package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/util"
	"errors"
	"sort"
	"strings"
	"time"
)

// 学生查看前后一周内的所有咨询
func (w *Workflow) GetReservationsByStudent(userId string, userType model.UserType) ([]*model.Reservation, error) {
	if userId == "" {
		return nil, errors.New("请先登录")
	} else if userType != model.STUDENT {
		return nil, errors.New("请重新登录")
	}
	student, err := w.model.GetStudentById(userId)
	if err != nil {
		return nil, errors.New("请先登录")
	} else if student.UserType != model.STUDENT {
		return nil, errors.New("请重新登录")
	}
	from := time.Now().AddDate(0, 0, -7)
	to := time.Now().AddDate(0, 0, 7).Add(-90 * time.Minute)
	reservations, err := w.model.GetReservationsBetweenTime(from, to)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	var result []*model.Reservation
	for _, r := range reservations {
		if r.Status == model.AVAILABLE && r.StartTime.Before(time.Now()) {
			continue
		} else if r.StudentId == student.Id.Hex() {
			if !r.TeacherFeedback.IsEmpty() && r.TeacherFeedback.Participants[0] == 0 {
				// 学生未参与的咨询不展示给学生（家长、老师或者辅导员参加）
				continue
			}
			result = append(result, r)
		} else if student.BindedTeacherId == "" || student.BindedTeacherId == r.TeacherId {
			result = append(result, r)
		}
		//} else if r.TeacherId == student.BindedTeacherId && r.Status == models.AVAILABLE {
		//	result = append(result, r)
		//} else if student.BindedTeacherId == "" && r.Status == models.AVAILABLE {
		//	result = append(result, r)
		//}
	}
	timedReservations, err := w.model.GetTimedReservationsAll()
	if err != nil {
		return result, nil
	}
	today := util.BeginOfDay(time.Now())
	for _, tr := range timedReservations {
		if tr.Status != model.AVAILABLE {
			continue
		}
		if len(student.BindedTeacherId) != 0 && !strings.EqualFold(student.BindedTeacherId, tr.TeacherId) {
			continue
		}
		minusWeekday := int(tr.Weekday - today.Weekday())
		if minusWeekday < 0 {
			minusWeekday += 7
		}
		date := today.AddDate(0, 0, minusWeekday)
		if util.ConcatTime(date, tr.StartTime).Before(time.Now()) {
			date = today.AddDate(0, 0, 7)
		}
		if !tr.Exceptions[date.Format("2006-01-02")] && !tr.Timed[date.Format("2006-01-02")] {
			result = append(result, tr.ToReservation(date))
		}
	}
	sort.Sort(model.ReservationSlice(result))
	return result, nil
}

// 咨询师查看负7天之后的所有咨询
func (w *Workflow) GetReservationsByTeacher(userId string, userType model.UserType) ([]*model.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != model.TEACHER {
		return nil, errors.New("权限不足")
	}
	teacher, err := w.model.GetTeacherById(userId)
	if err != nil {
		return nil, errors.New("请先登录")
	} else if teacher.UserType != model.TEACHER {
		return nil, errors.New("权限不足")
	}
	from := time.Now().AddDate(0, 0, -7)
	reservations, err := w.model.GetReservationsAfterTime(from)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	var result []*model.Reservation
	for _, r := range reservations {
		if r.Status == model.AVAILABLE && r.StartTime.Before(time.Now()) {
			continue
		} else if strings.EqualFold(r.TeacherId, teacher.Id.Hex()) {
			result = append(result, r)
		}
	}
	if timedReservations, err := w.model.GetTimedReservationsByTeacherId(teacher.Id.Hex()); err == nil {
		today := util.BeginOfDay(time.Now())
		for _, tr := range timedReservations {
			if tr.Status != model.AVAILABLE {
				continue
			}
			minusWeekday := int(tr.Weekday - today.Weekday())
			if minusWeekday < 0 {
				minusWeekday += 7
			}
			date := today.AddDate(0, 0, minusWeekday)
			if util.ConcatTime(date, tr.StartTime).Before(time.Now()) {
				date = today.AddDate(0, 0, 7)
			}
			if !tr.Exceptions[date.Format("2006-01-02")] && !tr.Timed[date.Format("2006-01-02")] {
				result = append(result, tr.ToReservation(date))
			}
			for i := 1; i <= 3; i++ {
				// 改变i的上阈值可以改变预设咨询的查看范围
				date = date.AddDate(0, 0, 7)
				if !tr.Exceptions[date.Format("2006-01-02")] && !tr.Timed[date.Format("2006-01-02")] {
					result = append(result, tr.ToReservation(date))
				}
			}
		}
	}
	sort.Sort(model.ReservationSlice(result))
	return result, nil
}

// 管理员查看负7天之后的所有咨询
func (w *Workflow) GetReservationsByAdmin(userId string, userType model.UserType) ([]*model.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, errors.New("权限不足")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	from := time.Now().AddDate(0, 0, -7)
	reservations, err := w.model.GetReservationsAfterTime(from)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	var result []*model.Reservation
	for _, r := range reservations {
		if r.Status == model.AVAILABLE && r.StartTime.Before(time.Now()) {
			continue
		}
		result = append(result, r)
	}
	if timedReservations, err := w.model.GetTimedReservationsAll(); err == nil {
		today := util.BeginOfDay(time.Now())
		for _, tr := range timedReservations {
			if tr.Status != model.AVAILABLE {
				continue
			}
			minusWeekday := int(tr.Weekday - today.Weekday())
			if minusWeekday < 0 {
				minusWeekday += 7
			}
			date := today.AddDate(0, 0, minusWeekday)
			if util.ConcatTime(date, tr.StartTime).Before(time.Now()) {
				date = today.AddDate(0, 0, 7)
			}
			if !tr.Exceptions[date.Format("2006-01-02")] && !tr.Timed[date.Format("2006-01-02")] {
				result = append(result, tr.ToReservation(date))
			}
			for i := 1; i <= 3; i++ {
				// 改变i的上阈值可以改变预设咨询的查看范围
				date = date.AddDate(0, 0, 7)
				if !tr.Exceptions[date.Format("2006-01-02")] && !tr.Timed[date.Format("2006-01-02")] {
					result = append(result, tr.ToReservation(date))
				}
			}
		}
	}
	sort.Sort(model.ReservationSlice(result))
	return result, nil
}

// 管理员查看指定日期的所有咨询
func (w *Workflow) GetReservationsDailyByAdmin(fromDate string, userId string, userType model.UserType) ([]*model.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, errors.New("权限不足")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	from, err := time.ParseInLocation("2006-01-02", fromDate, time.Local)
	if err != nil {
		return nil, errors.New("时间格式错误")
	}
	to := from.AddDate(0, 0, 1)
	reservations, err := w.model.GetReservationsBetweenTime(from, to)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	if timedReservations, err := w.model.GetTimedReservationsByWeekday(from.Weekday()); err == nil {
		for _, tr := range timedReservations {
			if tr.Status == model.AVAILABLE && !tr.Exceptions[fromDate] && !tr.Timed[fromDate] {
				reservations = append(reservations, tr.ToReservation(from))
			}
		}
	}
	sort.Sort(model.ReservationSlice(reservations))
	return reservations, nil
}

// 管理员通过咨询师工号查询咨询
func (w *Workflow) GetReservationsWithTeacherUsernameByAdmin(teacherUsername string, userId string, userType model.UserType) ([]*model.Reservation, error) {
	if userId == "" {
		return nil, errors.New("请先登录")
	} else if userType != model.ADMIN {
		return nil, errors.New("权限不足")
	} else if teacherUsername == "" {
		return nil, errors.New("咨询师工号为空")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil || admin.UserType != model.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	teacher, err := w.model.GetTeacherByUsername(teacherUsername)
	if err != nil {
		return nil, errors.New("咨询师不存在")
	}
	from := time.Now().AddDate(0, 0, -7)
	reservations, err := w.model.GetReservationsAfterTime(from)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	var result []*model.Reservation
	for _, r := range reservations {
		if r.Status == model.AVAILABLE && r.StartTime.Before(time.Now()) {
			continue
		}
		if r.TeacherId != teacher.Id.Hex() {
			continue
		}
		result = append(result, r)
	}
	if timedReservations, err := w.model.GetTimedReservationsByTeacherId(teacher.Id.Hex()); err == nil {
		today := util.BeginOfDay(time.Now())
		for _, tr := range timedReservations {
			if tr.Status != model.AVAILABLE {
				continue
			}
			minusWeekday := int(tr.Weekday - today.Weekday())
			if minusWeekday < 0 {
				minusWeekday += 7
			}
			date := today.AddDate(0, 0, minusWeekday)
			if util.ConcatTime(date, tr.StartTime).Before(time.Now()) {
				date = today.AddDate(0, 0, 7)
			}
			if !tr.Exceptions[date.Format("2006-01-02")] && !tr.Timed[date.Format("2006-01-02")] {
				result = append(result, tr.ToReservation(date))
			}
			for i := 1; i <= 3; i++ {
				// 改变i的上阈值可以改变预设咨询的查看范围
				date = date.AddDate(0, 0, 7)
				if !tr.Exceptions[date.Format("2006-01-02")] && !tr.Timed[date.Format("2006-01-02")] {
					result = append(result, tr.ToReservation(date))
				}
			}
		}
	}
	sort.Sort(model.ReservationSlice(result))
	return result, nil
}
