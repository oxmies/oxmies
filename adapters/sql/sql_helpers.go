package sqladapter

import (
	"database/sql"
	"errors"
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
	// Accept either *sql.Rows (from QueryContext) or *sql.Row
	rv := reflect.ValueOf(model)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("scanRow: model must be a non-nil pointer to struct")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("scanRow: model must point to a struct")
	}

	// We'll use sql.Rows to get column names and values
	var rows *sql.Rows
	switch r := row.(type) {
	case *sql.Rows:
		rows = r
	case *sql.Row:
		// sql.Row doesn't expose Columns; convert by querying a single row via Rows
		// but sql.Row can't be type asserted easily; fallthrough to error
		return errors.New("scanRow: *sql.Row is not supported; use QueryContext to pass *sql.Rows")
	default:
		return errors.New("scanRow: unsupported row type")
	}

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	// Prepare a slice of interfaces to scan into
	vals := make([]interface{}, len(cols))
	for i := range vals {
		var v interface{}
		vals[i] = &v
	}

	if err := rows.Scan(vals...); err != nil {
		return err
	}

	// Map column values to struct fields by column tag or field name
	colToVal := make(map[string]interface{})
	for i, c := range cols {
		// dereference scanned pointer
		vptr := vals[i].(*interface{})
		colToVal[c] = *vptr
	}

	typ := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := typ.Field(i)
		// ignore unexported fields
		if field.PkgPath != "" {
			continue
		}
		col := getColumnName(field)
		// try lowercase and original name variants
		if val, ok := colToVal[col]; ok {
			fv := rv.Field(i)
			if fv.CanSet() {
				v := reflect.ValueOf(val)
				if v.IsValid() {
					// handle nil database NULL
					if v.Kind() == reflect.Interface && v.IsNil() {
						// leave zero value
						continue
					}
					// attempt conversion
					if v.Type().AssignableTo(fv.Type()) {
						fv.Set(v)
					} else if v.Type().ConvertibleTo(fv.Type()) {
						fv.Set(v.Convert(fv.Type()))
					} else {
						// best-effort: try set via string when possible
						if s, ok := val.(string); ok {
							if fv.Kind() == reflect.String {
								fv.SetString(s)
							}
						}
					}
				}
			}
		} else if val, ok := colToVal[strings.ToLower(col)]; ok {
			fv := rv.Field(i)
			if fv.CanSet() {
				v := reflect.ValueOf(val)
				if v.IsValid() && v.Type().AssignableTo(fv.Type()) {
					fv.Set(v)
				} else if v.IsValid() && v.Type().ConvertibleTo(fv.Type()) {
					fv.Set(v.Convert(fv.Type()))
				}
			}
		}
	}

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
