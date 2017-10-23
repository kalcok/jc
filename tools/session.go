package tools

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"errors"
)

type SessionConf mgo.DialInfo

var (
	session *mgo.Session
)

func InitSession(conf *SessionConf) *mgo.Session {
	var err error
	session, err = mgo.DialWithInfo((*mgo.DialInfo)(conf))

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to DB server. %s", err))
	}
	return session
}

func CloseSession() {
	if session != nil {
		session.Close()
	}
}

func GetSessionCopy() (*mgo.Session, error) {
	var err error
	if session == nil {
		err = errors.New("Session is not initialized")
	}
	return session.Copy(), err
}

func GetSessionClone() (*mgo.Session, error) {
	var err error
	if session == nil {
		err = errors.New("Session is not initialized")
	}
	return session.Clone(), err
}