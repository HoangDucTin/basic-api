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

func Find(database, collection string, selector, result interface{}) error {
	s := cloneSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()
	return s.DB(database).C(collection).Find(selector).One(&result)
}

func FindAll(database, collection string, selector, result interface{}) error {
	s := cloneSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()
	return s.DB(database).C(collection).Find(selector).All(&result)
}

func Insert(database, collection string, data interface{}) error {
	s := cloneSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	return s.DB(database).C(collection).Insert(data)
}

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

func Remove(database, collection string, selector interface{}) error {
	s := cloneSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	return s.DB(database).C(collection).Remove(selector)
}

func RemoveAll(database, collection string, selector interface{}) error {
	s := cloneSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	_, err := s.DB(database).C(collection).RemoveAll(selector)
	return err
}

func Update(database, collection string, selector, updater interface{}) error {
	s := cloneSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	return s.DB(database).C(collection).Update(selector, updater)
}

func UpdateAll(database, collection string, selector, updater interface{}) error {
	s := cloneSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	_, err := s.DB(database).C(collection).UpdateAll(selector, updater)
	return err
}

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
