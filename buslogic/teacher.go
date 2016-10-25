package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"errors"
	"strconv"
	"time"
)

// 咨询师拉取反馈
func (w *Workflow) GetFeedbackByTeacher(reservationId string, sourceId string,
	userId string, userType int) (*model.Student, *model.Reservation, error) {
	if userId == "" {
		return nil, nil, errors.New("请先登录")
	} else if userType != model.USER_TYPE_TEACHER {
		return nil, nil, errors.New("权限不足")
	} else if reservationId == "" {
		return nil, nil, errors.New("咨询已下架")
	} else if reservationId == sourceId {
		return nil, nil, errors.New("咨询未被预约，不能反馈")
	}
	teacher, err := w.model.GetTeacherById(userId)
	if err != nil {
		return nil, nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != model.USER_TYPE_TEACHER {
		return nil, nil, errors.New("权限不足")
	}
	reservation, err := w.model.GetReservationById(reservationId)
	if err != nil || reservation.Status == model.RESERVATION_STATUS_DELETED {
		return nil, nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now()) {
		return nil, nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == model.RESERVATION_STATUS_AVAILABLE {
		return nil, nil, errors.New("咨询未被预约,不能反馈")
	} else if reservation.TeacherId != teacher.Id.Hex() {
		return nil, nil, errors.New("只能反馈本人开设的咨询")
	}
	student, err := w.model.GetStudentById(reservation.StudentId)
	if err != nil {
		return nil, nil, errors.New("获取数据失败")
	}
	return student, reservation, nil
}

// 咨询师提交反馈
func (w *Workflow) SubmitFeedbackByTeacher(reservationId string, sourceId string,
	category string, participants []int, emphasis string, severity []int, medicalDiagnosis []int, crisis []int,
	record string, crisisLevel string, userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, errors.New("请先登录")
	} else if userType != model.USER_TYPE_TEACHER {
		return nil, errors.New("权限不足")
	} else if reservationId == "" {
		return nil, errors.New("咨询已下架")
	} else if category == "" {
		return nil, errors.New("评估分类为空")
	} else if len(participants) != len(model.PARTICIPANTS) {
		return nil, errors.New("咨询参与者为空")
	} else if emphasis == "" {
		return nil, errors.New("重点明细为空")
	} else if len(severity) != len(model.SEVERITY) {
		return nil, errors.New("严重程度为空")
	} else if len(medicalDiagnosis) != len(model.MEDICAL_DIAGNOSIS) {
		return nil, errors.New("医疗诊断为空")
	} else if len(crisis) != len(model.CRISIS) {
		return nil, errors.New("危机情况为空")
	} else if record == "" {
		return nil, errors.New("咨询记录为空")
	} else if crisisLevel == "" {
		return nil, errors.New("危机等级为空")
	} else if reservationId == sourceId {
		return nil, errors.New("咨询未被预约，不能反馈")
	}
	emphasisInt, err := strconv.Atoi(emphasis)
	if err != nil || emphasisInt < 0 {
		return nil, errors.New("重点明细错误")
	}
	crisisLevelInt, err := strconv.Atoi(crisisLevel)
	if err != nil || crisisLevelInt < 0 {
		return nil, errors.New("危机等级错误")
	}
	teacher, err := w.model.GetTeacherById(userId)
	if err != nil {
		return nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != model.USER_TYPE_TEACHER {
		return nil, errors.New("权限不足")
	}
	reservation, err := w.model.GetReservationById(reservationId)
	if err != nil || reservation.Status == model.RESERVATION_STATUS_DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == model.RESERVATION_STATUS_AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	} else if reservation.TeacherId != teacher.Id.Hex() {
		return nil, errors.New("只能反馈本人开设的咨询")
	}
	sendFeedbackSMS := reservation.TeacherFeedback.IsEmpty() && reservation.StudentFeedback.IsEmpty()
	reservation.TeacherFeedback = model.TeacherFeedback{
		Category:         category,
		Participants:     participants,
		Emphasis:         emphasisInt,
		Severity:         severity,
		MedicalDiagnosis: medicalDiagnosis,
		Crisis:           crisis,
		Record:           record,
	}
	student, err := w.model.GetStudentById(reservation.StudentId)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	student.CrisisLevel = crisisLevelInt
	if w.model.UpsertReservation(reservation) != nil || w.model.UpsertStudent(student) != nil {
		return nil, errors.New("获取数据失败")
	}
	if sendFeedbackSMS && participants[0] > 0 {
		w.SendFeedbackSMS(reservation)
	}
	return reservation, nil
}

// 咨询师查看学生信息
func (w *Workflow) GetStudentInfoByTeacher(studentId string,
	userId string, userType int) (*model.Student, []*model.Reservation, error) {
	if userId == "" {
		return nil, nil, errors.New("请先登录")
	} else if userType != model.USER_TYPE_TEACHER {
		return nil, nil, errors.New("权限不足")
	} else if studentId == "" {
		return nil, nil, errors.New("咨询未被预约")
	}
	teacher, err := w.model.GetTeacherById(userId)
	if err != nil {
		return nil, nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != model.USER_TYPE_TEACHER {
		return nil, nil, errors.New("权限不足")
	}
	student, err := w.model.GetStudentById(studentId)
	if err != nil {
		return nil, nil, errors.New("学生未注册")
	}
	if student.BindedTeacherId != teacher.Id.Hex() {
		return nil, nil, errors.New("只能查看本人绑定的学生")
	}
	reservations, err := w.model.GetReservationsByStudentId(student.Id.Hex())
	if err != nil {
		return nil, nil, errors.New("获取数据失败")
	}
	return student, reservations, nil
}

// 咨询师查询学生信息
func (w *Workflow) QueryStudentInfoByTeacher(studentUsername string,
	userId string, userType int) (*model.Student, []*model.Reservation, error) {
	if userId == "" {
		return nil, nil, errors.New("请先登录")
	} else if userType != model.USER_TYPE_TEACHER {
		return nil, nil, errors.New("权限不足")
	} else if studentUsername == "" {
		return nil, nil, errors.New("学号为空")
	}
	teacher, err := w.model.GetTeacherById(userId)
	if err != nil {
		return nil, nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != model.USER_TYPE_TEACHER {
		return nil, nil, errors.New("权限不足")
	}
	student, err := w.model.GetStudentByUsername(studentUsername)
	if err != nil || student.UserType != model.USER_TYPE_STUDENT {
		return nil, nil, errors.New("学生未注册")
	}
	if student.BindedTeacherId != teacher.Id.Hex() {
		return nil, nil, errors.New("只能查看本人绑定的学生")
	}
	reservations, err := w.model.GetReservationsByStudentId(student.Id.Hex())
	if err != nil {
		return nil, nil, errors.New("获取数据失败")
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
