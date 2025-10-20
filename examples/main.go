package main

import (
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/oxmies/oxmies"
)

func main() {

	// Configuration for initializing Oxmies with a SQL adapter.
	cfg := map[string]any{
		"db": oxmies.SQLConfig{
			Driver:   "postgres",
			User:     "user",
			Password: "pass",
			Host:     "localhost",
			Port:     5432,
			DBName:   "testdb",
			SSLMode:  "disable",
			Params:   map[string]string{"search_path": "public"},
			Debug:    true,
			OxmiesDbConfig: oxmies.OxmiesDbConfig{
				Models: []any{&User{}},
			},
		},
	}
	oxmies.Initialize(cfg)

}
