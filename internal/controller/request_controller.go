package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/senoagung27/warehousex/internal/dto"
	"github.com/senoagung27/warehousex/internal/middleware"
	"github.com/senoagung27/warehousex/internal/service"
)

type RequestController struct {
	requestService service.RequestServiceInterface
}

func NewRequestController(requestService service.RequestServiceInterface) *RequestController {
	return &RequestController{requestService: requestService}
}

// CreateInbound godoc
// @Summary Create inbound request
// @Tags Requests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateRequestInput true "Create Inbound Input"
// @Success 201 {object} model.Request
// @Router /api/v1/requests/inbound [post]
func (ctrl *RequestController) CreateInbound(c *gin.Context) {
	var input dto.CreateRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	req, err := ctrl.requestService.CreateInbound(input, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "inbound request created",
		"data":    req,
	})
}

// CreateOutbound godoc
// @Summary Create outbound request
// @Tags Requests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateRequestInput true "Create Outbound Input"
// @Success 201 {object} model.Request
// @Router /api/v1/requests/outbound [post]
func (ctrl *RequestController) CreateOutbound(c *gin.Context) {
	var input dto.CreateRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	req, err := ctrl.requestService.CreateOutbound(input, userID)
	if err != nil {
		statusCode := http.StatusBadRequest
		if strings.Contains(err.Error(), "insufficient stock") {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "outbound request created",
		"data":    req,
	})
}

// GetAll godoc
// @Summary List all requests
// @Tags Requests
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param type query string false "Request type (INBOUND/OUTBOUND)"
// @Param status query string false "Request status"
// @Success 200 {array} model.Request
// @Router /api/v1/requests [get]
func (ctrl *RequestController) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	reqType := c.Query("type")
	status := c.Query("status")

	requests, total, err := ctrl.requestService.GetAll(page, limit, reqType, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  requests,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetByID godoc
// @Summary Get request by ID
// @Tags Requests
// @Security BearerAuth
// @Produce json
// @Param id path string true "Request ID"
// @Success 200 {object} model.Request
// @Router /api/v1/requests/{id} [get]
func (ctrl *RequestController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request ID"})
		return
	}

	req, err := ctrl.requestService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": req})
}

// Approve godoc
// @Summary Approve a request
// @Tags Requests
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Success 200 {object} model.Request
// @Router /api/v1/requests/{id}/approve [put]
func (ctrl *RequestController) Approve(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request ID"})
		return
	}

	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)

	req, err := ctrl.requestService.ApproveRequest(id, userID, userRole)
	if err != nil {
		statusCode := http.StatusBadRequest
		errMsg := err.Error()

		switch {
		case strings.Contains(errMsg, "insufficient permissions"):
			statusCode = http.StatusForbidden
		case strings.Contains(errMsg, "not found"):
			statusCode = http.StatusNotFound
		case strings.Contains(errMsg, "lock conflict"):
			statusCode = http.StatusConflict
		case strings.Contains(errMsg, "insufficient stock"):
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "request approved successfully",
		"data":    req,
	})
}

// Reject godoc
// @Summary Reject a request
// @Tags Requests
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Success 200 {object} model.Request
// @Router /api/v1/requests/{id}/reject [put]
func (ctrl *RequestController) Reject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request ID"})
		return
	}

	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)

	req, err := ctrl.requestService.RejectRequest(id, userID, userRole)
	if err != nil {
		statusCode := http.StatusBadRequest
		if strings.Contains(err.Error(), "insufficient permissions") {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "request rejected",
		"data":    req,
	})
}
