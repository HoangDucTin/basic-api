package sql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/tinwoan-go/basic-api/logger"
	"reflect"
)

const (
	sqlConnectionStringFormat = "server=%v;user id=%v;password=%v;"
)

var (
	db *sql.DB
)

// This function creates an instance
// of SQL database connection, based
// on given host, username and password.
func NewSql(host, username, password string) error {
	connString := fmt.Sprintf(sqlConnectionStringFormat, host, username, password)
	switch sqlDB, err := sql.Open("mssql", connString); {
	case err != nil:
		return err
	default:
		db = sqlDB
		return nil
	}
}

// This function closes the connection
// of SQL database, based on the
// current instance in application.
func Close() error {
	return db.Close()
}

// This function finds one latest record
// from given table within given database.
func Find(database, table string, response interface{}, condition interface{}) error {
	selectStatement := fmt.Sprintf("USE %v; SELECT TOP %v * FROM %v", database, 1, table)
	if condition != nil {
		b, err := json.Marshal(condition)
		if err != nil {
			return err
		}
		var conditionMap map[string]interface{}
		if err := json.Unmarshal(b, &conditionMap); err != nil {
			return err
		}
		selectStatement += fmt.Sprintf(" WHERE 1 = 1")
		for key, value := range conditionMap {
			if to := reflect.TypeOf(value); to.Kind() == reflect.String {
				value = fmt.Sprintf("'%v'", value)
			}
			selectStatement += fmt.Sprintf(" AND %v = %v", key, value)
		}
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

// This function finds all the records
// from given table within given database.
func FindAll(database, table string, response interface{}, condition interface{}) error {
	selectStatement := fmt.Sprintf("USE %v; SELECT * FROM %v", database, table)
	logger.Warn("condition: %+v", condition)
	if condition != nil {
		b, err := json.Marshal(condition)
		if err != nil {
			return err
		}
		var conditionMap map[string]interface{}
		if err := json.Unmarshal(b, &conditionMap); err != nil {
			return err
		}
		selectStatement += fmt.Sprintf(" WHERE 1 = 1")
		for key, value := range conditionMap {
			if to := reflect.TypeOf(value); to.Kind() == reflect.String {
				value = fmt.Sprintf("'%v'", value)
			}
			selectStatement += fmt.Sprintf(" AND %v = %v", key, value)
		}
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

// This function returns a number of
// latest records by given limit number
// from given table within given database.
func FindWithLimit(database, table string, limit int, response interface{}, condition interface{}) error {
	selectStatement := fmt.Sprintf("USE %v; SELECT TOP %v * FROM %v", database, limit, table)
	if condition != nil {
		b, err := json.Marshal(condition)
		if err != nil {
			return err
		}
		var conditionMap map[string]interface{}
		if err := json.Unmarshal(b, &conditionMap); err != nil {
			return err
		}
		selectStatement += fmt.Sprintf(" WHERE 1 = 1")
		for key, value := range conditionMap {
			if to := reflect.TypeOf(value); to.Kind() == reflect.String {
				value = fmt.Sprintf("'%v'", value)
			}
			selectStatement += fmt.Sprintf(" AND %v = %v", key, value)
		}
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

// This function updates the records
// satisfied the updater with the values
// from the selector within the
// given table inside the given
// database.
func Update(database, table string, selector interface{}, updater interface{}) error {
	if updater == nil {
		return errors.New("can not update record(s) without updater")
	}
	b, err := json.Marshal(updater)
	if err != nil {
		return err
	}
	var updaterMap map[string]interface{}
	if err := json.Unmarshal(b, &updaterMap); err != nil {
		return err
	}
	updateStatement := fmt.Sprintf("USE %v; UPDATE %v SET ", database, table)
	for key, value := range updaterMap {
		if to := reflect.TypeOf(value); to.Kind() == reflect.String {
			value = fmt.Sprintf("'%v'", value)
		}
		updateStatement += fmt.Sprintf(", %v = %v", key, value)
	}
	if selector != nil {
		updateStatement += " WHERE 1 = 1"
		b, err := json.Marshal(selector)
		if err != nil {
			return err
		}
		var selectorMap map[string]interface{}
		if err := json.Unmarshal(b, &selectorMap); err != nil {
			return err
		}
		for key, value := range selectorMap {
			if to := reflect.TypeOf(value); to.Kind() == reflect.String {
				value = fmt.Sprintf("'%v'", value)
			}
			updateStatement += fmt.Sprintf(" AND %v = %v", key, value)
		}
	}
	if _, err := db.Exec(updateStatement); err != nil {
		return err
	}
	return nil
}

// This function deletes the records
// selected by the give selector
// with the given table, inside the
// given database.
func Delete(database, table string, selector interface{}) error {
	if selector == nil {
		return errors.New("can not delete record(s) without selector")
	}
	b, err := json.Marshal(selector)
	if err != nil {
		return err
	}
	var selectorMap map[string]interface{}
	if err := json.Unmarshal(b, &selectorMap); err != nil {
		return err
	}
	deleteStatement := fmt.Sprintf("USE %v; DELETE FROM %v WHERE 1 = 1", database, table)
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

// This function inserts one record to
// the given table and database names.
func Insert(database, table string, data interface{}) error {
	if data == nil {
		return errors.New("can not insert nil data")
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	insertStatement := fmt.Sprintf("USE %v; INSERT INTO %v VALUES(",database, table)
	for _, value := range m {
		if to := reflect.TypeOf(value); to.Kind() == reflect.String {
			value = fmt.Sprintf("'%v'", value)
		}
		insertStatement += fmt.Sprintf("%v, ", value)
	}
	insertStatement = insertStatement[:len(insertStatement) - 2] + ")"
	if _, err := db.Exec(insertStatement); err != nil {
		return err
	}
	return nil
}

// This function inserts multiple
// records to the given table and
// database name.
func InsertMany(database, table string, data interface{}) error {
	if data == nil {
		return errors.New("can not insert nil data")
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	var m []map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	insertStatement := fmt.Sprintf("USE %v; INSERT INTO %v VALUES", database, table)
	for _, info := range m {
		insertStatement += "("
		for _, value := range info {
			if to := reflect.TypeOf(value); to.Kind() == reflect.String {
				value = fmt.Sprintf("'%v'", value)
			}
			insertStatement += fmt.Sprintf("%v, ", value)
		}
		insertStatement = insertStatement[:len(insertStatement) - 2] + "),"
	}
	insertStatement = insertStatement[:len(insertStatement)-1]
	if _, err := db.Exec(insertStatement); err != nil {
		return err
	}
	return nil
}