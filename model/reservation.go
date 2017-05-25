package model

import (
	re "github.com/shudiwsh2009/reservation_thxl_go/rerror"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"time"
)

const (
	RESERVATION_STATUS_AVAILABLE  = 1
	RESERVATION_STATUS_RESERVATED = 2
	RESERVATION_STATUS_FEEDBACK   = 3
	RESERVATION_STATUS_DELETED    = 4
	RESERVATION_STATUS_CLOSED     = 5

	RESERVATION_SOURCE_TIMETABLE   = 1
	RESERVATION_SOURCE_TEACHER_ADD = 2
	RESERVATION_SOURCE_ADMIN_ADD   = 3

	RESERVATION_STUDENT_FEEDBACK_SCORES_LENGTH = 5
)

// Index: student_id + status + start_time
// Index: start_time + end_time + status
// Index: start_time + status
type Reservation struct {
	Id              bson.ObjectId   `bson:"_id"`
	StartTime       time.Time       `bson:"start_time"`
	EndTime         time.Time       `bson:"end_time"`
	Status          int             `bson:"status"`
	Source          int             `bson:"source"`
	SourceId        string          `bson:"source_id"`
	IsAdminSet      bool            `bson:"is_admin_set"`
	SendSms         bool            `bson:"send_sms"`
	TeacherId       string          `bson:"teacher_id"`
	StudentId       string          `bson:"student_id"`
	StudentFeedback StudentFeedback `bson:"student_feedback"`
	TeacherFeedback TeacherFeedback `bson:"teacher_feedback"`
	CreatedAt       time.Time       `bson:"created_at"`
	UpdatedAt       time.Time       `bson:"updated_at"`
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
	Participants     []int  `bson:"participants"`      // deprecated
	Problem          string `bson:"problem"`           // deprecated
	Emphasis         int    `bson:"emphasis"`          // deprecated, 重点选项
	Severity         []int  `bson:"severity"`          // 严重程度
	MedicalDiagnosis []int  `bson:"medical_diagnosis"` // 疑似或明确的医疗诊断
	Crisis           []int  `bson:"crisis"`            // 危急情况
	HasCrisis        bool   `bson:"has_crisis"`        // 本次会谈是否有危机
	HasReservated    bool   `bson:"has_reservated"`    // 本次会谈是否有预约
	IsSendNotify     bool   `bson:"is_send_notify"`    // 是否发危机通告
	SchoolContact    string `bson:"school_contact"`    // 院系联系人
	Record           string `bson:"record"`
}

var (
	FeedbackSeverity         = [...]string{"缓考", "休学复学", "家属陪读", "家属不知情", "任何其他需要知会院系关注的原因", "试读期"}
	FeedbackMedicalDiagnosis = [...]string{"服药", "精神分裂", "双相情感障碍", "焦虑症（状态）", "抑郁症（状态）", "强迫症", "进食障碍", "失眠", "其他精神症状", "躯体疾病", "不遵医嘱"}
	FeedbackCrisis           = [...]string{"自伤", "伤害他人", "自杀念头", "自杀未遂"}
)

func (tf TeacherFeedback) IsEmpty() bool {
	return tf.Category == "" || len(tf.Severity) == 0 || len(tf.MedicalDiagnosis) == 0 || len(tf.Crisis) == 0 || tf.Record == ""
}

func (tf TeacherFeedback) GetServerityStr() string {
	var severity []string
	for i := 0; i < len(tf.Severity); i++ {
		if tf.Severity[i] > 0 {
			severity = append(severity, FeedbackSeverity[i])
		}
	}
	return strings.Join(severity, "、")
}

func (tf TeacherFeedback) GetMedicalDiagnosisStr() string {
	var medicalDiagnosis []string
	for i := 0; i < len(tf.MedicalDiagnosis); i++ {
		if tf.MedicalDiagnosis[i] > 0 {
			medicalDiagnosis = append(medicalDiagnosis, FeedbackMedicalDiagnosis[i])
		}
	}
	return strings.Join(medicalDiagnosis, "、")
}

func (tf TeacherFeedback) GetCrisisStr() string {
	var crisis []string
	for i := 0; i < len(tf.Crisis); i++ {
		if tf.Crisis[i] > 0 {
			crisis = append(crisis, FeedbackCrisis[i])
		}
	}
	return strings.Join(crisis, "、")
}

func (tf TeacherFeedback) GetEmphasisStr() string {
	severity := tf.GetServerityStr()
	medicalDiagnosis := tf.GetMedicalDiagnosisStr()
	crisis := tf.GetCrisisStr()
	var emphasis []string
	if severity != "" {
		emphasis = append(emphasis, severity)
	}
	if medicalDiagnosis != "" {
		emphasis = append(emphasis, medicalDiagnosis)
	}
	if crisis != "" {
		emphasis = append(emphasis, crisis)
	}
	return strings.Join(emphasis, "、")
}

func (tf TeacherFeedback) ToJson() map[string]interface{} {
	var feedback = make(map[string]interface{})
	feedback["category"] = tf.Category
	feedback["severity"] = tf.Severity
	feedback["medical_diagnosis"] = tf.MedicalDiagnosis
	feedback["crisis"] = tf.Crisis
	feedback["has_crisis"] = tf.HasCrisis
	feedback["has_reservated"] = tf.HasReservated
	feedback["is_send_notify"] = tf.IsSendNotify
	feedback["school_contact"] = tf.SchoolContact
	feedback["record"] = tf.Record
	return feedback
}

func (tf TeacherFeedback) ToStringJson() map[string]interface{} {
	var json = make(map[string]interface{})
	json["category"] = FeedbackAllCategoryMap[tf.Category]
	json["severity"] = tf.GetServerityStr()
	json["medical_diagnosis"] = tf.GetMedicalDiagnosisStr()
	json["crisis"] = tf.GetCrisisStr()
	json["has_crisis"] = tf.HasCrisis
	json["has_reservated"] = tf.HasReservated
	json["is_send_notify"] = tf.IsSendNotify
	json["school_contact"] = tf.SchoolContact
	json["record"] = tf.Record
	return json
}

func (m *MongoClient) InsertReservation(reservation *Reservation) error {
	now := time.Now()
	reservation.Id = bson.NewObjectId()
	reservation.CreatedAt = now
	reservation.UpdatedAt = now
	return dbReservation.Insert(reservation)
}

func (m *MongoClient) InsertReservationAndUpdateTimedReservation(reservation *Reservation, timedReservation *TimedReservation) error {
	err := m.InsertReservation(reservation)
	if err != nil {
		return err
	}
	return m.UpdateTimedReservation(timedReservation)
}

func (m *MongoClient) UpdateReservation(reservation *Reservation) error {
	reservation.UpdatedAt = time.Now()
	return dbReservation.UpdateId(reservation.Id, reservation)
}

func (m *MongoClient) UpdateReservationWithoutUpdatedTime(reservation *Reservation) error {
	return dbReservation.UpdateId(reservation.Id, reservation)
}

func (m *MongoClient) UpdateReservationAndTimedReservation(reservation *Reservation, timedReservation *TimedReservation) error {
	err := m.UpdateReservation(reservation)
	if err != nil {
		return err
	}
	return m.UpdateTimedReservation(timedReservation)
}

func (m *MongoClient) UpdateReservationAndStudent(reservation *Reservation, student *Student) error {
	err := m.UpdateReservation(reservation)
	if err != nil {
		return err
	}
	return m.UpdateStudent(student)
}

func (m *MongoClient) GetAllReservations() ([]*Reservation, error) {
	var reservations []*Reservation
	err := dbReservation.Find(bson.M{}).All(&reservations)
	return reservations, err
}

func (m *MongoClient) GetReservationById(id string) (*Reservation, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, re.NewRErrorCode("id is not valid", nil, re.ERROR_DATABASE)
	}
	var reservation Reservation
	err := dbReservation.FindId(bson.ObjectIdHex(id)).One(&reservation)
	if err == mgo.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &reservation, nil
	}
}

func (m *MongoClient) GetReservationsByStudentId(studentId string) ([]*Reservation, error) {
	var reservations []*Reservation
	err := dbReservation.Find(bson.M{"student_id": studentId,
		"status": bson.M{"$ne": RESERVATION_STATUS_DELETED}}).Sort("start_time").All(&reservations)
	return reservations, err
}

func (m *MongoClient) GetReservationsBetweenTime(start time.Time, end time.Time) ([]*Reservation, error) {
	var reservations []*Reservation
	err := dbReservation.Find(bson.M{"start_time": bson.M{"$gte": start, "$lt": end},
		"status": bson.M{"$ne": RESERVATION_STATUS_DELETED}}).Sort("start_time").All(&reservations)
	return reservations, err
}

func (m *MongoClient) GetReservatedReservationsBetweenTime(start time.Time, end time.Time) ([]*Reservation, error) {
	var reservations []*Reservation
	err := dbReservation.Find(bson.M{"start_time": bson.M{"$gte": start, "$lt": end},
		"status": RESERVATION_STATUS_RESERVATED}).Sort("start_time").All(&reservations)
	return reservations, err
}

func (m *MongoClient) GetReservationsAfterTime(start time.Time) ([]*Reservation, error) {
	var reservations []*Reservation
	err := dbReservation.Find(bson.M{"start_time": bson.M{"$gte": start},
		"status": bson.M{"$ne": RESERVATION_STATUS_DELETED}}).Sort("start_time").All(&reservations)
	return reservations, err
}

func (m *MongoClient) GetReservationsBySchoolContact(schoolContact string) ([]*Reservation, error) {
	var reservations []*Reservation
	err := dbReservation.Find(bson.M{"teacher_feedback.school_contact": schoolContact,
		"status": bson.M{"$ne": RESERVATION_STATUS_DELETED}}).Sort("start_time").All(&reservations)
	return reservations, err
}

var FeedbackFirstCategoryKeys = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "Y", "Z",
}

var FeedbackFirstCategoryMap = map[string]string{
	"":  "请选择",
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
	"Z": "Z 个体督导",
}

var FeedbackSecondCategoryMap = map[string]map[string]string{
	"": map[string]string{
		"": "请选择",
	},
	"A": map[string]string{
		"":   "请选择",
		"A1": "A1 学业成就困扰",
		"A2": "A2 专业认同困扰",
		"A3": "A3 缓考评估",
		"A4": "A4 休学复学评估",
	},
	"B": map[string]string{
		"":   "请选择",
		"B1": "B1 恋爱困扰",
		"B2": "B2 性困扰",
		"B3": "B3 性取向",
	},
	"C": map[string]string{
		"":   "请选择",
		"C1": "C1 同伴人际",
		"C2": "C2 家庭人际",
		"C3": "C3 与辅导员人际",
		"C4": "C4 与教师人际",
	},
	"D": map[string]string{
		"":   "请选择",
		"D1": "D1 就业困扰",
		"D2": "D2 事业探索",
		"D3": "D3 价值感与意义感",
		"D4": "D4 完美情结",
	},
	"E": map[string]string{
		"":   "请选择",
		"E1": "E1 焦虑情绪",
		"E2": "E2 抑郁情绪",
		"E3": "E3 焦虑抑郁情绪",
	},
	"F": map[string]string{
		"":   "请选择",
		"F1": "F1 睡眠问题",
		"F2": "F2 进食问题",
		"F3": "F3 身心问题",
		"F4": "F4 电脑依赖",
		"F5": "F5 强迫问题",
		"F6": "F6 品行问题",
	},
	"G": map[string]string{
		"":   "请选择",
		"G1": "G1 应激状态干预",
		"G2": "G2 精神障碍发作期干预",
		"G3": "G3 精神障碍康复期干预",
		"G4": "G4 创伤后干预",
	},
	"H": map[string]string{
		"":   "请选择",
		"H1": "H1 会商（与辅导员）",
		"H2": "H2 会商（与教师）",
		"H3": "H3 会商（与家属）",
		"H4": "H4 会商（与学生）",
		"H5": "H5 会商（与咨询师）",
		"H6": "H6 会商（联席会议）",
	},
	"I": map[string]string{
		"":   "请选择",
		"I1": "I1 人格测验与反馈",
		"I2": "I2 情绪测验与反馈",
		"I3": "I3 学业测验与反馈",
		"I4": "I4 职业测验与反馈",
		"I5": "I5 新生回访适应正常",
	},
	"J": map[string]string{
		"":   "请选择",
		"J1": "J1 躯体疾病转介",
		"J2": "J2 严重心理问题/精神疾病转介",
		"J3": "J3 转介至学习发展中心",
		"J4": "J4 转介至就业指导中心",
	},
	"Y": map[string]string{
		"":   "请选择",
		"Y1": "Y1 学习压力团体",
		"Y2": "Y2 人际关系团体",
		"Y3": "Y3 恋爱情感团体",
		"Y4": "Y4 辅导员团体",
	},
	"Z": map[string]string{
		"":   "请选择",
		"Z1": "Z1 个体心理督导",
	},
}

var FeedbackAllCategoryMap = map[string]string{
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
