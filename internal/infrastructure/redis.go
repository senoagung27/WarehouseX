package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/senoagung27/warehousex/internal/config"
	"go.uber.org/zap"
)

type RedisClient struct {
	Client *redis.Client
	log    *zap.Logger
}

func NewRedisClient(cfg *config.RedisConfig, log *zap.Logger) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info("Redis connected successfully",
		zap.String("addr", cfg.Addr()),
	)

	return &RedisClient{Client: client, log: log}, nil
}

// AcquireLock acquires a distributed lock for a given item
// SET lock:item:<id> <value> NX EX 10
func (r *RedisClient) AcquireLock(ctx context.Context, itemID uuid.UUID) (string, error) {
	lockKey := fmt.Sprintf("lock:item:%s", itemID.String())
	lockValue := uuid.New().String()

	ok, err := r.Client.SetNX(ctx, lockKey, lockValue, 10*time.Second).Result()
	if err != nil {
		return "", fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !ok {
		return "", fmt.Errorf("lock already held for item %s", itemID.String())
	}

	r.log.Debug("Lock acquired",
		zap.String("key", lockKey),
		zap.String("value", lockValue),
	)

	return lockValue, nil
}

// ReleaseLock releases the distributed lock (only if we own it)
func (r *RedisClient) ReleaseLock(ctx context.Context, itemID uuid.UUID, lockValue string) error {
	lockKey := fmt.Sprintf("lock:item:%s", itemID.String())

	// Lua script to ensure atomic check-and-delete
	script := redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`)

	result, err := script.Run(ctx, r.Client, []string{lockKey}, lockValue).Int64()
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}
	if result == 0 {
		r.log.Warn("Lock was not held or expired",
			zap.String("key", lockKey),
		)
	}

	return nil
}
