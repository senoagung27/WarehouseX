package repository

import (
	"fmt"

	"github.com/google/uuid"
	domainRepo "github.com/senoagung27/warehousex/internal/domain/repository"
	"github.com/senoagung27/warehousex/internal/model"
	"gorm.io/gorm"
)

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) domainRepo.AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(log *model.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *auditLogRepository) CreateWithTx(tx interface{}, log *model.AuditLog) error {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}
	return gormTx.Create(log).Error
}

func (r *auditLogRepository) FindAll(page, limit int, entityName string, entityID *uuid.UUID) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	query := r.db.Model(&model.AuditLog{})

	if entityName != "" {
		query = query.Where("entity = ?", entityName)
	}
	if entityID != nil {
		query = query.Where("entity_id = ?", *entityID)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	if err := query.Preload("User").
		Offset(offset).Limit(limit).Order("created_at DESC").
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
