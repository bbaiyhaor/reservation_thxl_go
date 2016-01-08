package models

import "gopkg.in/mgo.v2/bson"

type Archive struct {
	Id              bson.ObjectId `bson:"_id"`
	StudentUsername string        `bson:"student_username"` // Indexed
	ArchiveCategory string        `bson:"archive_category"`
	ArchiveNumber   string        `bson:"archive_number"`
}
