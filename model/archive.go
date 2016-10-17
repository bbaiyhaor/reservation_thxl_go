package model

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
)

type Archive struct {
	Id              bson.ObjectId `bson:"_id"`
	StudentUsername string        `bson:"student_username"` // Indexed
	ArchiveCategory string        `bson:"archive_category"`
	ArchiveNumber   string        `bson:"archive_number"`
}

func (m *Model) AddArchive(studentUsername string, archiveCategory string, archiveNumber string) (*Archive, error) {
	if len(studentUsername) == 0 || len(archiveCategory) == 0 || len(archiveNumber) == 0 {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("archive")
	newArchive := &Archive{
		Id:              bson.NewObjectId(),
		StudentUsername: studentUsername,
		ArchiveCategory: archiveCategory,
		ArchiveNumber:   archiveNumber,
	}
	if err := collection.Insert(newArchive); err != nil {
		return nil, err
	}
	return newArchive, nil
}

func (m *Model) GetArchiveByStudentUsername(studentUsername string) (*Archive, error) {
	if len(studentUsername) == 0 {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("archive")
	archive := &Archive{}
	if err := collection.Find(bson.M{"student_username": studentUsername}).One(archive); err != nil {
		return nil, err
	}
	return archive, nil
}

func (m *Model) GetArchiveByArchiveNumber(archiveNumber string) (*Archive, error) {
	if len(archiveNumber) == 0 {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("archive")
	archive := &Archive{}
	if err := collection.Find(bson.M{"archive_number": archiveNumber}).One(archive); err != nil {
		return nil, err
	}
	return archive, nil
}
