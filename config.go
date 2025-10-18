package ormies

import (
	"fmt"
	"net/url"
	"strings"
)

// -----------------------------
// SQL Configuration
// -----------------------------
type SQLConfig struct {
	Driver   string            // "postgres", "mysql", etc.
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
	SSLMode  string            // Optional: Postgres only
	Params   map[string]string // Extra query parameters
	Debug    bool              // Enable SQL logging
}

// DSN builds the connection string automatically
func (c SQLConfig) DSN() string {
	switch c.Driver {
	case "postgres":
		ssl := c.SSLMode
		if ssl == "" {
			ssl = "disable"
		}
		params := url.Values{}
		for k, v := range c.Params {
			params.Set(k, v)
		}
		paramStr := ""
		if len(params) > 0 {
			paramStr = "?" + params.Encode()
		}
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s%s",
			c.User, c.Password, c.Host, c.Port, c.DBName, paramStr,
		)

	case "mysql":
		// Format: user:pass@tcp(host:port)/dbname?param=value
		params := []string{}
		for k, v := range c.Params {
			params = append(params, fmt.Sprintf("%s=%s", k, v))
		}
		paramStr := ""
		if len(params) > 0 {
			paramStr = "?" + strings.Join(params, "&")
		}
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s%s",
			c.User, c.Password, c.Host, c.Port, c.DBName, paramStr,
		)

	default:
		panic("Unsupported SQL driver: " + c.Driver)
	}
}

// -----------------------------
// MongoDB Configuration
// -----------------------------
type MongoConfig struct {
	URI      string // e.g., "mongodb://localhost:27017"
	Database string // Database name
}

// -----------------------------
// Redis Configuration
// -----------------------------
type RedisConfig struct {
	URI string // e.g., "redis://localhost:6379"
}
