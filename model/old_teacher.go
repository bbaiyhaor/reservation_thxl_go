package model

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/utils"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type OldTeacher struct {
	Id                bson.ObjectId `bson:"_id"`
	CreateTime        time.Time     `bson:"create_time"`
	UpdateTime        time.Time     `bson:"update_time"`
	Username          string        `bson:"username"` // Indexed
	Password          string        `bson:"password"` // will be deprecated soon
	EncryptedPassword string        `bson:"encrypted_password"`
	UserType          int           `bson:"user_type"`
	Fullname          string        `bson:"fullname"`
	Mobile            string        `bson:"mobile"`
}

func (t *OldTeacher) ToTeacher() (*Teacher, error) {
	now := time.Now()
	teacher := &Teacher{
		Id:        t.Id,
		Username:  t.Username,
		UserType:  t.UserType,
		Fullname:  t.Fullname,
		Mobile:    t.Mobile,
		CreatedAt: t.CreateTime,
		UpdatedAt: t.UpdateTime,
	}
	if teacher.CreatedAt.Before(now.AddDate(-3, 0, 0)) {
		teacher.CreatedAt = now
	}
	if teacher.UpdatedAt.Before(now.AddDate(-3, 0, 0)) {
		teacher.UpdatedAt = now
	}
	if t.EncryptedPassword != "" {
		teacher.EncryptedPassword = t.EncryptedPassword
	} else {
		teacher.Password = t.Password
		teacher.PreInsert()
	}
	return teacher, nil
}

func (m *MongoClient) AddOldTeacher(username string, password string, fullname string, mobile string) (*OldTeacher, error) {
	if username == "" || password == "" || fullname == "" || mobile == "" {
		return nil, errors.New("字段不合法")
	}
	encryptedPassword, err := utils.EncryptPassword(password)
	if err != nil {
		return nil, errors.New("加密出错，请联系技术支持")
	}
	collection := m.mongo.C("teacher")
	newOldTeacher := &OldTeacher{
		Id:                bson.NewObjectId(),
		CreateTime:        time.Now(),
		UpdateTime:        time.Now(),
		Username:          username,
		EncryptedPassword: encryptedPassword,
		Fullname:          fullname,
		Mobile:            mobile,
		UserType:          USER_TYPE_TEACHER,
	}
	if err := collection.Insert(newOldTeacher); err != nil {
		return nil, err
	}
	return newOldTeacher, nil
}

func (m *MongoClient) UpsertOldTeacher(teacher *OldTeacher) error {
	if teacher == nil || !teacher.Id.Valid() {
		return errors.New("字段不合法")
	}
	collection := m.mongo.C("teacher")
	teacher.UpdateTime = time.Now()
	_, err := collection.UpsertId(teacher.Id, teacher)
	return err
}

func (m *MongoClient) GetOldTeacherById(id string) (*OldTeacher, error) {
	if id == "" || !bson.IsObjectIdHex(id) {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("teacher")
	var teacher OldTeacher
	if err := collection.FindId(bson.ObjectIdHex(id)).One(&teacher); err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (m *MongoClient) GetOldTeacherByUsername(username string) (*OldTeacher, error) {
	if username == "" {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("teacher")
	var teacher OldTeacher
	if err := collection.Find(bson.M{"username": username, "user_type": USER_TYPE_TEACHER}).One(&teacher); err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (m *MongoClient) GetOldTeacherByFullname(fullname string) (*OldTeacher, error) {
	if fullname == "" {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("teacher")
	var teacher OldTeacher
	if err := collection.Find(bson.M{"fullname": fullname}).One(&teacher); err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (m *MongoClient) GetOldTeacherByMobile(mobile string) (*OldTeacher, error) {
	if mobile == "" {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("teacher")
	var teacher OldTeacher
	if err := collection.Find(bson.M{"mobile": mobile}).One(&teacher); err != nil {
		return nil, err
	}
	return &teacher, nil
}
