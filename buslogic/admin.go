package buslogic

import (
	"errors"
	"fmt"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"github.com/shudiwsh2009/reservation_thxl_go/workflow"
	"sort"
	"strings"
	"time"
)

type AdminLogic struct {
}

// 管理员添加咨询
func (al *AdminLogic) AddReservationByAdmin(startTime string, endTime string, teacherUsername string,
	teacherFullname string, teacherMobile string, userId string, userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(startTime) == 0 {
		return nil, errors.New("开始时间为空")
	} else if len(endTime) == 0 {
		return nil, errors.New("结束时间为空")
	} else if len(teacherUsername) == 0 {
		return nil, errors.New("咨询师工号为空")
	} else if len(teacherFullname) == 0 {
		return nil, errors.New("咨询师姓名为空")
	} else if len(teacherMobile) == 0 {
		return nil, errors.New("咨询师手机号为空")
	} else if !utils.IsMobile(teacherMobile) {
		return nil, errors.New("咨询师手机号格式不正确")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	start, err := time.ParseInLocation(utils.TIME_PATTERN, startTime, utils.Location)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	end, err := time.ParseInLocation(utils.TIME_PATTERN, endTime, utils.Location)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}
	if start.After(end) {
		return nil, errors.New("开始时间不能晚于结束时间")
	}
	teacher, err := models.GetTeacherByUsername(teacherUsername)
	if err != nil {
		if teacher, err = models.AddTeacher(teacherUsername, TeacherDefaultPassword, teacherFullname, teacherMobile); err != nil {
			return nil, errors.New("获取数据失败")
		}
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if !strings.EqualFold(teacher.Fullname, teacherFullname) || !strings.EqualFold(teacher.Mobile, teacherMobile) {
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if models.UpsertTeacher(teacher) != nil {
			return nil, errors.New("获取数据失败")
		}
	}
	reservation, err := models.AddReservation(start, end, models.ADMIN_ADD, "", teacher.Id.Hex())
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	return reservation, nil
}

// 管理员编辑咨询
func (al *AdminLogic) EditReservationByAdmin(reservationId string, sourceId string, originalStartTime string,
	startTime string, endTime string, teacherUsername string, teacherFullname string, teacherMobile string,
	userId string, userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	} else if len(startTime) == 0 {
		return nil, errors.New("开始时间为空")
	} else if len(endTime) == 0 {
		return nil, errors.New("结束时间为空")
	} else if len(teacherUsername) == 0 {
		return nil, errors.New("咨询师工号为空")
	} else if len(teacherFullname) == 0 {
		return nil, errors.New("咨询师姓名为空")
	} else if len(teacherMobile) == 0 {
		return nil, errors.New("咨询师手机号为空")
	} else if !utils.IsMobile(teacherMobile) {
		return nil, errors.New("咨询师手机号格式不正确")
	} else if len(sourceId) != 0 {
		return nil, errors.New("请在安排表中编辑预设咨询")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("请在安排表中编辑预设咨询")
	} else if reservation.Status == models.RESERVATED {
		return nil, errors.New("不能编辑已被预约的咨询")
	}
	start, err := time.ParseInLocation(utils.TIME_PATTERN, startTime, utils.Location)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	end, err := time.ParseInLocation(utils.TIME_PATTERN, endTime, utils.Location)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}
	if start.After(end) {
		return nil, errors.New("开始时间不能晚于结束时间")
	} else if start.Before(utils.GetNow()) {
		return nil, errors.New("不能编辑已过期咨询")
	}
	teacher, err := models.GetTeacherByUsername(teacherUsername)
	if err != nil {
		if teacher, err = models.AddTeacher(teacherUsername, TeacherDefaultPassword, teacherFullname, teacherMobile); err != nil {
			return nil, errors.New("获取数据失败")
		}
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if !strings.EqualFold(teacher.Fullname, teacherFullname) || !strings.EqualFold(teacher.Mobile, teacherMobile) {
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if models.UpsertTeacher(teacher) != nil {
			return nil, errors.New("获取数据失败")
		}
	}
	reservation.StartTime = start
	reservation.EndTime = end
	reservation.TeacherId = teacher.Id.Hex()
	if err = models.UpsertReservation(reservation); err != nil {
		return nil, errors.New("获取数据失败")
	}
	return reservation, nil
}

// 管理员删除咨询
func (al *AdminLogic) RemoveReservationsByAdmin(reservationIds []string, sourceIds []string, startTimes []string,
	userId string, userType models.UserType) (int, error) {
	if len(userId) == 0 {
		return 0, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return 0, errors.New("权限不足")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return 0, errors.New("管理员账户出错,请联系技术支持")
	}
	removed := 0
	for index, reservationId := range reservationIds {
		if len(sourceIds[index]) == 0 {
			// Source为ADD，无SourceId：直接置为DELETED（TODO 目前不能删除已预约咨询）
			if reservation, err := models.GetReservationById(reservationId); err == nil && reservation.Status != models.RESERVATED {
				reservation.Status = models.DELETED
				if models.UpsertReservation(reservation) == nil {
					removed++
				}
			}
		} else if strings.EqualFold(reservationId, sourceIds[index]) {
			// Source为TIMETABLE且未预约，rId=sourceId：加入exception
			if timedReservation, err := models.GetTimedReservationById(sourceIds[index]); err == nil {
				if time, err := time.ParseInLocation(utils.TIME_PATTERN, startTimes[index], utils.Location); err == nil {
					date := time.Format(utils.DATE_PATTERN)
					timedReservation.Exceptions[date] = true
					if models.UpsertTimedReservation(timedReservation) == nil {
						removed++
					}
				}
			}
		} else {
			// Source为TIMETABLE且已预约，rId!=sourceId
			// TODO 目前不能删除已预约咨询
		}
	}
	return removed, nil
}

// 管理员取消预约
func (al *AdminLogic) CancelReservationsByAdmin(reservationIds []string, sourceIds []string,
	userId string, userType models.UserType) (int, error) {
	if len(userId) == 0 {
		return 0, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return 0, errors.New("权限不足")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return 0, errors.New("管理员账户出错,请联系技术支持")
	}
	removed := 0
	for index, reservationId := range reservationIds {
		if !strings.EqualFold(reservationId, sourceIds[index]) {
			// 1、Source为ADD，无SourceId：置为AVAILABLE
			// 2、Source为TIMETABLE且已预约：置为DELETED并去除timed
			if reservation, err := models.GetReservationById(reservationId); err == nil &&
				reservation.Status == models.RESERVATED && reservation.StartTime.After(utils.GetNow()) {
				if reservation.Source != models.TIMETABLE {
					// 1
					reservation.Status = models.AVAILABLE
					reservation.StudentId = ""
					reservation.StudentFeedback = models.StudentFeedback{}
					reservation.TeacherFeedback = models.TeacherFeedback{}
					if models.UpsertReservation(reservation) == nil {
						removed++
					}
				} else {
					// 2
					reservation.Status = models.DELETED
					if timedReservation, err := models.GetTimedReservationById(sourceIds[index]); err == nil {
						date := reservation.StartTime.In(utils.Location).Format(utils.DATE_PATTERN)
						delete(timedReservation.Timed, date)
						if models.UpsertReservation(reservation) == nil && models.UpsertTimedReservation(timedReservation) == nil {
							removed++
						}
					}
				}
			}
		}
	}
	return removed, nil
}

// 管理员拉取反馈
func (al *AdminLogic) GetFeedbackByAdmin(reservationId string, sourceId string,
	userId string, userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	} else if strings.EqualFold(reservationId, sourceId) {
		return nil, errors.New("咨询未被预约，不能反馈")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(utils.GetNow()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	}
	return reservation, nil
}

// 管理员提交反馈
func (al *AdminLogic) SubmitFeedbackByAdmin(reservationId string, sourceId string,
	category string, participants []int, problem string, record string,
	userId string, userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
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
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(utils.GetNow()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
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
		workflow.SendFeedbackSMS(reservation)
	}
	return reservation, nil
}

// 管理员指定某次预约的学生
func (al *AdminLogic) SetStudentByAdmin(reservationId string, sourceId string, startTime string, studentUsername string,
	userId string, userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	} else if len(studentUsername) == 0 {
		return nil, errors.New("学生学号为空")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := models.GetStudentByUsername(studentUsername)
	if err != nil {
		return nil, errors.New("学生未注册")
	}
	reservation := &models.Reservation{}
	if len(sourceId) == 0 {
		// Source为ADD，无SourceId：直接指定
		reservation, err = models.GetReservationById(reservationId)
		if err != nil || reservation.Status == models.DELETED {
			return nil, errors.New("咨询已下架")
		} else if reservation.StartTime.Before(utils.GetNow()) {
			return nil, errors.New("咨询已过期")
		} else if reservation.Status != models.AVAILABLE {
			return nil, errors.New("咨询已被预约")
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
		} else if start.Before(utils.GetNow()) {
			return nil, errors.New("咨询已过期")
		} else if !strings.EqualFold(start.Format(utils.CLOCK_PATTERN),
			timedReservation.StartTime.In(utils.Location).Format(utils.CLOCK_PATTERN)) {
			return nil, errors.New("开始时间不匹配")
		} else if timedReservation.Timed[start.Format(utils.DATE_PATTERN)] {
			return nil, errors.New("咨询已被预约")
		}
		end := utils.ConcatTime(start, timedReservation.EndTime)
		reservation, err = models.AddReservation(start, end, models.TIMETABLE, timedReservation.Id.Hex(),
			timedReservation.TeacherId)
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
	student.BindedTeacherId = reservation.TeacherId
	if err = models.UpsertStudent(student); err != nil {
		return nil, errors.New("获取数据失败")
	}
	// 更新咨询信息
	reservation.StudentId = student.Id.Hex()
	reservation.Status = models.RESERVATED
	if err = models.UpsertReservation(reservation); err != nil {
		return nil, errors.New("获取数据失败")
	}
	return reservation, nil
}

// 管理员查看学生信息
func (al *AdminLogic) GetStudentInfoByAdmin(studentId string,
	userId string, userType models.UserType) (*models.Student, []*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, nil, errors.New("权限不足")
	} else if len(studentId) == 0 {
		return nil, nil, errors.New("咨询未被预约，不能查看")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, nil, errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := models.GetStudentById(studentId)
	if err != nil {
		return nil, nil, errors.New("学生未注册")
	}
	reservations, err := models.GetReservationsByStudentId(student.Id.Hex())
	if err != nil {
		return nil, nil, errors.New("获取数据失败")
	}
	return student, reservations, nil
}

// 管理员导出学生信息
func (al *AdminLogic) ExportStudentByAdmin(studentId string, userId string, userType models.UserType) (string, error) {
	if len(userId) == 0 {
		return "", errors.New("请先登录")
	} else if userType != models.ADMIN {
		return "", errors.New("权限不足")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return "", errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := models.GetStudentById(studentId)
	if err != nil {
		return "", errors.New("学生未注册")
	}
	filename := "student_" + student.Username + "_" + utils.GetNow().Format(utils.DATE_PATTERN) + utils.CsvSuffix
	if err = workflow.ExportStudentInfo(student, filename); err != nil {
		return "", err
	}
	return "/" + utils.ExportFolder + filename, nil
}

// 管理员解绑学生的匹配咨询师
func (al *AdminLogic) UnbindStudentByAdmin(studentId string, userId string, userType models.UserType) (*models.Student, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := models.GetStudentById(studentId)
	if err != nil {
		return nil, errors.New("学生未注册")
	}
	student.BindedTeacherId = ""
	if err = models.UpsertStudent(student); err != nil {
		return nil, errors.New("获取数据失败")
	}
	return student, nil
}

// 管理员绑定学生的匹配咨询师
func (al *AdminLogic) BindStudentByAdmin(studentId string, teacherUsername string,
	userId string, userType models.UserType) (*models.Student, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(studentId) == 0 {
		return nil, errors.New("请先登录")
	} else if len(teacherUsername) == 0 {
		return nil, errors.New("咨询师工号为空")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := models.GetStudentById(studentId)
	if err != nil {
		return nil, errors.New("学生未注册")
	}
	teacher, err := models.GetTeacherByUsername(teacherUsername)
	if err != nil {
		return nil, errors.New("咨询师未注册")
	}
	student.BindedTeacherId = teacher.Id.Hex()
	if err = models.UpsertStudent(student); err != nil {
		return nil, errors.New("获取数据失败")
	}
	return student, nil
}

// 管理员查询学生信息
func (al *AdminLogic) QueryStudentInfoByAdmin(studentUsername string,
	userId string, userType models.UserType) (*models.Student, []*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, nil, errors.New("权限不足")
	} else if len(studentUsername) == 0 {
		return nil, nil, errors.New("学号为空")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, nil, errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := models.GetStudentByUsername(studentUsername)
	if err != nil {
		return nil, nil, errors.New("学生未注册")
	}
	reservations, err := models.GetReservationsByStudentId(student.Id.Hex())
	if err != nil {
		return nil, nil, errors.New("获取数据失败")
	}
	return student, reservations, nil
}

// 管理员导出当天时间表
func (al *AdminLogic) ExportTodayReservationTimetableByAdmin(userId string, userType models.UserType) (string, error) {
	if len(userId) == 0 {
		return "", errors.New("请先登录")
	} else if userType != models.ADMIN {
		return "", errors.New("权限不足")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return "", errors.New("管理员账户出错,请联系技术支持")
	}
	today := utils.GetToday()
	tomorrow := today.AddDate(0, 0, 1)
	reservations, err := models.GetReservationsBetweenTime(today, tomorrow)
	if err != nil {
		return "", errors.New("获取数据失败")
	}
	todayDate := today.Format(utils.DATE_PATTERN)
	if timedReservations, err := models.GetTimedReservationsByWeekday(today.Weekday()); err == nil {
		for _, tr := range timedReservations {
			if !tr.Exceptions[todayDate] && !tr.Timed[todayDate] {
				reservations = append(reservations, tr.ToReservation(today))
			}
		}
	}
	sort.Sort(models.ReservationSlice(reservations))
	filename := "timetable_" + todayDate + utils.CsvSuffix
	if len(reservations) == 0 {
		return "", errors.New("今日无咨询")
	}
	if err = workflow.ExportTodayReservationTimetable(reservations, filename); err != nil {
		return "", err
	}
	return "/" + utils.ExportFolder + filename, nil
}

// 查找咨询师
// 查找顺序:全名 > 工号 > 手机号
func (al *AdminLogic) SearchTeacherByAdmin(teacherFullname string, teacherUsername string, teacherMobile string,
	userId string, userType models.UserType) (*models.Teacher, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	if len(teacherFullname) != 0 {
		teacher, err := models.GetTeacherByFullname(teacherFullname)
		if err == nil {
			return teacher, nil
		}
	}
	if len(teacherUsername) != 0 {
		teacher, err := models.GetTeacherByUsername(teacherUsername)
		if err == nil {
			return teacher, nil
		}
	}
	if len(teacherMobile) != 0 {
		teacher, err := models.GetTeacherByMobile(teacherMobile)
		if err == nil {
			return teacher, nil
		}
	}
	return nil, errors.New("用户不存在")
}

type WorkLoad struct {
	TeacherId       string          `json:"teacher_id"`
	TeacherUsername string          `json:"teacher_username"`
	TeacherFullname string          `json:"teacher_fullname"`
	TeacherMobile   string          `json:"teacher_mobile"`
	Students        map[string]bool `json:"students"`
	Reservations    map[string]bool `json:"reservations"`
}

// 管理员统计咨询师工作量
func (al *AdminLogic) GetTeacherWorkloadByAdmin(fromDate string, toDate string,
	userId string, userType models.UserType) (map[string]WorkLoad, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(fromDate) == 0 {
		return nil, errors.New("开始日期为空")
	} else if len(toDate) == 0 {
		return nil, errors.New("结束日期为空")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	from, err := time.ParseInLocation(utils.DATE_PATTERN, fromDate, utils.Location)
	if err != nil {
		return nil, errors.New("开始日期格式错误")
	}
	to, err := time.ParseInLocation(utils.DATE_PATTERN, toDate, utils.Location)
	if err != nil {
		return nil, errors.New("结束日期格式错误")
	}
	to = to.AddDate(0, 0, 1)
	reservations, err := models.GetReservatedReservationsBetweenTime(from, to)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	workload := make(map[string]WorkLoad)
	for _, r := range reservations {
		if _, exist := workload[r.TeacherId]; !exist {
			teacher, err := models.GetTeacherById(r.TeacherId)
			if err != nil {
				continue
			}
			workload[r.TeacherId] = WorkLoad{
				TeacherId:       teacher.Id.Hex(),
				TeacherUsername: teacher.Username,
				TeacherFullname: teacher.Fullname,
				TeacherMobile:   teacher.Mobile,
				Students:        make(map[string]bool),
				Reservations:    make(map[string]bool),
			}
		}
		workload[r.TeacherId].Students[r.StudentId] = true
		workload[r.TeacherId].Reservations[r.Id.Hex()] = true
	}
	return workload, nil
}

// 管理员导出月报
func (al *AdminLogic) ExportMonthlyReportByAdmin(monthlyDate string, userId string, userType models.UserType) (string, error) {
	if len(userId) == 0 {
		return "", errors.New("请先登录")
	} else if userType != models.ADMIN {
		return "", errors.New("权限不足")
	} else if len(monthlyDate) == 0 {
		return "", errors.New("日期为空")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return "", errors.New("管理员账户出错,请联系技术支持")
	}
	date, err := time.ParseInLocation(utils.DATE_PATTERN, monthlyDate, utils.Location)
	if err != nil {
		return "", errors.New("日期格式错误")
	}
	from := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, utils.Location)
	to := from.AddDate(0, 1, 0)
	reservations, err := models.GetReservatedReservationsBetweenTime(from, to)
	if err != nil {
		return "", errors.New("获取数据失败")
	}
	filename := fmt.Sprintf("monthly_report_%d_%d%s", date.Year(), date.Month(), utils.CsvSuffix)
	if len(reservations) == 0 {
		return "", nil
	}
	if err = workflow.ExportMonthlyReport(reservations, filename); err != nil {
		return "", err
	}
	return "/" + utils.ExportFolder + filename, nil
}
