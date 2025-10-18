package ormies

import (
	"database/sql"
	"fmt"
	adapter "github.com/ormies/ormies/adapters"
	adaptersql "github.com/ormies/ormies/adapters/sql"
)

var adapters = make(map[string]adapter.DBAdapter)

// RegisterAdapter registers a named adapter
func RegisterAdapter(name string, adapter adapter.DBAdapter) {
	adapters[name] = adapter
}

// GetAdapter returns an adapter by name
func GetAdapter(name string) adapter.DBAdapter {
	adapter, ok := adapters[name]
	if !ok {
		panic(fmt.Sprintf("❌ Ormies: Adapter '%s' not found. Did you register it?", name))
	}
	return adapter
}

func InitSQL(cfg SQLConfig, name string) {
	dsn := cfg.DSN()
	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		panic(err)
	}
	adapter := adaptersql.NewSQLAdapter(db, cfg.Debug)
	RegisterAdapter(name, adapter)
	if cfg.Debug {
		fmt.Printf("✅ SQL Adapter '%s' initialized\n", name)
	}
}