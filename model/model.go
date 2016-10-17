package model

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/config"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

type Model struct {
	mongo *mgo.Database
}

func NewModel() *Model {
	var session *mgo.Session
	var err error
	if config.Instance().IsSmockServer() {
		session, err = mgo.Dial("127.0.0.1:27017")
	} else {
		mongoDbDialInfo := mgo.DialInfo{
			Addrs:    []string{config.Instance().MongoHost},
			Timeout:  60 * time.Second,
			Database: config.Instance().MongoDatabase,
			Username: config.Instance().MongoUser,
			Password: config.Instance().MongoPassword,
		}
		session, err = mgo.DialWithInfo(&mongoDbDialInfo)
	}
	if err != nil {
		log.Fatalf("连接数据库失败：%v", err)
	}
	//defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	ret := &Model{
		mongo: session.DB("reservation_thxl"),
	}
	return ret
}
