package sql

import (
	"errors"
	"reflect"

	// Native packages
	"fmt"

	// Third parties
	"github.com/jinzhu/gorm"
)

const (
	sqlConnectionStringFormat = "server=%s;user id=%s;password=%s;port=%s;database=%s"
)

var (
	db *gorm.DB

	// ErrNotSliceOrStructPtr is returned
	// when the parameter is not a pointer
	// of a slice or a struct to avoid panic.
	ErrNotSliceOrStructPtr = errors.New("not a pointer of a slice or a struct")

	// ErrNoSelector is returned when
	// the parameter is nil to avoid delete
	// all records.
	ErrNoSelector = errors.New("no selector")
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

// NewSQLClient creates a
// connection to SQL server
// base on the configuration
// passed into it. (Host,
// Username, Password, Port,
// Database name, and one of
// the most important things
// is the driver.)
func NewSQLClient(cfg Configs) error {
	connString := fmt.Sprintf(sqlConnectionStringFormat, cfg.Host, cfg.Username, cfg.Password, cfg.Port, cfg.Database)
	switch sqlDB, err := gorm.Open(cfg.Driver, connString); {
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

// Find selects the records
// base on the parameter
// 'condition' and push the
// data into the parameter
// 'result'.
func Find(table string, result interface{}, condition interface{}) error {
	if t := reflect.TypeOf(result).Kind(); t != reflect.Ptr {
		return ErrNotSliceOrStructPtr
	}
	if condition != nil {
		return db.Table(table).Where(condition).Find(result).Error
	}
	return db.Table(table).Find(result).Error
}

// Update updates all the
// records base on the
// parameter 'selector' for
// condition. If parameter
// 'selector' is nil, Update
// will updates all the records
// in table with data 'updater'.
func Update(table string, updater, selector interface{}) error {
	if selector != nil {
		return db.Table(table).Where(selector).UpdateColumns(updater).Error
	}
	return db.Table(table).UpdateColumns(updater).Error
}

// Delete removes all the records
// which satisfied the 'selector'.
func Delete(table string, selector interface{}) error {
	if selector != nil {
		return db.Table(table).Delete(selector).Error
	}
	return ErrNoSelector
}

// Insert accepts a pointer of
// a slice or a pointer of a
// struct. Insert creates all
// the records in the 'data'.
func Insert(table string, data interface{}) error {
	if t := reflect.TypeOf(data).Kind(); t != reflect.Ptr {
		return ErrNotSliceOrStructPtr
	}
	v := reflect.ValueOf(data).Elem()
	fmt.Printf("%+v\n", v)
	if vt := v.Kind(); vt == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			if err := db.Table(table).Create(v.Index(i).Interface()).Error; err != nil {
				return err
			}
		}
		return nil
	}
	if vt := v.Kind(); vt == reflect.Struct {
		return db.Table(table).Create(v).Error
	}
	return ErrNotSliceOrStructPtr
}