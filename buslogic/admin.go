package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	re "bitbucket.org/shudiwsh2009/reservation_thxl_go/rerror"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/utils"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

// 管理员添加咨询
func (w *Workflow) AddReservationByAdmin(startTime string, endTime string, teacherUsername string,
	teacherFullname string, teacherMobile string, force bool, userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if startTime == "" {
		return nil, re.NewRErrorCodeContext("start time is empty", nil, re.ERROR_MISSING_PARAM, "start_time")
	} else if endTime == "" {
		return nil, re.NewRErrorCodeContext("end time is empty", nil, re.ERROR_MISSING_PARAM, "end_time")
	} else if teacherUsername == "" {
		return nil, re.NewRErrorCodeContext("teacher username is empty", nil, re.ERROR_MISSING_PARAM, "teacher_username")
	} else if teacherFullname == "" {
		return nil, re.NewRErrorCodeContext("teacher fullname is empty", nil, re.ERROR_MISSING_PARAM, "teacher_fullname")
	} else if teacherMobile == "" {
		return nil, re.NewRErrorCodeContext("teacher mobile is empty", nil, re.ERROR_MISSING_PARAM, "teacher_mobile")
	} else if !utils.IsMobile(teacherMobile) {
		return nil, re.NewRErrorCode("teacher mobile format is wrong", nil, re.ERROR_FORMAT_MOBILE)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	start, err := time.ParseInLocation("2006-01-02 15:04", startTime, time.Local)
	if err != nil {
		return nil, re.NewRErrorCodeContext("start time format is wrong", err, re.ERROR_INVALID_PARAM, "start_time")
	}
	end, err := time.ParseInLocation("2006-01-02 15:04", endTime, time.Local)
	if err != nil {
		return nil, re.NewRErrorCodeContext("end time format is wrong", err, re.ERROR_INVALID_PARAM, "end_time")
	}
	if end.Before(start) {
		return nil, re.NewRErrorCode("start time cannot be after end time", nil, re.ERROR_ADMIN_EDIT_RESERVATION_END_TIME_BEFORE_START_TIME)
	}
	teacher, err := w.mongoClient.GetTeacherByUsername(teacherUsername)
	if err != nil {
		teacher := &model.Teacher{
			Username: teacherUsername,
			Password: TEACHER_DEFAULT_PASSWORD,
			Fullname: teacherFullname,
			Mobile:   teacherMobile,
		}
		if err = w.mongoClient.InsertTeacher(teacher); err != nil {
			return nil, re.NewRErrorCode("fail to insert new teacher", err, re.ERROR_DATABASE)
		}
	} else if teacher.UserType != model.USER_TYPE_TEACHER {
		return nil, re.NewRErrorCode("teacher has wrong user type", nil, re.ERROR_DATABASE)
	} else if teacher.Fullname != teacherFullname || teacher.Mobile != teacherMobile {
		if !force {
			return nil, re.NewRErrorCode("teacher info changes without force symbol", nil, re.CHECK)
		}
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if err = w.mongoClient.UpdateTeacher(teacher); err != nil {
			return nil, re.NewRErrorCode("fail to update teacher", err, re.ERROR_DATABASE)
		}
	}
	// 检查时间是否有冲突
	theDay := utils.BeginOfDay(start)
	nextDay := utils.BeginOfTomorrow(start)
	theDayReservations, err := w.mongoClient.GetReservationsBetweenTime(theDay, nextDay)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get the day reservations", err, re.ERROR_DATABASE)
	}
	for _, r := range theDayReservations {
		if r.TeacherId == teacher.Id.Hex() {
			if start.After(r.StartTime) && start.Before(r.EndTime) ||
				(end.After(r.StartTime) && end.Before(r.EndTime)) ||
				(!start.After(r.StartTime) && !end.Before(r.EndTime)) {
				return nil, re.NewRErrorCode("teacher time conflicts", nil, re.ERROR_ADMIN_EDIT_RESERVATION_TEACHER_TIME_CONFLICT)
			}
		}
	}
	theDayTimedReservations, err := w.mongoClient.GetTimedReservationsByWeekday(start.Weekday())
	if err != nil {
		return nil, re.NewRErrorCode("fail to get the day timetables", err, re.ERROR_DATABASE)
	}
	for _, tr := range theDayTimedReservations {
		if !tr.Exceptions[start.Format("2006-01-02")] && tr.TeacherId == teacher.Id.Hex() {
			r := tr.ToReservation(start)
			if start.After(r.StartTime) && start.Before(r.EndTime) ||
				(end.After(r.StartTime) && end.Before(r.EndTime)) ||
				(!start.After(r.StartTime) && !end.Before(r.EndTime)) {
				return nil, re.NewRErrorCode("teacher time conflicts", nil, re.ERROR_ADMIN_EDIT_RESERVATION_TEACHER_TIME_CONFLICT)
			}
		}
	}
	reservation := &model.Reservation{
		StartTime:       start,
		EndTime:         end,
		Status:          model.RESERVATION_STATUS_AVAILABLE,
		Source:          model.RESERVATION_SOURCE_ADMIN_ADD,
		TeacherId:       teacher.Id.Hex(),
		StudentFeedback: model.StudentFeedback{},
		TeacherFeedback: model.TeacherFeedback{},
	}
	if err = w.mongoClient.InsertReservation(reservation); err != nil {
		return nil, re.NewRErrorCode("fail to insert new reservation", err, re.ERROR_DATABASE)
	}
	return reservation, nil
}

// 管理员编辑咨询
func (w *Workflow) EditReservationByAdmin(reservationId string, sourceId string, originalStartTime string,
	startTime string, endTime string, teacherUsername string, teacherFullname string, teacherMobile string,
	force bool, userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation id is empty", nil, re.ERROR_MISSING_PARAM, "reservation_id")
	} else if startTime == "" {
		return nil, re.NewRErrorCodeContext("start time is empty", nil, re.ERROR_MISSING_PARAM, "start_time")
	} else if endTime == "" {
		return nil, re.NewRErrorCodeContext("end time is empty", nil, re.ERROR_MISSING_PARAM, "end_time")
	} else if teacherUsername == "" {
		return nil, re.NewRErrorCodeContext("teacher username is empty", nil, re.ERROR_MISSING_PARAM, "teacher_username")
	} else if teacherFullname == "" {
		return nil, re.NewRErrorCodeContext("teacher fullname is empty", nil, re.ERROR_MISSING_PARAM, "teacher_fullname")
	} else if teacherMobile == "" {
		return nil, re.NewRErrorCodeContext("teacher mobile is empty", nil, re.ERROR_MISSING_PARAM, "teacher_mobile")
	} else if !utils.IsMobile(teacherMobile) {
		return nil, re.NewRErrorCode("teacher mobile format is wrong", nil, re.ERROR_FORMAT_MOBILE)
	} else if sourceId != "" {
		return nil, re.NewRErrorCode("cannot edit timetable here", nil, re.ERROR_ADMIN_EDIT_TIMETABLE_IN_RESERVATION)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	reservation, err := w.mongoClient.GetReservationById(reservationId)
	if err != nil || reservation.Status == model.RESERVATION_STATUS_DELETED {
		return nil, re.NewRErrorCode("cannot get reservation", err, re.ERROR_DATABASE)
	} else if reservation.Status == model.RESERVATION_STATUS_RESERVATED {
		return nil, re.NewRErrorCode("cannot edit reservated reservation", nil, re.ERROR_ADMIN_EDIT_RESERVATED_RESERVATION)
	}
	start, err := time.ParseInLocation("2006-01-02 15:04", startTime, time.Local)
	if err != nil {
		return nil, re.NewRErrorCodeContext("start time format is wrong", err, re.ERROR_INVALID_PARAM, "start_time")
	}
	end, err := time.ParseInLocation("2006-01-02 15:04", endTime, time.Local)
	if err != nil {
		return nil, re.NewRErrorCodeContext("end time format is wrong", err, re.ERROR_INVALID_PARAM, "end_time")
	}
	if end.Before(start) {
		return nil, re.NewRErrorCode("start time cannot be after end time", nil, re.ERROR_ADMIN_EDIT_RESERVATION_END_TIME_BEFORE_START_TIME)
	} else if start.Before(time.Now()) {
		return nil, re.NewRErrorCode("cannot edit outdated reservation", nil, re.ERROR_ADMIN_EDIT_RESERVATION_OUTDATED)
	}
	teacher, err := w.mongoClient.GetTeacherByUsername(teacherUsername)
	if err != nil {
		teacher := &model.Teacher{
			Username: teacherUsername,
			Password: TEACHER_DEFAULT_PASSWORD,
			Fullname: teacherFullname,
			Mobile:   teacherMobile,
		}
		if err = w.mongoClient.InsertTeacher(teacher); err != nil {
			return nil, re.NewRErrorCode("fail to insert new teacher", err, re.ERROR_DATABASE)
		}
	} else if teacher.UserType != model.USER_TYPE_TEACHER {
		return nil, re.NewRErrorCode("teacher has wrong user type", nil, re.ERROR_DATABASE)
	} else if teacher.Fullname != teacherFullname || teacher.Mobile != teacherMobile {
		if !force {
			return nil, re.NewRErrorCode("teacher info changes without force symbol", nil, re.CHECK)
		}
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if err = w.mongoClient.UpdateTeacher(teacher); err != nil {
			return nil, re.NewRErrorCode("fail to update teacher", err, re.ERROR_DATABASE)
		}
	}
	// 检查时间是否有冲突
	theDay := utils.BeginOfDay(start)
	nextDay := utils.BeginOfTomorrow(start)
	theDayReservations, err := w.mongoClient.GetReservationsBetweenTime(theDay, nextDay)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get the day reservations", err, re.ERROR_DATABASE)
	}
	for _, r := range theDayReservations {
		if r.TeacherId == teacher.Id.Hex() {
			if start.After(r.StartTime) && start.Before(r.EndTime) ||
				(end.After(r.StartTime) && end.Before(r.EndTime)) ||
				(!start.After(r.StartTime) && !end.Before(r.EndTime)) {
				return nil, re.NewRErrorCode("teacher time conflicts", nil, re.ERROR_ADMIN_EDIT_RESERVATION_TEACHER_TIME_CONFLICT)
			}
		}
	}
	theDayTimedReservations, err := w.mongoClient.GetTimedReservationsByWeekday(start.Weekday())
	if err != nil {
		return nil, re.NewRErrorCode("fail to get the day timetables", err, re.ERROR_DATABASE)
	}
	for _, tr := range theDayTimedReservations {
		if !tr.Exceptions[start.Format("2006-01-02")] && tr.TeacherId == teacher.Id.Hex() {
			r := tr.ToReservation(start)
			if start.After(r.StartTime) && start.Before(r.EndTime) ||
				(end.After(r.StartTime) && end.Before(r.EndTime)) ||
				(!start.After(r.StartTime) && !end.Before(r.EndTime)) {
				return nil, re.NewRErrorCode("teacher time conflicts", nil, re.ERROR_ADMIN_EDIT_RESERVATION_TEACHER_TIME_CONFLICT)
			}
		}
	}
	reservation.StartTime = start
	reservation.EndTime = end
	reservation.TeacherId = teacher.Id.Hex()
	if err = w.mongoClient.UpdateReservation(reservation); err != nil {
		return nil, re.NewRErrorCode("fail to update reservation", err, re.ERROR_DATABASE)
	}
	return reservation, nil
}

// 管理员删除咨询
func (w *Workflow) RemoveReservationsByAdmin(reservationIds []string, sourceIds []string, startTimes []string,
	userId string, userType int) (int, error) {
	if userId == "" {
		return 0, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return 0, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return 0, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	removed := 0
	for index, reservationId := range reservationIds {
		if sourceIds[index] == "" {
			// Source为ADD，无SourceId：直接置为DELETED（TODO 目前不能删除已预约咨询）
			if reservation, err := w.mongoClient.GetReservationById(reservationId); err == nil &&
				(reservation.Source == model.RESERVATION_SOURCE_ADMIN_ADD || reservation.Source == model.RESERVATION_SOURCE_TEACHER_ADD) &&
				reservation.Status != model.RESERVATION_STATUS_RESERVATED {
				reservation.Status = model.RESERVATION_STATUS_DELETED
				if w.mongoClient.UpdateReservation(reservation) == nil {
					removed++
				}
			}
		} else if reservationId == sourceIds[index] {
			// Source为TIMETABLE且未预约，rId=sourceId：加入exception
			if timedReservation, err := w.mongoClient.GetTimedReservationById(sourceIds[index]); err == nil {
				if time, err := time.ParseInLocation("2006-01-02 15:04", startTimes[index], time.Local); err == nil {
					date := time.Format("2006-01-02")
					timedReservation.Exceptions[date] = true
					if w.mongoClient.UpdateTimedReservation(timedReservation) == nil {
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
	userId string, userType int) (int, error) {
	if userId == "" {
		return 0, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return 0, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return 0, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	canceled := 0
	for index, reservationId := range reservationIds {
		if reservationId != sourceIds[index] {
			// 1、Source为ADD，无SourceId：置为AVAILABLE
			// 2、Source为TIMETABLE且已预约：置为DELETED并去除timed
			if reservation, err := w.mongoClient.GetReservationById(reservationId); err == nil &&
				reservation.Status == model.RESERVATION_STATUS_RESERVATED { // && reservation.StartTime.After(time.Now()) {
				if reservation.Source == model.RESERVATION_SOURCE_ADMIN_ADD || reservation.Source == model.RESERVATION_SOURCE_TEACHER_ADD {
					// 1
					sendSms := reservation.SendSms
					studentId := reservation.StudentId
					reservation.Status = model.RESERVATION_STATUS_AVAILABLE
					reservation.StudentId = ""
					reservation.StudentFeedback = model.StudentFeedback{}
					reservation.TeacherFeedback = model.TeacherFeedback{}
					reservation.IsAdminSet = false
					reservation.SendSms = false
					if w.mongoClient.UpdateReservation(reservation) == nil {
						canceled++
						reservation.StudentId = studentId
						if sendSms {
							w.SendCancelSMS(reservation)
						}
					}
				} else if reservation.Source == model.RESERVATION_SOURCE_TIMETABLE {
					// 2
					reservation.Status = model.RESERVATION_STATUS_DELETED
					if timedReservation, err := w.mongoClient.GetTimedReservationById(sourceIds[index]); err == nil {
						date := reservation.StartTime.Format("2006-01-02")
						delete(timedReservation.Timed, date)
						if w.mongoClient.UpdateReservationAndTimedReservation(reservation, timedReservation) == nil {
							canceled++
							if reservation.SendSms {
								w.SendCancelSMS(reservation)
							}
						}
					}
				}
			}
		}
	}
	return canceled, nil
}

// 管理员拉取反馈
func (w *Workflow) GetFeedbackByAdmin(reservationId string, sourceId string,
	userId string, userType int) (*model.Student, *model.Reservation, error) {
	if userId == "" {
		return nil, nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if reservationId == "" {
		return nil, nil, re.NewRErrorCodeContext("reservation id is empty", nil, re.ERROR_MISSING_PARAM, "reservation_id")
	} else if reservationId == sourceId {
		return nil, nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ERROR_FEEDBACK_AVAILABLE_RESERVATION)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	reservation, err := w.mongoClient.GetReservationById(reservationId)
	if err != nil || reservation.Status == model.RESERVATION_STATUS_DELETED {
		return nil, nil, re.NewRErrorCode("fail to get reservation", err, re.ERROR_DATABASE)
	} else if reservation.StartTime.After(time.Now()) {
		return nil, nil, re.NewRErrorCode("cannot get feedback of future reservation", nil, re.ERROR_FEEDBACK_FUTURE_RESERVATION)
	} else if reservation.Status == model.RESERVATION_STATUS_AVAILABLE {
		return nil, nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ERROR_FEEDBACK_AVAILABLE_RESERVATION)
	}
	student, err := w.mongoClient.GetStudentById(reservation.StudentId)
	if err != nil {
		return nil, nil, re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
	}
	return student, reservation, nil
}

// 管理员提交反馈
func (w *Workflow) SubmitFeedbackByAdmin(reservationId string, sourceId string,
	category string, participants []int, emphasis string, severity []int, medicalDiagnosis []int, crisis []int,
	record string, crisisLevel string, userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation id is empty", nil, re.ERROR_MISSING_PARAM, "reservation_id")
	} else if category == "" {
		return nil, re.NewRErrorCodeContext("category is empty", nil, re.ERROR_MISSING_PARAM, "category")
	} else if len(participants) != len(model.PARTICIPANTS) {
		return nil, re.NewRErrorCodeContext("participants is not valid", nil, re.ERROR_INVALID_PARAM, "participants")
	} else if emphasis == "" {
		return nil, re.NewRErrorCodeContext("emphasis is empty", nil, re.ERROR_MISSING_PARAM, "emphasis")
	} else if len(severity) != len(model.SEVERITY) {
		return nil, re.NewRErrorCodeContext("severity is not valid", nil, re.ERROR_INVALID_PARAM, "severity")
	} else if len(medicalDiagnosis) != len(model.MEDICAL_DIAGNOSIS) {
		return nil, re.NewRErrorCodeContext("medical_diagnosis is not valid", nil, re.ERROR_INVALID_PARAM, "medical_diagnosis")
	} else if len(crisis) != len(model.CRISIS) {
		return nil, re.NewRErrorCodeContext("crisis is not valid", nil, re.ERROR_INVALID_PARAM, "crisis")
	} else if record == "" {
		return nil, re.NewRErrorCodeContext("record is empty", nil, re.ERROR_MISSING_PARAM, "record")
	} else if crisisLevel == "" {
		return nil, re.NewRErrorCodeContext("crisis_level is empty", nil, re.ERROR_MISSING_PARAM, "crisis_level")
	} else if reservationId == sourceId {
		return nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ERROR_FEEDBACK_AVAILABLE_RESERVATION)
	}
	emphasisInt, err := strconv.Atoi(emphasis)
	if err != nil || emphasisInt < 0 {
		return nil, re.NewRErrorCodeContext("emphasis is not valid", err, re.ERROR_INVALID_PARAM, "emphasis")
	}
	crisisLevelInt, err := strconv.Atoi(crisisLevel)
	if err != nil || crisisLevelInt < 0 {
		return nil, re.NewRErrorCodeContext("crisis_level is not valid", err, re.ERROR_INVALID_PARAM, "crisis_level")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	reservation, err := w.mongoClient.GetReservationById(reservationId)
	if err != nil || reservation.Status == model.RESERVATION_STATUS_DELETED {
		return nil, re.NewRErrorCode("fail to get reservation", err, re.ERROR_DATABASE)
	} else if reservation.StartTime.After(time.Now()) {
		return nil, re.NewRErrorCode("cannot get feedback of future reservation", nil, re.ERROR_FEEDBACK_FUTURE_RESERVATION)
	} else if reservation.Status == model.RESERVATION_STATUS_AVAILABLE {
		return nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ERROR_FEEDBACK_AVAILABLE_RESERVATION)
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
	student, err := w.mongoClient.GetStudentById(reservation.StudentId)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
	}
	student.CrisisLevel = crisisLevelInt
	if err = w.mongoClient.UpdateReservationAndStudent(reservation, student); err != nil {
		return nil, re.NewRErrorCode("fail to update reservation and student", err, re.ERROR_DATABASE)
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
	userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation id is empty", nil, re.ERROR_MISSING_PARAM, "reservation_id")
	} else if studentUsername == "" {
		return nil, re.NewRErrorCodeContext("student_username is empty", nil, re.ERROR_MISSING_PARAM, "student_username")
	} else if fullname == "" {
		return nil, re.NewRErrorCodeContext("fullname is empty", nil, re.ERROR_MISSING_PARAM, "fullname")
	} else if gender == "" {
		return nil, re.NewRErrorCodeContext("gender is empty", nil, re.ERROR_MISSING_PARAM, "gender")
	} else if birthday == "" {
		return nil, re.NewRErrorCodeContext("birthday is empty", nil, re.ERROR_MISSING_PARAM, "birthday")
	} else if school == "" {
		return nil, re.NewRErrorCodeContext("school is empty", nil, re.ERROR_MISSING_PARAM, "school")
	} else if grade == "" {
		return nil, re.NewRErrorCodeContext("grade is empty", nil, re.ERROR_MISSING_PARAM, "grade")
	} else if currentAddress == "" {
		return nil, re.NewRErrorCodeContext("current_address is empty", nil, re.ERROR_MISSING_PARAM, "current_address")
	} else if familyAddress == "" {
		return nil, re.NewRErrorCodeContext("family_address is empty", nil, re.ERROR_MISSING_PARAM, "family_address")
	} else if mobile == "" {
		return nil, re.NewRErrorCodeContext("mobile is empty", nil, re.ERROR_MISSING_PARAM, "mobile")
	} else if email == "" {
		return nil, re.NewRErrorCodeContext("email is empty", nil, re.ERROR_MISSING_PARAM, "email")
	} else if problem == "" {
		return nil, re.NewRErrorCodeContext("problem is empty", nil, re.ERROR_MISSING_PARAM, "problem")
	} else if !utils.IsMobile(mobile) {
		return nil, re.NewRErrorCode("mobile format is wrong", nil, re.ERROR_FORMAT_MOBILE)
	} else if !utils.IsEmail(email) {
		return nil, re.NewRErrorCode("email format is wrong", nil, re.ERROR_FORMAT_EMAIL)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	student, err := w.mongoClient.GetStudentByUsername(studentUsername)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get student", err, re.ERROR_NO_STUDENT)
	}
	var reservation *model.Reservation
	if sourceId == "" {
		// Source为ADD，无SourceId：直接指定
		reservation, err = w.mongoClient.GetReservationById(reservationId)
		if err != nil || reservation.Status == model.RESERVATION_STATUS_DELETED {
			return nil, re.NewRErrorCode("fail to get reservation", err, re.ERROR_DATABASE)
			//		} else if reservation.StartTime.Before(time.Now()) {
			//			// 允许指定过期咨询，作为补录（网页正常情况不显示过期咨询，要通过查询咨询的方式来补录）
			//			return nil, errors.New("咨询已过期")
		} else if reservation.Status != model.RESERVATION_STATUS_AVAILABLE {
			return nil, re.NewRErrorCode("cannot set reservated reservation", nil, re.ERROR_ADMIN_SET_RESERVATED_RESERVATION)
		}
	} else if reservationId == sourceId {
		// Source为TIMETABLE且未被预约
		timedReservation, err := w.mongoClient.GetTimedReservationById(sourceId)
		if err != nil || timedReservation.Status == model.RESERVATION_STATUS_DELETED {
			return nil, re.NewRErrorCode("fail to get timetable", err, re.ERROR_DATABASE)
		}
		start, err := time.ParseInLocation("2006-01-02 15:04", startTime, time.Local)
		if err != nil {
			return nil, re.NewRErrorCodeContext("start is not valid", err, re.ERROR_INVALID_PARAM, "start_time")
			//		} else if start.Before(time.Now()) {
			//			// 允许指定过期咨询，作为补录（网页正常情况不显示过期咨询，要通过查询咨询的方式来补录）
			//			return nil, errors.New("咨询已过期")
		} else if start.Format("15:04") != timedReservation.StartTime.Format("15:04") {
			return nil, re.NewRErrorCode("start time mismatch", nil, re.ERROR_START_TIME_MISMATCH)
		} else if timedReservation.Timed[start.Format("2006-01-02")] {
			return nil, re.NewRErrorCode("cannot set reservated reservation", nil, re.ERROR_ADMIN_SET_RESERVATED_RESERVATION)
		}
		end := utils.ConcatTime(start, timedReservation.EndTime)
		reservation = &model.Reservation{
			StartTime:       start,
			EndTime:         end,
			Status:          model.RESERVATION_STATUS_AVAILABLE,
			Source:          model.RESERVATION_SOURCE_TIMETABLE,
			SourceId:        timedReservation.Id.Hex(),
			TeacherId:       timedReservation.TeacherId,
			StudentFeedback: model.StudentFeedback{},
			TeacherFeedback: model.TeacherFeedback{},
		}
		timedReservation.Timed[start.Format("2006-01-02")] = true
		if err = w.mongoClient.InsertReservationAndUpdateTimedReservation(reservation, timedReservation); err != nil {
			return nil, re.NewRErrorCode("fail to insert reservation and update timetable", err, re.ERROR_DATABASE)
		}
	} else {
		return nil, re.NewRErrorCode("cannot set reservated reservation", nil, re.ERROR_ADMIN_SET_RESERVATED_RESERVATION)
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
	student.ParentInfo = model.ParentInfo{
		FatherAge:      fatherAge,
		FatherJob:      fatherJob,
		FatherEdu:      fatherEdu,
		MotherAge:      motherAge,
		MotherJob:      motherJob,
		MotherEdu:      motherEdu,
		ParentMarriage: parentMarriage,
	}
	student.Significant = siginificant
	student.Problem = problem
	student.BindedTeacherId = reservation.TeacherId
	// 更新咨询信息
	reservation.StudentId = student.Id.Hex()
	reservation.IsAdminSet = true
	reservation.SendSms = sendSms
	reservation.Status = model.RESERVATION_STATUS_RESERVATED
	if err = w.mongoClient.UpdateReservationAndStudent(reservation, student); err != nil {
		return nil, re.NewRErrorCode("fail to update reservation and student", err, re.ERROR_DATABASE)
	}
	// send success sms
	if sendSms {
		w.SendSuccessSMS(reservation)
	}
	return reservation, nil
}

// 管理员查看学生信息
func (w *Workflow) GetStudentInfoByAdmin(studentId string,
	userId string, userType int) (*model.Student, []*model.Reservation, error) {
	if userId == "" {
		return nil, nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if studentId == "" {
		return nil, nil, re.NewRErrorCodeContext("student id is empty", nil, re.ERROR_MISSING_PARAM, "student_id")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	student, err := w.mongoClient.GetStudentById(studentId)
	if err != nil || student.UserType != model.USER_TYPE_STUDENT {
		return nil, nil, re.NewRErrorCode("fail to get student", err, re.ERROR_NO_STUDENT)
	}
	reservations, err := w.mongoClient.GetReservationsByStudentId(student.Id.Hex())
	if err != nil {
		return nil, nil, re.NewRErrorCode("fail to get reservations", err, re.ERROR_DATABASE)
	}
	return student, reservations, nil
}

// 管理员更新学生危机等级
func (w *Workflow) UpdateStudentCrisisLevelByAdmin(studentId string, crisisLevel string,
	userId string, userType int) (*model.Student, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if studentId == "" {
		return nil, re.NewRErrorCodeContext("student id is empty", nil, re.ERROR_MISSING_PARAM, "student_id")
	} else if crisisLevel == "" {
		return nil, re.NewRErrorCodeContext("crisis_level is empty", nil, re.ERROR_MISSING_PARAM, "crisis_level")
	}
	crisisLevelInt, err := strconv.Atoi(crisisLevel)
	if err != nil || crisisLevelInt < 0 {
		return nil, re.NewRErrorCodeContext("crisis_level is not valid", err, re.ERROR_INVALID_PARAM, "crisis_level")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	student, err := w.mongoClient.GetStudentById(studentId)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
	}
	student.CrisisLevel = crisisLevelInt
	if err = w.mongoClient.UpdateStudent(student); err != nil {
		return nil, re.NewRErrorCode("fail to update student", err, re.ERROR_DATABASE)
	}
	return student, nil
}

// 管理员更新学生档案编号
func (w *Workflow) UpdateStudentArchiveNumberByAdmin(studentId string, archiveCategory string, archiveNumber string,
	userId string, userType int) (*model.Student, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if studentId == "" {
		return nil, re.NewRErrorCodeContext("student id is empty", nil, re.ERROR_MISSING_PARAM, "student_id")
	} else if archiveCategory == "" {
		return nil, re.NewRErrorCodeContext("archive_category is empty", nil, re.ERROR_MISSING_PARAM, "archive_category")
	} else if archiveNumber == "" {
		return nil, re.NewRErrorCodeContext("archive_number is empty", nil, re.ERROR_MISSING_PARAM, "archive_number")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	student, err := w.mongoClient.GetStudentById(studentId)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
	}
	archiveStudent, err := w.mongoClient.GetStudentByArchiveCategoryAndNumber(archiveCategory, archiveNumber)
	if err == nil && archiveStudent.Id.Valid() && archiveStudent.Id.Hex() != student.Id.Hex() {
		return nil, re.NewRErrorCode("archive number already exists", nil, re.ERROR_ADMIN_ARCHIVE_NUMBER_ALREADY_EXIST)
	}
	archive, err := w.mongoClient.GetArchiveByArchiveCategoryAndNumber(archiveCategory, archiveNumber)
	if err == nil && archive.Id.Valid() && archive.StudentUsername != student.Username {
		return nil, re.NewRErrorCode("archive number already exists", nil, re.ERROR_ADMIN_ARCHIVE_NUMBER_ALREADY_EXIST)
	}
	student.ArchiveCategory = archiveCategory
	student.ArchiveNumber = archiveNumber
	if err = w.mongoClient.UpdateStudent(student); err != nil {
		return nil, re.NewRErrorCode("fail to update student", err, re.ERROR_DATABASE)
	}
	return student, nil
}

// 管理员重置学生密码
func (w *Workflow) ResetStudentPasswordByAdmin(studentId string, password string,
	userId string, userType int) (*model.Student, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if studentId == "" {
		return nil, re.NewRErrorCodeContext("student id is empty", nil, re.ERROR_MISSING_PARAM, "student_id")
	} else if password == "" {
		return nil, re.NewRErrorCodeContext("password is empty", nil, re.ERROR_MISSING_PARAM, "password")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	student, err := w.mongoClient.GetStudentById(studentId)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
	}
	student.Password = password
	student.PreInsert()
	if err = w.mongoClient.UpdateStudent(student); err != nil {
		return nil, re.NewRErrorCode("fail to update student", err, re.ERROR_DATABASE)
	}
	w.ClearUserLoginRedisKey(student.Id.Hex(), student.UserType)
	return student, nil
}

// 管理员删除学生账户
func (w *Workflow) DeleteStudentAccountByAdmin(studentId string, userId string, userType int) error {
	if userId == "" {
		return re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if studentId == "" {
		return re.NewRErrorCodeContext("student id is empty", nil, re.ERROR_MISSING_PARAM, "student_id")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	student, err := w.mongoClient.GetStudentById(studentId)
	if err != nil || student.UserType != model.USER_TYPE_STUDENT {
		return re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
	}
	student.UserType = model.USER_TYPE_UNKNOWN
	if err = w.mongoClient.UpdateStudent(student); err != nil {
		return re.NewRErrorCode("fail to update student", err, re.ERROR_DATABASE)
	}
	return nil
}

// 管理员导出学生信息
func (w *Workflow) ExportStudentByAdmin(studentId string, userId string, userType int) (string, error) {
	if userId == "" {
		return "", re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return "", re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if studentId == "" {
		return "", re.NewRErrorCodeContext("student id is empty", nil, re.ERROR_MISSING_PARAM, "student_id")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return "", re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	student, err := w.mongoClient.GetStudentById(studentId)
	if err != nil {
		return "", re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
	}
	if student.ArchiveNumber == "" {
		return "", re.NewRErrorCode("lack archive number", nil, re.ERROR_ADMIN_EXPORT_STUDENT_NO_ARCHIVE_NUMBER)
	}
	path := filepath.Join(utils.EXPORT_FOLDER, fmt.Sprintf("student_%s_%s_%s%s", student.ArchiveNumber, student.Username, time.Now().Format("2006-01-02"), utils.CSV_FILE_SUFFIX))
	if err = w.ExportStudentInfoToFile(student, path); err != nil {
		return "", err
	}
	return path, nil
}

// 管理员解绑学生的匹配咨询师
func (w *Workflow) UnbindStudentByAdmin(studentId string, userId string, userType int) (*model.Student, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if studentId == "" {
		return nil, re.NewRErrorCodeContext("student id is empty", nil, re.ERROR_MISSING_PARAM, "student_id")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	student, err := w.mongoClient.GetStudentById(studentId)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
	}
	student.BindedTeacherId = ""
	if err = w.mongoClient.UpdateStudent(student); err != nil {
		return nil, re.NewRErrorCode("fail to update student", err, re.ERROR_DATABASE)
	}
	return student, nil
}

// 管理员绑定学生的匹配咨询师
func (w *Workflow) BindStudentByAdmin(studentId string, teacherUsername string,
	userId string, userType int) (*model.Student, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if studentId == "" {
		return nil, re.NewRErrorCodeContext("student id is empty", nil, re.ERROR_MISSING_PARAM, "student_id")
	} else if teacherUsername == "" {
		return nil, re.NewRErrorCodeContext("teacher username is empty", nil, re.ERROR_MISSING_PARAM, "teacher_username")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	student, err := w.mongoClient.GetStudentById(studentId)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
	}
	teacher, err := w.mongoClient.GetTeacherByUsername(teacherUsername)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ERROR_DATABASE)
	}
	student.BindedTeacherId = teacher.Id.Hex()
	if err = w.mongoClient.UpdateStudent(student); err != nil {
		return nil, re.NewRErrorCode("fail to update student", err, re.ERROR_DATABASE)
	}
	return student, nil
}

// 管理员查询学生信息
func (w *Workflow) QueryStudentInfoByAdmin(studentUsername string,
	userId string, userType int) (*model.Student, []*model.Reservation, error) {
	if userId == "" {
		return nil, nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if studentUsername == "" {
		return nil, nil, re.NewRErrorCodeContext("student username is empty", nil, re.ERROR_MISSING_PARAM, "student_username")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	student, err := w.mongoClient.GetStudentByUsername(studentUsername)
	if err != nil || student.UserType != model.USER_TYPE_STUDENT {
		return nil, nil, re.NewRErrorCode("fail to get student", err, re.ERROR_NO_STUDENT)
	}
	reservations, err := w.mongoClient.GetReservationsByStudentId(student.Id.Hex())
	if err != nil {
		return nil, nil, re.NewRErrorCode("fail to get reservations", err, re.ERROR_DATABASE)
	}
	return student, reservations, nil
}

// 管理员重置咨询师密码
func (w *Workflow) ResetTeacherPasswordByAdmin(teacherUsername string, teacherFullname string, teacherMobile string, password string,
	userId string, userType int) (*model.Teacher, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if teacherUsername == "" {
		return nil, re.NewRErrorCodeContext("teacher_username id is empty", nil, re.ERROR_MISSING_PARAM, "teacher_username")
	} else if teacherFullname == "" {
		return nil, re.NewRErrorCodeContext("teacher_fullname id is empty", nil, re.ERROR_MISSING_PARAM, "teacher_fullname")
	} else if teacherMobile == "" {
		return nil, re.NewRErrorCodeContext("teacher_mobile id is empty", nil, re.ERROR_MISSING_PARAM, "teacher_mobile")
	} else if password == "" {
		return nil, re.NewRErrorCodeContext("password is empty", nil, re.ERROR_MISSING_PARAM, "password")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	teacher, err := w.mongoClient.GetTeacherByUsername(teacherUsername)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ERROR_NO_USER)
	} else if teacherFullname != teacher.Fullname || teacherMobile != teacher.Mobile {
		return nil, re.NewRErrorCode("teacher fullname/mobile mismatch", nil, re.ERROR_LOGIN_PWDCHANGE_INFO_MISMATCH)
	}
	teacher.Password = password
	teacher.PreInsert()
	if err = w.mongoClient.UpdateTeacher(teacher); err != nil {
		return nil, re.NewRErrorCode("fail to update teacher", err, re.ERROR_DATABASE)
	}
	w.ClearUserLoginRedisKey(teacher.Id.Hex(), teacher.UserType)
	return teacher, nil
}

// 管理员导出当天时间表
func (w *Workflow) ExportTodayReservationTimetableByAdmin(userId string, userType int) (string, error) {
	if userId == "" {
		return "", re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return "", re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return "", re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	today := utils.BeginOfDay(time.Now())
	tomorrow := today.AddDate(0, 0, 1)
	reservations, err := w.mongoClient.GetReservationsBetweenTime(today, tomorrow)
	if err != nil {
		return "", re.NewRErrorCode("fail to get reservations", err, re.ERROR_DATABASE)
	}
	todayDate := today.Format("2006-01-02")
	if timedReservations, err := w.mongoClient.GetTimedReservationsByWeekday(today.Weekday()); err == nil {
		for _, tr := range timedReservations {
			if !tr.Exceptions[todayDate] && !tr.Timed[todayDate] {
				reservations = append(reservations, tr.ToReservation(today))
			}
		}
	}
	sort.Sort(ByStartTimeOfReservation(reservations))
	if len(reservations) == 0 {
		return "", re.NewRErrorCode("no reservations today", nil, re.ERROR_ADMIN_NO_RESERVATIONS_TODAY)
	}
	path := filepath.Join(utils.EXPORT_FOLDER, fmt.Sprintf("timetable_%s%s", todayDate, utils.CSV_FILE_SUFFIX))
	if err = w.ExportTodayReservationTimetableToFile(reservations, path); err != nil {
		return "", err
	}
	return path, nil
}

// 查找咨询师
// 查找顺序:全名 > 工号 > 手机号
func (w *Workflow) SearchTeacherByAdmin(teacherFullname string, teacherUsername string, teacherMobile string,
	userId string, userType int) (*model.Teacher, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	if teacherFullname != "" {
		teacher, err := w.mongoClient.GetTeacherByFullname(teacherFullname)
		if err == nil {
			return teacher, nil
		}
	}
	if teacherUsername != "" {
		teacher, err := w.mongoClient.GetTeacherByUsername(teacherUsername)
		if err == nil {
			return teacher, nil
		}
	}
	if teacherMobile != "" {
		teacher, err := w.mongoClient.GetTeacherByMobile(teacherMobile)
		if err == nil {
			return teacher, nil
		}
	}
	return nil, re.NewRErrorCode("fail to get teacher", nil, re.ERROR_NO_USER)
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
	userId string, userType int) (map[string]WorkLoad, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if fromDate == "" {
		return nil, re.NewRErrorCodeContext("from_date is empty", nil, re.ERROR_MISSING_PARAM, "from_date")
	} else if toDate == "" {
		return nil, re.NewRErrorCodeContext("to_date is empty", nil, re.ERROR_MISSING_PARAM, "to_date")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	from, err := time.ParseInLocation("2006-01-02", fromDate, time.Local)
	if err != nil {
		return nil, re.NewRErrorCodeContext("from_date is not valid", err, re.ERROR_INVALID_PARAM, "from_date")
	}
	to, err := time.ParseInLocation("2006-01-02", toDate, time.Local)
	if err != nil {
		return nil, re.NewRErrorCodeContext("to_date is not valid", err, re.ERROR_INVALID_PARAM, "to_date")
	}
	to = to.AddDate(0, 0, 1)
	reservations, err := w.mongoClient.GetReservatedReservationsBetweenTime(from, to)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get reservations", err, re.ERROR_DATABASE)
	}
	workload := make(map[string]WorkLoad)
	for _, r := range reservations {
		if _, exist := workload[r.TeacherId]; !exist {
			teacher, err := w.mongoClient.GetTeacherById(r.TeacherId)
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
func (w *Workflow) ExportReportFormByAdmin(fromDate string, toDate string, userId string, userType int) (string, error) {
	if userId == "" {
		return "", re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return "", re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if fromDate == "" {
		return "", re.NewRErrorCodeContext("from_date is empty", nil, re.ERROR_MISSING_PARAM, "from_date")
	} else if toDate == "" {
		return "", re.NewRErrorCodeContext("to_date is empty", nil, re.ERROR_MISSING_PARAM, "to_date")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return "", re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	from, err := time.ParseInLocation("2006-01-02", fromDate, time.Local)
	if err != nil {
		return "", re.NewRErrorCodeContext("from_date is not valid", err, re.ERROR_INVALID_PARAM, "from_date")
	}
	to, err := time.ParseInLocation("2006-01-02", toDate, time.Local)
	if err != nil {
		return "", re.NewRErrorCodeContext("to_date is not valid", err, re.ERROR_INVALID_PARAM, "to_date")
	}
	to = to.AddDate(0, 0, 1)
	reservations, err := w.mongoClient.GetReservatedReservationsBetweenTime(from, to)
	if err != nil {
		return "", re.NewRErrorCode("fail to get reservations", err, re.ERROR_DATABASE)
	}
	if len(reservations) == 0 {
		return "", nil
	}
	path := filepath.Join(utils.EXPORT_FOLDER, fmt.Sprintf("monthly_report_%s_%s%s", fromDate, toDate, utils.CSV_FILE_SUFFIX))
	if err = w.ExportReportFormToFile(reservations, path); err != nil {
		return "", err
	}
	return path, nil
}

// 管理员导出报表
func (w *Workflow) ExportReportMonthlyByAdmin(monthlyDate string, userId string, userType int) (string, string, error) {
	if userId == "" {
		return "", "", re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return "", "", re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if monthlyDate == "" {
		return "", "", re.NewRErrorCodeContext("monthly_date is empty", nil, re.ERROR_MISSING_PARAM, "monthly_date")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return "", "", re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	date, err := time.ParseInLocation("2006-01-02", monthlyDate, time.Local)
	if err != nil {
		return "", "", re.NewRErrorCodeContext("date is not valid", err, re.ERROR_INVALID_PARAM, "date")
	}
	from := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
	to := from.AddDate(0, 1, 0)
	reservations, err := w.mongoClient.GetReservatedReservationsBetweenTime(from, to)
	if err != nil {
		return "", "", re.NewRErrorCode("fail to get reservations", err, re.ERROR_DATABASE)
	}
	if len(reservations) == 0 {
		return "", "", nil
	}
	reportPath := filepath.Join(utils.EXPORT_FOLDER, fmt.Sprintf("monthly_report_%d_%d%s", date.Year(), date.Month(), utils.CSV_FILE_SUFFIX))
	keyCasePath := filepath.Join(utils.EXPORT_FOLDER, fmt.Sprintf("monthly_key_case_%d_%d%s", date.Year(), date.Month(), utils.CSV_FILE_SUFFIX))
	if err = w.ExportReportFormToFile(reservations, reportPath); err != nil {
		return "", "", err
	}
	//if err = workflow.ExportKeyCaseReport(reservations, keyCasePath); err != nil {
	//	return "", "", err
	//}
	return reportPath, keyCasePath, nil
}

func (w *Workflow) WrapAdmin(admin *model.Admin) map[string]interface{} {
	var result = make(map[string]interface{})
	if admin == nil {
		return result
	}
	result["id"] = admin.Id.Hex()
	result["username"] = admin.Username
	result["user_type"] = admin.UserType
	result["fullname"] = admin.Fullname
	result["mobile"] = admin.Mobile
	return result
}
