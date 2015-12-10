package buslogic

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"strings"
	"time"
)

// 管理员添加时间表
func (al *AdminLogic) AddTimetableByAdmin(weekday time.Weekday, startTime string, endTime string,
	teacherUsername string, teacherFullname string, teacherMobile string,
	userId string, userType models.UserType) (*models.TimedReservation, error) {
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
		return nil, errors.New("管理员账户出错，请联系技术支持")
	}
	start, err := time.ParseInLocation(utils.CLOCK_PATTERN, startTime, utils.Location)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	end, err := time.ParseInLocation(utils.CLOCK_PATTERN, endTime, utils.Location)
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
		return nil, errors.New("咨询师权限不足")
	} else if !strings.EqualFold(teacher.Fullname, teacherFullname) || !strings.EqualFold(teacher.Mobile, teacherMobile) {
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if err = models.UpsertTeacher(teacherFullname); err != nil {
			return nil, errors.New("获取数据失败")
		}
	}
	timedReservation, err := models.AddTimedReservation(weekday, start, end, teacher.Id)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	return timedReservation, nil
}

// 管理员编辑时间表
func (al *AdminLogic) EditTimetableByAdmin(timedReservationId string, weekday time.Weekday,
	startTime string, endTime string, teacherUsername string, teacherFullname string, teacherMobile string,
	userId string, userType models.UserType) (*models.TimedReservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(timedReservationId) == 0 {
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
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错，请联系技术支持")
	}
	timedReservation, err := models.GetTimedReservationById(timedReservationId)
	if err != nil || timedReservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	}
	start, err := time.ParseInLocation(utils.CLOCK_PATTERN, startTime, utils.Location)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	end, err := time.ParseInLocation(utils.CLOCK_PATTERN, endTime, utils.Location)
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
		return nil, errors.New("咨询师权限不足")
	} else if !strings.EqualFold(teacher.Fullname, teacherFullname) || !strings.EqualFold(teacher.Mobile, teacherMobile) {
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if err = models.UpsertTeacher(teacherFullname); err != nil {
			return nil, errors.New("获取数据失败")
		}
	}
	timedReservation.Weekday = weekday
	timedReservation.StartTime = start
	timedReservation.EndTime = end
	timedReservation.TeacherId = teacher.Id
	if err = models.UpsertTimedReservation(timedReservation); err != nil {
		return nil, errors.New("获取数据失败")
	}
	return timedReservation, nil
}

// 管理员删除时间表
func (al *AdminLogic) RemoveTimetablesByAdmin(timedReservationIds []string, userId string, userType models.UserType) (int, error) {
	if len(userId) == 0 {
		return 0, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return 0, errors.New("权限不足")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错，请联系技术支持")
	}
	removed := 0
	for _, id := range timedReservationIds {
		if timedReservation, err := models.GetTimedReservationById(id); err == nil {
			timedReservation.Status = models.DELETED
			if err = models.UpsertTimedReservation(timedReservation); err == nil {
				removed++
			}
		}
	}
	return removed, nil
}
