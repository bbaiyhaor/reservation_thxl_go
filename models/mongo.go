package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	Mongo *mgo.Database
)

/**
Student
*/

func AddStudent(username string, password string) (*Student, error) {
	collection := Mongo.C("student")
	newStudent := &Student{
		Id:       bson.NewObjectId(),
		Username: username,
		Password: password,
	}
	if err := collection.Insert(newStudent); err != nil {
		return nil, err
	}
	return newStudent, nil
}

func UpsertStudent(student *Student) error {
	collection := Mongo.C("student")
	_, err := collection.UpdateId(student.Id, student)
	return err
}

func GetStudentById(studentId string) (*Student, error) {
	collection := Mongo.C("student")
	student := &Student{}
	if err := collection.FindId(bson.ObjectIdHex(studentId)).One(student); err != nil {
		return nil, err
	}
	return student, nil
}

func GetStudentByUsername(username string) (*Student, error) {
	collection := Mongo.C("student")
	student := &Student{}
	if err := collection.Find(bson.M{"username": username}).One(student); err != nil {
		return nil, err
	}
	return student, nil
}

/**
Teacher
*/

func AddSimpleTeacher(username string, password string) (*Teacher, error) {
	collection := Mongo.C("teacher")
	newTeacher := &Teacher{
		Id:       bson.NewObjectId(),
		Username: username,
		Password: password,
		UserType: TEACHER,
	}
	if err := collection.Insert(newTeacher); err != nil {
		return nil, err
	}
	return newTeacher, nil
}

func AddFullTeacher(username string, password string, fullname string, mobile string) (*Teacher, error) {
	collection := Mongo.C("teacher")
	newTeacher := &Teacher{
		Id:       bson.NewObjectId(),
		Username: username,
		Password: password,
		Fullname: fullname,
		Mobile:   mobile,
		UserType: TEACHER,
	}
	if err := collection.Insert(newTeacher); err != nil {
		return nil, err
	}
	return newTeacher, nil
}

func UpsertTeacher(teacher *Teacher) error {
	collection := Mongo.C("teacher")
	_, err := collection.UpsertId(teacher.Id, teacher)
	return err
}

func GetTeacherById(teacherId string) (*Teacher, error) {
	collection := Mongo.C("teacher")
	teacher := &Teacher{}
	if err := collection.FindId(bson.ObjectIdHex(teacherId)).One(teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func GetTeacherByUsername(username string) (*Teacher, error) {
	collection := Mongo.C("teacher")
	teacher := &Teacher{}
	if err := collection.Find(bson.M{"username": username}).One(teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func GetTeacherByFullname(fullname string) (*Teacher, error) {
	collection := Mongo.C("teacher")
	teacher := &Teacher{}
	if err := collection.Find(bson.M{"fullname": fullname}).One(teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func GetTeacherByMobile(mobile string) (*Teacher, error) {
	collection := Mongo.C("teacher")
	teacher := &Teacher{}
	if err := collection.Find(bson.M{"mobile": mobile}).One(teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

/**
Reservation
*/

func AddReservation(startTime time.Time, endTime time.Time, teacherFullname string, teacherUsername string,
	teacherMobile string) (*Reservation, error) {
	collection := Mongo.C("reservation")
	newReservation := &Reservation{
		Id:              bson.NewObjectId(),
		StartTime:       startTime,
		EndTime:         endTime,
		Status:          AVAILABLE,
		TeacherFullname: teacherFullname,
		TeacherUsername: teacherUsername,
		TeacherMobile:   teacherMobile,
		StudentFeedback: StudentFeedback{},
		TeacherFeedback: TeacherFeedback{},
	}
	if err := collection.Insert(newReservation); err != nil {
		return nil, err
	}
	return newReservation, nil
}

func UpsertReservation(reservation *Reservation) error {
	collection := Mongo.C("reservation")
	_, err := collection.UpsertId(reservation.Id, reservation)
	return err
}

func GetReservationById(id string) (*Reservation, error) {
	collection := Mongo.C("reservation")
	reservation := &Reservation{}
	if err := collection.FindId(bson.ObjectIdHex(id)).One(reservation); err != nil {
		return nil, err
	}
	return reservation, nil
}

func GetReservationsByStudentUsername(username string) ([]*Reservation, error) {
	collection := Mongo.C("reservation")
	var reservations []*Reservation
	if err := collection.Find(bson.M{"student_username": username}).All(&reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func GetReservationsBetweenTime(from time.Time, to time.Time) ([]*Reservation, error) {
	collection := Mongo.C("reservation")
	var reservations []*Reservation
	if err := collection.Find(bson.M{"start_time": bson.M{"$gte": from, "$lte": to},
		"status": bson.M{"$ne": DELETED}}).Sort("start_time").All(&reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func GetReservationsAfterTime(from time.Time) ([]*Reservation, error) {
	collection := Mongo.C("reservation")
	var reservations []*Reservation
	if err := collection.Find(bson.M{"start_time": bson.M{"$gte": from},
		"status": bson.M{"$ne": DELETED}}).Sort("start_time").All(&reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}
