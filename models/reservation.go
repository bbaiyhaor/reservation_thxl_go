package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type ReservationStatus int

const (
	AVAILABLE ReservationStatus = 1 + iota
	RESERVATED
	FEEDBACK
	DELETED
)

var reservationStatuses = [...]string{
	"AVAILABLE",
	"RESERVATED",
	"FEEDBACK",
	"DELETED",
}

func (rs ReservationStatus) String() string {
	return reservationStatuses[rs-1]
}

type StudentFeedback struct {
	Scores []int `bson:"scores"`
}

func (sf StudentFeedback) IsEmpty() bool {
	return sf.Scores == nil || len(sf.Scores) == 0
}

type TeacherFeedback struct {
	Problem string `bson:"problem"`
	Record  string `bson:"record"`
}

func (tf TeacherFeedback) IsEmpty() bool {
	return len(tf.Problem) == 0 || len(tf.Record) == 0
}

type Reservation struct {
	Id bson.ObjectId `bson:"_id"`
	// Indexed
	StartTime       time.Time         `bson:"start_time"`
	EndTime         time.Time         `bson:"end_time"`
	Status          ReservationStatus `bson:"status"`
	TeacherUsername string            `bson:"teacher_username"`
	TeacherFullname string            `bson:"teacher_fullname"`
	TeacherMobile   string            `bson:"teacher_mobile"`
	// Indexed
	StudentUsername string          `bson:"student_username"`
	StudentFullname string          `bson:"student_fullname"`
	StudentMobile   string          `bson:"student_mobile"`
	StudentFeedback StudentFeedback `bson:"student_feedback"`
	TeacherFeedback TeacherFeedback `bson:"teacher_feedback"`
}
