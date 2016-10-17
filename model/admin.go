package model

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Admin struct {
	Id         bson.ObjectId `bson:"_id"`
	CreateTime time.Time     `bson:"create_time"`
	UpdateTime time.Time     `bson:"update_time"`
	Username   string        `bson:"username"` // Indexed
	Password   string        `bson:"password"`
	UserType   UserType      `bson:"user_type"`
}

func (m *Model) AddAdmin(username string, password string) (*Admin, error) {
	if len(username) == 0 || len(password) == 0 {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("admin")
	newAdmin := &Admin{
		Id:         bson.NewObjectId(),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		Username:   username,
		Password:   password,
		UserType:   ADMIN,
	}
	if err := collection.Insert(newAdmin); err != nil {
		return nil, err
	}
	return newAdmin, nil
}

func (m *Model) UpsertAdmin(admin *Admin) error {
	if admin == nil || !admin.Id.Valid() {
		return errors.New("字段不合法")
	}
	collection := m.mongo.C("admin")
	admin.UpdateTime = time.Now()
	_, err := collection.UpsertId(admin.Id, admin)
	return err
}

func (m *Model) GetAdminById(adminId string) (*Admin, error) {
	if len(adminId) == 0 || !bson.IsObjectIdHex(adminId) {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("admin")
	admin := &Admin{}
	if err := collection.FindId(bson.ObjectIdHex(adminId)).One(admin); err != nil {
		return nil, err
	}
	return admin, nil
}

func (m *Model) GetAdminByUsername(username string) (*Admin, error) {
	if len(username) == 0 {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("admin")
	admin := &Admin{}
	if err := collection.Find(bson.M{"username": username}).One(admin); err != nil {
		return nil, err
	}
	return admin, nil
}
