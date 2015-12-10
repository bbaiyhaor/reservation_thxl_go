package models

import (
	"fmt"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type TimedReservation struct {
	Id         bson.ObjectId     `bson:"_id"`
	Weekday    time.Weekday      `bson:"weekday"`
	StartTime  time.Time         `bson:"start_time"`
	EndTime    time.Time         `bson:"end_time"`
	Status     ReservationStatus `bson:"status"`
	TeacherId  string            `bson:"teacher_id"`
	Exceptions map[string]bool   `bson:"exception_map"` // exceptional dates
	Timed      map[string]bool   `bson:"timed_map"`     // timed dates
}

func (tr TimedReservation) ToReservation(date time.Time) *Reservation {
	return &Reservation{
		Id:              tr.Id,
		StartTime:       utils.ConcatTime(date, tr.StartTime),
		EndTime:         utils.ConcatTime(date, tr.EndTime),
		Status:          AVAILABLE,
		Source:          TIMETABLE,
		SourceId:        fmt.Sprintf("%x", string(tr.Id)),
		TeacherId:       tr.TeacherId,
		StudentId:       "",
		StudentFeedback: StudentFeedback{},
		TeacherFeedback: TeacherFeedback{},
	}
}
