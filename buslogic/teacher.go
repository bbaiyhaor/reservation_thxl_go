package buslogic

import (
	"github.com/shudiwsh2009/reservation_thxl_go/model"
	re "github.com/shudiwsh2009/reservation_thxl_go/rerror"
	"strconv"
	"strings"
	"time"
)

// 咨询师拉取反馈
func (w *Workflow) GetFeedbackByTeacher(reservationId string, sourceId string,
	userId string, userType int) (*model.Student, *model.Reservation, error) {
	if userId == "" {
		return nil, nil, re.NewRErrorCode("teacher not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_TEACHER {
		return nil, nil, re.NewRErrorCode("user is not teacher", nil, re.ERROR_NOT_AUTHORIZED)
	} else if reservationId == "" {
		return nil, nil, re.NewRErrorCodeContext("reservation id is empty", nil, re.ERROR_MISSING_PARAM, "reservation_id")
	} else if reservationId == sourceId {
		return nil, nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ERROR_FEEDBACK_AVAILABLE_RESERVATION)
	}
	teacher, err := w.mongoClient.GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.USER_TYPE_TEACHER {
		return nil, nil, re.NewRErrorCode("fail to get teacher", err, re.ERROR_DATABASE)
	}
	reservation, err := w.mongoClient.GetReservationById(reservationId)
	if err != nil || reservation == nil || reservation.Status == model.RESERVATION_STATUS_DELETED {
		return nil, nil, re.NewRErrorCode("fail to get reservation", err, re.ERROR_DATABASE)
	} else if reservation.StartTime.After(time.Now()) {
		return nil, nil, re.NewRErrorCode("cannot get feedback of future reservation", nil, re.ERROR_FEEDBACK_FUTURE_RESERVATION)
	} else if reservation.Status == model.RESERVATION_STATUS_AVAILABLE {
		return nil, nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ERROR_FEEDBACK_AVAILABLE_RESERVATION)
	} else if reservation.TeacherId != teacher.Id.Hex() {
		return nil, nil, re.NewRErrorCode("cannot get feedback of other one's reservation", nil, re.ERROR_FEEDBACK_OTHER_RESERVATION)
	}
	student, err := w.mongoClient.GetStudentById(reservation.StudentId)
	if err != nil || student == nil || student.UserType != model.USER_TYPE_STUDENT {
		return nil, nil, re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
	}
	return student, reservation, nil
}

// 咨询师提交反馈
func (w *Workflow) SubmitFeedbackByTeacher(reservationId string, sourceId string, category string, severity []int,
	medicalDiagnosis []int, crisis []int, hasCrisis bool, hasReservated bool, isSendNotify bool, schoolContact string,
	record string, crisisLevel string, userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("teacher not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_TEACHER {
		return nil, re.NewRErrorCode("user is not teacher", nil, re.ERROR_NOT_AUTHORIZED)
	} else if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation id is empty", nil, re.ERROR_MISSING_PARAM, "reservation_id")
	} else if category == "" {
		return nil, re.NewRErrorCodeContext("category is empty", nil, re.ERROR_MISSING_PARAM, "category")
	} else if len(severity) != len(model.FeedbackSeverity) {
		return nil, re.NewRErrorCodeContext("severity is not valid", nil, re.ERROR_INVALID_PARAM, "severity")
	} else if len(medicalDiagnosis) != len(model.FeedbackMedicalDiagnosis) {
		return nil, re.NewRErrorCodeContext("medical_diagnosis is not valid", nil, re.ERROR_INVALID_PARAM, "medical_diagnosis")
	} else if len(crisis) != len(model.FeedbackCrisis) {
		return nil, re.NewRErrorCodeContext("crisis is not valid", nil, re.ERROR_INVALID_PARAM, "crisis")
	} else if record == "" {
		return nil, re.NewRErrorCodeContext("record is empty", nil, re.ERROR_MISSING_PARAM, "record")
	} else if crisisLevel == "" {
		return nil, re.NewRErrorCodeContext("crisis_level is empty", nil, re.ERROR_MISSING_PARAM, "crisis_level")
	} else if reservationId == sourceId {
		return nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ERROR_FEEDBACK_AVAILABLE_RESERVATION)
	}
	crisisLevelInt, err := strconv.Atoi(crisisLevel)
	if err != nil || crisisLevelInt < 0 {
		return nil, re.NewRErrorCodeContext("crisis_level is not valid", err, re.ERROR_INVALID_PARAM, "crisis_level")
	}
	teacher, err := w.mongoClient.GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.USER_TYPE_TEACHER {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ERROR_DATABASE)
	}
	reservation, err := w.mongoClient.GetReservationById(reservationId)
	if err != nil || reservation == nil || reservation.Status == model.RESERVATION_STATUS_DELETED {
		return nil, re.NewRErrorCode("fail to get reservation", err, re.ERROR_DATABASE)
	} else if reservation.StartTime.After(time.Now()) {
		return nil, re.NewRErrorCode("cannot get feedback of future reservation", nil, re.ERROR_FEEDBACK_FUTURE_RESERVATION)
	} else if reservation.Status == model.RESERVATION_STATUS_AVAILABLE {
		return nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ERROR_FEEDBACK_AVAILABLE_RESERVATION)
	} else if reservation.TeacherId != teacher.Id.Hex() {
		return nil, re.NewRErrorCode("cannot get feedback of other one's reservation", nil, re.ERROR_FEEDBACK_OTHER_RESERVATION)
	}
	sendFeedbackSMS := reservation.TeacherFeedback.IsEmpty() && reservation.StudentFeedback.IsEmpty()
	reservation.TeacherFeedback = model.TeacherFeedback{
		Category:         category,
		Severity:         severity,
		MedicalDiagnosis: medicalDiagnosis,
		Crisis:           crisis,
		HasCrisis:        hasCrisis,
		HasReservated:    hasReservated,
		IsSendNotify:     isSendNotify,
		SchoolContact:    schoolContact,
		Record:           record,
	}
	student, err := w.mongoClient.GetStudentById(reservation.StudentId)
	if err != nil || student == nil || student.UserType != model.USER_TYPE_STUDENT {
		return nil, re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
	}
	student.CrisisLevel = crisisLevelInt
	if err = w.mongoClient.UpdateReservationAndStudent(reservation, student); err != nil {
		return nil, re.NewRErrorCode("fail to update reservation and student", err, re.ERROR_DATABASE)
	}
	if sendFeedbackSMS && !strings.HasPrefix(category, "H") {
		w.SendFeedbackSMS(reservation)
	}
	return reservation, nil
}

// 咨询师查看学生信息
func (w *Workflow) GetStudentInfoByTeacher(studentId string,
	userId string, userType int) (*model.Student, []*model.Reservation, error) {
	if userId == "" {
		return nil, nil, re.NewRErrorCode("teacher not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_TEACHER {
		return nil, nil, re.NewRErrorCode("user is not teacher", nil, re.ERROR_NOT_AUTHORIZED)
	} else if studentId == "" {
		return nil, nil, re.NewRErrorCodeContext("student id is empty", nil, re.ERROR_MISSING_PARAM, "student_id")
	}
	teacher, err := w.mongoClient.GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.USER_TYPE_TEACHER {
		return nil, nil, re.NewRErrorCode("fail to get teacher", err, re.ERROR_DATABASE)
	}
	student, err := w.mongoClient.GetStudentById(studentId)
	if err != nil || student == nil || student.UserType != model.USER_TYPE_STUDENT {
		return nil, nil, re.NewRErrorCode("fail to get student", err, re.ERROR_NO_STUDENT)
	}
	if student.BindedTeacherId != teacher.Id.Hex() {
		return nil, nil, re.NewRErrorCode("cannot view other one's student", nil, re.ERROR_TEACHER_VIEW_OTHER_STUDENT)
	}
	reservations, err := w.mongoClient.GetReservationsByStudentId(student.Id.Hex())
	if err != nil {
		return nil, nil, re.NewRErrorCode("fail to get reservations", err, re.ERROR_DATABASE)
	}
	return student, reservations, nil
}

// 咨询师查询学生信息
func (w *Workflow) QueryStudentInfoByTeacher(studentUsername string,
	userId string, userType int) (*model.Student, []*model.Reservation, error) {
	if userId == "" {
		return nil, nil, re.NewRErrorCode("teacher not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_TEACHER {
		return nil, nil, re.NewRErrorCode("user is not teacher", nil, re.ERROR_NOT_AUTHORIZED)
	} else if studentUsername == "" {
		return nil, nil, re.NewRErrorCodeContext("student username is empty", nil, re.ERROR_MISSING_PARAM, "student_username")
	}
	teacher, err := w.mongoClient.GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.USER_TYPE_TEACHER {
		return nil, nil, re.NewRErrorCode("fail to get teacher", err, re.ERROR_DATABASE)
	}
	student, err := w.mongoClient.GetStudentByUsername(studentUsername)
	if err != nil || student == nil || student.UserType != model.USER_TYPE_STUDENT {
		return nil, nil, re.NewRErrorCode("fail to get student", err, re.ERROR_NO_STUDENT)
	}
	if student.BindedTeacherId != teacher.Id.Hex() {
		return nil, nil, re.NewRErrorCode("cannot view other one's student", nil, re.ERROR_TEACHER_VIEW_OTHER_STUDENT)
	}
	reservations, err := w.mongoClient.GetReservationsByStudentId(student.Id.Hex())
	if err != nil {
		return nil, nil, re.NewRErrorCode("fail to get reservations", err, re.ERROR_DATABASE)
	}
	return student, reservations, nil
}

func (w *Workflow) WrapTeacher(teacher *model.Teacher) map[string]interface{} {
	var result = make(map[string]interface{})
	if teacher == nil {
		return result
	}
	result["id"] = teacher.Id.Hex()
	result["username"] = teacher.Username
	result["user_type"] = teacher.UserType
	result["fullname"] = teacher.Fullname
	result["mobile"] = teacher.Mobile
	return result
}
