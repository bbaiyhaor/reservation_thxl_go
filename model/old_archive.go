package model

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
)

type OldArchive struct {
	Id                 bson.ObjectId `bson:"_id"`
	StudentUsername    string        `bson:"student_username"` // Indexed
	OldArchiveCategory string        `bson:"archive_category"`
	OldArchiveNumber   string        `bson:"archive_number"`
}

func (m *MongoClient) AddOldArchive(studentUsername string, archiveCategory string, archiveNumber string) (*OldArchive, error) {
	if studentUsername == "" || archiveCategory == "" || archiveNumber == "" {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("archive")
	newOldArchive := &OldArchive{
		Id:                 bson.NewObjectId(),
		StudentUsername:    studentUsername,
		OldArchiveCategory: archiveCategory,
		OldArchiveNumber:   archiveNumber,
	}
	if err := collection.Insert(newOldArchive); err != nil {
		return nil, err
	}
	return newOldArchive, nil
}

func (m *MongoClient) GetOldArchiveByStudentUsername(studentUsername string) (*OldArchive, error) {
	if studentUsername == "" {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("archive")
	var archive OldArchive
	if err := collection.Find(bson.M{"student_username": studentUsername}).One(&archive); err != nil {
		return nil, err
	}
	return &archive, nil
}

func (m *MongoClient) GetOldArchiveByOldArchiveNumber(archiveNumber string) (*OldArchive, error) {
	if archiveNumber == "" {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("archive")
	var archive OldArchive
	if err := collection.Find(bson.M{"archive_number": archiveNumber}).One(&archive); err != nil {
		return nil, err
	}
	return &archive, nil
}
