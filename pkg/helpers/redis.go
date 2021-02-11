package helpers

import (
	"context"
	"fmt"

	"redismigrator/pkg/model"

	redis "github.com/go-redis/redis/v8"
)

// Ctx ...
var Ctx = context.Background()

// GetDBSize returns the keys length
func GetDBSize(client *redis.Client, keyFilter string) (int64, error) {
	res := client.DBSize(Ctx)
	return res.Val(), res.Err()
}

// NewRedisClient creates a new client
func NewRedisClient(cfg model.RedisServerConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port),
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DBIndex,
	})
}
