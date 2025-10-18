package sqladapter

import (
	"context"
	"database/sql"
	"log"
	"github.com/ormies/ormies/adapters"
)

type SQLAdapter struct {
	DB    *sql.DB
	Debug bool
}

// Constructor
func NewSQLAdapter(db *sql.DB, debug bool) adapters.DBAdapter {
	return &SQLAdapter{
		DB:    db,
		Debug: debug,
	}
}

// ----------------- DBAdapter Methods -----------------
func (a *SQLAdapter) Insert(ctx context.Context, model interface{}) error {
	// Call helper to build query and args
	query, args := buildInsertQuery(model)
	if a.Debug {
		log.Println("SQL Insert:", query, args)
	}
	_, err := a.DB.ExecContext(ctx, query, args...)
	return err
}

func (a *SQLAdapter) Update(ctx context.Context, model interface{}) error {
	query, args := buildUpdateQuery(model)
	if a.Debug {
		log.Println("SQL Update:", query, args)
	}
	_, err := a.DB.ExecContext(ctx, query, args...)
	return err
}

func (a *SQLAdapter) FindByID(ctx context.Context, model interface{}, id any) error {
	query := buildSelectByIDQuery(model)
	row := a.DB.QueryRowContext(ctx, query, id)
	return scanRow(row, model)
}

func (a *SQLAdapter) Delete(ctx context.Context, model interface{}) error {
	query, args := buildDeleteQuery(model)
	if a.Debug {
		log.Println("SQL Delete:", query, args)
	}
	_, err := a.DB.ExecContext(ctx, query, args...)
	return err
}
