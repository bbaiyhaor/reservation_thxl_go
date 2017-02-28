package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/config"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	re "bitbucket.org/shudiwsh2009/reservation_thxl_go/rerror"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/utils"
	"fmt"
	"github.com/mijia/sweb/log"
	"github.com/scorredoira/email"
	"net/mail"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// 管理员查看时间表
func (w *Workflow) ViewTimetableByAdmin(userId string, userType int) (map[time.Weekday][]*model.TimedReservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	timedReservations := make(map[time.Weekday][]*model.TimedReservation)
	for i := time.Sunday; i <= time.Saturday; i++ {
		if trs, err := w.mongoClient.GetTimedReservationsByWeekday(i); err == nil {
			sort.Sort(ByWeekdayOfTimedReservation(trs))
			timedReservations[i] = trs
		}
	}
	return timedReservations, nil
}

// 管理员添加时间表
func (w *Workflow) AddTimetableByAdmin(weekday string, startClock string, endClock string,
	teacherUsername string, teacherFullname string, teacherMobile string, force bool,
	userId string, userType int) (*model.TimedReservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if startClock == "" {
		return nil, re.NewRErrorCodeContext("start clock is empty", nil, re.ERROR_MISSING_PARAM, "start_clock")
	} else if endClock == "" {
		return nil, re.NewRErrorCodeContext("end clock is empty", nil, re.ERROR_MISSING_PARAM, "end_clock")
	} else if teacherUsername == "" {
		return nil, re.NewRErrorCodeContext("teacher username is empty", nil, re.ERROR_MISSING_PARAM, "teacher_username")
	} else if teacherFullname == "" {
		return nil, re.NewRErrorCodeContext("teacher fullname is empty", nil, re.ERROR_MISSING_PARAM, "teacher_fullname")
	} else if teacherMobile == "" {
		return nil, re.NewRErrorCodeContext("teacher mobile is empty", nil, re.ERROR_MISSING_PARAM, "teacher_mobile")
	} else if !utils.IsMobile(teacherMobile) {
		return nil, re.NewRErrorCode("mobile format is wrong", nil, re.ERROR_FORMAT_MOBILE)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	week, err := utils.StringToWeekday(weekday)
	if err != nil {
		return nil, re.NewRErrorCode("weekday format is wrong", nil, re.ERROR_FORMAT_WEEKDAY)
	}
	start, err := time.ParseInLocation("2006-01-02 15:04", "2006-01-02 "+startClock, time.Local)
	if err != nil {
		return nil, re.NewRErrorCodeContext("start time format is wrong", err, re.ERROR_INVALID_PARAM, "start_time")
	}
	end, err := time.ParseInLocation("2006-01-02 15:04", "2006-01-02 "+endClock, time.Local)
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
	timedReservation := &model.TimedReservation{
		Weekday:    week,
		StartTime:  start,
		EndTime:    end,
		Status:     model.RESERVATION_STATUS_CLOSED,
		TeacherId:  teacher.Id.Hex(),
		Exceptions: make(map[string]bool),
		Timed:      make(map[string]bool),
	}
	if err = w.mongoClient.InsertTimedReservation(timedReservation); err != nil {
		return nil, re.NewRErrorCode("fail to insert new timetable", err, re.ERROR_DATABASE)
	}
	return timedReservation, nil
}

// 管理员编辑时间表
func (w *Workflow) EditTimetableByAdmin(timedReservationId string, weekday string,
	startClock string, endClock string, teacherUsername string, teacherFullname string, teacherMobile string,
	force bool, userId string, userType int) (*model.TimedReservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if timedReservationId == "" {
		return nil, re.NewRErrorCodeContext("timed reservation id is empty", nil, re.ERROR_MISSING_PARAM, "timed_reservation_id")
	} else if startClock == "" {
		return nil, re.NewRErrorCodeContext("start clock is empty", nil, re.ERROR_MISSING_PARAM, "start_clock")
	} else if endClock == "" {
		return nil, re.NewRErrorCodeContext("end clock is empty", nil, re.ERROR_MISSING_PARAM, "end_clock")
	} else if teacherUsername == "" {
		return nil, re.NewRErrorCodeContext("teacher username is empty", nil, re.ERROR_MISSING_PARAM, "teacher_username")
	} else if teacherFullname == "" {
		return nil, re.NewRErrorCodeContext("teacher fullname is empty", nil, re.ERROR_MISSING_PARAM, "teacher_fullname")
	} else if teacherMobile == "" {
		return nil, re.NewRErrorCodeContext("teacher mobile is empty", nil, re.ERROR_MISSING_PARAM, "teacher_mobile")
	} else if !utils.IsMobile(teacherMobile) {
		return nil, re.NewRErrorCode("mobile format is wrong", nil, re.ERROR_FORMAT_MOBILE)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	timedReservation, err := w.mongoClient.GetTimedReservationById(timedReservationId)
	if err != nil || timedReservation.Status == model.RESERVATION_STATUS_DELETED {
		return nil, re.NewRErrorCode("fail to get timetable", err, re.ERROR_DATABASE)
	}
	week, err := utils.StringToWeekday(weekday)
	if err != nil {
		return nil, re.NewRErrorCode("weekday format is wrong", nil, re.ERROR_FORMAT_WEEKDAY)
	}
	start, err := time.ParseInLocation("2006-01-02 15:04", "2006-01-02 "+startClock, time.Local)
	if err != nil {
		return nil, re.NewRErrorCodeContext("start time format is wrong", err, re.ERROR_INVALID_PARAM, "start_time")
	}
	end, err := time.ParseInLocation("2006-01-02 15:04", "2006-01-02 "+endClock, time.Local)
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
	timedReservation.Weekday = week
	timedReservation.StartTime = start
	timedReservation.EndTime = end
	timedReservation.Status = model.RESERVATION_STATUS_CLOSED
	timedReservation.TeacherId = teacher.Id.Hex()
	if err = w.mongoClient.UpdateTimedReservation(timedReservation); err != nil {
		return nil, re.NewRErrorCode("fail to insert new timetable", err, re.ERROR_DATABASE)
	}
	return timedReservation, nil
}

// 管理员删除时间表
func (w *Workflow) RemoveTimetablesByAdmin(timedReservationIds []string, userId string, userType int) (int, error) {
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
	for _, id := range timedReservationIds {
		if timedReservation, err := w.mongoClient.GetTimedReservationById(id); err == nil {
			timedReservation.Status = model.RESERVATION_STATUS_DELETED
			if err = w.mongoClient.UpdateTimedReservation(timedReservation); err == nil {
				removed++
			}
		}
	}
	return removed, nil
}

// 管理员开启时间表
func (w *Workflow) OpenTimetablesByAdmin(timedReservationIds []string, userId string, userType int) (int, error) {
	if userId == "" {
		return 0, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return 0, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return 0, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	opened := 0
	for _, id := range timedReservationIds {
		if timedReservation, err := w.mongoClient.GetTimedReservationById(id); err == nil {
			if timedReservation.Status == model.RESERVATION_STATUS_CLOSED {
				timedReservation.Status = model.RESERVATION_STATUS_AVAILABLE
				if err = w.mongoClient.UpdateTimedReservation(timedReservation); err == nil {
					opened++
				}
			}
		}
	}
	return opened, nil
}

// 管理员关闭时间表
func (w *Workflow) CloseTimetablesByAdmin(timedReservationIds []string, userId string, userType int) (int, error) {
	if userId == "" {
		return 0, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return 0, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return 0, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	closed := 0
	for _, id := range timedReservationIds {
		if timedReservation, err := w.mongoClient.GetTimedReservationById(id); err == nil {
			if timedReservation.Status == model.RESERVATION_STATUS_AVAILABLE {
				timedReservation.Status = model.RESERVATION_STATUS_CLOSED
				if err = w.mongoClient.UpdateTimedReservation(timedReservation); err == nil {
					closed++
				}
			}
		}
	}
	return closed, nil
}

func (w *Workflow) ExportTodayReservationTimetableToFile(reservations []*model.Reservation, path string) error {
	data := make([][]string, 0)
	today := utils.BeginOfDay(time.Now())
	data = append(data, []string{today.Format("2006-01-02")})
	data = append(data, []string{"时间", "咨询师", "学生姓名", "学生学号", "联系方式"})
	for _, r := range reservations {
		teacher, err := w.mongoClient.GetTeacherById(r.TeacherId)
		if err != nil {
			continue
		}
		if student, err := w.mongoClient.GetStudentById(r.StudentId); err == nil {
			data = append(data, []string{r.StartTime.Format("15:04") + " - " + r.EndTime.Format("15:04"),
				teacher.Fullname, student.Fullname, student.Username, student.Mobile})
		} else {
			data = append(data, []string{r.StartTime.Format("15:04") + " - " + r.EndTime.Format("15:04"),
				teacher.Fullname, "", "", ""})
		}
	}
	if err := utils.WriteToCSV(data, path); err != nil {
		return err
	}
	return nil
}

// 每天8:00发送当天咨询安排表邮件
func (w *Workflow) SendTodayTimetableMail(mailTo string) {
	today := utils.BeginOfDay(time.Now())
	tomorrow := today.AddDate(0, 0, 1)
	reservations, err := w.mongoClient.GetReservationsBetweenTime(today, tomorrow)
	if err != nil {
		log.Errorf("%v", err)
		return
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
	path := filepath.Join(utils.EXPORT_FOLDER, fmt.Sprintf("timetable_%s%s", todayDate, utils.CSV_FILE_SUFFIX))
	if err = w.ExportTodayReservationTimetableToFile(reservations, path); err != nil {
		log.Errorf("%v", err)
		return
	}
	// email
	title := fmt.Sprintf("【心理发展中心】%s咨询安排表", todayDate)
	m := email.NewMessage(title, title)
	m.From = mail.Address{Name: "", Address: config.Instance().SMTPUser}
	m.To = strings.Split(mailTo, ",")
	m.Attach(path)
	if err := SendEmail(m); err != nil {
		log.Errorf("发送邮件失败：%v", err)
		return
	}
	log.Infof("发送邮件成功")
}

type ByWeekdayOfTimedReservation []*model.TimedReservation

func (ts ByWeekdayOfTimedReservation) Len() int {
	return len(ts)
}

func (ts ByWeekdayOfTimedReservation) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

func (ts ByWeekdayOfTimedReservation) Less(i, j int) bool {
	if ts[i].Weekday != ts[j].Weekday {
		return ts[i].Weekday < ts[j].Weekday
	} else if !ts[i].StartTime.Equal(ts[j].StartTime) {
		return ts[i].StartTime.Before(ts[j].StartTime)
	}
	return ts[i].TeacherId < ts[j].TeacherId
}

func (w *Workflow) WrapTimedReservation(timedReservation *model.TimedReservation) map[string]interface{} {
	var result = make(map[string]interface{})
	if timedReservation == nil {
		return result
	}
	result["id"] = timedReservation.Id.Hex()
	result["weekday"] = timedReservation.Weekday
	result["start_clock"] = timedReservation.StartTime.Format("15:04")
	result["end_clock"] = timedReservation.EndTime.Format("15:04")
	result["status"] = timedReservation.Status
	if timedReservation.TeacherId != "" {
		result["teacher_id"] = timedReservation.TeacherId
		if teacher, err := w.mongoClient.GetTeacherById(timedReservation.TeacherId); err == nil {
			result["teacher_username"] = teacher.Username
			result["teacher_fullname"] = teacher.Fullname
			result["teacher_mobile"] = teacher.Mobile
		}
	}
	return result
}
