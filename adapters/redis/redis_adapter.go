package redisadapter

import (
	"context"
	"errors"

	"github.com/oxmies/oxmies/adapters"
)

// RedisAdapter is a minimal placeholder implementing adapters.DBAdapter.
// Replace with a real implementation when integrating Redis.
type RedisAdapter struct{}

func NewRedisAdapter() adapters.DBAdapter {
	return &RedisAdapter{}
}

func (r *RedisAdapter) Insert(ctx context.Context, model interface{}) error {
	return errors.New("redis adapter: Insert not implemented")
}
func (r *RedisAdapter) Update(ctx context.Context, model interface{}) error {
	return errors.New("redis adapter: Update not implemented")
}
func (r *RedisAdapter) FindByID(ctx context.Context, model interface{}, id any) error {
	return errors.New("redis adapter: FindByID not implemented")
}
func (r *RedisAdapter) Delete(ctx context.Context, model interface{}) error {
	return errors.New("redis adapter: Delete not implemented")
}

// AdapterType returns the adapter type for this adapter
func (r *RedisAdapter) AdapterType() adapters.AdapterType {
	return adapters.Redis
}
