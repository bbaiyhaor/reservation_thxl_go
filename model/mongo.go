package model

import (
	"github.com/mijia/sweb/log"
	"github.com/shudiwsh2009/reservation_thxl_go/config"
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
		//mongoDbDialInfo := mgo.DialInfo{
		//	Addrs:    []string{"123.56.190.91:20171"},
		//	Timeout:  60 * time.Second,
		//	Database: "reservation_thxl",
		//	Username: "reservation_readonly",
		//	Password: "CczzgMv^Cfr9Gh*G^K6nz6G@xaQ/(yYYcDXFu3Eq*ryW6xg2JN3++8D7aFc]2DcB",
		//}
		//session, err = mgo.DialWithInfo(&mongoDbDialInfo)
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

func (m *MongoClient) EnsureAllIndexes() error {
	var err error
	err = dbStudent.EnsureIndex(mgo.Index{
		Key: []string{"username", "user_type"},
	})
	if err != nil {
		return err
	}
	err = dbStudent.EnsureIndex(mgo.Index{
		Key: []string{"archive_category", "archive_number", "user_type"},
	})
	if err != nil {
		return err
	}
	err = dbStudent.EnsureIndex(mgo.Index{
		Key: []string{"binded_teacher_id", "user_type"},
	})
	if err != nil {
		return err
	}

	err = dbTeacher.EnsureIndex(mgo.Index{
		Key: []string{"username", "user_type"},
	})
	if err != nil {
		return err
	}
	err = dbTeacher.EnsureIndex(mgo.Index{
		Key: []string{"fullname", "user_type"},
	})
	if err != nil {
		return err
	}
	err = dbTeacher.EnsureIndex(mgo.Index{
		Key: []string{"mobile", "user_type"},
	})
	if err != nil {
		return err
	}

	err = dbAdmin.EnsureIndex(mgo.Index{
		Key: []string{"username", "user_type"},
	})
	if err != nil {
		return err
	}

	err = dbArchive.EnsureIndex(mgo.Index{
		Key: []string{"student_username"},
	})
	if err != nil {
		return err
	}
	err = dbArchive.EnsureIndex(mgo.Index{
		Key: []string{"archive_category", "archive_number"},
	})
	if err != nil {
		return err
	}

	err = dbReservation.EnsureIndex(mgo.Index{
		Key: []string{"student_id", "status", "start_time"},
	})
	if err != nil {
		return err
	}
	err = dbReservation.EnsureIndex(mgo.Index{
		Key: []string{"start_time", "end_time", "status"},
	})
	if err != nil {
		return err
	}
	err = dbReservation.EnsureIndex(mgo.Index{
		Key: []string{"start_time", "status"},
	})
	if err != nil {
		return err
	}
	err = dbReservation.EnsureIndex(mgo.Index{
		Key: []string{"teacher_feedback.school_contact", "status", "start_time"},
	})
	if err != nil {
		return err
	}

	err = dbTimetable.EnsureIndex(mgo.Index{
		Key: []string{"status"},
	})
	if err != nil {
		return err
	}
	err = dbTimetable.EnsureIndex(mgo.Index{
		Key: []string{"weekday", "status"},
	})
	if err != nil {
		return err
	}
	err = dbTimetable.EnsureIndex(mgo.Index{
		Key: []string{"teacher_id", "status"},
	})
	if err != nil {
		return err
	}

	return nil
}

// DANGER!!!
func (m *MongoClient) DropAllIndexes() error {
	for _, coll := range []*mgo.Collection{dbStudent, dbTeacher, dbAdmin, dbArchive, dbReservation, dbTimetable} {
		err := m.DropIndexes(coll)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MongoClient) DropIndexes(coll *mgo.Collection) error {
	indexes, err := coll.Indexes()
	if err != nil {
		return err
	}
	for _, index := range indexes {
		if index.Name == "_id_" {
			continue
		}
		err = coll.DropIndexName(index.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
