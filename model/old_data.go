package model

import "gopkg.in/mgo.v2/bson"

func (m *MongoClient) GetAllOldAdmins() ([]*OldAdmin, error) {
	collection := m.mongo.C("admin")
	var admins []*OldAdmin
	err := collection.Find(bson.M{}).All(&admins)
	return admins, err
}

func (m *MongoClient) GetAllOldTeachers() ([]*OldTeacher, error) {
	collection := m.mongo.C("teacher")
	var teachers []*OldTeacher
	err := collection.Find(bson.M{}).All(&teachers)
	return teachers, err
}

func (m *MongoClient) GetAllOldStudents() ([]*OldStudent, error) {
	collection := m.mongo.C("student")
	var students []*OldStudent
	err := collection.Find(bson.M{}).All(&students)
	return students, err
}

func (m *MongoClient) GetAllOldArchives() ([]*OldArchive, error) {
	collection := m.mongo.C("archive")
	var archives []*OldArchive
	err := collection.Find(bson.M{}).All(&archives)
	return archives, err
}

func (m *MongoClient) GetAllOldReservations() ([]*OldReservation, error) {
	collection := m.mongo.C("reservation")
	var reservations []*OldReservation
	err := collection.Find(bson.M{}).All(&reservations)
	return reservations, err
}

func (m *MongoClient) GetAllOldTimetables() ([]*OldTimedReservation, error) {
	collection := m.mongo.C("timetable")
	var timedReservations []*OldTimedReservation
	err := collection.Find(bson.M{}).All(&timedReservations)
	return timedReservations, err
}
