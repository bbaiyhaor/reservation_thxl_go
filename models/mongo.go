package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	Mongo *mgo.Database
)

// 为保证一致性，models中不用冗余存储，采用给字段构建索引的办法加快查询速度，以观后效

/**
Student
*/

func AddStudent(username string, password string) (*Student, error) {
	collection := Mongo.C("student")
	newStudent := &Student{
		Id:       bson.NewObjectId(),
		Username: username,
		Password: password,
		UserType: STUDENT,
	}
	if err := collection.Insert(newStudent); err != nil {
		return nil, err
	}
	return newStudent, nil
}

func UpsertStudent(student *Student) error {
	collection := Mongo.C("student")
	_, err := collection.UpsertId(student.Id, student)
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

func AddTeacher(username string, password string, fullname string, mobile string) (*Teacher, error) {
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
Admin
*/

func AddAdmin(username string, password string) (*Admin, error) {
	collection := Mongo.C("admin")
	newAdmin := &Admin{
		Id:       bson.NewObjectId(),
		Username: username,
		Password: password,
		UserType: ADMIN,
	}
	if err := collection.Insert(newAdmin); err != nil {
		return nil, err
	}
	return newAdmin, nil
}

func UpsertAdmin(admin *Admin) error {
	collection := Mongo.C("admin")
	_, err := collection.UpsertId(admin.Id, admin)
	return err
}

func GetAdminById(adminId string) (*Admin, error) {
	collection := Mongo.C("admin")
	admin := &Admin{}
	if err := collection.FindId(bson.ObjectIdHex(adminId)).One(admin); err != nil {
		return nil, err
	}
	return admin, nil
}

func GetAdminByUsername(username string) (*Admin, error) {
	collection := Mongo.C("admin")
	admin := &Admin{}
	if err := collection.Find(bson.M{"username": username}).One(admin); err != nil {
		return nil, err
	}
	return admin, nil
}

/**
Reservation
*/

func AddReservation(startTime time.Time, endTime time.Time, source ReservationSource, sourceId string,
	teacherId string) (*Reservation, error) {
	collection := Mongo.C("reservation")
	newReservation := &Reservation{
		Id:              bson.NewObjectId(),
		StartTime:       startTime,
		EndTime:         endTime,
		Status:          AVAILABLE,
		Source:          source,
		SourceId:        sourceId,
		TeacherId:       teacherId,
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

func GetReservationsByStudentId(studentId string) ([]*Reservation, error) {
	collection := Mongo.C("reservation")
	var reservations []*Reservation
	if err := collection.Find(bson.M{"student_id": studentId,
		"status": bson.M{"$ne": DELETED}}).Sort("start_time").All(&reservations); err != nil {
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

func GetReservatedReservationsBetweenTime(from time.Time, to time.Time) ([]*Reservation, error) {
	collection := Mongo.C("reservation")
	var reservations []*Reservation
	if err := collection.Find(bson.M{"start_time": bson.M{"$gte": from, "$lte": to},
		"status": RESERVATED}).Sort("start_time").All(&reservations); err != nil {
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

/**
TimedReservation
*/

func AddTimedReservation(weekday time.Weekday, startTime time.Time, endTime time.Time, teacherId string) (*TimedReservation, error) {
	collection := Mongo.C("timetable")
	timedReservation := &TimedReservation{
		Id:         bson.NewObjectId(),
		Weekday:    weekday,
		StartTime:  startTime,
		EndTime:    endTime,
		Status:     AVAILABLE,
		TeacherId:  teacherId,
		Exceptions: make(map[string]bool),
	}
	if err := collection.Insert(timedReservation); err != nil {
		return nil, err
	}
	return timedReservation, nil
}

func UpsertTimedReservation(timedReservation *TimedReservation) error {
	collection := Mongo.C("timetable")
	_, err := collection.UpsertId(timedReservation.Id, timedReservation)
	return err
}

func GetTimedReservationById(timedReservtionId string) (*TimedReservation, error) {
	collection := Mongo.C("timetable")
	timedReservation := &TimedReservation{}
	if err := collection.FindId(bson.ObjectIdHex(timedReservtionId)).One(timedReservation); err != nil {
		return nil, err
	}
	return timedReservation, nil
}

func GetTimedReservationsAll() ([]*TimedReservation, error) {
	collection := Mongo.C("timetable")
	var timedReservations []*TimedReservation
	if err := collection.Find(bson.M{"status": bson.M{"$ne": DELETED}}).All(&timedReservations); err != nil {
		return nil, err
	}
	return timedReservations, nil
}

func GetTimedReservationsByWeekday(weekday time.Weekday) ([]*TimedReservation, error) {
	collection := Mongo.C("timetable")
	var timedReservations []*TimedReservation
	if err := collection.Find(bson.M{"weekday": weekday,
		"status": bson.M{"$ne": DELETED}}).Sort("start_time").All(&timedReservations); err != nil {
		return nil, err
	}
	return timedReservations, nil
}

func GetTimedReservationsByTeacherId(teacherId string) ([]*TimedReservation, error) {
	collection := Mongo.C("timetable")
	var timedReservations []*TimedReservation
	if err := collection.Find(bson.M{"teacher_id": teacherId,
		"status": bson.M{"$ne": DELETED}}).All(&timedReservations); err != nil {
		return nil, err
	}
	return timedReservations, nil
}
