package sql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	// This import is for the driver to
	// connect Microsoft SQL server.
	_ "github.com/denisenkom/go-mssqldb"

	"github.com/tinwoan-go/basic-api/utils"
)

const (
	sqlConnectionStringFormat = "server=%v;user id=%v;password=%v;database=%v;"
)

var (
	db          *sql.DB
	// ErrNilParam shows that the
	// passed in parameter is nil.
	ErrNilParam = errors.New("nil parameter")
)

// NewSQL creates an instance
// of SQL database connection, based
// on given host, username, password
// and database name.
func NewSQL(host, username, password, database string) error {
	connString := fmt.Sprintf(sqlConnectionStringFormat, host, username, password, database)
	switch sqlDB, err := sql.Open("mssql", connString); {
	case err != nil:
		return err
	default:
		db = sqlDB
		return nil
	}
}

// Close closes the connection
// of SQL database, based on the
// current instance in application.
func Close() error {
	return db.Close()
}

// Find finds one latest record
// from given table within given database.
// Parameter 'condition' can either be a
// struct with json tag, or a
// map[string]interface{}.
// The parameter 'response' MUST be a struct
// with json tags attached to fields.
func Find(table string, response interface{}, condition interface{}) error {
	selectStatement := fmt.Sprintf("SELECT TOP %v * FROM %v", 1, table)
	switch conditionMap, err := BuildMap(condition); {
	case err != nil:
		if err != ErrNilParam {
			return err
		}
	case conditionMap != nil:
		selectStatement += fmt.Sprintf(" WHERE 1 = 1")
		for key, value := range conditionMap {
			if to := reflect.TypeOf(value); to.Kind() == reflect.String {
				value = fmt.Sprintf("'%v'", value)
			}
			selectStatement += fmt.Sprintf(" AND %v = %v", key, value)
		}
		defer utils.MapDestructor(conditionMap)
	}

	statement, err := db.Prepare(selectStatement)
	if err != nil {
		return err
	}
	rows, err := statement.Query()
	if err != nil {
		return err
	}
	var objects map[string]interface{}
	defer utils.MapDestructor(objects)
	for rows.Next() {
		columns, err := rows.ColumnTypes()
		if err != nil {
			return err
		}

		values := make([]interface{}, len(columns))
		object := map[string]interface{}{}
		for i, column := range columns {
			object[column.Name()] = reflect.New(column.ScanType()).Interface()
			values[i] = object[column.Name()]
		}

		err = rows.Scan(values...)
		if err != nil {
			return err
		}

		objects = object
		break
	}
	b, err := json.Marshal(objects)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &response); err != nil {
		return err
	}
	return nil
}

// FindAll finds all the records
// from given table within given database.
// Parameter 'condition' can either be a
// struct with json tags or a
// map[string]interface{}
// Parameter 'response' MUST be a struct
// with json tags.
func FindAll(table string, response interface{}, condition interface{}) error {
	selectStatement := fmt.Sprintf("SELECT * FROM %v", table)
	switch conditionMap, err := BuildMap(condition); {
	case err != nil:
		if err != ErrNilParam {
			return err
		}
	case conditionMap != nil:
		selectStatement += fmt.Sprintf(" WHERE 1 = 1")
		for key, value := range conditionMap {
			if to := reflect.TypeOf(value); to.Kind() == reflect.String {
				value = fmt.Sprintf("'%v'", value)
			}
			selectStatement += fmt.Sprintf(" AND %v = %v", key, value)
		}
		defer utils.MapDestructor(conditionMap)
	}

	statement, err := db.Prepare(selectStatement)
	if err != nil {
		return err
	}
	rows, err := statement.Query()
	if err != nil {
		return err
	}
	var objects []map[string]interface{}
	defer utils.MapDestructor(objects)
	for rows.Next() {
		columns, err := rows.ColumnTypes()
		if err != nil {
			return err
		}

		values := make([]interface{}, len(columns))
		object := map[string]interface{}{}
		for i, column := range columns {
			object[column.Name()] = reflect.New(column.ScanType()).Interface()
			values[i] = object[column.Name()]
		}

		err = rows.Scan(values...)
		if err != nil {
			return err
		}

		objects = append(objects, object)
	}
	b, err := json.Marshal(objects)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &response); err != nil {
		return err
	}
	return nil
}

// FindWithLimit returns a number of
// latest records by given limit number
// from given table within given database.
// Parameter 'condition' can either be a
// struct with json tags or a
// map[string]interface{}.
// Parameter 'response' MUST be a struct
// with json tags.
func FindWithLimit(table string, limit int, response interface{}, condition interface{}) error {
	if limit < 1 {
		return nil
	}
	selectStatement := fmt.Sprintf("SELECT TOP %v * FROM %v", limit, table)
	switch conditionMap, err := BuildMap(condition); {
	case err != nil:
		if err != ErrNilParam {
			return err
		}
	case conditionMap != nil:
		selectStatement += fmt.Sprintf(" WHERE 1 = 1")
		for key, value := range conditionMap {
			if to := reflect.TypeOf(value); to.Kind() == reflect.String {
				value = fmt.Sprintf("'%v'", value)
			}
			selectStatement += fmt.Sprintf(" AND %v = %v", key, value)
		}
		defer utils.MapDestructor(conditionMap)
	}

	statement, err := db.Prepare(selectStatement)
	if err != nil {
		return err
	}
	rows, err := statement.Query()
	if err != nil {
		return err
	}
	var objects []map[string]interface{}
	defer utils.MapDestructor(objects)
	for rows.Next() {
		columns, err := rows.ColumnTypes()
		if err != nil {
			return err
		}

		values := make([]interface{}, len(columns))
		object := map[string]interface{}{}
		for i, column := range columns {
			object[column.Name()] = reflect.New(column.ScanType()).Interface()
			values[i] = object[column.Name()]
		}

		err = rows.Scan(values...)
		if err != nil {
			return err
		}

		objects = append(objects, object)
	}
	b, err := json.Marshal(objects)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &response); err != nil {
		return err
	}
	return nil
}

// Update updates the records
// satisfied the updater with the values
// from the selector within the
// given table inside the given
// database.
// Parameter 'selector' can either be
// a struct with json tags or a
// map[string]interface{}, but MUST
// be not nil.
// Parameter 'updater' can either be a
// struct with json tags or a
// map[string]interface{}, but MUST
// be not nil.
func Update(table string, selector interface{}, updater interface{}) error {
	updaterMap, err := BuildMap(updater)
	if err != nil {
		if err == ErrNilParam {
			return errors.New("can not update record(s) without updater")
		}
		return err
	}
	defer utils.MapDestructor(updaterMap)
	selectorMap, err := BuildMap(selector)
	if err != nil {
		if err == ErrNilParam {
			return errors.New("can not update record(s) without selector")
		}
		return err
	}
	defer utils.MapDestructor(selectorMap)
	updateStatement := fmt.Sprintf("UPDATE %v SET ", table)
	for key, value := range updaterMap {
		if to := reflect.TypeOf(value); to.Kind() == reflect.String {
			value = fmt.Sprintf("'%v'", value)
		}
		updateStatement += fmt.Sprintf("%v = %v, ", key, value)
	}
	updateStatement = updateStatement[:len(updateStatement)-2] + " WHERE 1 = 1"

	for key, value := range selectorMap {
		if to := reflect.TypeOf(value); to.Kind() == reflect.String {
			value = fmt.Sprintf("'%v'", value)
		}
		updateStatement += fmt.Sprintf(" AND %v = %v", key, value)
	}

	if _, err := db.Exec(updateStatement); err != nil {
		return err
	}
	return nil
}

// Delete deletes the records
// selected by the give selector
// with the given table, inside the
// given database.
// Parameter 'selector' can either be
// a struct with json tags or a
// map[string]interface{}, but MUST
// be not nil.
func Delete(table string, selector interface{}) error {
	selectorMap, err := BuildMap(selector)
	if err != nil {
		if err == ErrNilParam {
			return errors.New("can not delete record(s) without selector")
		}
		return err
	}
	defer utils.MapDestructor(selectorMap)
	deleteStatement := fmt.Sprintf("DELETE FROM %v WHERE 1 = 1", table)
	for key, value := range selectorMap {
		if to := reflect.TypeOf(value); to.Kind() == reflect.String {
			value = fmt.Sprintf("'%v'", value)
		}
		deleteStatement += fmt.Sprintf(" AND %v = %v", key, value)
	}
	if _, err := db.Exec(deleteStatement); err != nil {
		return err
	}
	return nil
}

// Insert inserts one record to
// the given table and database names.
// Parameter 'data' can either be a struct
// with json tags or a map[string]interface{}
// but MUST be not nil.
// In case parameter 'data' is a struct with
// json tags, the order of the fields MUST
// be the same order of fields in the
// destination table.
func Insert(table string, data interface{}) error {
	insertMap, err := BuildMap(data)
	if err != nil {
		if err == ErrNilParam {
			return errors.New("can not insert nil data")
		}
		return err
	}
	defer utils.MapDestructor(insertMap)
	insertStatement := fmt.Sprintf("INSERT INTO %v VALUES(", table)
	for _, value := range insertMap {
		if to := reflect.TypeOf(value); to.Kind() == reflect.String {
			value = fmt.Sprintf("'%v'", value)
		}
		insertStatement += fmt.Sprintf("%v, ", value)
	}
	insertStatement = insertStatement[:len(insertStatement)-2] + ")"
	if _, err := db.Exec(insertStatement); err != nil {
		return err
	}
	return nil
}

// InsertMany inserts multiple
// records to the given table and
// database name.
// Parameter 'data' can either be a struct
// with json tags or a map[string]interface{}
// but MUST be not nil.
// In case parameter 'data' is a struct with
// json tags, the order of the fields MUST
// be the same order of fields in the
// destination table.
func InsertMany(table string, data interface{}) error {
	insertMap, err := BuildListMap(data)
	if err != nil {
		if err == ErrNilParam {
			return errors.New("can not insert nil data")
		}
		return err
	}
	defer utils.MapDestructor(insertMap)
	insertStatement := fmt.Sprintf("INSERT INTO %v VALUES", table)
	for _, info := range insertMap {
		insertStatement += "("
		for _, value := range info {
			if to := reflect.TypeOf(value); to.Kind() == reflect.String {
				value = fmt.Sprintf("'%v'", value)
			}
			insertStatement += fmt.Sprintf("%v, ", value)
		}
		insertStatement = insertStatement[:len(insertStatement)-2] + "),"
	}
	insertStatement = insertStatement[:len(insertStatement)-1]
	if _, err := db.Exec(insertStatement); err != nil {
		return err
	}
	return nil
}

// BuildMap returns a map[string]interface{}
// built from parameter 'parameter'.
func BuildMap(parameter interface{}) (map[string]interface{}, error) {
	if parameter == nil {
		return nil, ErrNilParam
	}
	var parameterMap map[string]interface{}
	switch m, ok := parameter.(map[string]interface{}); {
	case ok:
		return m, nil
	default:
		b, err := json.Marshal(parameter)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(b, &parameterMap); err != nil {
			return nil, err
		}
		return parameterMap, nil
	}
}

// BuildListMap returns a list of
// map[string]interface{} built from
// parameter 'parameter.
func BuildListMap(parameter interface{}) ([]map[string]interface{}, error) {
	if parameter == nil {
		return nil, ErrNilParam
	}
	var parameterListMap []map[string]interface{}
	switch lm, ok := parameter.([]map[string]interface{}); {
	case ok:
		return lm, nil
	default:
		b, err := json.Marshal(parameter)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(b, &parameterListMap); err != nil {
			return nil, err
		}
		return parameterListMap, nil
	}
}
