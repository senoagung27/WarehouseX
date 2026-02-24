package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/senoagung27/warehousex/internal/controller"
	"github.com/senoagung27/warehousex/internal/middleware"
)

type Router struct {
	Engine              *gin.Engine
	authController      *controller.AuthController
	inventoryController *controller.InventoryController
	requestController   *controller.RequestController
	auditController     *controller.AuditController
	jwtSecret           string
}

func NewRouter(
	authController *controller.AuthController,
	inventoryController *controller.InventoryController,
	requestController *controller.RequestController,
	auditController *controller.AuditController,
	jwtSecret string,
	ginMode string,
) *Router {
	gin.SetMode(ginMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	r := &Router{
		Engine:              engine,
		authController:      authController,
		inventoryController: inventoryController,
		requestController:   requestController,
		auditController:     auditController,
		jwtSecret:           jwtSecret,
	}

	r.setupRoutes()
	return r
}

func (r *Router) setupRoutes() {
	// Health check
	r.Engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "WarehouseX",
			"version": "1.0.0",
		})
	})

	// ========== API v1 ==========
	v1 := r.Engine.Group("/api/v1")

	// --- Auth (public) ---
	auth := v1.Group("/auth")
	{
		auth.POST("/register", r.authController.Register)
		auth.POST("/login", r.authController.Login)
	}

	// --- Protected routes ---
	protected := v1.Group("")
	protected.Use(middleware.JWTAuth(r.jwtSecret))

	// --- Inventory ---
	inventory := protected.Group("/inventory")
	{
		inventory.GET("", r.inventoryController.GetAll)
		inventory.GET("/:id", r.inventoryController.GetByID)
		inventory.POST("", middleware.RequireRole("admin"), r.inventoryController.Create)
		inventory.PUT("/:id", middleware.RequireRole("admin"), r.inventoryController.Update)
	}

	// --- Requests (Inbound / Outbound) ---
	requests := protected.Group("/requests")
	{
		requests.GET("", r.requestController.GetAll)
		requests.GET("/:id", r.requestController.GetByID)
		requests.POST("/inbound", middleware.RequireRole("staff"), r.requestController.CreateInbound)
		requests.POST("/outbound", middleware.RequireRole("staff"), r.requestController.CreateOutbound)
		requests.PUT("/:id/approve", middleware.RequireRoles("supervisor", "admin"), r.requestController.Approve)
		requests.PUT("/:id/reject", middleware.RequireRoles("supervisor", "admin"), r.requestController.Reject)
	}

	// --- Audit Logs ---
	auditLogs := protected.Group("/audit-logs")
	{
		auditLogs.GET("", middleware.RequireRole("auditor"), r.auditController.GetAll)
	}
}
