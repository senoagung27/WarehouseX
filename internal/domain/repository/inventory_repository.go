package repository

import (
	"github.com/google/uuid"
	"github.com/senoagung27/warehousex/internal/model"
)

type InventoryRepository interface {
	Create(item *model.Inventory) error
	FindByID(id uuid.UUID) (*model.Inventory, error)
	FindAll(page, limit int) ([]model.Inventory, int64, error)
	Update(item *model.Inventory) error
	// FindByIDForUpdate uses SELECT FOR UPDATE row-level locking
	FindByIDForUpdate(tx interface{}, id uuid.UUID) (*model.Inventory, error)
	// UpdateWithTx updates within an existing transaction
	UpdateWithTx(tx interface{}, item *model.Inventory) error
}
