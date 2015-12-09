package models

import (
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
}
