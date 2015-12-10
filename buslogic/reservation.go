package buslogic

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"strings"
	"time"
	"sort"
)

type ReservationLogic struct {
}

// 学生查看前后一周内的所有咨询
func (rl *ReservationLogic) GetReservationsByStudent(userId string, userType models.UserType) ([]*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.STUDENT {
		return nil, errors.New("请重新登录")
	}
	student, err := models.GetStudentById(userId)
	if err != nil {
		return nil, errors.New("请先登录")
	} else if student.UserType != models.STUDENT {
		return nil, errors.New("请重新登录")
	}
	from := utils.GetNow().AddDate(0, 0, -7)
	to := utils.GetNow().AddDate(0, 0, 7)
	reservations, err := models.GetReservationsBetweenTime(from, to)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	var result []*models.Reservation
	for _, r := range reservations {
		if r.Status == models.AVAILABLE && r.StartTime.Before(utils.GetNow()) {
			continue
		} else if strings.EqualFold(r.StudentId, student.Id.Hex()) {
			if !r.TeacherFeedback.IsEmpty() && r.TeacherFeedback.Participants[0] == 0 {
				// 学生未参与的咨询不展示给学生（家长、老师或者辅导员参加）
				continue
			}
			result = append(result, r)
		} else if strings.EqualFold(r.TeacherId, student.BindedTeacherId) && r.Status == models.AVAILABLE {
			result = append(result, r)
		} else if len(student.BindedTeacherId) == 0 && r.Status == models.AVAILABLE {
			result = append(result, r)
		}
	}
	timedReservations, err := models.GetTimedReservationsAll()
	if err != nil {
		return result, nil
	}
	today := utils.GetToday()
	for _, tr := range timedReservations {
		if len(student.BindedTeacherId) != 0 && !strings.EqualFold(student.BindedTeacherId, tr.TeacherId) {
			continue
		}
		minusWeekday := int(tr.Weekday - today.Weekday())
		if minusWeekday < 0 {
			minusWeekday = 7 - minusWeekday
		}
		date := today.AddDate(0, 0, minusWeekday)
		if utils.ConcatTime(date, tr.StartTime).Before(utils.GetNow()) {
			date = today.AddDate(0, 0, 7)
		}
		if !tr.Exceptions[date.Format(utils.DATE_PATTERN)] && !tr.Timed[date.Format(utils.DATE_PATTERN)] {
			result = append(result, tr.ToReservation(date))
		}
	}
	sort.Sort(models.ReservationSlice(result))
	return result, nil
}

// 咨询师查看负7天之后的所有咨询
func (rl *ReservationLogic) GetReservationsByTeacher(userId string, userType models.UserType) ([]*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return nil, errors.New("权限不足")
	}
	teacher, err := models.GetTeacherById(userId)
	if err != nil {
		return nil, errors.New("请先登录")
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	}
	from := utils.GetNow().AddDate(0, 0, -7)
	reservations, err := models.GetReservationsAfterTime(from)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	var result []*models.Reservation
	for _, r := range reservations {
		if r.Status == models.AVAILABLE && r.StartTime.Before(utils.GetNow()) {
			continue
		} else if strings.EqualFold(r.TeacherId, teacher.Id.Hex()) {
			result = append(result, r)
		}
	}
	return result, nil
}

// 管理员查看负7天之后的所有咨询
func (rl *ReservationLogic) GetReservationsByAdmin(userId string, userType models.UserType) ([]*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	}
	admin, err := models.GetTeacherById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	from := utils.GetNow().AddDate(0, 0, -7)
	reservations, err := models.GetReservationsAfterTime(from)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	var result []*models.Reservation
	for _, r := range reservations {
		if r.Status == models.AVAILABLE && r.StartTime.Before(utils.GetNow()) {
			continue
		}
		result = append(result, r)
	}
	timedReservations, err := models.GetTimedReservationsAll()
	if err != nil {
		return result, nil
	}
	today := utils.GetToday()
	for _, tr := range timedReservations {
		minusWeekday := int(tr.Weekday - today.Weekday())
		if minusWeekday < 0 {
			minusWeekday = 7 - minusWeekday
		}
		date := today.AddDate(0, 0, minusWeekday)
		if utils.ConcatTime(date, tr.StartTime).Before(utils.GetNow()) {
			date = today.AddDate(0, 0, 7)
		}
		if tr.Exceptions[date.Format(utils.DATE_PATTERN)] || tr.Timed[date.Format(utils.DATE_PATTERN)] {
			result = append(result, tr.ToReservation(date))
		}
		for i := 1; i <= 3; i++ {
			// 改变i的上阈值可以改变预设咨询的查看范围
			date = date.AddDate(0, 0, 7)
			if tr.Exceptions[date.Format(utils.DATE_PATTERN)] || tr.Timed[date.Format(utils.DATE_PATTERN)] {
				result = append(result, tr.ToReservation(date))
			}
			result = append(result, tr.ToReservation(date.AddDate(0, 0, 7)))
		}
	}
	sort.Sort(models.ReservationSlice(result))
	return result, nil
}

// 管理员查看指定日期后30天内的所有咨询（看不到预设咨询）
func (rl *ReservationLogic) GetReservationsMonthlyByAdmin(from string, userId string, userType models.UserType) ([]*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	}
	admin, err := models.GetTeacherById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	fromDate, err := time.ParseInLocation(utils.DATE_PATTERN, from, utils.Location)
	if err != nil {
		return nil, errors.New("时间格式错误")
	}
	toDate := fromDate.AddDate(0, 0, 30)
	reservations, err := models.GetReservationsBetweenTime(fromDate, toDate)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	var result []*models.Reservation
	for _, r := range reservations {
		if r.Status == models.AVAILABLE && r.StartTime.Before(time.Now().In(utils.Location)) {
			continue
		}
		result = append(result, r)
	}
	return result, nil
}
