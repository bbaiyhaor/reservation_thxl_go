package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/util"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 管理员添加咨询
func (w *Workflow) AddReservationByAdmin(startTime string, endTime string, teacherUsername string,
	teacherFullname string, teacherMobile string, force bool, userId string, userType model.UserType) (*model.Reservation, error) {
	if userId == "" {
		return nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, errors.New("权限不足")
	} else if startTime == "" {
		return nil, errors.New("开始时间为空")
	} else if endTime == "" {
		return nil, errors.New("结束时间为空")
	} else if teacherUsername == "" {
		return nil, errors.New("咨询师工号为空")
	} else if teacherFullname == "" {
		return nil, errors.New("咨询师姓名为空")
	} else if teacherMobile == "" {
		return nil, errors.New("咨询师手机号为空")
	} else if !util.IsMobile(teacherMobile) {
		return nil, errors.New("咨询师手机号格式不正确")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	start, err := time.ParseInLocation("2006-01-02 15:04", startTime, time.Local)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	end, err := time.ParseInLocation("2006-01-02 15:04", endTime, time.Local)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}
	if start.After(end) {
		return nil, errors.New("开始时间不能晚于结束时间")
	}
	teacher, err := w.model.GetTeacherByUsername(teacherUsername)
	if err != nil {
		if teacher, err = w.model.AddTeacher(teacherUsername, TeacherDefaultPassword, teacherFullname, teacherMobile); err != nil {
			return nil, errors.New("获取数据失败")
		}
	} else if teacher.UserType != model.TEACHER {
		return nil, errors.New("权限不足")
	} else if !strings.EqualFold(teacher.Fullname, teacherFullname) || !strings.EqualFold(teacher.Mobile, teacherMobile) {
		if !force {
			return nil, errors.New(model.CHECK_MESSAGE)
		}
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if w.model.UpsertTeacher(teacher) != nil {
			return nil, errors.New("获取数据失败")
		}
	}
	reservation, err := w.model.AddReservation(start, end, model.ADMIN_ADD, "", teacher.Id.Hex())
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	return reservation, nil
}

// 管理员编辑咨询
func (w *Workflow) EditReservationByAdmin(reservationId string, sourceId string, originalStartTime string,
	startTime string, endTime string, teacherUsername string, teacherFullname string, teacherMobile string,
	force bool, userId string, userType model.UserType) (*model.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
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
	} else if !util.IsMobile(teacherMobile) {
		return nil, errors.New("咨询师手机号格式不正确")
	} else if len(sourceId) != 0 {
		return nil, errors.New("请在安排表中编辑预设咨询")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	reservation, err := w.model.GetReservationById(reservationId)
	if err != nil || reservation.Status == model.DELETED {
		return nil, errors.New("请在安排表中编辑预设咨询")
	} else if reservation.Status == model.RESERVATED {
		return nil, errors.New("不能编辑已被预约的咨询")
	}
	start, err := time.ParseInLocation("2006-01-02 15:04", startTime, time.Local)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	end, err := time.ParseInLocation("2006-01-02 15:04", endTime, time.Local)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}
	if start.After(end) {
		return nil, errors.New("开始时间不能晚于结束时间")
	} else if start.Before(time.Now()) {
		return nil, errors.New("不能编辑已过期咨询")
	}
	teacher, err := w.model.GetTeacherByUsername(teacherUsername)
	if err != nil {
		if teacher, err = w.model.AddTeacher(teacherUsername, TeacherDefaultPassword, teacherFullname, teacherMobile); err != nil {
			return nil, errors.New("获取数据失败")
		}
	} else if teacher.UserType != model.TEACHER {
		return nil, errors.New("权限不足")
	} else if !strings.EqualFold(teacher.Fullname, teacherFullname) || !strings.EqualFold(teacher.Mobile, teacherMobile) {
		if !force {
			return nil, errors.New(model.CHECK_MESSAGE)
		}
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if w.model.UpsertTeacher(teacher) != nil {
			return nil, errors.New("获取数据失败")
		}
	}
	reservation.StartTime = start
	reservation.EndTime = end
	reservation.TeacherId = teacher.Id.Hex()
	if err = w.model.UpsertReservation(reservation); err != nil {
		return nil, errors.New("获取数据失败")
	}
	return reservation, nil
}

// 管理员删除咨询
func (w *Workflow) RemoveReservationsByAdmin(reservationIds []string, sourceIds []string, startTimes []string,
	userId string, userType model.UserType) (int, error) {
	if len(userId) == 0 {
		return 0, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return 0, errors.New("权限不足")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return 0, errors.New("管理员账户出错,请联系技术支持")
	}
	removed := 0
	for index, reservationId := range reservationIds {
		if len(sourceIds[index]) == 0 {
			// Source为ADD，无SourceId：直接置为DELETED（TODO 目前不能删除已预约咨询）
			if reservation, err := w.model.GetReservationById(reservationId); err == nil && reservation.Status != model.RESERVATED {
				reservation.Status = model.DELETED
				if w.model.UpsertReservation(reservation) == nil {
					removed++
				}
			}
		} else if strings.EqualFold(reservationId, sourceIds[index]) {
			// Source为TIMETABLE且未预约，rId=sourceId：加入exception
			if timedReservation, err := w.model.GetTimedReservationById(sourceIds[index]); err == nil {
				if time, err := time.ParseInLocation("2006-01-02 15:04", startTimes[index], time.Local); err == nil {
					date := time.Format("2006-01-02")
					timedReservation.Exceptions[date] = true
					if w.model.UpsertTimedReservation(timedReservation) == nil {
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
func (w *Workflow) CancelReservationsByAdmin(reservationIds []string, sourceIds []string,
	userId string, userType model.UserType) (int, error) {
	if userId == "" {
		return 0, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return 0, errors.New("权限不足")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return 0, errors.New("管理员账户出错,请联系技术支持")
	}
	removed := 0
	for index, reservationId := range reservationIds {
		if !strings.EqualFold(reservationId, sourceIds[index]) {
			// 1、Source为ADD，无SourceId：置为AVAILABLE
			// 2、Source为TIMETABLE且已预约：置为DELETED并去除timed
			if reservation, err := w.model.GetReservationById(reservationId); err == nil &&
				reservation.Status == model.RESERVATED { // && reservation.StartTime.After(time.Now()) {
				if reservation.Source != model.TIMETABLE {
					// 1
					sendSms := reservation.SendSms
					reservation.Status = model.AVAILABLE
					studentId := reservation.StudentId
					reservation.StudentId = ""
					reservation.StudentFeedback = model.StudentFeedback{}
					reservation.TeacherFeedback = model.TeacherFeedback{}
					reservation.IsAdminSet = false
					reservation.SendSms = false
					if w.model.UpsertReservation(reservation) == nil {
						removed++
						reservation.StudentId = studentId
						if sendSms {
							w.SendCancelSMS(reservation)
						}
					}
				} else {
					// 2
					reservation.Status = model.DELETED
					if timedReservation, err := w.model.GetTimedReservationById(sourceIds[index]); err == nil {
						date := reservation.StartTime.Format("2006-01-02")
						delete(timedReservation.Timed, date)
						if w.model.UpsertReservation(reservation) == nil && w.model.UpsertTimedReservation(timedReservation) == nil {
							removed++
							if reservation.SendSms {
								w.SendCancelSMS(reservation)
							}
						}
					}
				}
			}
		}
	}
	return removed, nil
}

// 管理员拉取反馈
func (w *Workflow) GetFeedbackByAdmin(reservationId string, sourceId string,
	userId string, userType model.UserType) (*model.Student, *model.Reservation, error) {
	if len(userId) == 0 {
		return nil, nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, nil, errors.New("咨询已下架")
	} else if strings.EqualFold(reservationId, sourceId) {
		return nil, nil, errors.New("咨询未被预约，不能反馈")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, nil, errors.New("管理员账户出错,请联系技术支持")
	}
	reservation, err := w.model.GetReservationById(reservationId)
	if err != nil || reservation.Status == model.DELETED {
		return nil, nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now()) {
		return nil, nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == model.AVAILABLE {
		return nil, nil, errors.New("咨询未被预约,不能反馈")
	}
	student, err := w.model.GetStudentById(reservation.StudentId)
	if err != nil {
		return nil, nil, errors.New("获取数据失败")
	}
	return student, reservation, nil
}

// 管理员提交反馈
func (w *Workflow) SubmitFeedbackByAdmin(reservationId string, sourceId string,
	category string, participants []int, emphasis string, severity []int, medicalDiagnosis []int, crisis []int,
	record string, crisisLevel string, userId string, userType model.UserType) (*model.Reservation, error) {
	if userId == "" {
		return nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
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
	} else if len(record) == 0 {
		return nil, errors.New("咨询记录为空")
	} else if crisisLevel == "" {
		return nil, errors.New("危机等级为空")
	} else if strings.EqualFold(reservationId, sourceId) {
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
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	reservation, err := w.model.GetReservationById(reservationId)
	if err != nil || reservation.Status == model.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == model.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
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

// 管理员指定某次预约的学生
func (w *Workflow) SetStudentByAdmin(reservationId string, sourceId string, startTime string, studentUsername string,
	fullname string, gender string, birthday string, school string, grade string, currentAddress string,
	familyAddress string, mobile string, email string, experienceTime string, experienceLocation string,
	experienceTeacher string, fatherAge string, fatherJob string, fatherEdu string, motherAge string, motherJob string,
	motherEdu string, parentMarriage string, siginificant string, problem string, sendSms bool,
	userId string, userType model.UserType) (*model.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	} else if len(studentUsername) == 0 {
		return nil, errors.New("学生学号为空")
	} else if len(fullname) == 0 {
		return nil, errors.New("姓名为空")
	} else if len(gender) == 0 {
		return nil, errors.New("性别为空")
	} else if len(birthday) == 0 {
		return nil, errors.New("出生日期为空")
	} else if len(school) == 0 {
		return nil, errors.New("院系为空")
	} else if len(grade) == 0 {
		return nil, errors.New("年纪为空")
	} else if len(currentAddress) == 0 {
		return nil, errors.New("现在住址为空")
	} else if len(familyAddress) == 0 {
		return nil, errors.New("家庭住址为空")
	} else if len(mobile) == 0 {
		return nil, errors.New("手机号为空")
	} else if len(email) == 0 {
		return nil, errors.New("邮箱为空")
	} else if len(problem) == 0 {
		return nil, errors.New("问题为空")
	} else if !util.IsMobile(mobile) {
		return nil, errors.New("手机号格式不正确")
	} else if !util.IsEmail(email) {
		return nil, errors.New("邮箱格式不正确")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := w.model.GetStudentByUsername(studentUsername)
	if err != nil {
		return nil, errors.New("学生未注册")
	}
	var reservation *model.Reservation
	if sourceId == "" {
		// Source为ADD，无SourceId：直接指定
		reservation, err = w.model.GetReservationById(reservationId)
		if err != nil || reservation.Status == model.DELETED {
			return nil, errors.New("咨询已下架")
			//		} else if reservation.StartTime.Before(time.Now()) {
			//			// 允许指定过期咨询，作为补录（网页正常情况不显示过期咨询，要通过查询咨询的方式来补录）
			//			return nil, errors.New("咨询已过期")
		} else if reservation.Status != model.AVAILABLE {
			return nil, errors.New("咨询已被预约")
		}
	} else if reservationId == sourceId {
		// Source为TIMETABLE且未被预约
		timedReservation, err := w.model.GetTimedReservationById(sourceId)
		if err != nil || timedReservation.Status == model.DELETED {
			return nil, errors.New("咨询已下架")
		}
		start, err := time.ParseInLocation("2006-01-02 15:04", startTime, time.Local)
		if err != nil {
			return nil, errors.New("开始时间格式错误")
			//		} else if start.Before(time.Now()) {
			//			// 允许指定过期咨询，作为补录（网页正常情况不显示过期咨询，要通过查询咨询的方式来补录）
			//			return nil, errors.New("咨询已过期")
		} else if start.Format("15:04") != timedReservation.StartTime.Format("15:04") {
			return nil, errors.New("开始时间不匹配")
		} else if timedReservation.Timed[start.Format("2006-01-02")] {
			return nil, errors.New("咨询已被预约")
		}
		end := util.ConcatTime(start, timedReservation.EndTime)
		reservation, err = w.model.AddReservation(start, end, model.TIMETABLE, timedReservation.Id.Hex(),
			timedReservation.TeacherId)
		if err != nil {
			return nil, errors.New("获取数据失败")
		}
		timedReservation.Timed[start.Format("2006-01-02")] = true
		if w.model.UpsertTimedReservation(timedReservation) != nil {
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
	if w.model.UpsertStudent(student) != nil {
		return nil, errors.New("获取数据失败")
	}
	// 更新咨询信息
	reservation.StudentId = student.Id.Hex()
	reservation.IsAdminSet = true
	reservation.SendSms = sendSms
	reservation.Status = model.RESERVATED
	if err = w.model.UpsertReservation(reservation); err != nil {
		return nil, errors.New("获取数据失败")
	}
	// send success sms
	if sendSms {
		w.SendSuccessSMS(reservation)
	}
	return reservation, nil
}

// 管理员查看学生信息
func (w *Workflow) GetStudentInfoByAdmin(studentId string,
	userId string, userType model.UserType) (*model.Student, []*model.Reservation, error) {
	if len(userId) == 0 {
		return nil, nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, nil, errors.New("权限不足")
	} else if len(studentId) == 0 {
		return nil, nil, errors.New("咨询未被预约，不能查看")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, nil, errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := w.model.GetStudentById(studentId)
	if err != nil || student.UserType != model.STUDENT {
		return nil, nil, errors.New("学生未注册")
	}
	reservations, err := w.model.GetReservationsByStudentId(student.Id.Hex())
	if err != nil {
		return nil, nil, errors.New("获取数据失败")
	}
	return student, reservations, nil
}

// 管理员更新学生危机等级
func (w *Workflow) UpdateStudentCrisisLevelByAdmin(studentId string, crisisLevel string,
	userId string, userType model.UserType) (*model.Student, error) {
	if userId == "" {
		return nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, errors.New("权限不足")
	} else if studentId == "" {
		return nil, errors.New("学生未注册")
	} else if crisisLevel == "" {
		return nil, errors.New("危机等级为空")
	}
	crisisLevelInt, err := strconv.Atoi(crisisLevel)
	if err != nil || crisisLevelInt < 0 {
		return nil, errors.New("危机等级错误")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := w.model.GetStudentById(studentId)
	if err != nil {
		return nil, errors.New("学生未注册")
	}
	student.CrisisLevel = crisisLevelInt
	if err := w.model.UpsertStudent(student); err != nil {
		return nil, errors.New("获取数据失败")
	}
	return student, nil
}

// 管理员更新学生档案编号
func (w *Workflow) UpdateStudentArchiveNumberByAdmin(studentId string, archiveCategory string, archiveNumber string,
	userId string, userType model.UserType) (*model.Student, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(studentId) == 0 {
		return nil, errors.New("学生未注册")
	} else if len(archiveCategory) == 0 {
		return nil, errors.New("档案分类为空")
	} else if len(archiveNumber) == 0 {
		return nil, errors.New("档案编号为空")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := w.model.GetStudentById(studentId)
	if err != nil {
		return nil, errors.New("学生未注册")
	}
	archiveStudent, err := w.model.GetStudentByArchiveNumber(archiveNumber)
	if err == nil && archiveStudent.Id.Valid() && !strings.EqualFold(archiveStudent.Id.Hex(), student.Id.Hex()) {
		return nil, errors.New("档案号已存在，请重新分配")
	}
	archive, err := w.model.GetArchiveByArchiveNumber(archiveNumber)
	if err == nil && archive.Id.Valid() && !strings.EqualFold(archive.StudentUsername, student.Username) {
		return nil, errors.New("档案号已存在，请重新分配")
	}
	student.ArchiveCategory = archiveCategory
	student.ArchiveNumber = archiveNumber
	if err := w.model.UpsertStudent(student); err != nil {
		return nil, errors.New("获取数据失败")
	}
	return student, nil
}

// 管理员重置学生密码
func (w *Workflow) ResetStudentPasswordByAdmin(studentId string, password string,
	userId string, userType model.UserType) (*model.Student, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(studentId) == 0 {
		return nil, errors.New("学生未注册")
	} else if len(password) == 0 {
		return nil, errors.New("密码为空")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := w.model.GetStudentById(studentId)
	if err != nil {
		return nil, errors.New("学生未注册")
	}
	//student.Password = password
	encryptedPassword, err := util.EncryptPassword(password)
	if err != nil {
		return nil, errors.New("加密出错，请联系技术支持")
	}
	student.EncryptedPassword = encryptedPassword
	if err := w.model.UpsertStudent(student); err != nil {
		return nil, errors.New("获取数据失败")
	}
	return student, nil
}

// 管理员删除学生账户
func (w *Workflow) DeleteStudentAccountByAdmin(studentId string, userId string, userType model.UserType) error {
	if len(userId) == 0 {
		return errors.New("请先登录")
	} else if userType != model.ADMIN {
		return errors.New("权限不足")
	} else if len(studentId) == 0 {
		return errors.New("学生未注册")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := w.model.GetStudentById(studentId)
	if err != nil || student.UserType != model.STUDENT {
		return errors.New("学生未注册")
	}
	student.UserType = model.UNKNOWN
	if err := w.model.UpsertStudent(student); err != nil {
		return errors.New("获取数据失败")
	}
	return nil
}

// 管理员导出学生信息
func (w *Workflow) ExportStudentByAdmin(studentId string, userId string, userType model.UserType) (string, error) {
	if len(userId) == 0 {
		return "", errors.New("请先登录")
	} else if userType != model.ADMIN {
		return "", errors.New("权限不足")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return "", errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := w.model.GetStudentById(studentId)
	if err != nil {
		return "", errors.New("学生未注册")
	}
	if len(student.ArchiveNumber) == 0 {
		return "", errors.New("请先分配档案号")
	}
	filename := "student_" + student.ArchiveNumber + "_" + student.Username + "_" +
		time.Now().Format("2006-01-02") + util.CsvSuffix
	if err = w.ExportStudentInfoToFile(student, filename); err != nil {
		return "", err
	}
	return "/" + util.ExportFolder + filename, nil
}

// 管理员解绑学生的匹配咨询师
func (w *Workflow) UnbindStudentByAdmin(studentId string, userId string, userType model.UserType) (*model.Student, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, errors.New("权限不足")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := w.model.GetStudentById(studentId)
	if err != nil {
		return nil, errors.New("学生未注册")
	}
	student.BindedTeacherId = ""
	if err = w.model.UpsertStudent(student); err != nil {
		return nil, errors.New("获取数据失败")
	}
	return student, nil
}

// 管理员绑定学生的匹配咨询师
func (w *Workflow) BindStudentByAdmin(studentId string, teacherUsername string,
	userId string, userType model.UserType) (*model.Student, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(studentId) == 0 {
		return nil, errors.New("请先登录")
	} else if len(teacherUsername) == 0 {
		return nil, errors.New("咨询师工号为空")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := w.model.GetStudentById(studentId)
	if err != nil {
		return nil, errors.New("学生未注册")
	}
	teacher, err := w.model.GetTeacherByUsername(teacherUsername)
	if err != nil {
		return nil, errors.New("咨询师未注册")
	}
	student.BindedTeacherId = teacher.Id.Hex()
	if err = w.model.UpsertStudent(student); err != nil {
		return nil, errors.New("获取数据失败")
	}
	return student, nil
}

// 管理员查询学生信息
func (w *Workflow) QueryStudentInfoByAdmin(studentUsername string,
	userId string, userType model.UserType) (*model.Student, []*model.Reservation, error) {
	if len(userId) == 0 {
		return nil, nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, nil, errors.New("权限不足")
	} else if len(studentUsername) == 0 {
		return nil, nil, errors.New("学号为空")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, nil, errors.New("管理员账户出错,请联系技术支持")
	}
	student, err := w.model.GetStudentByUsername(studentUsername)
	if err != nil || student.UserType != model.STUDENT {
		return nil, nil, errors.New("学生未注册")
	}
	reservations, err := w.model.GetReservationsByStudentId(student.Id.Hex())
	if err != nil {
		return nil, nil, errors.New("获取数据失败")
	}
	return student, reservations, nil
}

// 管理员导出当天时间表
func (w *Workflow) ExportTodayReservationTimetableByAdmin(userId string, userType model.UserType) (string, error) {
	if len(userId) == 0 {
		return "", errors.New("请先登录")
	} else if userType != model.ADMIN {
		return "", errors.New("权限不足")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return "", errors.New("管理员账户出错,请联系技术支持")
	}
	today := util.BeginOfDay(time.Now())
	tomorrow := today.AddDate(0, 0, 1)
	reservations, err := w.model.GetReservationsBetweenTime(today, tomorrow)
	if err != nil {
		return "", errors.New("获取数据失败")
	}
	todayDate := today.Format("2006-01-02")
	if timedReservations, err := w.model.GetTimedReservationsByWeekday(today.Weekday()); err == nil {
		for _, tr := range timedReservations {
			if !tr.Exceptions[todayDate] && !tr.Timed[todayDate] {
				reservations = append(reservations, tr.ToReservation(today))
			}
		}
	}
	sort.Sort(model.ReservationSlice(reservations))
	filename := "timetable_" + todayDate + util.CsvSuffix
	if len(reservations) == 0 {
		return "", errors.New("今日无咨询")
	}
	if err = w.ExportTodayReservationTimetableToFile(reservations, filename); err != nil {
		return "", err
	}
	return "/" + util.ExportFolder + filename, nil
}

// 查找咨询师
// 查找顺序:全名 > 工号 > 手机号
func (w *Workflow) SearchTeacherByAdmin(teacherFullname string, teacherUsername string, teacherMobile string,
	userId string, userType model.UserType) (*model.Teacher, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, errors.New("权限不足")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	if len(teacherFullname) != 0 {
		teacher, err := w.model.GetTeacherByFullname(teacherFullname)
		if err == nil {
			return teacher, nil
		}
	}
	if len(teacherUsername) != 0 {
		teacher, err := w.model.GetTeacherByUsername(teacherUsername)
		if err == nil {
			return teacher, nil
		}
	}
	if len(teacherMobile) != 0 {
		teacher, err := w.model.GetTeacherByMobile(teacherMobile)
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
func (w *Workflow) GetTeacherWorkloadByAdmin(fromDate string, toDate string,
	userId string, userType model.UserType) (map[string]WorkLoad, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(fromDate) == 0 {
		return nil, errors.New("开始日期为空")
	} else if len(toDate) == 0 {
		return nil, errors.New("结束日期为空")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	from, err := time.ParseInLocation("2006-01-02", fromDate, time.Local)
	if err != nil {
		return nil, errors.New("开始日期格式错误")
	}
	to, err := time.ParseInLocation("2006-01-02", toDate, time.Local)
	if err != nil {
		return nil, errors.New("结束日期格式错误")
	}
	to = to.AddDate(0, 0, 1)
	reservations, err := w.model.GetReservatedReservationsBetweenTime(from, to)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	workload := make(map[string]WorkLoad)
	for _, r := range reservations {
		if _, exist := workload[r.TeacherId]; !exist {
			teacher, err := w.model.GetTeacherById(r.TeacherId)
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

// 管理员导出报表
func (w *Workflow) ExportReportFormByAdmin(fromDate string, toDate string, userId string, userType model.UserType) (string, error) {
	if len(userId) == 0 {
		return "", errors.New("请先登录")
	} else if userType != model.ADMIN {
		return "", errors.New("权限不足")
	} else if len(fromDate) == 0 {
		return "", errors.New("开始日期为空")
	} else if len(toDate) == 0 {
		return "", errors.New("结束日期为空")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return "", errors.New("管理员账户出错,请联系技术支持")
	}
	from, err := time.ParseInLocation("2006-01-02", fromDate, time.Local)
	if err != nil {
		return "", errors.New("开始日期格式错误")
	}
	to, err := time.ParseInLocation("2006-01-02", toDate, time.Local)
	if err != nil {
		return "", errors.New("结束日期格式错误")
	}
	to = to.AddDate(0, 0, 1)
	reservations, err := w.model.GetReservatedReservationsBetweenTime(from, to)
	if err != nil {
		return "", errors.New("获取数据失败")
	}
	filename := fmt.Sprintf("monthly_report_%s_%s%s", fromDate, toDate, util.CsvSuffix)
	if len(reservations) == 0 {
		return "", nil
	}
	if err = w.ExportReportFormToFile(reservations, filename); err != nil {
		return "", err
	}
	return "/" + util.ExportFolder + filename, nil
}

// 管理员导出报表
func (w *Workflow) ExportReportMonthlyByAdmin(monthlyDate string, userId string, userType model.UserType) (string, string, error) {
	if len(userId) == 0 {
		return "", "", errors.New("请先登录")
	} else if userType != model.ADMIN {
		return "", "", errors.New("权限不足")
	} else if len(monthlyDate) == 0 {
		return "", "", errors.New("开始日期为空")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return "", "", errors.New("管理员账户出错,请联系技术支持")
	}
	date, err := time.ParseInLocation("2006-01-02", monthlyDate, time.Local)
	if err != nil {
		return "", "", errors.New("开始日期格式错误")
	}
	from := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
	to := from.AddDate(0, 1, 0)
	reservations, err := w.model.GetReservatedReservationsBetweenTime(from, to)
	if err != nil {
		return "", "", errors.New("获取数据失败")
	}
	reportFilename := fmt.Sprintf("monthly_report_%d_%d%s", date.Year(), date.Month(), util.CsvSuffix)
	keyCaseFilename := fmt.Sprintf("monthly_key_case_%d_%d%s", date.Year(), date.Month(), util.CsvSuffix)
	if len(reservations) == 0 {
		return "", "", nil
	}
	if err = w.ExportReportFormToFile(reservations, reportFilename); err != nil {
		return "", "", err
	}
	//if err = workflow.ExportKeyCaseReport(reservations, keyCaseFilename); err != nil {
	//	return "", "", err
	//}
	return "/" + util.ExportFolder + reportFilename, "/" + util.ExportFolder + keyCaseFilename, nil
}
