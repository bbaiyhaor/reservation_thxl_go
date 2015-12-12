package models

import (
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
		SourceId:        tr.Id.Hex(),
		TeacherId:       tr.TeacherId,
		StudentId:       "",
		StudentFeedback: StudentFeedback{},
		TeacherFeedback: TeacherFeedback{},
	}
}

type TimedReservationSlice []*TimedReservation

func (ts TimedReservationSlice) Len() int {
	return len(ts)
}

func (ts TimedReservationSlice) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

func (ts TimedReservationSlice) Less(i, j int) bool {
	if ts[i].Weekday != ts[j].Weekday {
		return ts[i].Weekday < ts[j].Weekday
	}
	return ts[i].StartTime.Before(ts[j].StartTime)
}
