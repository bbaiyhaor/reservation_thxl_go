package models

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UserType int

const (
	UNKNOWN UserType = iota
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

func (ut UserType) IntStr() string {
	return fmt.Sprintf("%d", ut)
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
	Id               bson.ObjectId `bson:"_id"`
	CreateTime       time.Time     `bson:"create_time"`
	UpdateTime       time.Time     `bson:"update_time"`
	Username         string        `bson:"username"` // Indexed
	Password         string        `bson:"password"`
	UserType         UserType      `bson:"user_type"`
	BindedTeacherId  string        `bson:"binded_teacher_id"` // Indexed
	ArchiveCategory  string        `bson:"archive_category"`
	ArchiveNumber    string        `bson:"archive_number"` // Indexed
	CrisisLevel      int           `bson:"crisis_level"`
	KeyCase          []int         `bson:"key_case"`          // deprecated
	MedicalDiagnosis []int         `bson:"medical_diagnosis"` // deprecated
	Fullname         string        `bson:"fullname"`
	Gender           string        `bson:"gender"`
	Birthday         string        `bson:"birthday"`
	School           string        `bson:"school"`
	Grade            string        `bson:"grade"`
	CurrentAddress   string        `bson:"current_address"`
	FamilyAddress    string        `bson:"family_address"`
	Mobile           string        `bson:"mobile"`
	Email            string        `bson:"email"`
	Experience       Experience    `bson:"experience"`
	FatherAge        string        `bson:"father_age"`
	FatherJob        string        `bson:"father_job"`
	FatherEdu        string        `bson:"father_edu"`
	MotherAge        string        `bson:"mother_age"`
	MotherJob        string        `bson:"mother_job"`
	MotherEdu        string        `bson:"mother_edu"`
	ParentMarriage   string        `bson:"parent_marriage"`
	Significant      string        `bson:"significant"`
	Problem          string        `bson:"problem"`
}

type Teacher struct {
	Id         bson.ObjectId `bson:"_id"`
	CreateTime time.Time     `bson:"create_time"`
	UpdateTime time.Time     `bson:"update_time"`
	Username   string        `bson:"username"` // Indexed
	Password   string        `bson:"password"`
	Fullname   string        `bson:"fullname"`
	Mobile     string        `bson:"mobile"`
	UserType   UserType      `bson:"user_type"`
}

type Admin struct {
	Id         bson.ObjectId `bson:"_id"`
	CreateTime time.Time     `bson:"create_time"`
	UpdateTime time.Time     `bson:"update_time"`
	Username   string        `bson:"username"` // Indexed
	Password   string        `bson:"password"`
	UserType   UserType      `bson:"user_type"`
}
