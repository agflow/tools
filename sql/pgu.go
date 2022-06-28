package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/agflow/tools/log"
	"github.com/agflow/tools/types"
)

// Select runs query on database with arguments and saves result on dest variable
func Select(db *sql.DB, dest interface{}, query string, args ...interface{}) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer func() { log.IfErrorDiffNil(rows.Close()) }()

	return scanAll(rows, dest, false)
}

func scannerInterface() reflect.Type {
	return reflect.TypeOf((*sql.Scanner)(nil)).Elem()
}

func structOnlyError(t reflect.Type) error {
	isStruct := t.Kind() == reflect.Struct
	isScanner := reflect.PtrTo(t).Implements(scannerInterface())
	if !isStruct {
		return fmt.Errorf("expected %s but got %s", reflect.Struct, t.Kind())
	}
	if isScanner {
		return fmt.Errorf(
			"structs can expects a struct dest but the provided struct type %s implements scanner",
			t.Name())
	}
	return fmt.Errorf("expected a struct, but struct %s has no exported fields", t.Name())
}

func isScannable(t reflect.Type) bool {
	if reflect.PtrTo(t).Implements(scannerInterface()) {
		return true
	}
	if t.Kind() != reflect.Struct {
		return true
	}

	// it's not important that we use the right mapper for this particular object,
	// we're only concerned on how many exported fields this struct has
	return t.NumField() == 0
}

func processScan(rows *sql.Rows, isPtr bool, base reflect.Type, direct reflect.Value) error {
	var vp reflect.Value
	for rows.Next() {
		vp = reflect.New(base)

		if err := rows.Scan(vp.Interface()); err != nil {
			return err
		}

		if isPtr {
			direct.Set(reflect.Append(direct, vp))
		} else {
			direct.Set(reflect.Append(direct, reflect.Indirect(vp)))
		}
	}
	return rows.Err()
}

func processNotScan(
	rows *sql.Rows, isPtr bool, base reflect.Type, direct reflect.Value, columns []string,
) error {
	var v, vp reflect.Value
	values := make([]interface{}, len(columns))
	for rows.Next() {
		// create a new struct type (which returns PtrTo) and indirect it
		vp = reflect.New(base)
		v = reflect.Indirect(vp)

		valuesByFields(v, values, columns)

		// scan into the struct field pointers and append to our results
		if err := rows.Scan(values...); err != nil {
			return err
		}

		if isPtr {
			direct.Set(reflect.Append(direct, vp))
		} else {
			direct.Set(reflect.Append(direct, v))
		}
	}
	return rows.Err()
}

func scanAll(rows *sql.Rows, dest interface{}, structOnly bool) error {
	value := reflect.ValueOf(dest)

	// json.Unmarshal returns errors for these
	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if value.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}
	direct := reflect.Indirect(value)

	slice, err := types.BaseType(value.Type(), reflect.Slice)
	if err != nil {
		return err
	}

	isPtr := slice.Elem().Kind() == reflect.Ptr
	base := types.DeRef(slice.Elem())
	scannable := isScannable(base)

	if structOnly {
		return structOnlyError(base)
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	if scannable {
		return processScan(rows, isPtr, base, direct)
	}
	return processNotScan(rows, isPtr, base, direct, columns)
}

func findByDBTag(v reflect.Value, s string) (interface{}, bool) {
	for i := 0; i < v.NumField(); i++ {
		f := reflect.Indirect(v).Field(i)
		dbTag, ok := v.Type().Field(i).Tag.Lookup("db")
		if !ok && f.Kind() == reflect.Struct {
			if value, ok := findByDBTag(f, s); ok {
				return value, ok
			}
			continue
		}
		if dbTag == s {
			return f.Addr().Interface(), true
		}
	}
	return nil, false
}

func valuesByFields(v reflect.Value, values []interface{}, columns []string) {
	v = reflect.Indirect(v)
	for i := 0; i < len(columns); i++ {
		values[i], _ = findByDBTag(v, columns[i])
	}
}
