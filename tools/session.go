package tools

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"errors"
)

type SessionConf struct {
	Hosts    []string
	Database string
}

var (
	session *mgo.Session
)

func InitSession(conf SessionConf) *mgo.Session {
	var err error
	dialInfo := mgo.DialInfo{
		Addrs:    conf.Hosts,
		Database: conf.Database,
	}

	session, err = mgo.DialWithInfo(&dialInfo)

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to DB server. %s", err))
	}
	return session
}

func GetSession() (*mgo.Session, error) {
	var err error
	if session == nil {
		err = errors.New("Session is not initialized")
	}
	return session, err
}
