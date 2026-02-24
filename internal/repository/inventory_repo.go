package repository

import (
	"fmt"

	"github.com/google/uuid"
	domainRepo "github.com/senoagung27/warehousex/internal/domain/repository"
	"github.com/senoagung27/warehousex/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) domainRepo.InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) Create(item *model.Inventory) error {
	return r.db.Create(item).Error
}

func (r *inventoryRepository) FindByID(id uuid.UUID) (*model.Inventory, error) {
	var item model.Inventory
	if err := r.db.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *inventoryRepository) FindAll(page, limit int) ([]model.Inventory, int64, error) {
	var items []model.Inventory
	var total int64

	r.db.Model(&model.Inventory{}).Count(&total)

	offset := (page - 1) * limit
	if err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *inventoryRepository) Update(item *model.Inventory) error {
	return r.db.Save(item).Error
}

func (r *inventoryRepository) FindByIDForUpdate(tx interface{}, id uuid.UUID) (*model.Inventory, error) {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("invalid transaction type")
	}

	var item model.Inventory
	if err := gormTx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *inventoryRepository) UpdateWithTx(tx interface{}, item *model.Inventory) error {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}
	return gormTx.Save(item).Error
}
