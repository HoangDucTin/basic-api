package sql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	sqlConnectionStringFormat = "server=%s;user id=%s;password=%s;port=%s;database=%s"
)

var (
	db *sqlx.DB
)

// Configs contains the configuration
// for opening connection to SQL server.
type Configs struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

// NewSQL creates an instance
// of SQL database connection, based
// on given host, username, password
// and database name.
func NewSQL(cfg Configs) error {
	connString := fmt.Sprintf(sqlConnectionStringFormat, cfg.Host, cfg.Username, cfg.Password, cfg.Port, cfg.Database)
	switch sqlDB, err := sqlx.Connect(cfg.Driver, connString); {
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

func exec(statement string) error {
	_, err := db.Exec(statement)
	return err
}

// Find will return one latest record
// base on the Select statement.
func Find(statement string, response interface{}) error {
	return db.Get(response, statement)
}

// FindAll will return all the records
// base on the Select statement.
func FindAll(statement string, response interface{}) error {
	return db.Select(response, statement)
}

// Execute will execute the statement
// for Insert, Update, Delete.
func Execute(statement string) error {
	return exec(statement)
}
