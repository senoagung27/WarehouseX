package service

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/senoagung27/warehousex/internal/domain/repository"
	"github.com/senoagung27/warehousex/internal/dto"
	"github.com/senoagung27/warehousex/internal/model"
	"go.uber.org/zap"
)

var _ InventoryServiceInterface = (*InventoryService)(nil)

type InventoryService struct {
	inventoryRepo repository.InventoryRepository
	auditRepo     repository.AuditLogRepository
	log           *zap.Logger
}

func NewInventoryService(
	inventoryRepo repository.InventoryRepository,
	auditRepo repository.AuditLogRepository,
	log *zap.Logger,
) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
		auditRepo:     auditRepo,
		log:           log,
	}
}

func (s *InventoryService) Create(input dto.CreateInventoryInput, userID uuid.UUID) (*model.Inventory, error) {
	item := &model.Inventory{
		ID:       uuid.New(),
		ItemName: input.ItemName,
		SKU:      input.SKU,
		Quantity: input.Quantity,
		Unit:     input.Unit,
		Version:  1,
	}

	if err := s.inventoryRepo.Create(item); err != nil {
		return nil, fmt.Errorf("failed to create inventory item: %w", err)
	}

	afterJSON, _ := json.Marshal(item)
	_ = s.auditRepo.Create(&model.AuditLog{
		ID:         uuid.New(),
		Entity:     "inventory",
		EntityID:   item.ID,
		Action:     "CREATE",
		UserID:     userID,
		AfterValue: afterJSON,
	})

	s.log.Info("Inventory item created",
		zap.String("item_id", item.ID.String()),
		zap.String("item_name", item.ItemName),
	)

	return item, nil
}

func (s *InventoryService) GetByID(id uuid.UUID) (*model.Inventory, error) {
	return s.inventoryRepo.FindByID(id)
}

func (s *InventoryService) GetAll(page, limit int) ([]model.Inventory, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.inventoryRepo.FindAll(page, limit)
}

func (s *InventoryService) Update(id uuid.UUID, input dto.UpdateInventoryInput, userID uuid.UUID) (*model.Inventory, error) {
	item, err := s.inventoryRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("inventory item not found: %w", err)
	}

	beforeJSON, _ := json.Marshal(item)

	if input.ItemName != "" {
		item.ItemName = input.ItemName
	}
	if input.SKU != "" {
		item.SKU = input.SKU
	}
	if input.Unit != "" {
		item.Unit = input.Unit
	}

	if err := s.inventoryRepo.Update(item); err != nil {
		return nil, fmt.Errorf("failed to update inventory item: %w", err)
	}

	afterJSON, _ := json.Marshal(item)
	_ = s.auditRepo.Create(&model.AuditLog{
		ID:          uuid.New(),
		Entity:      "inventory",
		EntityID:    item.ID,
		Action:      "UPDATE",
		UserID:      userID,
		BeforeValue: beforeJSON,
		AfterValue:  afterJSON,
	})

	s.log.Info("Inventory item updated",
		zap.String("item_id", item.ID.String()),
	)

	return item, nil
}
