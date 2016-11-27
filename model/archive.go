package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Index: student_username
// Index: archive_category + archive_number
type Archive struct {
	Id              bson.ObjectId `bson:"_id"`
	StudentUsername string        `bson:"student_username"`
	ArchiveCategory string        `bson:"archive_category"`
	ArchiveNumber   string        `bson:"archive_number"`
	CreatedAt       time.Time     `bson:"created_at"`
	UpdatedAt       time.Time     `bson:"updated_at"`
}

func (m *MongoClient) InsertArchive(archive *Archive) error {
	now := time.Now()
	archive.CreatedAt = now
	archive.UpdatedAt = now
	return dbArchive.Insert(archive)
}

func (m *MongoClient) UpdateArchive(archive *Archive) error {
	archive.UpdatedAt = time.Now()
	return dbArchive.UpdateId(archive.Id, archive)
}

func (m *MongoClient) UpdateArchiveWithoutTime(archive *Archive) error {
	return dbArchive.UpdateId(archive.Id, archive)
}

func (m *MongoClient) CountByStudentUsername(studentUsername string) (int, error) {
	return dbArchive.Find(bson.M{"student_username": studentUsername}).Count()
}

func (m *MongoClient) GetArchiveByStudentUsername(studentUsername string) (*Archive, error) {
	var archive Archive
	err := dbArchive.Find(bson.M{"student_username": studentUsername}).One(&archive)
	return &archive, err
}

func (m *MongoClient) GetArchiveByArchiveCategoryAndNumber(archiveCategory string, archiveNumber string) (*Archive, error) {
	var archive Archive
	err := dbArchive.Find(bson.M{"archive_category": archiveCategory, "archive_number": archiveNumber}).One(&archive)
	return &archive, err
}
