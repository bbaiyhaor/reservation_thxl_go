package model

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

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

type Experience struct {
	Time     string `bson:"time"`
	Location string `bson:"location"`
	Teacher  string `bson:"teacher"`
}

func (e Experience) IsEmpty() bool {
	return len(e.Time) == 0 && len(e.Location) == 0 && len(e.Teacher) == 0
}

func (m *Model) AddStudent(username string, password string) (*Student, error) {
	if len(username) == 0 || len(password) == 0 {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("student")
	newStudent := &Student{
		Id:               bson.NewObjectId(),
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
		Username:         username,
		Password:         password,
		UserType:         STUDENT,
		CrisisLevel:      0,
		KeyCase:          make([]int, 5),
		MedicalDiagnosis: make([]int, 8),
	}
	if err := collection.Insert(newStudent); err != nil {
		return nil, err
	}
	return newStudent, nil
}

func (m *Model) UpsertStudent(student *Student) error {
	if student == nil || !student.Id.Valid() {
		return errors.New("字段不合法")
	}
	collection := m.mongo.C("student")
	student.UpdateTime = time.Now()
	_, err := collection.UpsertId(student.Id, student)
	return err
}

func (m *Model) GetStudentById(studentId string) (*Student, error) {
	if studentId == "" || !bson.IsObjectIdHex(studentId) {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("student")
	student := &Student{}
	if err := collection.FindId(bson.ObjectIdHex(studentId)).One(student); err != nil {
		return nil, err
	}
	return student, nil
}

func (m *Model) GetStudentByUsername(username string) (*Student, error) {
	if len(username) == 0 {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("student")
	student := &Student{}
	if err := collection.Find(bson.M{"username": username, "user_type": STUDENT}).One(student); err != nil {
		return nil, err
	}
	return student, nil
}

func (m *Model) GetStudentByArchiveNumber(archiveNumber string) (*Student, error) {
	if len(archiveNumber) == 0 {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("student")
	student := &Student{}
	if err := collection.Find(bson.M{"archive_number": archiveNumber}).One(student); err != nil {
		return nil, err
	}
	return student, nil
}
