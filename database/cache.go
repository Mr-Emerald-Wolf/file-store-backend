package database

import (
	"context"
	"log"
	"time"

	"github.com/mr-emerald-wolf/21BCE0665_Backend/config"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

var RedisClient *RedisRepository

func NewRepository(redisConfig config.RedisConfig) {

	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.REDIS_HOST,
		Password: redisConfig.REDIS_PASSWORD,
		DB:       0,
	})

	// Verify the connection to Redis
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis Init Failed: %v", err.Error())
	}

	RedisClient = &RedisRepository{client}
}

func (r *RedisRepository) Set(key, value string, time time.Duration) error {
	ctx := context.Background()
	err := r.client.Set(ctx, key, value, time).Err()
	return err
}

func (r *RedisRepository) Get(key string) (string, error) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, key).Result()
	return val, err
}

func (r *RedisRepository) Delete(key string) error {
	ctx := context.Background()
	err := r.client.Del(ctx, key).Err()
	if err == redis.Nil {
		return nil
	}
	return err
}
