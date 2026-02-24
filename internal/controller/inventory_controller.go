package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/senoagung27/warehousex/internal/dto"
	"github.com/senoagung27/warehousex/internal/middleware"
	"github.com/senoagung27/warehousex/internal/service"
)

type InventoryController struct {
	inventoryService service.InventoryServiceInterface
}

func NewInventoryController(inventoryService service.InventoryServiceInterface) *InventoryController {
	return &InventoryController{inventoryService: inventoryService}
}

// Create godoc
// @Summary Create a new inventory item
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateInventoryInput true "Create Inventory Input"
// @Success 201 {object} model.Inventory
// @Router /api/v1/inventory [post]
func (ctrl *InventoryController) Create(c *gin.Context) {
	var input dto.CreateInventoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	item, err := ctrl.inventoryService.Create(input, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "inventory item created",
		"data":    item,
	})
}

// GetAll godoc
// @Summary List all inventory items
// @Tags Inventory
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {array} model.Inventory
// @Router /api/v1/inventory [get]
func (ctrl *InventoryController) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	items, total, err := ctrl.inventoryService.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetByID godoc
// @Summary Get inventory item by ID
// @Tags Inventory
// @Security BearerAuth
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} model.Inventory
// @Router /api/v1/inventory/{id} [get]
func (ctrl *InventoryController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	item, err := ctrl.inventoryService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// Update godoc
// @Summary Update an inventory item
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Item ID"
// @Param input body dto.UpdateInventoryInput true "Update Inventory Input"
// @Success 200 {object} model.Inventory
// @Router /api/v1/inventory/{id} [put]
func (ctrl *InventoryController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	var input dto.UpdateInventoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	item, err := ctrl.inventoryService.Update(id, input, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "inventory item updated",
		"data":    item,
	})
}
