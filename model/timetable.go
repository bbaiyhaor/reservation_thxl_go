package model

import (
	re "bitbucket.org/shudiwsh2009/reservation_thxl_go/rerror"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Index: status
// Index: weekday + status
// Index: teacher_id + status
type TimedReservation struct {
	Id         bson.ObjectId   `bson:"_id"`
	Weekday    time.Weekday    `bson:"weekday"`
	StartTime  time.Time       `bson:"start_time"`
	EndTime    time.Time       `bson:"end_time"`
	Status     int             `bson:"status"`
	TeacherId  string          `bson:"teacher_id"`
	Exceptions map[string]bool `bson:"exception_map"`
	Timed      map[string]bool `bson:"timed_map"`
	CreatedAt  time.Time       `bson:"created_at"`
	UpdatedAt  time.Time       `bson:"updated_at"`
}

func (m *MongoClient) InsertTimedReservation(timedReservation *TimedReservation) error {
	now := time.Now()
	timedReservation.CreatedAt = now
	timedReservation.UpdatedAt = now
	return dbTimetable.Insert(timedReservation)
}

func (m *MongoClient) UpdateTimedReservation(timedReservation *TimedReservation) error {
	timedReservation.UpdatedAt = time.Now()
	return dbTimetable.UpdateId(timedReservation.Id, timedReservation)
}

func (m *MongoClient) UpdateTimedReservationWithoutTime(timedReservation *TimedReservation) error {
	return dbTimetable.UpdateId(timedReservation.Id, timedReservation)
}

func (m *MongoClient) GetTimedReservationById(id string) (*TimedReservation, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, re.NewRErrorCode("id is not valid", nil, re.ERROR_DATABASE)
	}
	var timedReservation TimedReservation
	err := dbTimetable.FindId(bson.ObjectIdHex(id)).One(&timedReservation)
	return &timedReservation, err
}

func (m *MongoClient) GetAllTimedReservations() ([]*TimedReservation, error) {
	var timedReservations []*TimedReservation
	err := dbTimetable.Find(bson.M{"status": bson.M{"$ne": RESERVATION_STATUS_DELETED}}).All(&timedReservations)
	return timedReservations, err
}

func (m *MongoClient) GetTimedReservationsByWeekday(weekday time.Weekday) ([]*TimedReservation, error) {
	var timedReservations []*TimedReservation
	err := dbTimetable.Find(bson.M{"weekday": weekday,
		"status": bson.M{"$ne": RESERVATION_STATUS_DELETED}}).All(&timedReservations)
	return timedReservations, err
}

func (m *MongoClient) GetTimedReservationsByTeacherId(teacherId string) ([]*TimedReservation, error) {
	var timedReservations []*TimedReservation
	err := dbTimetable.Find(bson.M{"teacher_id": teacherId,
		"status": bson.M{"$ne": RESERVATION_STATUS_DELETED}}).All(&timedReservations)
	return timedReservations, err
}

func (tr TimedReservation) ToReservation(date time.Time) *Reservation {
	return &Reservation{
		Id:              tr.Id,
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
