package model

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

const (
	RESERVATION_STATUS_AVAILABLE = 1 + iota
	RESERVATION_STATUS_RESERVATED
	RESERVATION_STATUS_FEEDBACK
	RESERVATION_STATUS_DELETED
	RESERVATION_STATUS_CLOSED

	RESERVATION_SOURCE_TIMETABLE = 1 + iota
	RESERVATION_SOURCE_TEACHER_ADD
	RESERVATION_SOURCE_ADMIN_ADD

	RESERVATION_STUDENT_FEEDBACK_SCORES_LENGTH = 5
)

type Reservation struct {
	Id              bson.ObjectId   `bson:"_id"`
	CreateTime      time.Time       `bson:"create_time"`
	UpdateTime      time.Time       `bson:"update_time"`
	StartTime       time.Time       `bson:"start_time"` // indexed
	EndTime         time.Time       `bson:"end_time"`
	Status          int             `bson:"status"`
	Source          int             `bson:"source"`
	SourceId        string          `bson:"source_id"`
	IsAdminSet      bool            `bson:"is_admin_set"`
	SendSms         bool            `bson:"send_sms"`
	TeacherId       string          `bson:"teacher_id"` // indexed
	StudentId       string          `bson:"student_id"` // indexed
	StudentFeedback StudentFeedback `bson:"student_feedback"`
	TeacherFeedback TeacherFeedback `bson:"teacher_feedback"`
}

type StudentFeedback struct {
	Scores []int `bson:"scores"`
}

func (sf StudentFeedback) IsEmpty() bool {
	return len(sf.Scores) == 0
}

func (sf StudentFeedback) ToStringJson() map[string]interface{} {
	var json = make(map[string]interface{})
	scores := ""
	for _, s := range sf.Scores {
		scores += strconv.Itoa(s) + " "
	}
	json["scores"] = scores
	return json
}

type TeacherFeedback struct {
	Category         string `bson:"category"`
	Participants     []int  `bson:"participants"`
	Problem          string `bson:"problem"`           // deprecated
	Emphasis         int    `bson:"emphasis"`          // 重点选项
	Severity         []int  `bson:"severity"`          // 严重程度
	MedicalDiagnosis []int  `bson:"medical_diagnosis"` // 疑似或明确的医疗诊断
	Crisis           []int  `bson:"crisis"`            // 危急情况
	Record           string `bson:"record"`
}

var (
	PARTICIPANTS      = [...]string{"学生", "家长", "教师", "辅导员", "其他"}
	SEVERITY          = [...]string{"缓考", "休学复学", "家属陪读", "家属不知情", "任何其他需要知会院系关注的原因"}
	MEDICAL_DIAGNOSIS = [...]string{"服药", "精神分裂", "双相情感障碍", "焦虑症（状态）", "抑郁症（状态）", "强迫症", "进食障碍", "失眠", "其他精神症状", "躯体疾病", "不遵医嘱"}
	CRISIS            = [...]string{"自伤", "伤害他人", "自杀念头", "自杀未遂"}
)

func (tf TeacherFeedback) IsEmpty() bool {
	return tf.Category == "" || len(tf.Participants) == 0 || len(tf.Severity) == 0 ||
		len(tf.MedicalDiagnosis) == 0 || len(tf.Crisis) == 0 || tf.Record == ""
}

func (tf TeacherFeedback) ToJson() map[string]interface{} {
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

func (tf TeacherFeedback) ToStringJson() map[string]interface{} {
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

func (m *Model) AddReservation(startTime time.Time, endTime time.Time, source int, sourceId string, teacherId string) (*Reservation, error) {
	collection := m.mongo.C("reservation")
	newReservation := &Reservation{
		Id:              bson.NewObjectId(),
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
		StartTime:       startTime,
		EndTime:         endTime,
		Status:          RESERVATION_STATUS_AVAILABLE,
		Source:          source,
		SourceId:        sourceId,
		TeacherId:       teacherId,
		StudentFeedback: StudentFeedback{},
		TeacherFeedback: TeacherFeedback{},
	}
	if err := collection.Insert(newReservation); err != nil {
		return nil, err
	}
	return newReservation, nil
}

func (m *Model) UpsertReservation(reservation *Reservation) error {
	if reservation == nil || !reservation.Id.Valid() {
		return errors.New("字段不合法")
	}
	collection := m.mongo.C("reservation")
	reservation.UpdateTime = time.Now()
	_, err := collection.UpsertId(reservation.Id, reservation)
	return err
}

func (m *Model) GetReservationById(id string) (*Reservation, error) {
	if id == "" || !bson.IsObjectIdHex(id) {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("reservation")
	var reservation *Reservation
	if err := collection.FindId(bson.ObjectIdHex(id)).One(reservation); err != nil {
		return nil, err
	}
	return reservation, nil
}

func (m *Model) GetReservationsByStudentId(studentId string) ([]*Reservation, error) {
	if studentId == "" || !bson.IsObjectIdHex(studentId) {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("reservation")
	var reservations []*Reservation
	if err := collection.Find(bson.M{"student_id": studentId,
		"status": bson.M{"$ne": RESERVATION_STATUS_DELETED}}).Sort("start_time").All(&reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func (m *Model) GetReservationsBetweenTime(from time.Time, to time.Time) ([]*Reservation, error) {
	collection := m.mongo.C("reservation")
	var reservations []*Reservation
	if err := collection.Find(bson.M{"start_time": bson.M{"$gte": from, "$lte": to},
		"status": bson.M{"$ne": RESERVATION_STATUS_DELETED}}).Sort("start_time").All(&reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func (m *Model) GetReservatedReservationsBetweenTime(from time.Time, to time.Time) ([]*Reservation, error) {
	collection := m.mongo.C("reservation")
	var reservations []*Reservation
	if err := collection.Find(bson.M{"start_time": bson.M{"$gte": from, "$lte": to},
		"status": RESERVATION_STATUS_RESERVATED}).Sort("start_time").All(&reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func (m *Model) GetReservationsAfterTime(from time.Time) ([]*Reservation, error) {
	collection := m.mongo.C("reservation")
	var reservations []*Reservation
	if err := collection.Find(bson.M{"start_time": bson.M{"$gte": from},
		"status": bson.M{"$ne": RESERVATION_STATUS_DELETED}}).Sort("start_time").All(&reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

var FeedbackFirstCategory = map[string]interface{}{
	"A": "A 学业问题",
	"B": "B 情感问题",
	"C": "C 人际问题",
	"D": "D 发展问题",
	"E": "E 情绪问题",
	"F": "F 身心与行为问题",
	"G": "G 危机干预",
	"H": "H 会商",
	"I": "I 心理测试与回访",
	"J": "J 转介",
	"Y": "Y 团体辅导",
	"Z": "Z 个体心理督导",
}

var FeedbackSecondCategory = map[string]interface{}{
	"A": map[string]interface{}{
		"A1": "A1 学业成就困扰",
		"A2": "A2 专业认同困扰",
		"A3": "A3 缓考评估",
		"A4": "A4 休学复学评估",
	},
	"B": map[string]interface{}{
		"B1": "B1 恋爱困扰",
		"B2": "B2 性困扰",
		"B3": "B3 性取向",
	},
	"C": map[string]interface{}{
		"C1": "C1 同伴人际",
		"C2": "C2 家庭人际",
		"C3": "C3 与辅导员人际",
		"C4": "C4 与教师人际",
	},
	"D": map[string]interface{}{
		"D1": "D1 就业困扰",
		"D2": "D2 事业探索",
		"D3": "D3 价值感与意义感",
		"D4": "D4 完美情结",
	},
	"E": map[string]interface{}{
		"E1": "E1 焦虑情绪",
		"E2": "E2 抑郁情绪",
		"E3": "E3 焦虑抑郁情绪",
	},
	"F": map[string]interface{}{
		"F1": "F1 睡眠问题",
		"F2": "F2 进食问题",
		"F3": "F3 身心问题",
		"F4": "F4 电脑依赖",
		"F5": "F5 强迫问题",
		"F6": "F6 品行问题",
	},
	"G": map[string]interface{}{
		"G1": "G1 应激状态干预",
		"G2": "G2 精神障碍发作期干预",
		"G3": "G3 精神障碍康复期干预",
		"G4": "G4 创伤后干预",
	},
	"H": map[string]interface{}{
		"H1": "H1 会商（与辅导员）",
		"H2": "H2 会商（与教师）",
		"H3": "H3 会商（与家属）",
		"H4": "H4 会商（与学生）",
		"H5": "H5 会商（与咨询师）",
		"H6": "H6 会商（联席会议）",
	},
	"I": map[string]interface{}{
		"I1": "I1 人格测验与反馈",
		"I2": "I2 情绪测验与反馈",
		"I3": "I3 学业测验与反馈",
		"I4": "I4 职业测验与反馈",
		"I5": "I5 新生回访适应正常",
	},
	"J": map[string]interface{}{
		"J1": "J1 躯体疾病转介",
		"J2": "J2 严重心理问题/精神疾病转介",
		"J3": "J3 转介至学习发展中心",
		"J4": "J4 转介至就业指导中心",
	},
	"Y": map[string]interface{}{
		"Y1": "Y1 学习压力团体",
		"Y2": "Y2 人际关系团体",
		"Y3": "Y3 恋爱情感团体",
		"Y4": "Y4 辅导员团体",
	},
	"Z": map[string]interface{}{
		"Z1": "Z1 个体心理督导",
	},
}

var FeedbackAllCategory = map[string]string{
	"A1": "A1 学业成就困扰",
	"A2": "A2 专业认同困扰",
	"A3": "A3 缓考评估",
	"A4": "A4 休学复学评估",
	"B1": "B1 恋爱困扰",
	"B2": "B2 性困扰",
	"B3": "B3 性取向",
	"C1": "C1 同伴人际",
	"C2": "C2 家庭人际",
	"C3": "C3 与辅导员人际",
	"C4": "C4 与教师人际",
	"D1": "D1 就业困扰",
	"D2": "D2 事业探索",
	"D3": "D3 价值感与意义感",
	"D4": "D4 完美情结",
	"E1": "E1 焦虑情绪",
	"E2": "E2 抑郁情绪",
	"E3": "E3 焦虑抑郁情绪",
	"F1": "F1 睡眠问题",
	"F2": "F2 进食问题",
	"F3": "F3 身心问题",
	"F4": "F4 电脑依赖",
	"F5": "F5 强迫问题",
	"F6": "F6 品行问题",
	"G1": "G1 应激状态干预",
	"G2": "G2 精神障碍发作期干预",
	"G3": "G3 精神障碍康复期干预",
	"G4": "G4 创伤后干预",
	"H1": "H1 会商（与辅导员）",
	"H2": "H2 会商（与教师）",
	"H3": "H3 会商（与家属）",
	"H4": "H4 会商（与学生）",
	"H5": "H5 会商（与咨询师）",
	"H6": "H6 会商（联席会议）",
	"I1": "I1 人格测验与反馈",
	"I2": "I2 情绪测验与反馈",
	"I3": "I3 学业测验与反馈",
	"I4": "I4 职业测验与反馈",
	"I5": "I5 新生回访适应正常",
	"J1": "J1 躯体疾病转介",
	"J2": "J2 严重心理问题/精神疾病转介",
	"J3": "J3 转介至学习发展中心",
	"J4": "J4 转介至就业指导中心",
	"Y1": "Y1 学习压力团体",
	"Y2": "Y2 人际关系团体",
	"Y3": "Y3 恋爱情感团体",
	"Y4": "Y4 辅导员团体",
	"Z1": "Z1 个体心理督导",
}
