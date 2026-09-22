package services

import (
	"context"

	"github.com/redis/go-redis/v9"
	"lumi.yellowlabs.space/internal/config"
)

var Redis *redis.Client

func InitRedis() error {
	Redis = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.Redis.URL,
		Password: config.AppConfig.Redis.Password,
	})

	return Redis.Ping(context.Background()).Err()
}

func CloseRedis() error {
	if Redis != nil {
		return Redis.Close()
	}
	return nil
}
