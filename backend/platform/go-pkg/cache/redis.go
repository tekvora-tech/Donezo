package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient wrapper Redis untuk cache
type RedisClient struct {
	client *redis.Client
}

// NewRedis buat instance cache
func NewRedis(addr string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisClient{client: client}
}

// Set simpan data ke cache
func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

// Get ambil data dari cache
func (c *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// Delete hapus key dari cache
func (c *RedisClient) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists cek key ada di cache
func (c *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	return n > 0, err
}

// FlushDB hapus semua data (hati-hati!)
func (c *RedisClient) FlushDB(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

// Close tutup koneksi
func (c *RedisClient) Close() error {
	return c.client.Close()
}

// Ping cek koneksi
func (c *RedisClient) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}