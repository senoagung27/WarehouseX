package infrastructure

import (
	"fmt"

	"github.com/senoagung27/warehousex/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(cfg *config.DatabaseConfig, log *zap.Logger) (*gorm.DB, error) {
	dsn := cfg.DSN()

	gormConfig := &gorm.Config{}
	if cfg.SSLMode != "disable" {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	log.Info("Database connected successfully",
		zap.String("host", cfg.Host),
		zap.String("dbname", cfg.DBName),
	)

	return db, nil
}
