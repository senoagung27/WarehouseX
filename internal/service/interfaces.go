package service

import (
	"github.com/google/uuid"
	"github.com/senoagung27/warehousex/internal/dto"
	"github.com/senoagung27/warehousex/internal/model"
)

// AuthServiceInterface defines the contract for authentication operations
type AuthServiceInterface interface {
	Register(input dto.RegisterInput) (*dto.AuthResponse, error)
	Login(input dto.LoginInput) (*dto.AuthResponse, error)
}

// InventoryServiceInterface defines the contract for inventory operations
type InventoryServiceInterface interface {
	Create(input dto.CreateInventoryInput, userID uuid.UUID) (*model.Inventory, error)
	GetByID(id uuid.UUID) (*model.Inventory, error)
	GetAll(page, limit int) ([]model.Inventory, int64, error)
	Update(id uuid.UUID, input dto.UpdateInventoryInput, userID uuid.UUID) (*model.Inventory, error)
}

// RequestServiceInterface defines the contract for request operations
type RequestServiceInterface interface {
	CreateInbound(input dto.CreateRequestInput, userID uuid.UUID) (*model.Request, error)
	CreateOutbound(input dto.CreateRequestInput, userID uuid.UUID) (*model.Request, error)
	ApproveRequest(requestID uuid.UUID, approverID uuid.UUID, approverRole string) (*model.Request, error)
	RejectRequest(requestID uuid.UUID, approverID uuid.UUID, approverRole string) (*model.Request, error)
	GetByID(id uuid.UUID) (*model.Request, error)
	GetAll(page, limit int, reqType, status string) ([]model.Request, int64, error)
}

// AuditServiceInterface defines the contract for audit log operations
type AuditServiceInterface interface {
	GetAll(page, limit int, entityName string, entityID *uuid.UUID) ([]model.AuditLog, int64, error)
}
