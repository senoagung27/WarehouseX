package repository

import (
	"github.com/google/uuid"
	"github.com/senoagung27/warehousex/internal/model"
)

type AuditLogRepository interface {
	Create(log *model.AuditLog) error
	CreateWithTx(tx interface{}, log *model.AuditLog) error
	FindAll(page, limit int, entityName string, entityID *uuid.UUID) ([]model.AuditLog, int64, error)
}
