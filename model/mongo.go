package model

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/config"
	"github.com/mijia/sweb/log"
	"gopkg.in/mgo.v2"
	"time"
)

var (
	dbStudent     *mgo.Collection
	dbTeacher     *mgo.Collection
	dbAdmin       *mgo.Collection
	dbArchive     *mgo.Collection
	dbReservation *mgo.Collection
	dbTimetable   *mgo.Collection
)

type MongoClient struct {
	mongo *mgo.Database
}

func NewMongoClient() *MongoClient {
	var session *mgo.Session
	var err error
	if config.Instance().IsSmockServer() {
		session, err = mgo.Dial("127.0.0.1:27017")
	} else {
		mongoDbDialInfo := mgo.DialInfo{
			Addrs:    []string{config.Instance().MongoHost},
			Timeout:  60 * time.Second,
			Database: config.Instance().MongoAuthDatabase,
			Username: config.Instance().MongoAuthUser,
			Password: config.Instance().MongoAuthPassword,
		}
		session, err = mgo.DialWithInfo(&mongoDbDialInfo)
	}
	if err != nil {
		log.Fatalf("连接数据库失败：%v", err)
	}
	//defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	mongo := session.DB(config.Instance().MongoDatabase)
	dbStudent = mongo.C("student")
	dbTeacher = mongo.C("teacher")
	dbAdmin = mongo.C("admin")
	dbArchive = mongo.C("archive")
	dbReservation = mongo.C("reservation")
	dbTimetable = mongo.C("timetable")
	ret := &MongoClient{
		mongo: mongo,
	}
	return ret
}
