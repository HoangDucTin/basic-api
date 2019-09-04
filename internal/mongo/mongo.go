package mongo

import (
	"errors"
	"strings"
	"time"

	"github.com/globalsign/mgo"
)

var session *mgo.Session

const InitializedError = "MongoDB connection has not been initialized"

func Setup(addresses, database, username, password string, timeout time.Duration) error {
	dialInfo := mgo.DialInfo{
		Addrs:    strings.Split(addresses, ","),
		Database: database,
		Username: username,
		Password: password,
		Timeout:  timeout,
	}

	if s, err := mgo.DialWithInfo(&dialInfo);
		err != nil {
		return err
	} else {
		s.SetMode(mgo.Monotonic, true)
		session = s
	}

	return nil
}

func Close() {
	if session != nil {
		session.Close()
		session = nil
	}
}

func newSession() *mgo.Session {
	if session == nil {
		return nil
	}

	return session.Copy()
}

func Search(database, collection string, query, result interface{}) error {
	s := newSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	return s.DB(database).C(collection).Find(query).One(result)
}

func SearchAll(database, collection string, query, result interface{}) error {
	s := newSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	return s.DB(database).C(collection).Find(query).All(result)
}

func Insert(database, collection string, data interface{}) error {
	s := newSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	return s.DB(database).C(collection).Insert(data)
}

func InsertAll(database, collection string, list []interface{}) error {
	s := newSession()
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

func Remove(database, collection string, query interface{}) error {
	s := newSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	return s.DB(database).C(collection).Remove(query)
}

func RemoveAll(database, collection string, query interface{}) error {
	s := newSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	_, err := s.DB(database).C(collection).RemoveAll(query)
	return err
}

func Update(database, collection string, query, into interface{}) error {
	s := newSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	return s.DB(database).C(collection).Update(query, into)
}

func UpdateAll(database, collection string, query, into interface{}) error {
	s := newSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	_, err := s.DB(database).C(collection).UpdateAll(query, into)
	return err
}

func Change(database, collection string, query, into, result interface{}) error {
	s := newSession()
	if s == nil {
		return errors.New(InitializedError)
	}
	defer s.Close()

	change := mgo.Change{Update: into, ReturnNew: true}
	_, err := s.DB(database).C(collection).Find(query).Apply(change, result)
	return err
}

// End-of-file
