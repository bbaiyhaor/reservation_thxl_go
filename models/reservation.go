package models

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"time"
)

type ReservationStatus int

const (
	AVAILABLE ReservationStatus = 1 + iota
	RESERVATED
	FEEDBACK
	DELETED
	CLOSED
)

var reservationStatuses = [...]string{
	"AVAILABLE",
	"RESERVATED",
	"FEEDBACK",
	"DELETED",
	"CLOSED",
}

func (rs ReservationStatus) String() string {
	return reservationStatuses[rs-1]
}

type ReservationSource int

const (
	TIMETABLE ReservationSource = 1 + iota
	TEACHER_ADD
	ADMIN_ADD
)

var reservationSources = [...]string{
	"TIMETABLE",
	"TEACHER",
	"ADMIN",
}

func (rs ReservationSource) String() string {
	return reservationSources[rs-1]
}

type StudentFeedback struct {
	Scores []int `bson:"scores"`
}

func (sf StudentFeedback) IsEmpty() bool {
	return sf.Scores == nil || len(sf.Scores) == 0
}

func (sf StudentFeedback) ToJson() map[string]interface{} {
	var json = make(map[string]interface{})
	scores := ""
	for _, s := range sf.Scores {
		scores += strconv.Itoa(s) + " "
	}
	json["scores"] = scores
	return json
}

type TeacherFeedback struct {
	Category     string `bson:"category"`
	Participants []int  `bson:"participants"`
	Problem      string `bson:"problem"`
	Record       string `bson:"record"`
}

var Reservation_Participants = [...]string{
	"学生", "家长", "教师", "辅导员", "其他",
}

func (tf TeacherFeedback) IsEmpty() bool {
	return len(tf.Category) == 0 || tf.Participants == nil || len(tf.Participants) != len(Reservation_Participants) ||
		len(tf.Problem) == 0 || len(tf.Record) == 0
}

func (tf TeacherFeedback) ToJson() map[string]interface{} {
	var json = make(map[string]interface{})
	json["category"] = FeedbackAllCategory[tf.Category]
	participants := ""
	if len(tf.Participants) == len(Reservation_Participants) {
		for i := 0; i < len(tf.Participants); i++ {
			if tf.Participants[i] > 0 {
				participants += Reservation_Participants[i] + " "
			}
		}
	}
	json["participants"] = participants
	json["problem"] = tf.Problem
	json["record"] = tf.Record
	return json
}

type Reservation struct {
	Id              bson.ObjectId     `bson:"_id"`
	CreateTime      time.Time         `bson:"create_time"`
	UpdateTime      time.Time         `bson:"update_time"`
	StartTime       time.Time         `bson:"start_time"` // indexed
	EndTime         time.Time         `bson:"end_time"`
	Status          ReservationStatus `bson:"status"`
	Source          ReservationSource `bson:"source"`
	SourceId        string            `bson:"source_id"`
	TeacherId       string            `bson:"teacher_id"` // indexed
	StudentId       string            `bson:"student_id"` // indexed
	StudentFeedback StudentFeedback   `bson:"student_feedback"`
	TeacherFeedback TeacherFeedback   `bson:"teacher_feedback"`
}

type ReservationSlice []*Reservation

func (rs ReservationSlice) Len() int {
	return len(rs)
}

func (rs ReservationSlice) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs ReservationSlice) Less(i, j int) bool {
	if rs[i].StartTime.Equal(rs[j].StartTime) {
		return strings.Compare(rs[i].TeacherId, rs[j].TeacherId) < 0
	}
	return rs[i].StartTime.Before(rs[j].StartTime)
}

var FeedbackFirstCategory = map[string]interface{}{
	"A": "A 学业问题",
	"B": "B 情感问题",
	"C": "C 人际问题",
	"D": "D 发展问题",
	"E": "E 情绪问题",
	"F": "F 身心与行为问题",
	"G": "G 危机干预",
	"H": "H 心理测验",
	"I": "I 其他",
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
		"H1": "H1 人格测验与反馈",
		"H2": "H2 情绪测验与反馈",
		"H3": "H3 学业测验与反馈",
		"H4": "H4 职业测验与反馈",
	},
	"I": map[string]interface{}{
		"I1": "I1 躯体疾病转介",
		"I2": "I2 严重心理问题转介",
		"I3": "I3 转介至学习发展中心",
		"I4": "I4 转介至就业指导中心",
		"I5": "I5 反映学生情况",
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
	"H1": "H1 人格测验与反馈",
	"H2": "H2 情绪测验与反馈",
	"H3": "H3 学业测验与反馈",
	"H4": "H4 职业测验与反馈",
	"I1": "I1 躯体疾病转介",
	"I2": "I2 严重心理问题转介",
	"I3": "I3 转介至学习发展中心",
	"I4": "I4 转介至就业指导中心",
	"I5": "I5 反映学生情况",
	"Y1": "Y1 学习压力团体",
	"Y2": "Y2 人际关系团体",
	"Y3": "Y3 恋爱情感团体",
	"Y4": "Y4 辅导员团体",
	"Z1": "Z1 个体心理督导",
}

const CHECK_MESSAGE = "CHECK"
