package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/config"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/util"
	"errors"
	"fmt"
	"github.com/scorredoira/email"
	"log"
	"net/mail"
	"sort"
	"strings"
	"time"
)

// 管理员查看时间表
func (w *Workflow) ViewTimetableByAdmin(userId string, userType model.UserType) (map[time.Weekday][]*model.TimedReservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, errors.New("权限不足")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, errors.New("管理员账户出错，请联系技术支持")
	}
	timedReservations := make(map[time.Weekday][]*model.TimedReservation)
	for i := time.Sunday; i <= time.Saturday; i++ {
		if trs, err := w.model.GetTimedReservationsByWeekday(i); err == nil {
			sort.Sort(model.TimedReservationSlice(trs))
			timedReservations[i] = trs
		}
	}
	return timedReservations, nil
}

// 管理员添加时间表
func (w *Workflow) AddTimetableByAdmin(weekday string, startClock string, endClock string,
	teacherUsername string, teacherFullname string, teacherMobile string, force bool,
	userId string, userType model.UserType) (*model.TimedReservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(startClock) == 0 {
		return nil, errors.New("开始时间为空")
	} else if len(endClock) == 0 {
		return nil, errors.New("结束时间为空")
	} else if len(teacherUsername) == 0 {
		return nil, errors.New("咨询师工号为空")
	} else if len(teacherFullname) == 0 {
		return nil, errors.New("咨询师姓名为空")
	} else if len(teacherMobile) == 0 {
		return nil, errors.New("咨询师手机号为空")
	} else if !util.IsMobile(teacherMobile) {
		return nil, errors.New("咨询师手机号格式不正确")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, errors.New("管理员账户出错，请联系技术支持")
	}
	week, err := util.StringToWeekday(weekday)
	if err != nil {
		return nil, errors.New("星期格式错误")
	}
	start, err := time.ParseInLocation("2006-01-02 15:04", "2006-01-02 "+startClock, time.Local)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	end, err := time.ParseInLocation("2006-01-02 15:04", "2006-01-02 "+endClock, time.Local)
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
		return nil, errors.New("咨询师权限不足")
	} else if !strings.EqualFold(teacher.Fullname, teacherFullname) || !strings.EqualFold(teacher.Mobile, teacherMobile) {
		if !force {
			return nil, errors.New(model.CHECK_MESSAGE)
		}
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if err = w.model.UpsertTeacher(teacher); err != nil {
			return nil, errors.New("获取数据失败")
		}
	}
	timedReservation, err := w.model.AddTimedReservation(week, start, end, teacher.Id.Hex())
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	return timedReservation, nil
}

// 管理员编辑时间表
func (w *Workflow) EditTimetableByAdmin(timedReservationId string, weekday string,
	startClock string, endClock string, teacherUsername string, teacherFullname string, teacherMobile string,
	force bool, userId string, userType model.UserType) (*model.TimedReservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(timedReservationId) == 0 {
		return nil, errors.New("咨询已下架")
	} else if len(startClock) == 0 {
		return nil, errors.New("开始时间为空")
	} else if len(endClock) == 0 {
		return nil, errors.New("结束时间为空")
	} else if len(teacherUsername) == 0 {
		return nil, errors.New("咨询师工号为空")
	} else if len(teacherFullname) == 0 {
		return nil, errors.New("咨询师姓名为空")
	} else if len(teacherMobile) == 0 {
		return nil, errors.New("咨询师手机号为空")
	} else if !util.IsMobile(teacherMobile) {
		return nil, errors.New("咨询师手机号格式不正确")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, errors.New("管理员账户出错，请联系技术支持")
	}
	timedReservation, err := w.model.GetTimedReservationById(timedReservationId)
	if err != nil || timedReservation.Status == model.DELETED {
		return nil, errors.New("咨询已下架")
	}
	week, err := util.StringToWeekday(weekday)
	if err != nil {
		return nil, errors.New("星期格式错误")
	}
	start, err := time.ParseInLocation("2006-01-02 15:04", "2006-01-02"+startClock, time.Local)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	end, err := time.ParseInLocation("2006-01-02 15:04", "2006-01-02"+endClock, time.Local)
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
		return nil, errors.New("咨询师权限不足")
	} else if !strings.EqualFold(teacher.Fullname, teacherFullname) || !strings.EqualFold(teacher.Mobile, teacherMobile) {
		if !force {
			return nil, errors.New(model.CHECK_MESSAGE)
		}
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if err = w.model.UpsertTeacher(teacher); err != nil {
			return nil, errors.New("获取数据失败")
		}
	}
	timedReservation.Weekday = week
	timedReservation.StartTime = start
	timedReservation.EndTime = end
	timedReservation.Status = model.CLOSED
	timedReservation.TeacherId = teacher.Id.Hex()
	if err = w.model.UpsertTimedReservation(timedReservation); err != nil {
		return nil, errors.New("获取数据失败")
	}
	return timedReservation, nil
}

// 管理员删除时间表
func (w *Workflow) RemoveTimetablesByAdmin(timedReservationIds []string, userId string, userType model.UserType) (int, error) {
	if len(userId) == 0 {
		return 0, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return 0, errors.New("权限不足")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return 0, errors.New("管理员账户出错，请联系技术支持")
	}
	removed := 0
	for _, id := range timedReservationIds {
		if timedReservation, err := w.model.GetTimedReservationById(id); err == nil {
			timedReservation.Status = model.DELETED
			if err = w.model.UpsertTimedReservation(timedReservation); err == nil {
				removed++
			}
		}
	}
	return removed, nil
}

// 管理员开启时间表
func (w *Workflow) OpenTimetablesByAdmin(timedReservationIds []string, userId string, userType model.UserType) (int, error) {
	if len(userId) == 0 {
		return 0, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return 0, errors.New("权限不足")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return 0, errors.New("管理员账户出错，请联系技术支持")
	}
	opened := 0
	for _, id := range timedReservationIds {
		if timedReservation, err := w.model.GetTimedReservationById(id); err == nil {
			if timedReservation.Status == model.CLOSED {
				timedReservation.Status = model.AVAILABLE
				if err = w.model.UpsertTimedReservation(timedReservation); err == nil {
					opened++
				}
			}
		}
	}
	return opened, nil
}

// 管理员关闭时间表
func (w *Workflow) CloseTimetablesByAdmin(timedReservationIds []string, userId string, userType model.UserType) (int, error) {
	if len(userId) == 0 {
		return 0, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return 0, errors.New("权限不足")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return 0, errors.New("管理员账户出错，请联系技术支持")
	}
	closed := 0
	for _, id := range timedReservationIds {
		if timedReservation, err := w.model.GetTimedReservationById(id); err == nil {
			if timedReservation.Status == model.AVAILABLE {
				timedReservation.Status = model.CLOSED
				if err = w.model.UpsertTimedReservation(timedReservation); err == nil {
					closed++
				}
			}
		}
	}
	return closed, nil
}

func (w *Workflow) ExportTodayReservationTimetableToFile(reservations []*model.Reservation, filename string) error {
	data := make([][]string, 0)
	today := util.BeginOfDay(time.Now())
	data = append(data, []string{today.Format("2006-01-02")})
	data = append(data, []string{"时间", "咨询师", "学生姓名", "联系方式"})
	for _, r := range reservations {
		teacher, err := w.model.GetTeacherById(r.TeacherId)
		if err != nil {
			continue
		}
		if student, err := w.model.GetStudentById(r.StudentId); err == nil {
			data = append(data, []string{r.StartTime.Format("15:04") + " - " + r.EndTime.Format("15:04"),
				teacher.Fullname, student.Fullname, student.Mobile})
		} else {
			data = append(data, []string{r.StartTime.Format("15:04") + " - " + r.EndTime.Format("15:04"),
				teacher.Fullname, "", ""})
		}
	}
	if err := util.WriteToCSV(data, filename); err != nil {
		return err
	}
	return nil
}

// 每天8:00发送当天咨询安排表邮件
func (w *Workflow) SendTodayTimetableMail(mailTo string) {
	today := util.BeginOfDay(time.Now())
	tomorrow := today.AddDate(0, 0, 1)
	reservations, err := w.model.GetReservationsBetweenTime(today, tomorrow)
	if err != nil {
		log.Printf("%v", err)
		return
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
	if err = w.ExportTodayReservationTimetableToFile(reservations, filename); err != nil {
		log.Printf("%v", err)
		return
	}
	// email
	title := fmt.Sprintf("【心理发展中心】%s咨询安排表", todayDate)
	m := email.NewMessage(title, title)
	m.From = mail.Address{Name: "", Address: config.Instance().SMTPUser}
	m.To = strings.Split(mailTo, ",")
	m.Attach(fmt.Sprintf("%s%s", util.ExportFolder, filename))
	if err := util.SendEmail(m); err != nil {
		log.Printf("发送邮件失败：%v", err)
		return
	}
	log.Printf("发送邮件成功")
}
