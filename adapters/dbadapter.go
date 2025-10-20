package adapters

import (
	"context"
)

type DBAdapter interface {
	Insert(ctx context.Context, model interface{}) error
	Update(ctx context.Context, model interface{}) error
	FindByID(ctx context.Context, model interface{}, id any) error
	Delete(ctx context.Context, model interface{}) error
	AdapterType() AdapterType
}

type AdapterType string

const (
	SQL   AdapterType = "sql"
	Redis AdapterType = "redis"
)

