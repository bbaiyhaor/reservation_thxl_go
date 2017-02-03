package model

import (
	re "bitbucket.org/shudiwsh2009/reservation_thxl_go/rerror"
	"crypto/sha256"
	"encoding/base64"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	USER_TYPE_UNKNOWN = 0
	USER_TYPE_STUDENT = 1
	USER_TYPE_TEACHER = 2
	USER_TYPE_ADMIN   = 3

	USER_GENDER_MALE   = "男"
	USER_GENDER_FEMALE = "女"
)

// Index: username + user_type
// Index: archive_category + archive_number + user_type
// Index: binded_teacher_id + user_type
type Student struct {
	Id                bson.ObjectId `bson:"_id"`
	Username          string        `bson:"username"`
	Password          string        `bson:"password"`
	Salt              string        `bson:"salt"`
	EncryptedPassword string        `bson:"encrypted_password"`
	UserType          int           `bson:"user_type"`
	BindedTeacherId   string        `bson:"binded_teacher_id"`
	ArchiveCategory   string        `bson:"archive_category"`
	ArchiveNumber     string        `bson:"archive_number"`
	CrisisLevel       int           `bson:"crisis_level"`
	Fullname          string        `bson:"fullname"`
	Gender            string        `bson:"gender"`
	Birthday          string        `bson:"birthday"`
	School            string        `bson:"school"`
	Grade             string        `bson:"grade"`
	CurrentAddress    string        `bson:"current_address"`
	FamilyAddress     string        `bson:"family_address"`
	Mobile            string        `bson:"mobile"`
	Email             string        `bson:"email"`
	Experience        Experience    `bson:"experience"`
	ParentInfo        ParentInfo    `bson:"parent_info"`
	Significant       string        `bson:"significant"`
	Problem           string        `bson:"problem"`
	CreatedAt         time.Time     `bson:"created_at"`
	UpdatedAt         time.Time     `bson:"updated_at"`
}

type Experience struct {
	Time     string `bson:"time"`
	Location string `bson:"location"`
	Teacher  string `bson:"teacher"`
}

func (e Experience) IsEmpty() bool {
	return e.Time == "" && e.Location == "" && e.Teacher == ""
}

type ParentInfo struct {
	FatherAge      string `bson:"father_age"`
	FatherJob      string `bson:"father_job"`
	FatherEdu      string `bson:"father_edu"`
	MotherAge      string `bson:"mother_age"`
	MotherJob      string `bson:"mother_job"`
	MotherEdu      string `bson:"mother_edu"`
	ParentMarriage string `bson:"parent_marriage"`
}

func (student *Student) PreInsert() error {
	salt := EncodePassword("salt", strconv.Itoa(rand.Int()))
	student.Salt = salt[:16]
	student.Password = EncodePassword(student.Salt, student.Password)
	student.Username = strings.TrimSpace(student.Username)
	student.UserType = USER_TYPE_STUDENT
	return nil
}

func (m *MongoClient) InsertStudent(student *Student) error {
	student.PreInsert()
	now := time.Now()
	student.Id = bson.NewObjectId()
	student.CreatedAt = now
	student.UpdatedAt = now
	return dbStudent.Insert(student)
}

func (m *MongoClient) UpdateStudent(student *Student) error {
	student.UpdatedAt = time.Now()
	return dbStudent.UpdateId(student.Id, student)
}

func (m *MongoClient) UpdateStudentWithoutTime(student *Student) error {
	return dbStudent.UpdateId(student.Id, student)
}

func (m *MongoClient) GetAllStudents() ([]*Student, error) {
	var students []*Student
	err := dbStudent.Find(bson.M{}).All(&students)
	return students, err
}

func (m *MongoClient) GetStudentById(id string) (*Student, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, re.NewRErrorCode("id is not valid", nil, re.ERROR_DATABASE)
	}
	var student Student
	err := dbStudent.FindId(bson.ObjectIdHex(id)).One(&student)
	return &student, err
}

func (m *MongoClient) GetStudentByUsername(username string) (*Student, error) {
	var student Student
	err := dbStudent.Find(bson.M{"username": username, "user_type": USER_TYPE_STUDENT}).One(&student)
	return &student, err
}

func (m *MongoClient) GetStudentByArchiveCategoryAndNumber(archiveCategory string, archiveNumber string) (*Student, error) {
	var student Student
	err := dbStudent.Find(bson.M{"archive_category": archiveCategory, "archive_number": archiveNumber, "user_type": USER_TYPE_STUDENT}).One(&student)
	return &student, err
}

func (m *MongoClient) GetStudentsByBindedTeacherId(teacherId string) ([]*Student, error) {
	var students []*Student
	err := dbStudent.Find(bson.M{"binded_teacher_id": teacherId, "user_type": USER_TYPE_STUDENT}).All(&students)
	return students, err
}

// Index: username + user_type
// Index: fullname + user_type
// Index: mobile + user_type
type Teacher struct {
	Id                bson.ObjectId `bson:"_id"`
	Username          string        `bson:"username"` // Indexed
	Password          string        `bson:"password"`
	Salt              string        `bson:"salt"`
	EncryptedPassword string        `bson:"encrypted_password"`
	UserType          int           `bson:"user_type"`
	Fullname          string        `bson:"fullname"`
	Mobile            string        `bson:"mobile"`
	CreatedAt         time.Time     `bson:"created_at"`
	UpdatedAt         time.Time     `bson:"updated_at"`
}

func (teacher *Teacher) PreInsert() error {
	salt := EncodePassword("salt", strconv.Itoa(rand.Int()))
	teacher.Salt = salt[:16]
	teacher.Password = EncodePassword(teacher.Salt, teacher.Password)
	teacher.Username = strings.TrimSpace(teacher.Username)
	teacher.UserType = USER_TYPE_TEACHER
	return nil
}

func (m *MongoClient) InsertTeacher(teacher *Teacher) error {
	teacher.PreInsert()
	now := time.Now()
	teacher.Id = bson.NewObjectId()
	teacher.CreatedAt = now
	teacher.UpdatedAt = now
	return dbTeacher.Insert(teacher)
}

func (m *MongoClient) UpdateTeacher(teacher *Teacher) error {
	teacher.UpdatedAt = time.Now()
	return dbTeacher.UpdateId(teacher.Id, teacher)
}

func (m *MongoClient) UpdateTeacherWithoutTime(teacher *Teacher) error {
	return dbTeacher.UpdateId(teacher.Id, teacher)
}

func (m *MongoClient) GetTeacherById(id string) (*Teacher, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, re.NewRErrorCode("id is not valid", nil, re.ERROR_DATABASE)
	}
	var teacher Teacher
	err := dbTeacher.FindId(bson.ObjectIdHex(id)).One(&teacher)
	return &teacher, err
}

func (m *MongoClient) GetTeacherByUsername(username string) (*Teacher, error) {
	var teacher Teacher
	err := dbTeacher.Find(bson.M{"username": username, "user_type": USER_TYPE_TEACHER}).One(&teacher)
	return &teacher, err
}

func (m *MongoClient) GetTeacherByFullname(fullname string) (*Teacher, error) {
	var teacher Teacher
	err := dbTeacher.Find(bson.M{"fullname": fullname, "user_type": USER_TYPE_TEACHER}).One(&teacher)
	return &teacher, err
}

func (m *MongoClient) GetTeacherByMobile(mobile string) (*Teacher, error) {
	var teacher Teacher
	err := dbTeacher.Find(bson.M{"mobile": mobile, "user_type": USER_TYPE_TEACHER}).One(&teacher)
	return &teacher, err
}

// Index: username + user_type
type Admin struct {
	Id                bson.ObjectId `bson:"_id"`
	Username          string        `bson:"username"`
	Password          string        `bson:"password"`
	Salt              string        `bson:"salt"`
	EncryptedPassword string        `bson:"encrypted_password"`
	UserType          int           `bson:"user_type"`
	Fullname          string        `bson:"fullname"`
	Mobile            string        `bson:"mobile"`
	CreatedAt         time.Time     `bson:"created_at"`
	UpdatedAt         time.Time     `bson:"updated_at"`
}

func (admin *Admin) PreInsert() error {
	salt := EncodePassword("salt", strconv.Itoa(rand.Int()))
	admin.Salt = salt[:16]
	admin.Password = EncodePassword(admin.Salt, admin.Password)
	admin.Username = strings.TrimSpace(admin.Username)
	admin.UserType = USER_TYPE_ADMIN
	return nil
}

func (m *MongoClient) InsertAdmin(admin *Admin) error {
	admin.PreInsert()
	now := time.Now()
	admin.Id = bson.NewObjectId()
	admin.CreatedAt = now
	admin.UpdatedAt = now
	return dbAdmin.Insert(admin)
}

func (m *MongoClient) UpdateAdmin(admin *Admin) error {
	admin.UpdatedAt = time.Now()
	return dbAdmin.UpdateId(admin.Id, admin)
}

func (m *MongoClient) UpdateAdminWithoutTime(admin *Admin) error {
	return dbAdmin.UpdateId(admin.Id, admin)
}

func (m *MongoClient) GetAdminById(id string) (*Admin, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, re.NewRErrorCode("id is not valid", nil, re.ERROR_DATABASE)
	}
	var admin Admin
	err := dbAdmin.FindId(bson.ObjectIdHex(id)).One(&admin)
	return &admin, err
}

func (m *MongoClient) GetAdminByUsername(username string) (*Admin, error) {
	var admin Admin
	err := dbAdmin.Find(bson.M{"username": username, "user_type": USER_TYPE_ADMIN}).One(&admin)
	return &admin, err
}

func EncodePassword(salt, passwd string) string {
	h := sha256.New()
	h.Write([]byte(passwd + salt))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
