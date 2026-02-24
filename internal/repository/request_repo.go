package repository

import (
	"fmt"

	"github.com/google/uuid"
	domainRepo "github.com/senoagung27/warehousex/internal/domain/repository"
	"github.com/senoagung27/warehousex/internal/model"
	"gorm.io/gorm"
)

type requestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) domainRepo.RequestRepository {
	return &requestRepository{db: db}
}

func (r *requestRepository) Create(req *model.Request) error {
	return r.db.Create(req).Error
}

func (r *requestRepository) FindByID(id uuid.UUID) (*model.Request, error) {
	var req model.Request
	if err := r.db.Preload("Item").Preload("Creator").Preload("Approver").
		Where("id = ?", id).First(&req).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *requestRepository) FindAll(page, limit int, filters map[string]interface{}) ([]model.Request, int64, error) {
	var requests []model.Request
	var total int64

	query := r.db.Model(&model.Request{})

	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	if err := query.Preload("Item").Preload("Creator").Preload("Approver").
		Offset(offset).Limit(limit).Order("created_at DESC").
		Find(&requests).Error; err != nil {
		return nil, 0, err
	}
	return requests, total, nil
}

func (r *requestRepository) Update(req *model.Request) error {
	return r.db.Save(req).Error
}

func (r *requestRepository) FindByIDWithTx(tx interface{}, id uuid.UUID) (*model.Request, error) {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("invalid transaction type")
	}

	var req model.Request
	if err := gormTx.Where("id = ?", id).First(&req).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *requestRepository) UpdateWithTx(tx interface{}, req *model.Request) error {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}
	return gormTx.Save(req).Error
}
