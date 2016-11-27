package model

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type OldAdmin struct {
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

func (a *OldAdmin) ToAdmin() (*Admin, error) {
	now := time.Now()
	admin := &Admin{
		Id:        a.Id,
		Username:  a.Username,
		UserType:  a.UserType,
		Fullname:  a.Fullname,
		Mobile:    a.Mobile,
		CreatedAt: a.CreateTime,
		UpdatedAt: a.UpdateTime,
	}
	if admin.CreatedAt.Before(now.AddDate(-3, 0, 0)) {
		admin.CreatedAt = now
	}
	if admin.UpdatedAt.Before(now.AddDate(-3, 0, 0)) {
		admin.UpdatedAt = now
	}
	if a.EncryptedPassword != "" {
		admin.EncryptedPassword = a.EncryptedPassword
	} else {
		admin.Password = a.Password
		admin.PreInsert()
	}
	return admin, nil
}

func (m *MongoClient) UpsertOldAdmin(admin *OldAdmin) error {
	if admin == nil || !admin.Id.Valid() {
		return errors.New("字段不合法")
	}
	collection := m.mongo.C("admin")
	admin.UpdateTime = time.Now()
	_, err := collection.UpsertId(admin.Id, admin)
	return err
}

func (m *MongoClient) GetOldAdminById(id string) (*OldAdmin, error) {
	if id == "" || !bson.IsObjectIdHex(id) {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("admin")
	var admin OldAdmin
	if err := collection.FindId(bson.ObjectIdHex(id)).One(&admin); err != nil {
		return nil, err
	}
	return &admin, nil
}

func (m *MongoClient) GetOldAdminByUsername(username string) (*OldAdmin, error) {
	if len(username) == 0 {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("admin")
	var admin OldAdmin
	if err := collection.Find(bson.M{"username": username, "user_type": USER_TYPE_ADMIN}).One(&admin); err != nil {
		return nil, err
	}
	return &admin, nil
}
