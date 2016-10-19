package model

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Admin struct {
	Id                bson.ObjectId `bson:"_id"`
	CreateTime        time.Time     `bson:"create_time"`
	UpdateTime        time.Time     `bson:"update_time"`
	Username          string        `bson:"username"` // Indexed
	Password          string        `bson:"password"` // will be deprecated soon
	EncryptedPassword string        `bson:"encrypted_password"`
	Fullname          string        `bson:"fullname"`
	Mobile            string        `bson:"mobile"`
	UserType          int           `bson:"user_type"`
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
	if err := collection.Find(bson.M{"username": username, "user_type": USER_TYPE_ADMIN}).One(admin); err != nil {
		return nil, err
	}
	return admin, nil
}
