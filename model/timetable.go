package model

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/utils"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type TimedReservation struct {
	Id         bson.ObjectId   `bson:"_id"`
	CreateTime time.Time       `bson:"create_time"`
	UpdateTime time.Time       `bson:"update_time"`
	Weekday    time.Weekday    `bson:"weekday"`
	StartTime  time.Time       `bson:"start_time"`
	EndTime    time.Time       `bson:"end_time"`
	Status     int             `bson:"status"`
	TeacherId  string          `bson:"teacher_id"`
	Exceptions map[string]bool `bson:"exception_map"` // exceptional dates
	Timed      map[string]bool `bson:"timed_map"`     // timed dates
}

func (m *Model) AddTimedReservation(weekday time.Weekday, startTime time.Time, endTime time.Time, teacherId string) (*TimedReservation, error) {
	collection := m.mongo.C("timetable")
	timedReservation := &TimedReservation{
		Id:         bson.NewObjectId(),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		Weekday:    weekday,
		StartTime:  startTime,
		EndTime:    endTime,
		Status:     RESERVATION_STATUS_CLOSED,
		TeacherId:  teacherId,
		Exceptions: make(map[string]bool),
		Timed:      make(map[string]bool),
	}
	if err := collection.Insert(timedReservation); err != nil {
		return nil, err
	}
	return timedReservation, nil
}

func (m *Model) UpsertTimedReservation(timedReservation *TimedReservation) error {
	if timedReservation == nil || !timedReservation.Id.Valid() {
		return errors.New("字段不合法")
	}
	collection := m.mongo.C("timetable")
	timedReservation.UpdateTime = time.Now()
	_, err := collection.UpsertId(timedReservation.Id, timedReservation)
	return err
}

func (m *Model) GetTimedReservationById(timedReservtionId string) (*TimedReservation, error) {
	if len(timedReservtionId) == 0 || !bson.IsObjectIdHex(timedReservtionId) {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("timetable")
	timedReservation := &TimedReservation{}
	if err := collection.FindId(bson.ObjectIdHex(timedReservtionId)).One(timedReservation); err != nil {
		return nil, err
	}
	return timedReservation, nil
}

func (m *Model) GetTimedReservationsAll() ([]*TimedReservation, error) {
	collection := m.mongo.C("timetable")
	var timedReservations []*TimedReservation
	if err := collection.Find(bson.M{"status": bson.M{"$ne": RESERVATION_STATUS_DELETED}}).All(&timedReservations); err != nil {
		return nil, err
	}
	return timedReservations, nil
}

func (m *Model) GetTimedReservationsByWeekday(weekday time.Weekday) ([]*TimedReservation, error) {
	collection := m.mongo.C("timetable")
	var timedReservations []*TimedReservation
	if err := collection.Find(bson.M{"weekday": weekday,
		"status": bson.M{"$ne": RESERVATION_STATUS_DELETED}}).All(&timedReservations); err != nil {
		return nil, err
	}
	return timedReservations, nil
}

func (m *Model) GetTimedReservationsByTeacherId(teacherId string) ([]*TimedReservation, error) {
	if len(teacherId) == 0 || !bson.IsObjectIdHex(teacherId) {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("timetable")
	var timedReservations []*TimedReservation
	if err := collection.Find(bson.M{"teacher_id": teacherId,
		"status": bson.M{"$ne": RESERVATION_STATUS_DELETED}}).All(&timedReservations); err != nil {
		return nil, err
	}
	return timedReservations, nil
}

func (tr TimedReservation) ToReservation(date time.Time) *Reservation {
	return &Reservation{
		Id:              tr.Id,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
		StartTime:       utils.ConcatTime(date, tr.StartTime),
		EndTime:         utils.ConcatTime(date, tr.EndTime),
		Status:          RESERVATION_STATUS_AVAILABLE,
		Source:          RESERVATION_SOURCE_TIMETABLE,
		SourceId:        tr.Id.Hex(),
		TeacherId:       tr.TeacherId,
		StudentId:       "",
		StudentFeedback: StudentFeedback{},
		TeacherFeedback: TeacherFeedback{},
	}
}
