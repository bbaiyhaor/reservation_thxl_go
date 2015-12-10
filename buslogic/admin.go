package buslogic

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
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
	reservation, err := models.AddReservation(start, end, models.ADMIN_ADD, "", teacher.Id)
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
	} else if end.Before(time.Now().In(utils.Location)) {
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
	reservation.TeacherId = teacher.Id
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
				reservation.Status == models.RESERVATED && reservation.StartTime.After(time.Now().In(utils.Location)) {
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
						date := reservation.StartTime.Format(utils.DATE_PATTERN)
						delete(timedReservation, date)
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
	} else if reservation.StartTime.After(time.Now().Local()) {
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
	} else if reservation.StartTime.After(time.Now().Local()) {
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
		utils.SendFeedbackSMS(reservation)
	}
	return reservation, nil
}

// 管理员查看学生信息
func (al *AdminLogic) GetStudentInfoByAdmin(reservationId string, sourceId string,
	userId string, userType models.UserType) (*models.Student, []*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, nil, errors.New("咨询已下架")
	} else if strings.EqualFold(reservationId, sourceId) {
		return nil, nil, errors.New("咨询未被预约，无法查看")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, nil, errors.New("管理员账户出错,请联系技术支持")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, nil, errors.New("咨询已下架")
	} else if reservation.Status == models.AVAILABLE {
		return nil, nil, errors.New("咨询未被预约,无法查看")
	}
	student, err := models.GetStudentById(reservation.StudentId)
	if err != nil {
		return nil, nil, errors.New("咨询已失效")
	}
	reservations, err := models.GetReservationsByStudentId(student.Id)
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
	filename := "student_" + student.Username + "_" + time.Now().In(utils.Location).Format(utils.DATE_PATTERN) + utils.ExcelSuffix
	if err = utils.ExportStudent(student, filename); err != nil {
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
	student.BindedTeacherId = teacher.Id
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
	reservations, err := models.GetReservationsByStudentId(student.Id)
	if err != nil {
		return nil, nil, errors.New("获取数据失败")
	}
	return student, reservations, nil
}

// 管理员导出一周时间表
func (al *AdminLogic) ExportReservationTimetable(fromTime string, userId string, userType models.UserType) (string, error) {
	if len(userId) == 0 {
		return "", errors.New("请先登录")
	} else if userType != models.ADMIN {
		return "", errors.New("权限不足")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return "", errors.New("管理员账户出错,请联系技术支持")
	}
	from, err := time.ParseInLocation(utils.DATE_PATTERN, fromTime, utils.Location)
	if err != nil {
		return "", errors.New("开始时间格式错误")
	}
	to := from.AddDate(0, 0, 7)
	reservations, err := models.GetReservatedReservationsBetweenTime(from, to)
	if err != nil {
		return "", errors.New("获取数据失败")
	}
	filename := "timetable_" + fromTime + utils.ExcelSuffix
	if len(reservations) == 0 {
		return "", nil
	}
	if err = utils.ExportReservationTimetable(reservations, filename); err != nil {
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

// TODO 管理员统计咨询师工作量

// TODO 管理员导出月报

// TODO 管理员指定某次预约的学生
