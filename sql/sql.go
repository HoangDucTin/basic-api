package sql

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
)

const (
	sqlConnectionStringFormat = "server=%v;user id=%v;password=%v;"
)

var (
	db *sql.DB
)

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

func Find(database, table string, response []interface{}) error {
	selectStatement := fmt.Sprintf("USE %v; SELECT * FROM %v LIMIT %v", database, table, 1)
	statement, err := db.Prepare(selectStatement)
	if err != nil {
		return err
	}
	rows, err := statement.Query()
	if err != nil {
		return err
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	columnsNumber := len(columns)
	for rows.Next() {
		row := make([]interface{}, columnsNumber)
		if err := rows.Scan(&row); err != nil {
			return err
		}
		response = append(response, row)
	}
	return nil
}

func FindAll(database, table string, response []interface{}) error {
	selectStatement := fmt.Sprintf("USE %v; SELECT * FROM %v", database, table)
	statement, err := db.Prepare(selectStatement)
	if err != nil {
		return err
	}
	rows, err := statement.Query()
	if err != nil {
		return err
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	columnsNumber := len(columns)
	for rows.Next() {
		row := make([]interface{}, columnsNumber)
		if err := rows.Scan(&row); err != nil {
			return err
		}
		response = append(response, row)
	}
	return nil
}

func FindWithLimit(database, table string, limit int, response []interface{}) error {
	selectStatement := fmt.Sprintf("USE %v; SELECT * FROM %v LIMIT %v", database, table, limit)
	statement, err := db.Prepare(selectStatement)
	if err != nil {
		return err
	}
	rows, err := statement.Query()
	if err != nil {
		return err
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	columnsNumber := len(columns)
	for rows.Next() {
		row := make([]interface{}, columnsNumber)
		if err := rows.Scan(&row); err != nil {
			return err
		}
		response = append(response, row)
	}
	return nil
}
