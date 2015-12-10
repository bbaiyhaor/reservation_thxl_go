package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"strings"
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

type ReservationSource int

const (
	TIMETABLE ReservationSource = 1 + iota
	TEACHER_ADD
	ADMIN_ADD
)

var reservationSources = [...]string{
	"TIMETABLE",
	"TEACHER",
	"ADMIN",
}

func (rs ReservationSource) String() string {
	return reservationSources[rs-1]
}

type StudentFeedback struct {
	Scores []string `bson:"scores"`
}

func (sf StudentFeedback) IsEmpty() bool {
	return sf.Scores == nil || len(sf.Scores) == 0
}

type TeacherFeedback struct {
	Category     string `bson:"category"`
	Participants []int  `bson:"participants"`
	Problem      string `bson:"problem"`
	Record       string `bson:"record"`
}

func (tf TeacherFeedback) IsEmpty() bool {
	return len(tf.Category) == 0 || tf.Participants == nil || len(tf.Problem) == 0 || len(tf.Record) == 0
}

type Reservation struct {
	Id              bson.ObjectId     `bson:"_id"`
	StartTime       time.Time         `bson:"start_time"` // indexed
	EndTime         time.Time         `bson:"end_time"`
	Status          ReservationStatus `bson:"status"`
	Source          ReservationSource `bson:"source"`
	SourceId        string            `bson:"source_id"`
	TeacherId       string            `bson:"teacher_id"` // indexed
	StudentId       string            `bson:"student_id"` // indexed
	StudentFeedback StudentFeedback   `bson:"student_feedback"`
	TeacherFeedback TeacherFeedback   `bson:"teacher_feedback"`
}

type ReservationSlice []*Reservation

func (rs ReservationSlice) Len() int {
	return len(rs)
}

func (rs ReservationSlice) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs ReservationSlice) Less(i, j int) bool {
	if rs[i].StartTime.Before(rs[j].StartTime) {
		return true
	} else {
		return strings.Compare(rs[i].TeacherId, rs[j].TeacherId) < 0
	}
}