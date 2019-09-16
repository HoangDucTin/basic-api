package mongo

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/globalsign/mgo"
)

// Configs contains the configuration
// for opening connection to MongoDB.
type Configs struct {
	Addresses string
	Database  string
	Username  string
	Password  string
	Timeout   time.Duration
}

var (
	session *mgo.Session
	// ErrInitialized is returned when
	// cannot clone from the current
	// MongoDB session.
	ErrInitialized = errors.New("MongoDB connection has not been initialized")
	ErrNotSlice    = errors.New("the list is not a slice")
)

// NewMongoClient creates an instance
// of MongoDB session based on given
// addresses, database name, username
// and password. The timeout duration
// for connecting is given.
// (Notice: You should not set the
// timeout duration too long.)
func NewMongoClient(cfg Configs) error {
	dialInfo := mgo.DialInfo{
		Addrs:    strings.Split(cfg.Addresses, ","),
		Database: cfg.Database,
		Username: cfg.Username,
		Password: cfg.Password,
		Timeout:  cfg.Timeout,
	}

	switch s, err := mgo.DialWithInfo(&dialInfo); {
	case err != nil:
		return err
	default:
		s.SetMode(mgo.Monotonic, true)
		session = s
		return nil
	}
}

// Close closes the session
// to connect with MongoDB.
func Close() {
	if session != nil {
		session.Close()
		session = nil
		return
	}
}

func cloneSession() *mgo.Session {
	if session == nil {
		return nil
	}
	return session.Copy()
}

// Find finds one specific record
// based on given selector, in side
// the collection within the database
// that are given the names.
func Find(database, collection string, selector, result interface{}) error {
	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()
	return s.DB(database).C(collection).Find(selector).One(&result)
}

// FindAll finds every record
// based on given selector, in side
// the collection within the database,
// which are given the names.
func FindAll(database, collection string, selector, result interface{}) error {
	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()
	return s.DB(database).C(collection).Find(selector).All(&result)
}

// Insert inserts one record
// with the given data into the collection
// within the database, which are given
// names.
func Insert(database, collection string, data interface{}) error {
	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()

	return s.DB(database).C(collection).Insert(data)
}

// InsertAll inserts all the records
// given by the list into the collection
// within the database, which are given
// names.
func InsertAll(database, collection string, list interface{}) error {
	slice := reflect.ValueOf(list)
	if slice.Kind() != reflect.Slice {
		return ErrNotSlice
	}

	ret := make([]interface{}, slice.Len())

	for i:=0; i<slice.Len(); i++ {
		ret[i] = slice.Index(i).Interface()
	}

	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()

	bulk := s.DB(database).C(collection).Bulk()
	for _, item := range ret {
		bulk.Insert(item)
	}

	_, err := bulk.Run()
	return err
}

// Remove removes latest record
// with the given selector from collection
// in database, which are given names.
func Remove(database, collection string, selector interface{}) error {
	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()

	return s.DB(database).C(collection).Remove(selector)
}

// RemoveAll removes all the records
// with the given selector from collection
// in database, which are given names.
func RemoveAll(database, collection string, selector interface{}) error {
	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()

	_, err := s.DB(database).C(collection).RemoveAll(selector)
	return err
}

// Update updates latest record
// selected by selector, with new data
// is updater, within collection in
// database, which are given names.
func Update(database, collection string, selector, updater interface{}) error {
	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()

	return s.DB(database).C(collection).Update(selector, updater)
}

// UpdateAll updates all the records
// selected by selector, with new data
// is updater, within collection in
// database, which are given names.
func UpdateAll(database, collection string, selector, updater interface{}) error {
	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()

	_, err := s.DB(database).C(collection).UpdateAll(selector, updater)
	return err
}

// Change updates the latest record
// selected by the selector, with the given
// new data, within collection in database,
// which are given names. This function
// also returns a new version of the record
// after update to database.
func Change(database, collection string, selector, new, result interface{}) error {
	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()

	change := mgo.Change{Update: new, ReturnNew: true}
	_, err := s.DB(database).C(collection).Find(selector).Apply(change, result)
	return err
}

// Delete finds a single
// document (first one)
// matching the provided
// selector document and
// removes it from the
// database.
func Delete(database, collection string, selector interface{}) error {
	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()

	return s.DB(database).C(collection).Remove(selector)
}

// DeleteAll finds all documents
// matching the provided selector
// document and removes them
// from the database.
func DeleteAll(database, collection string, selector interface{}) error {
	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()
	_, err := s.DB(database).C(collection).RemoveAll(selector)
	return err
}

// End-of-file
