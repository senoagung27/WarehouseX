package service

import (
	"github.com/google/uuid"
	"github.com/senoagung27/warehousex/internal/domain/repository"
	"github.com/senoagung27/warehousex/internal/model"
	"go.uber.org/zap"
)

var _ AuditServiceInterface = (*AuditService)(nil)

type AuditService struct {
	auditRepo repository.AuditLogRepository
	log       *zap.Logger
}

func NewAuditService(auditRepo repository.AuditLogRepository, log *zap.Logger) *AuditService {
	return &AuditService{
		auditRepo: auditRepo,
		log:       log,
	}
}

func (s *AuditService) GetAll(page, limit int, entityName string, entityID *uuid.UUID) ([]model.AuditLog, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	return s.auditRepo.FindAll(page, limit, entityName, entityID)
}
