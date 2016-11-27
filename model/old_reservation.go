package model

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

type OldReservation struct {
	Id                 bson.ObjectId      `bson:"_id"`
	CreateTime         time.Time          `bson:"create_time"`
	UpdateTime         time.Time          `bson:"update_time"`
	StartTime          time.Time          `bson:"start_time"` // indexed
	EndTime            time.Time          `bson:"end_time"`
	Status             int                `bson:"status"`
	Source             int                `bson:"source"`
	SourceId           string             `bson:"source_id"`
	IsAdminSet         bool               `bson:"is_admin_set"`
	SendSms            bool               `bson:"send_sms"`
	TeacherId          string             `bson:"teacher_id"` // indexed
	StudentId          string             `bson:"student_id"` // indexed
	OldStudentFeedback OldStudentFeedback `bson:"student_feedback"`
	OldTeacherFeedback OldTeacherFeedback `bson:"teacher_feedback"`
}

type OldStudentFeedback struct {
	Scores []int `bson:"scores"`
}

func (sf OldStudentFeedback) IsEmpty() bool {
	return len(sf.Scores) == 0
}

func (sf OldStudentFeedback) ToStringJson() map[string]interface{} {
	var json = make(map[string]interface{})
	scores := ""
	for _, s := range sf.Scores {
		scores += strconv.Itoa(s) + " "
	}
	json["scores"] = scores
	return json
}

type OldTeacherFeedback struct {
	Category         string `bson:"category"`
	Participants     []int  `bson:"participants"`
	Problem          string `bson:"problem"`           // deprecated
	Emphasis         int    `bson:"emphasis"`          // 重点选项
	Severity         []int  `bson:"severity"`          // 严重程度
	MedicalDiagnosis []int  `bson:"medical_diagnosis"` // 疑似或明确的医疗诊断
	Crisis           []int  `bson:"crisis"`            // 危急情况
	Record           string `bson:"record"`
}

func (tf OldTeacherFeedback) IsEmpty() bool {
	return tf.Category == "" || len(tf.Participants) == 0 || len(tf.Severity) == 0 ||
		len(tf.MedicalDiagnosis) == 0 || len(tf.Crisis) == 0 || tf.Record == ""
}

func (tf OldTeacherFeedback) ToJson() map[string]interface{} {
	var feedback = make(map[string]interface{})
	feedback["category"] = tf.Category
	if len(tf.Participants) != len(PARTICIPANTS) {
		feedback["participants"] = make([]int, len(PARTICIPANTS))
	} else {
		feedback["participants"] = tf.Participants
	}
	feedback["emphasis"] = tf.Emphasis
	if len(tf.Severity) != len(SEVERITY) {
		feedback["severity"] = make([]int, len(SEVERITY))
	} else {
		feedback["severity"] = tf.Severity
	}
	if len(tf.MedicalDiagnosis) != len(MEDICAL_DIAGNOSIS) {
		feedback["medical_diagnosis"] = make([]int, len(MEDICAL_DIAGNOSIS))
	} else {
		feedback["medical_diagnosis"] = tf.MedicalDiagnosis
	}
	if len(tf.Crisis) != len(CRISIS) {
		feedback["crisis"] = make([]int, len(CRISIS))
	} else {
		feedback["crisis"] = tf.Crisis
	}
	feedback["record"] = tf.Record
	return feedback
}

func (tf OldTeacherFeedback) ToStringJson() map[string]interface{} {
	var json = make(map[string]interface{})
	json["category"] = FeedbackAllCategory[tf.Category]
	var participants string
	if len(tf.Participants) == len(PARTICIPANTS) {
		for i := 0; i < len(tf.Participants); i++ {
			if tf.Participants[i] > 0 {
				participants += PARTICIPANTS[i] + " "
			}
		}
	}
	json["participants"] = participants
	json["emphasis"] = strconv.Itoa(tf.Emphasis)
	var severity string
	if len(tf.Severity) == len(SEVERITY) {
		for i := 0; i < len(tf.Severity); i++ {
			if tf.Severity[i] > 0 {
				severity += SEVERITY[i] + " "
			}
		}
	}
	json["severity"] = severity
	var medicalDiagnosis string
	if len(tf.MedicalDiagnosis) == len(MEDICAL_DIAGNOSIS) {
		for i := 0; i < len(tf.MedicalDiagnosis); i++ {
			if tf.MedicalDiagnosis[i] > 0 {
				medicalDiagnosis += MEDICAL_DIAGNOSIS[i] + " "
			}
		}
	}
	json["medical_diagnosis"] = medicalDiagnosis
	var crisis string
	if len(tf.Crisis) == len(CRISIS) {
		for i := 0; i < len(tf.Crisis); i++ {
			if tf.Crisis[i] > 0 {
				crisis += CRISIS[i] + " "
			}
		}
	}
	json["crisis"] = crisis
	json["record"] = tf.Record
	return json
}

func (m *MongoClient) AddOldReservation(startTime time.Time, endTime time.Time, source int, sourceId string, teacherId string) (*OldReservation, error) {
	collection := m.mongo.C("reservation")
	newOldReservation := &OldReservation{
		Id:                 bson.NewObjectId(),
		CreateTime:         time.Now(),
		UpdateTime:         time.Now(),
		StartTime:          startTime,
		EndTime:            endTime,
		Status:             RESERVATION_STATUS_AVAILABLE,
		Source:             source,
		SourceId:           sourceId,
		TeacherId:          teacherId,
		OldStudentFeedback: OldStudentFeedback{},
		OldTeacherFeedback: OldTeacherFeedback{},
	}
	if err := collection.Insert(newOldReservation); err != nil {
		return nil, err
	}
	return newOldReservation, nil
}

func (m *MongoClient) UpsertOldReservation(reservation *OldReservation) error {
	if reservation == nil || !reservation.Id.Valid() {
		return errors.New("字段不合法")
	}
	collection := m.mongo.C("reservation")
	reservation.UpdateTime = time.Now()
	_, err := collection.UpsertId(reservation.Id, reservation)
	return err
}

func (m *MongoClient) GetOldReservationById(id string) (*OldReservation, error) {
	if id == "" || !bson.IsObjectIdHex(id) {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("reservation")
	var reservation OldReservation
	if err := collection.FindId(bson.ObjectIdHex(id)).One(&reservation); err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (m *MongoClient) GetOldReservationsByStudentId(studentId string) ([]*OldReservation, error) {
	if studentId == "" || !bson.IsObjectIdHex(studentId) {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("reservation")
	var reservations []*OldReservation
	if err := collection.Find(bson.M{"student_id": studentId,
		"status": bson.M{"$ne": RESERVATION_STATUS_DELETED}}).Sort("start_time").All(&reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func (m *MongoClient) GetOldReservationsBetweenTime(from time.Time, to time.Time) ([]*OldReservation, error) {
	collection := m.mongo.C("reservation")
	var reservations []*OldReservation
	if err := collection.Find(bson.M{"start_time": bson.M{"$gte": from, "$lte": to},
		"status": bson.M{"$ne": RESERVATION_STATUS_DELETED}}).Sort("start_time").All(&reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func (m *MongoClient) GetReservatedOldReservationsBetweenTime(from time.Time, to time.Time) ([]*OldReservation, error) {
	collection := m.mongo.C("reservation")
	var reservations []*OldReservation
	if err := collection.Find(bson.M{"start_time": bson.M{"$gte": from, "$lte": to},
		"status": RESERVATION_STATUS_RESERVATED}).Sort("start_time").All(&reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func (m *MongoClient) GetOldReservationsAfterTime(from time.Time) ([]*OldReservation, error) {
	collection := m.mongo.C("reservation")
	var reservations []*OldReservation
	if err := collection.Find(bson.M{"start_time": bson.M{"$gte": from},
		"status": bson.M{"$ne": RESERVATION_STATUS_DELETED}}).Sort("start_time").All(&reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}
