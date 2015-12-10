package buslogic

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"strings"
	"time"
)

type StudentLogic struct {
}

// 学生预约咨询
func (sl *StudentLogic) MakeReservationByStudent(reservationId string, sourceId string, startTime string,
	fullname string, gender string, birthday string, school string, grade string, currentAddress string,
	familyAddress string, mobile string, email string, experienceTime string, experienceLocation string,
	experienceTeacher string, fatherAge string, fatherJob string, fatherEdu string, motherAge string, motherJob string,
	motherEdu string, parentMarriage string, siginificant string, problem string,
	userId string, userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.STUDENT {
		return nil, errors.New("请重新登录")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	} else if len(fullname) == 0 {
		return nil, errors.New("姓名为空")
	} else if len(gender) == 0 {
		return nil, errors.New("性别为空")
	} else if len(birthday) == 0 {
		return nil, errors.New("生日为空")
	} else if len(school) == 0 {
		return nil, errors.New("院系为空")
	} else if len(grade) == 0 {
		return nil, errors.New("年纪为空")
	} else if len(mobile) == 0 {
		return nil, errors.New("手机号为空")
	} else if len(email) == 0 {
		return nil, errors.New("邮箱为空")
	} else if len(problem) == 0 {
		return nil, errors.New("问题为空")
	} else if !utils.IsMobile(mobile) {
		return nil, errors.New("手机号格式不正确")
	} else if !utils.IsEmail(email) {
		return nil, errors.New("邮箱格式不正确")
	}
	student, err := models.GetStudentById(userId)
	if err != nil {
		return nil, errors.New("请先登录")
	} else if student.UserType != models.STUDENT {
		return nil, errors.New("请重新登录")
	}
	studentReservations, err := models.GetReservationsByStudentId(student.Id)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	for _, r := range studentReservations {
		if r.Status == models.RESERVATED && r.StartTime.After(time.Now().In(utils.Location)) {
			return nil, errors.New("你好！你已有一个咨询预约，请完成这次咨询后再预约下一次，或致电62782007取消已有预约。")
		}
	}
	reservation := *models.Reservation{}
	if len(sourceId) == 0 {
		// Source为ADD，无SourceId：直接预约
		reservation, err := models.GetReservationById(reservationId)
		if err != nil || reservation.Status == models.DELETED {
			return nil, errors.New("咨询已下架")
		} else if reservation.StartTime.Before(time.Now().In(utils.Location)) {
			return nil, errors.New("咨询已过期")
		} else if reservation.Status != models.AVAILABLE {
			return nil, errors.New("咨询已被预约")
		} else if len(student.BindedTeacherId) != 0 && !strings.EqualFold(student.BindedTeacherId, reservation.TeacherId) {
			return nil, errors.New("只能预约匹配咨询师")
		}
	} else if strings.EqualFold(reservationId, sourceId) {
		// Source为TIMETABLE且未被预约
		timedReservation, err := models.GetTimedReservationById(sourceId)
		if err != nil || timedReservation.Status == models.DELETED {
			return nil, errors.New("咨询已下架")
		}
		start, err := time.ParseInLocation(utils.TIME_PATTERN, startTime, utils.Location)
		if err != nil {
			return nil, errors.New("开始时间格式错误")
		} else if start.Before(time.Now().In(utils.Location)) {
			return nil, errors.New("咨询已过期")
		} else if !strings.EqualFold(start.Format(utils.CLOCK_PATTERN),
			timedReservation.StartTime.In(utils.Location).Format(utils.CLOCK_PATTERN)) {
			return nil, errors.New("开始时间不匹配")
		} else if timedReservation.Timed[start.Format(utils.DATE_PATTERN)] {
			return nil, errors.New("咨询已被预约")
		} else if len(student.BindedTeacherId) != 0 && !strings.EqualFold(student.BindedTeacherId, timedReservation.TeacherId) {
			return nil, errors.New("只能预约匹配咨询师")
		}
		end := utils.ConcatTime(start, timedReservation.EndTime)
		reservation, err = models.AddReservation(start, end, models.TIMETABLE, timedReservation.Id, timedReservation.TeacherId)
		if err != nil {
			return nil, errors.New("获取数据失败")
		}
		timedReservation.Timed[start.Format(utils.DATE_PATTERN)] = true
		if models.UpsertTimedReservation(timedReservation) != nil {
			return nil, errors.New("获取数据失败")
		}
	} else {
		return nil, errors.New("咨询已被预约")
	}
	// 更新学生信息
	student.Fullname = fullname
	student.Gender = gender
	student.Birthday = birthday
	student.School = school
	student.Grade = grade
	student.CurrentAddress = currentAddress
	student.FamilyAddress = familyAddress
	student.Mobile = mobile
	student.Email = email
	student.Experience.Time = experienceTime
	student.Experience.Location = experienceLocation
	student.Experience.Teacher = experienceTeacher
	student.FatherAge = fatherAge
	student.FatherJob = fatherJob
	student.FatherEdu = fatherEdu
	student.MotherAge = motherAge
	student.MotherJob = motherJob
	student.MotherEdu = motherEdu
	student.ParentMarriage = parentMarriage
	student.Significant = siginificant
	student.Problem = problem
	student.BindedTeacherId = reservation.TeacherId
	if models.UpsertStudent(student) != nil {
		return nil, errors.New("获取数据失败")
	}
	// 更新咨询信息
	reservation.StudentId = student.Id
	reservation.Status = models.RESERVATED
	if models.UpsertReservation(reservation) != nil {
		return nil, errors.New("获取数据失败")
	}
	// send success sms
	utils.SendSuccessSMS(reservation)
	return reservation, nil
}

// 学生拉取反馈
func (sl *StudentLogic) GetFeedbackByStudent(reservationId string, sourceId string,
	userId string, userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.STUDENT {
		return nil, errors.New("请重新登录")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	} else if strings.EqualFold(reservationId, sourceId) {
		return nil, errors.New("咨询未被预约，不能反馈")
	}
	student, err := models.GetStudentById(userId)
	if err != nil {
		return nil, errors.New("请先登录")
	} else if student.UserType != models.STUDENT {
		return nil, errors.New("请重新登录")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now().In(utils.Location)) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	} else if !strings.EqualFold(reservation.StudentId, student.Id) {
		return nil, errors.New("只能反馈本人预约的咨询")
	}
	return reservation, nil
}

// 学生反馈
func (sl *StudentLogic) SubmitFeedbackByStudent(reservationId string, sourceId string, scores []string,
	userId string, userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.STUDENT {
		return nil, errors.New("请重新登录")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	} else if len(scores) == 0 {
		return nil, errors.New("反馈为空")
	} else if strings.EqualFold(reservationId, sourceId) {
		return nil, errors.New("咨询未被预约，不能反馈")
	}
	student, err := models.GetStudentById(userId)
	if err != nil {
		return nil, errors.New("请先登录")
	} else if student.UserType != models.STUDENT {
		return nil, errors.New("请重新登录")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now().In(utils.Location)) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	} else if !strings.EqualFold(reservation.StudentId, student.Id) {
		return nil, errors.New("只能反馈本人预约的咨询")
	}
	reservation.StudentFeedback = models.StudentFeedback{
		Scores: scores,
	}
	if models.UpsertReservation(reservation) != nil {
		return nil, errors.New("获取数据失败")
	}
	return reservation, nil
}
