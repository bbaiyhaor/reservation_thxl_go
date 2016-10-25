package model

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/utils"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Teacher struct {
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

func (m *Model) AddTeacher(username string, password string, fullname string, mobile string) (*Teacher, error) {
	if username == "" || password == "" || fullname == "" || mobile == "" {
		return nil, errors.New("字段不合法")
	}
	encryptedPassword, err := utils.EncryptPassword(password)
	if err != nil {
		return nil, errors.New("加密出错，请联系技术支持")
	}
	collection := m.mongo.C("teacher")
	newTeacher := &Teacher{
		Id:                bson.NewObjectId(),
		CreateTime:        time.Now(),
		UpdateTime:        time.Now(),
		Username:          username,
		EncryptedPassword: encryptedPassword,
		Fullname:          fullname,
		Mobile:            mobile,
		UserType:          USER_TYPE_TEACHER,
	}
	if err := collection.Insert(newTeacher); err != nil {
		return nil, err
	}
	return newTeacher, nil
}

func (m *Model) UpsertTeacher(teacher *Teacher) error {
	if teacher == nil || !teacher.Id.Valid() {
		return errors.New("字段不合法")
	}
	collection := m.mongo.C("teacher")
	teacher.UpdateTime = time.Now()
	_, err := collection.UpsertId(teacher.Id, teacher)
	return err
}

func (m *Model) GetTeacherById(id string) (*Teacher, error) {
	if id == "" || !bson.IsObjectIdHex(id) {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("teacher")
	var teacher *Teacher
	if err := collection.FindId(bson.ObjectIdHex(id)).One(teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func (m *Model) GetTeacherByUsername(username string) (*Teacher, error) {
	if username == "" {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("teacher")
	var teacher *Teacher
	if err := collection.Find(bson.M{"username": username, "user_type": USER_TYPE_TEACHER}).One(teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func (m *Model) GetTeacherByFullname(fullname string) (*Teacher, error) {
	if fullname == "" {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("teacher")
	var teacher *Teacher
	if err := collection.Find(bson.M{"fullname": fullname}).One(teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func (m *Model) GetTeacherByMobile(mobile string) (*Teacher, error) {
	if mobile == "" {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("teacher")
	var teacher *Teacher
	if err := collection.Find(bson.M{"mobile": mobile}).One(teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}
