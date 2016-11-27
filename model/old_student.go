package model

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/utils"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type OldStudent struct {
	Id                bson.ObjectId `bson:"_id"`
	CreateTime        time.Time     `bson:"create_time"`
	UpdateTime        time.Time     `bson:"update_time"`
	Username          string        `bson:"username"` // Indexed
	Password          string        `bson:"password"` // will be deprecated soon
	EncryptedPassword string        `bson:"encrypted_password"`
	UserType          int           `bson:"user_type"`
	BindedTeacherId   string        `bson:"binded_teacher_id"` // Indexed
	ArchiveCategory   string        `bson:"archive_category"`
	ArchiveNumber     string        `bson:"archive_number"` // Indexed
	CrisisLevel       int           `bson:"crisis_level"`
	KeyCase           []int         `bson:"key_case"`          // deprecated
	MedicalDiagnosis  []int         `bson:"medical_diagnosis"` // deprecated
	Fullname          string        `bson:"fullname"`
	Gender            string        `bson:"gender"`
	Birthday          string        `bson:"birthday"`
	School            string        `bson:"school"`
	Grade             string        `bson:"grade"`
	CurrentAddress    string        `bson:"current_address"`
	FamilyAddress     string        `bson:"family_address"`
	Mobile            string        `bson:"mobile"`
	Email             string        `bson:"email"`
	OldExperience     OldExperience `bson:"experience"`
	FatherAge         string        `bson:"father_age"`
	FatherJob         string        `bson:"father_job"`
	FatherEdu         string        `bson:"father_edu"`
	MotherAge         string        `bson:"mother_age"`
	MotherJob         string        `bson:"mother_job"`
	MotherEdu         string        `bson:"mother_edu"`
	ParentMarriage    string        `bson:"parent_marriage"`
	Significant       string        `bson:"significant"`
	Problem           string        `bson:"problem"`
}

type OldExperience struct {
	Time     string `bson:"time"`
	Location string `bson:"location"`
	Teacher  string `bson:"teacher"`
}

func (e OldExperience) IsEmpty() bool {
	return e.Time == "" && e.Location == "" && e.Teacher == ""
}

func (m *MongoClient) AddOldStudent(username string, password string) (*OldStudent, error) {
	if username == "" || password == "" {
		return nil, errors.New("字段不合法")
	}
	encryptedPassword, err := utils.EncryptPassword(password)
	if err != nil {
		return nil, errors.New("加密出错，请联系技术支持")
	}
	collection := m.mongo.C("student")
	newOldStudent := &OldStudent{
		Id:                bson.NewObjectId(),
		CreateTime:        time.Now(),
		UpdateTime:        time.Now(),
		Username:          username,
		EncryptedPassword: encryptedPassword,
		UserType:          USER_TYPE_STUDENT,
		CrisisLevel:       0,
		KeyCase:           make([]int, 5),
		MedicalDiagnosis:  make([]int, 8),
	}
	if err := collection.Insert(newOldStudent); err != nil {
		return nil, err
	}
	return newOldStudent, nil
}

func (m *MongoClient) UpsertOldStudent(student *OldStudent) error {
	if student == nil || !student.Id.Valid() {
		return errors.New("字段不合法")
	}
	collection := m.mongo.C("student")
	student.UpdateTime = time.Now()
	_, err := collection.UpsertId(student.Id, student)
	return err
}

func (m *MongoClient) GetOldStudentById(id string) (*OldStudent, error) {
	if id == "" || !bson.IsObjectIdHex(id) {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("student")
	var student OldStudent
	if err := collection.FindId(bson.ObjectIdHex(id)).One(&student); err != nil {
		return nil, err
	}
	return &student, nil
}

func (m *MongoClient) GetOldStudentByUsername(username string) (*OldStudent, error) {
	if username == "" {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("student")
	var student OldStudent
	if err := collection.Find(bson.M{"username": username, "user_type": USER_TYPE_STUDENT}).One(&student); err != nil {
		return nil, err
	}
	return &student, nil
}

func (m *MongoClient) GetOldStudentByArchiveNumber(archiveNumber string) (*OldStudent, error) {
	if archiveNumber == "" {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("student")
	var student OldStudent
	if err := collection.Find(bson.M{"archive_number": archiveNumber}).One(&student); err != nil {
		return nil, err
	}
	return &student, nil
}
