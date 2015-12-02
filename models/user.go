package models

import (
	"gopkg.in/mgo.v2/bson"
)

type UserType int

const (
	UNKNWONUSER UserType = iota
	STUDENT
	TEACHER
	ADMIN
)

var userTypes = [...]string{
	"UNKNOWN",
	"STUDENT",
	"TEACHER",
	"ADMIN",
}

func (ut UserType) String() string {
	return userTypes[ut]
}

type Experience struct {
	Time     string `bson:"time"`
	Location string `bson:"location"`
	Teacher  string `bson:"teacher"`
}

func (e Experience) IsEmpty() bool {
	return len(e.Time) == 0 && len(e.Location) == 0 && len(e.Teacher) == 0
}

type Student struct {
	Id bson.ObjectId `bson:"_id"`
	// Indexed
	Username string   `bson:"username"`
	Password string   `bson:"password"`
	UserType UserType `bson:"user_type"`
	// Indexed
	BindedTeacher string `bson:"binded_teacher"`

	Fullname       string     `bson:"fullname"`
	Gender         string     `bson:"gender"`
	Birthday       string     `bson:"birthday"`
	School         string     `bson:"school"`
	Grade          string     `bson:"grade"`
	CurrentAddress string     `bson:"current_address"`
	FamilyAddress  string     `bson:"family_address"`
	Mobile         string     `bson:"mobile"`
	Email          string     `bson:"email"`
	Experience     Experience `bson:"experience"`
	FatherAge      string     `bson:"father_age"`
	FatherJob      string     `bson:"father_job"`
	FatherEdu      string     `bson:"father_edu"`
	MotherAge      string     `bson:"mother_age"`
	MotherJob      string     `bson:"mother_job"`
	MotherEdu      string     `bson:"mother_edu"`
	ParentMarriage string     `bson:"parent_marriage"`
	Significant    string     `bson:"significant"`
	Problem        string     `bson:"problem"`
}

type Teacher struct {
	Id bson.ObjectId `bson:"_id"`
	// Indexed
	Username string   `bson:"username"`
	Password string   `bson:"password"`
	Fullname string   `bson:"fullname"`
	Mobile   string   `bson:"mobile"`
	UserType UserType `bson:"user_type"`
}
