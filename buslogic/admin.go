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
	admin, err := models.GetTeacherById(userId)
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
		if teacher, err = models.AddFullTeacher(teacherUsername, TeacherDefaultPassword, teacherFullname,
			teacherMobile); err != nil {
			return nil, errors.New("获取数据失败")
		}
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if !strings.EqualFold(teacher.Fullname, teacherFullname) || !strings.EqualFold(teacher.Mobile, teacherMobile) {
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if err = models.UpsertTeacher(teacher); err != nil {
			return nil, errors.New("数据获取失败")
		}
	}
	reservation, err := models.AddReservation(start, end, teacher.Fullname, teacher.Username, teacher.Mobile)
	if err != nil {
		return nil, errors.New("数据获取失败")
	}
	return reservation, nil
}

// 管理员编辑咨询
func (al *AdminLogic) EditReservationByAdmin(reservationId string, startTime string, endTime string,
	teacherUsername string, teacherFullname string, teacherMobile string, userId string,
	userType models.UserType) (*models.Reservation, error) {
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
	}
	admin, err := models.GetTeacherById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
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
		if teacher, err = models.AddFullTeacher(teacherUsername, TeacherDefaultPassword, teacherFullname,
			teacherMobile); err != nil {
			return nil, errors.New("获取数据失败")
		}
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if !strings.EqualFold(teacher.Fullname, teacherFullname) || !strings.EqualFold(teacher.Mobile, teacherMobile) {
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if err = models.UpsertTeacher(teacher); err != nil {
			return nil, errors.New("数据获取失败")
		}
	}
	reservation.StartTime = start
	reservation.EndTime = end
	reservation.TeacherUsername = teacher.Username
	reservation.TeacherFullname = teacher.Fullname
	reservation.TeacherMobile = teacher.Mobile
	if err = models.UpsertReservation(reservation); err != nil {
		return nil, errors.New("数据获取失败")
	}
	return reservation, nil
}

// 管理员删除咨询
func (al *AdminLogic) RemoveReservationsByAdmin(reservationIds []string, userId string, userType models.UserType) (int, error) {
	if len(userId) == 0 {
		return 0, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return 0, errors.New("权限不足")
	}
	admin, err := models.GetTeacherById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return 0, errors.New("管理员账户出错,请联系技术支持")
	}
	removed := 0
	for _, reservationId := range reservationIds {
		if reservation, err := models.GetReservationById(reservationId); err == nil {
			reservation.Status = models.DELETED
			models.UpsertReservation(reservation)
			removed++
		}
	}
	return removed, nil
}

// 管理员取消预约
func (al *AdminLogic) CancelReservationsByAdmin(reservationIds []string, userId string, userType models.UserType) (int, error) {
	if len(userId) == 0 {
		return 0, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return 0, errors.New("权限不足")
	}
	admin, err := models.GetTeacherById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return 0, errors.New("管理员账户出错,请联系技术支持")
	}
	removed := 0
	for _, reservationId := range reservationIds {
		reseravtion, err := models.GetReservationById(reservationId)
		if err != nil || reseravtion.Status == models.DELETED {
			continue
		}
		if reseravtion.Status == models.RESERVATED && reseravtion.StartTime.After(time.Now().Local()) {
			reseravtion.Status = models.AVAILABLE
			reseravtion.StudentUsername = ""
			reseravtion.StudentFullname = ""
			reseravtion.StudentMobile = ""
			reseravtion.StudentFeedback = models.StudentFeedback{}
			reseravtion.TeacherFeedback = models.TeacherFeedback{}
			models.UpsertReservation(reseravtion)
			removed++
		}
	}
	return removed, nil
}

// 管理员拉取反馈
func (al *AdminLogic) GetFeedbackByAdmin(reservationId string, userId string, userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	}
	admin, err := models.GetTeacherById(userId)
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
func (al *AdminLogic) SubmitFeedbackByAdmin(reservationId string, problem string, record string,
	userId string, userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	} else if len(problem) == 0 {
		return nil, errors.New("问题评估为空")
	} else if len(record) == 0 {
		return nil, errors.New("咨询记录为空")
	}
	admin, err := models.GetTeacherById(userId)
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
		Problem: problem,
		Record:  record,
	}
	if err = models.UpsertReservation(reservation); err != nil {
		return nil, errors.New("数据获取失败")
	}
	if sendFeedbackSMS {
		utils.SendFeedbackSMS(reservation)
	}
	return reservation, nil
}

// 管理员查看学生信息
func (al *AdminLogic) GetStudentInfoByAdmin(reservationId string, userId string, userType models.UserType) (*models.Student, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	}
	admin, err := models.GetTeacherById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,无法查看")
	}
	student, err := models.GetStudentByUsername(reservation.StudentUsername)
	if err != nil {
		return nil, errors.New("咨询已失效")
	}
	return student, nil
}

// 管理员导出学生信息
func (al *AdminLogic) ExportStudentByAdmin(studentUsername string, userId string, userType models.UserType) (string, error) {
	if len(userId) == 0 {
		return "", errors.New("请先登录")
	} else if userType != models.ADMIN {
		return "", errors.New("权限不足")
	}
	admin, err := models.GetTeacherById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return "", errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := models.GetStudentByUsername(studentUsername)
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
func (al *AdminLogic) UnbindStudentByAdmin(studentUsername string, userId string, userType models.UserType) (*models.Student, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	}
	admin, err := models.GetTeacherById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := models.GetStudentByUsername(studentUsername)
	if err != nil {
		return nil, errors.New("学生未注册")
	}
	student.BindedTeacher = ""
	if err = models.UpsertStudent(student); err != nil {
		return nil, errors.New("获取数据失败")
	}
	return student, nil
}

// 管理员绑定学生的匹配咨询师
func (al *AdminLogic) BindStudentByAdmin(studentUsername string, teacherUsername string,
	userId string, userType models.UserType) (*models.Student, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	}
	admin, err := models.GetTeacherById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := models.GetStudentByUsername(studentUsername)
	if err != nil {
		return nil, errors.New("学生未注册")
	}
	teacher, err := models.GetTeacherByUsername(teacherUsername)
	if err != nil {
		return nil, errors.New("咨询师未注册")
	}
	student.BindedTeacher = teacher.Username
	if err = models.UpsertStudent(student); err != nil {
		return nil, errors.New("获取数据失败")
	}
	return student, nil
}

// 管理员查询学生信息
func (al *AdminLogic) QueryStudentInfoByAdmin(studentUsername string, userId string, userType models.UserType) (*models.Student, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(studentUsername) == 0 {
		return nil, errors.New("学号为空")
	}
	admin, err := models.GetTeacherById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := models.GetStudentByUsername(studentUsername)
	if err != nil {
		return nil, errors.New("学生未注册")
	}
	return student, nil
}

// 管理员导出一周时间表
func (al *AdminLogic) ExportReservationTimetable(fromTime string, userId string, userType models.UserType) (string, error) {
	if len(userId) == 0 {
		return "", errors.New("请先登录")
	} else if userType != models.ADMIN {
		return "", errors.New("权限不足")
	}
	admin, err := models.GetTeacherById(userId)
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
		return "", errors.New("数据获取失败")
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
func (al *AdminLogic) SearchTeacherByAdmin(teacherFullname string, teacherUsername string, teacherMobile string, userId string, userType models.UserType) (*models.Teacher, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	}
	admin, err := models.GetTeacherById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	if !strings.EqualFold(teacherFullname, "") {
		teacher, err := models.GetTeacherByFullname(teacherFullname)
		if err == nil {
			return teacher, nil
		}
	}
	if !strings.EqualFold(teacherUsername, "") {
		teacher, err := models.GetTeacherByUsername(teacherUsername)
		if err == nil {
			return teacher, nil
		}
	}
	if !strings.EqualFold(teacherMobile, "") {
		teacher, err := models.GetTeacherByMobile(teacherMobile)
		if err == nil {
			return teacher, nil
		}
	}
	return nil, errors.New("用户不存在")
}
