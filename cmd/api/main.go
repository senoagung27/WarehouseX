package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/senoagung27/warehousex/internal/config"
	"github.com/senoagung27/warehousex/internal/controller"
	"github.com/senoagung27/warehousex/internal/infrastructure"
	"github.com/senoagung27/warehousex/internal/repository"
	"github.com/senoagung27/warehousex/internal/router"
	"github.com/senoagung27/warehousex/internal/service"
	"go.uber.org/zap"
)

func main() {
	// ========== Logger ==========
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// ========== Config ==========
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// ========== Database ==========
	db, err := infrastructure.NewDatabase(&cfg.Database, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	logger.Info("Database connected")

	// ========== Redis ==========
	redisClient, err := infrastructure.NewRedisClient(&cfg.Redis, logger)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Info("Redis connected")

	// ========== Repositories ==========
	userRepo := repository.NewUserRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	requestRepo := repository.NewRequestRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)

	// ========== Services ==========
	authService := service.NewAuthService(userRepo, cfg.JWT, logger)
	inventoryService := service.NewInventoryService(inventoryRepo, auditLogRepo, logger)
	requestService := service.NewRequestService(requestRepo, inventoryRepo, auditLogRepo, redisClient, db, logger)
	auditService := service.NewAuditService(auditLogRepo, logger)

	// ========== Controllers ==========
	authController := controller.NewAuthController(authService)
	inventoryController := controller.NewInventoryController(inventoryService)
	requestController := controller.NewRequestController(requestService)
	auditController := controller.NewAuditController(auditService)

	// ========== Router ==========
	r := router.NewRouter(
		authController,
		inventoryController,
		requestController,
		auditController,
		cfg.JWT.Secret,
		cfg.Server.GinMode,
	)

	// ========== Server ==========
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      r.Engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Starting WarehouseX API server",
			zap.String("port", cfg.Server.Port),
			zap.String("mode", cfg.Server.GinMode),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// ========== Graceful Shutdown ==========
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	// Close Redis
	if err := redisClient.Client.Close(); err != nil {
		logger.Error("Failed to close Redis connection", zap.Error(err))
	}

	// Close DB
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}

	logger.Info("Server exited gracefully")
}
