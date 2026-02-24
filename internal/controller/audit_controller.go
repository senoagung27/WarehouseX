package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/senoagung27/warehousex/internal/service"
)

type AuditController struct {
	auditService service.AuditServiceInterface
}

func NewAuditController(auditService service.AuditServiceInterface) *AuditController {
	return &AuditController{auditService: auditService}
}

// GetAll godoc
// @Summary List audit logs
// @Tags Audit
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param entity query string false "Entity type filter"
// @Param entity_id query string false "Entity ID filter"
// @Success 200 {array} model.AuditLog
// @Router /api/v1/audit-logs [get]
func (ctrl *AuditController) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	entityName := c.Query("entity")
	entityIDStr := c.Query("entity_id")

	var entityID *uuid.UUID
	if entityIDStr != "" {
		id, err := uuid.Parse(entityIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id"})
			return
		}
		entityID = &id
	}

	logs, total, err := ctrl.auditService.GetAll(page, limit, entityName, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}
