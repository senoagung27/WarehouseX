package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/senoagung27/warehousex/internal/domain/repository"
	"github.com/senoagung27/warehousex/internal/dto"
	"github.com/senoagung27/warehousex/internal/infrastructure"
	"github.com/senoagung27/warehousex/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var _ RequestServiceInterface = (*RequestService)(nil)

type RequestService struct {
	requestRepo   repository.RequestRepository
	inventoryRepo repository.InventoryRepository
	auditRepo     repository.AuditLogRepository
	redisClient   *infrastructure.RedisClient
	db            *gorm.DB
	log           *zap.Logger
}

func NewRequestService(
	requestRepo repository.RequestRepository,
	inventoryRepo repository.InventoryRepository,
	auditRepo repository.AuditLogRepository,
	redisClient *infrastructure.RedisClient,
	db *gorm.DB,
	log *zap.Logger,
) *RequestService {
	return &RequestService{
		requestRepo:   requestRepo,
		inventoryRepo: inventoryRepo,
		auditRepo:     auditRepo,
		redisClient:   redisClient,
		db:            db,
		log:           log,
	}
}

func (s *RequestService) CreateInbound(input dto.CreateRequestInput, userID uuid.UUID) (*model.Request, error) {
	itemID, err := uuid.Parse(input.ItemID)
	if err != nil {
		return nil, errors.New("invalid item ID")
	}

	_, err = s.inventoryRepo.FindByID(itemID)
	if err != nil {
		return nil, errors.New("inventory item not found")
	}

	req := &model.Request{
		ID:        uuid.New(),
		Type:      model.RequestTypeInbound,
		Status:    model.StatusPending,
		ItemID:    itemID,
		Quantity:  input.Quantity,
		Notes:     input.Notes,
		CreatedBy: userID,
	}

	if err := s.requestRepo.Create(req); err != nil {
		return nil, fmt.Errorf("failed to create inbound request: %w", err)
	}

	afterJSON, _ := json.Marshal(req)
	_ = s.auditRepo.Create(&model.AuditLog{
		ID:         uuid.New(),
		Entity:     "request",
		EntityID:   req.ID,
		Action:     "CREATE_INBOUND",
		UserID:     userID,
		AfterValue: afterJSON,
	})

	s.log.Info("Inbound request created",
		zap.String("request_id", req.ID.String()),
		zap.String("item_id", itemID.String()),
		zap.Int("quantity", input.Quantity),
	)

	return req, nil
}

func (s *RequestService) CreateOutbound(input dto.CreateRequestInput, userID uuid.UUID) (*model.Request, error) {
	itemID, err := uuid.Parse(input.ItemID)
	if err != nil {
		return nil, errors.New("invalid item ID")
	}

	item, err := s.inventoryRepo.FindByID(itemID)
	if err != nil {
		return nil, errors.New("inventory item not found")
	}

	if item.Quantity < input.Quantity {
		return nil, fmt.Errorf("insufficient stock: available %d, requested %d", item.Quantity, input.Quantity)
	}

	req := &model.Request{
		ID:        uuid.New(),
		Type:      model.RequestTypeOutbound,
		Status:    model.StatusPending,
		ItemID:    itemID,
		Quantity:  input.Quantity,
		Notes:     input.Notes,
		CreatedBy: userID,
	}

	if err := s.requestRepo.Create(req); err != nil {
		return nil, fmt.Errorf("failed to create outbound request: %w", err)
	}

	afterJSON, _ := json.Marshal(req)
	_ = s.auditRepo.Create(&model.AuditLog{
		ID:         uuid.New(),
		Entity:     "request",
		EntityID:   req.ID,
		Action:     "CREATE_OUTBOUND",
		UserID:     userID,
		AfterValue: afterJSON,
	})

	s.log.Info("Outbound request created",
		zap.String("request_id", req.ID.String()),
		zap.String("item_id", itemID.String()),
		zap.Int("quantity", input.Quantity),
	)

	return req, nil
}

func (s *RequestService) ApproveRequest(requestID uuid.UUID, approverID uuid.UUID, approverRole string) (*model.Request, error) {
	if !model.CanApprove(approverRole) {
		return nil, errors.New("insufficient permissions to approve requests")
	}

	req, err := s.requestRepo.FindByID(requestID)
	if err != nil {
		return nil, errors.New("request not found")
	}

	if !model.ValidTransition(req.Status, model.StatusApproved) {
		return nil, fmt.Errorf("cannot approve request with status: %s", req.Status)
	}

	if req.CreatedBy == approverID {
		return nil, errors.New("cannot approve your own request")
	}

	ctx := context.Background()

	if req.Type == model.RequestTypeOutbound {
		return s.processOutboundApproval(ctx, req, approverID)
	}

	return s.processInboundApproval(req, approverID)
}

func (s *RequestService) processInboundApproval(req *model.Request, approverID uuid.UUID) (*model.Request, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.inventoryRepo.FindByIDForUpdate(tx, req.ItemID)
		if err != nil {
			return fmt.Errorf("inventory item not found: %w", err)
		}

		beforeJSON, _ := json.Marshal(item)

		item.Quantity += req.Quantity
		item.Version++

		if err := s.inventoryRepo.UpdateWithTx(tx, item); err != nil {
			return fmt.Errorf("failed to update inventory: %w", err)
		}

		req.Status = model.StatusCompleted
		req.ApprovedBy = &approverID

		if err := s.requestRepo.UpdateWithTx(tx, req); err != nil {
			return fmt.Errorf("failed to update request: %w", err)
		}

		afterJSON, _ := json.Marshal(item)
		if err := s.auditRepo.CreateWithTx(tx, &model.AuditLog{
			ID:          uuid.New(),
			Entity:      "inventory",
			EntityID:    item.ID,
			Action:      "INBOUND_APPROVED",
			UserID:      approverID,
			BeforeValue: beforeJSON,
			AfterValue:  afterJSON,
		}); err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.log.Info("Inbound request approved",
		zap.String("request_id", req.ID.String()),
		zap.String("approver_id", approverID.String()),
	)

	return s.requestRepo.FindByID(req.ID)
}

func (s *RequestService) processOutboundApproval(ctx context.Context, req *model.Request, approverID uuid.UUID) (*model.Request, error) {
	lockValue, err := s.redisClient.AcquireLock(ctx, req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("lock conflict: %w", err)
	}
	defer func() {
		_ = s.redisClient.ReleaseLock(ctx, req.ItemID, lockValue)
	}()

	err = s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.inventoryRepo.FindByIDForUpdate(tx, req.ItemID)
		if err != nil {
			return fmt.Errorf("inventory item not found: %w", err)
		}

		if item.Quantity < req.Quantity {
			return fmt.Errorf("insufficient stock: available %d, requested %d", item.Quantity, req.Quantity)
		}

		beforeJSON, _ := json.Marshal(item)

		item.Quantity -= req.Quantity
		item.Version++

		if err := s.inventoryRepo.UpdateWithTx(tx, item); err != nil {
			return fmt.Errorf("failed to update inventory: %w", err)
		}

		req.Status = model.StatusCompleted
		req.ApprovedBy = &approverID

		if err := s.requestRepo.UpdateWithTx(tx, req); err != nil {
			return fmt.Errorf("failed to update request: %w", err)
		}

		afterJSON, _ := json.Marshal(item)
		if err := s.auditRepo.CreateWithTx(tx, &model.AuditLog{
			ID:          uuid.New(),
			Entity:      "inventory",
			EntityID:    item.ID,
			Action:      "OUTBOUND_APPROVED",
			UserID:      approverID,
			BeforeValue: beforeJSON,
			AfterValue:  afterJSON,
		}); err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.log.Info("Outbound request approved",
		zap.String("request_id", req.ID.String()),
		zap.String("approver_id", approverID.String()),
	)

	return s.requestRepo.FindByID(req.ID)
}

func (s *RequestService) RejectRequest(requestID uuid.UUID, approverID uuid.UUID, approverRole string) (*model.Request, error) {
	if !model.CanApprove(approverRole) {
		return nil, errors.New("insufficient permissions to reject requests")
	}

	req, err := s.requestRepo.FindByID(requestID)
	if err != nil {
		return nil, errors.New("request not found")
	}

	if !model.ValidTransition(req.Status, model.StatusRejected) {
		return nil, fmt.Errorf("cannot reject request with status: %s", req.Status)
	}

	req.Status = model.StatusRejected
	req.ApprovedBy = &approverID

	if err := s.requestRepo.Update(req); err != nil {
		return nil, fmt.Errorf("failed to reject request: %w", err)
	}

	afterJSON, _ := json.Marshal(req)
	_ = s.auditRepo.Create(&model.AuditLog{
		ID:         uuid.New(),
		Entity:     "request",
		EntityID:   req.ID,
		Action:     "REJECTED",
		UserID:     approverID,
		AfterValue: afterJSON,
	})

	s.log.Info("Request rejected",
		zap.String("request_id", req.ID.String()),
		zap.String("approver_id", approverID.String()),
	)

	return s.requestRepo.FindByID(req.ID)
}

func (s *RequestService) GetByID(id uuid.UUID) (*model.Request, error) {
	return s.requestRepo.FindByID(id)
}

func (s *RequestService) GetAll(page, limit int, reqType, status string) ([]model.Request, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	filters := make(map[string]interface{})
	if reqType != "" {
		filters["type"] = reqType
	}
	if status != "" {
		filters["status"] = status
	}

	return s.requestRepo.FindAll(page, limit, filters)
}
