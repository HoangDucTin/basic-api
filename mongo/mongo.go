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
	// session holds the session
	// of connection to MongoDB
	// and can be used only inside
	// this package.
	session *mgo.Session

	// ErrInitialized is returned when
	// cannot clone from the current
	// MongoDB session.
	ErrInitialized = errors.New("MongoDB connection has not been initialized")

	// ErrNotSlice is returned when
	// the argument is not a slice.
	ErrNotSlice = errors.New("the list is not a slice")

	// ErrNotSliceAddress is returned
	// when the argument is not a slice
	// address to avoid panic.
	ErrNotSliceAddress = errors.New("argument must be a slice address")

	// ErrSliceOrPointerOfSliceOnly is
	// returned when the argument is
	// neither a slice or pointer of
	// a slice.
	ErrSliceOrPointerOfSliceOnly = errors.New("only accept slice or pointer of slice")

	// ErrStructOrPointerOfStructOnly is
	// returned when the argument is
	// neither a struct or pointer
	// of a struct.
	ErrStructOrPointerOfStructOnly = errors.New("only accept struct or pointer of struct")

	// ErrNotPointerOfStruct is returned
	// when the argument is not pointer
	// of a struct.
	ErrNotPointerOfStruct = errors.New("argument must be a pointer of struct")
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

// Find queries from 'collection'
// the first record that satisfied
// the 'selector'. The result will
// be put into 'result'.
// Find accepts empty 'database'.
// In the case of empty 'database',
// it will consider using database
// when initiate connection.
func Find(database, collection string, selector, result interface{}) error {
	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()
	return s.DB(database).C(collection).Find(selector).One(result)
}

// FindAll queries from 'collection'
// all the records that satisfied
// the 'selector'. The results will
// be put into 'result'. FindAll
// accepts 'result' as slice address
// only (Based on the specification
// of mgo.Collection.Find().All()).
// FindAll accepts empty 'database'.
// In the case of empty 'database',
// it will consider using database
// when initiate connection.
func FindAll(database, collection string, selector, result interface{}) error {
	if reflect.TypeOf(result).Kind() != reflect.Ptr ||
		reflect.TypeOf(result).Elem().Kind() != reflect.Slice {
		return ErrNotSliceAddress
	}
	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()
	return s.DB(database).C(collection).Find(selector).All(result)
}

// Insert creates a record in
// 'collection' with the value
// of 'data'.
// Insert accepts empty 'database'.
// In the case of empty 'database',
// it will consider using database
// when initiate connection.
func Insert(database, collection string, data interface{}) error {
	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()
	return s.DB(database).C(collection).Insert(data)
}

// InsertAll creates records in
// 'collection' with the values
// from 'list'.
// Arguments 'list' can only be
// a slice or pointer of a slice.
// InsertAll accepts empty 'database'.
// In the case of empty 'database',
// it will consider using database
// when initiate connection.
func InsertAll(database, collection string, list interface{}) error {
	slice := reflect.ValueOf(list)
	if slice.Kind() != reflect.Slice {
		if reflect.TypeOf(list).Kind() != reflect.Ptr ||
			reflect.TypeOf(list).Kind() == reflect.Ptr &&
				reflect.TypeOf(list).Elem().Kind() != reflect.Slice {
			return ErrSliceOrPointerOfSliceOnly
		}
	}
	ret := make([]interface{}, slice.Len())
	for i := 0; i < slice.Len(); i++ {
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

// Remove deletes the first record
// that satisfied 'selector' from
// 'collection'.
// Remove accepts empty 'database'.
// In the case of empty 'database',
// it will consider using database
// when initiate connection.
func Remove(database, collection string, selector interface{}) error {
	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()
	return s.DB(database).C(collection).Remove(selector)
}

// RemoveAll deletes all the records
// that satisfied 'selector' from
// 'collection'.
// RemoveAll accepts empty 'database'.
// In the case of empty 'database',
// it will consider using database
// when initiate connection.
func RemoveAll(database, collection string, selector interface{}) error {
	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()
	_, err := s.DB(database).C(collection).RemoveAll(selector)
	return err
}

// Update updates the first record
// in 'collection' that satisfied
// 'selector', with the new data
// 'updater'.
// Update accepts empty 'database'.
// In the case of empty 'database',
// it will consider using database
// when initiate connection.
func Update(database, collection string, selector, updater interface{}) error {
	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()
	return s.DB(database).C(collection).Update(selector, updater)
}

// UpdateAll updates all the records
// in 'collection' that satisfied
// 'selector', with the new data
// 'updater'.
// UpdateAll accepts empty 'database'.
// In the case of empty 'database',
// it will consider using database
// when initiate connection.
func UpdateAll(database, collection string, selector, updater interface{}) error {
	s := cloneSession()
	if s == nil {
		return ErrInitialized
	}
	defer s.Close()
	_, err := s.DB(database).C(collection).UpdateAll(selector, updater)
	return err
}

// Change updates the first record
// in 'collection' that satisfied
// 'selector', with new data 'new'
// into target 'result'.
// Change apply the 'new' data to
// 'result' and keep it as the
// latest version of the record.
// Change accepts empty 'database'.
// In the case of empty 'database',
// it will consider using database
// when initiate connection.
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

// End-of-file
