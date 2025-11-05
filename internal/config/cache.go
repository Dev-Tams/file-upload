package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Cache *redis.Client

func InitRedis() {
	ctx := context.Background()



	Cache = redis.NewClient(&redis.Options{
		Addr:     Config.RedisAddr,
		Password: Config.RedisPassword,
		DB:       Config.RedisDB,
		MaxRetries: 2,
	})

	if err := Cache.Ping(ctx).Err(); err != nil {
		log.Printf(" Redis not available: %v", err)
		Cache = nil // fallback to no-cache mode
	}
	log.Printf("✅ Connected to Redis at %s [DB=%d]", Config.RedisAddr, Config.RedisDB)
}
