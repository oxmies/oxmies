package sqladapter

import (
	"fmt"
	"reflect"
	"strings"
)

// -----------------------------
// Helpers for building SQL queries
// -----------------------------

// buildInsertQuery generates an INSERT query from a struct
func buildInsertQuery(model interface{}) (string, []interface{}) {
	val := reflect.ValueOf(model).Elem()
	typ := val.Type()

	var columns []string
	var placeholders []string
	var args []interface{}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("orm")
		if strings.Contains(tag, "primary_key") {
			continue // skip primary key (auto-increment)
		}

		column := getColumnName(field)
		columns = append(columns, column)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)+1))
		args = append(args, val.Field(i).Interface())
	}

	table := getTableName(model)
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	return query, args
}

// buildUpdateQuery generates an UPDATE query from a struct
func buildUpdateQuery(model interface{}) (string, []interface{}) {
	val := reflect.ValueOf(model).Elem()
	typ := val.Type()

	var sets []string
	var args []interface{}
	var id interface{}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("orm")
		column := getColumnName(field)
		value := val.Field(i).Interface()

		if strings.Contains(tag, "primary_key") {
			id = value
			continue
		}

		sets = append(sets, fmt.Sprintf("%s=$%d", column, len(args)+1))
		args = append(args, value)
	}

	args = append(args, id)
	table := getTableName(model)
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id=$%d",
		table,
		strings.Join(sets, ", "),
		len(args),
	)

	return query, args
}

// buildSelectByIDQuery generates a SELECT query to find a record by ID
func buildSelectByIDQuery(model interface{}) string {
	table := getTableName(model)
	return fmt.Sprintf("SELECT * FROM %s WHERE id=$1", table)
}

// buildDeleteQuery generates a DELETE query from a struct
func buildDeleteQuery(model interface{}) (string, []interface{}) {
	val := reflect.ValueOf(model).Elem()
	typ := val.Type()
	var id interface{}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("orm")
		if strings.Contains(tag, "primary_key") {
			id = val.Field(i).Interface()
			break
		}
	}

	table := getTableName(model)
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", table)
	return query, []interface{}{id}
}

// scanRow maps sql.Row to struct fields (simple version)
func scanRow(row interface{}, model interface{}) error {
	// This needs sql.Row type; leave for SQLAdapter implementation
	return nil
}

// -----------------------------
// Reflection Helpers
// -----------------------------

// getColumnName reads `orm:"column:name"` tag or defaults to field name
func getColumnName(field reflect.StructField) string {
	tag := field.Tag.Get("orm")
	for _, part := range strings.Split(tag, ",") {
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}
	return field.Name
}

// getTableName defaults to struct type name
func getTableName(model interface{}) string {
	typ := reflect.TypeOf(model).Elem()
	return strings.ToLower(typ.Name())
}
