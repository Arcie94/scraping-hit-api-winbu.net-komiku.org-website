package cache

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache is a Redis-based cache implementation
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache() *RedisCache {
	// Get Redis config from environment
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "", // no password
		DB:       0,  // default DB
	})

	return &RedisCache{
		client: client,
		ctx:    context.Background(),
	}
}

// Set stores a value in Redis with TTL
func (r *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	// Serialize value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, key, data, ttl).Err()
}

// Get retrieves a value from Redis
func (r *RedisCache) Get(key string) (interface{}, bool) {
	data, err := r.client.Get(r.ctx, key).Bytes()
	if err == redis.Nil {
		return nil, false // Key doesn't exist
	}
	if err != nil {
		return nil, false
	}

	// Deserialize JSON
	var value interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, false
	}

	return value, true
}

// Delete removes a key from Redis
func (r *RedisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Clear removes all keys (use with caution!)
func (r *RedisCache) Clear() error {
	return r.client.FlushDB(r.ctx).Err()
}

// Ping tests Redis connectivity
func (r *RedisCache) Ping() error {
	return r.client.Ping(r.ctx).Err()
}

// Close closes the Redis connection
func (r *RedisCache) Close() error {
	return r.client.Close()
}
