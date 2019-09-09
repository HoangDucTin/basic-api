package mongo

import (
	"errors"
	"github.com/globalsign/mgo"
	"github.com/tinwoan-go/basic-api/logger"
	"strings"
	"time"
)

var session *mgo.Session

const InitializedError = "MongoDB connection has not been initialized"

// This function creates an instance
// of MongoDB session based on given
// addresses, database name, username
// and password. The timeout duration
// for connecting is given.
// (Notice: You should not set the
// timeout duration too long.)
func NewMongoClient(addresses, database, username, password string, timeout time.Duration) error {
	dialInfo := mgo.DialInfo{
		Addrs:    strings.Split(addresses, ","),
		Database: database,
		Username: username,
		Password: password,
		Timeout:  timeout,
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

// This function closes the session
// to connect with MongoDB.
func Close() {
	if session != nil {
		session.Close()
		session = nil
		return
	}
	logger.Warn("Session does not exist or already been closed")
}

func cloneSession() *mgo.Session {
	if session == nil {
		return nil
	}
	return session.Copy()
}

// This function finds one specific record
// based on given selector, in side
// the collection within the database
// that are given the names.
func Find(database, collection string, selector, result interface{}) error {
	s := cloneSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()
	return s.DB(database).C(collection).Find(selector).One(&result)
}

// This function finds every record
// based on given selector, in side
// the collection within the database,
// which are given the names.
func FindAll(database, collection string, selector, result interface{}) error {
	s := cloneSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()
	return s.DB(database).C(collection).Find(selector).All(&result)
}

// This function inserts one record
// with the given data into the collection
// within the database, which are given
// names.
func Insert(database, collection string, data interface{}) error {
	s := cloneSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	return s.DB(database).C(collection).Insert(data)
}

// This function inserts all the records
// given by the list into the collection
// within the database, which are given
// names.
func InsertAll(database, collection string, list []interface{}) error {
	s := cloneSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	bulk := s.DB(database).C(collection).Bulk()
	for _, item := range list {
		bulk.Insert(item)
	}

	_, err := bulk.Run()
	return err
}

// This function removes latest record
// with the given selector from collection
// in database, which are given names.
func Remove(database, collection string, selector interface{}) error {
	s := cloneSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	return s.DB(database).C(collection).Remove(selector)
}

// This function removes all the records
// with the given selector from collection
// in database, which are given names.
func RemoveAll(database, collection string, selector interface{}) error {
	s := cloneSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	_, err := s.DB(database).C(collection).RemoveAll(selector)
	return err
}

// This function updates latest record
// selected by selector, with new data
// is updater, within collection in
// database, which are given names.
func Update(database, collection string, selector, updater interface{}) error {
	s := cloneSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	return s.DB(database).C(collection).Update(selector, updater)
}

// This function updates all the records
// selected by selector, with new data
// is updater, within collection in
// database, which are given names.
func UpdateAll(database, collection string, selector, updater interface{}) error {
	s := cloneSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	_, err := s.DB(database).C(collection).UpdateAll(selector, updater)
	return err
}

// This function updates the latest record
// selected by the selector, with the given
// new data, within collection in database,
// which are given names. This function
// also returns a new version of the record
// after update to database.
func Change(database, collection string, selector, new, result interface{}) error {
	s := cloneSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	change := mgo.Change{Update: new, ReturnNew: true}
	_, err := s.DB(database).C(collection).Find(selector).Apply(change, result)
	return err
}

// End-of-file
