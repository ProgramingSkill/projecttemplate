package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"strings"
	"time"
)

type MgoDB struct {
	Session *mgo.Session
}

var mgoDB = new(MgoDB)

func checkMongoReachable(err error) {
	if err == nil {
		return
	}
	if strings.Contains(err.Error(), mgo.ErrNotFound.Error()) || strings.Contains(err.Error(), mgo.ErrCursor.Error()) || strings.Contains(err.Error(), "E11000") {
		return
	}
	if strings.Contains(err.Error(), "no reachable") {
		return
	}
	if strings.Contains(err.Error(), "timeout") {
		return
	}
	return
}

func mongoInit() {
	var mongoPara []string
	for _, value := range Config.MongoDB {
		one := fmt.Sprintf("%s:%s",
			value.IP,
			value.Port,
		)
		mongoPara = append(mongoPara, one)
	}
	mongoUser := Config.MongoDB["1"].User
	mongoPasswd := Config.MongoDB["1"].Password
	mongoPoolLimit := Config.MongoDB["1"].MongoPoolLimit
	if mongoPoolLimit <= 0 {
		mongoPoolLimit = 64
	}
	mgoDB.mgoConn(mongoPara, mongoUser, mongoPasswd, mongoPoolLimit)
	return
}

func (m *MgoDB) mgoConn(addrs []string, user, passwd string, mongoPoolLimit int) {

	dialInfo := &mgo.DialInfo{
		Addrs:    addrs,
		Username: user,
		Password: passwd,
		Timeout:  5 * time.Second,
	}
	for {
		session, err := mgo.DialWithInfo(dialInfo)
		if err != nil {
			Warningf("mgo failed:", err)
			time.Sleep(1 * time.Second)
			continue
		}
		session.SetPoolLimit(mongoPoolLimit)

		m.Session = session
		Info("mgo conn ok")
		break

	}
	return
}
