package buslogic

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"strings"
	"github.com/shudiwsh2009/reservation_thxl_go/sms"
)

type TeacherLogic struct {
}

// 咨询师拉取反馈
func (tl *TeacherLogic) GetFeedbackByTeacher(reservationId string, sourceId string,
	userId string, userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	} else if strings.EqualFold(reservationId, sourceId) {
		return nil, errors.New("咨询未被预约，不能反馈")
	}
	teacher, err := models.GetTeacherById(userId)
	if err != nil {
		return nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(utils.GetNow()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	} else if !strings.EqualFold(reservation.TeacherId, teacher.Id.Hex()) {
		return nil, errors.New("只能反馈本人开设的咨询")
	}
	return reservation, nil
}

// 咨询师提交反馈
func (tl *TeacherLogic) SubmitFeedbackByTeacher(reservationId string, sourceId string,
	category string, participants []int, problem string, record string,
	userId string, userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	} else if len(category) == 0 {
		return nil, errors.New("评估分类为空")
	} else if len(participants) != 4 {
		return nil, errors.New("咨询参与者为空")
	} else if len(problem) == 0 {
		return nil, errors.New("问题评估为空")
	} else if len(record) == 0 {
		return nil, errors.New("咨询记录为空")
	} else if strings.EqualFold(reservationId, sourceId) {
		return nil, errors.New("咨询未被预约，不能反馈")
	}
	teacher, err := models.GetTeacherById(userId)
	if err != nil {
		return nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(utils.GetNow()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	} else if !strings.EqualFold(reservation.TeacherId, teacher.Id.Hex()) {
		return nil, errors.New("只能反馈本人开设的咨询")
	}
	if reservation.TeacherFeedback.IsEmpty() && reservation.StudentFeedback.IsEmpty() {
		sms.SendFeedbackSMS(reservation)
	}
	sendFeedbackSMS := reservation.TeacherFeedback.IsEmpty() && reservation.StudentFeedback.IsEmpty()
	reservation.TeacherFeedback = models.TeacherFeedback{
		Category:     category,
		Participants: participants,
		Problem:      problem,
		Record:       record,
	}
	if err = models.UpsertReservation(reservation); err != nil {
		return nil, errors.New("获取数据失败")
	}
	if sendFeedbackSMS && participants[0] > 0 {
		sms.SendFeedbackSMS(reservation)
	}
	return reservation, nil
}

// 咨询师查看学生信息
func (tl *TeacherLogic) GetStudentInfoByTeacher(studentId string,
	userId string, userType models.UserType) (*models.Student, []*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return nil, nil, errors.New("权限不足")
	} else if len(studentId) == 0 {
		return nil, nil, errors.New("咨询未被预约")
	}
	teacher, err := models.GetTeacherById(userId)
	if err != nil {
		return nil, nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != models.TEACHER {
		return nil, nil, errors.New("权限不足")
	}
	student, err := models.GetStudentById(studentId)
	if err != nil {
		return nil, nil, errors.New("学生未注册")
	}
	if !strings.EqualFold(student.BindedTeacherId, teacher.Id.Hex()) {
		return nil, nil, errors.New("只能查看本人绑定的学生")
	}
	reservations, err := models.GetReservationsByStudentId(student.Id.Hex())
	if err != nil {
		return nil, nil, errors.New("获取数据失败")
	}
	return student, reservations, nil
}

// 咨询师查询学生信息
func (tl *TeacherLogic) QueryStudentInfoByTeacher(studentUsername string,
	userId string, userType models.UserType) (*models.Student, []*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return nil, nil, errors.New("权限不足")
	} else if len(studentUsername) == 0 {
		return nil, nil, errors.New("学号为空")
	}
	teacher, err := models.GetTeacherById(userId)
	if err != nil {
		return nil, nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != models.TEACHER {
		return nil, nil, errors.New("权限不足")
	}
	student, err := models.GetStudentByUsername(studentUsername)
	if err != nil {
		return nil, nil, errors.New("学生未注册")
	}
	if !strings.EqualFold(student.BindedTeacherId, teacher.Id.Hex()) {
		return nil, nil, errors.New("只能查看本人绑定的学生")
	}
	reservations, err := models.GetReservationsByStudentId(student.Id.Hex())
	if err != nil {
		return nil, nil, errors.New("获取数据失败")
	}
	return student, reservations, nil
}
