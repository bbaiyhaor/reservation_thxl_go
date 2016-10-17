package service

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"net/http"
	"strings"
)

func (s *Service) ViewTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	var result = map[string]interface{}{"state": "SUCCESS"}

	timedReservations, err := s.w.ViewTimetableByAdmin(userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	var timetable = make(map[string]interface{})
	for weekday, trs := range timedReservations {
		var array = make([]interface{}, 0)
		for _, tr := range trs {
			trJson := make(map[string]interface{})
			trJson["timed_reservation_id"] = tr.Id.Hex()
			trJson["weekday"] = tr.Weekday
			trJson["start_clock"] = tr.StartTime.Format("15:04")
			trJson["end_clock"] = tr.EndTime.Format("15:04")
			trJson["status"] = tr.Status.String()
			trJson["teacher_id"] = tr.TeacherId
			if teacher, err := s.w.GetTeacherById(tr.TeacherId); err == nil {
				trJson["teacher_username"] = teacher.Username
				trJson["teacher_fullname"] = teacher.Fullname
				trJson["teacher_mobile"] = teacher.Mobile
			}
			array = append(array, trJson)
		}
		timetable[weekday.String()] = array
	}
	result["timed_reservations"] = timetable

	return result
}

func (s *Service) AddTimedReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	weekday := r.PostFormValue("weekday")
	startTime := r.PostFormValue("start_clock")
	endTime := r.PostFormValue("end_clock")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")
	force := strings.EqualFold(r.PostFormValue("force"), "FORCE")

	var result = map[string]interface{}{"state": "SUCCESS"}

	var timedReservationJson = make(map[string]interface{})
	timedReservation, err := s.w.AddTimetableByAdmin(weekday, startTime, endTime, teacherUsername, teacherFullname,
		teacherMobile, force, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	timedReservationJson["timed_reservation_id"] = timedReservation.Id.Hex()
	timedReservationJson["weekday"] = timedReservation.Weekday
	timedReservationJson["start_clock"] = timedReservation.StartTime.Format("15:04")
	timedReservationJson["end_clock"] = timedReservation.EndTime.Format("15:04")
	timedReservationJson["status"] = timedReservation.Status.String()
	timedReservationJson["teacher_id"] = timedReservation.TeacherId
	if teacher, err := s.w.GetTeacherById(timedReservation.TeacherId); err == nil {
		timedReservationJson["teacher_username"] = teacher.Username
		timedReservationJson["teacher_fullname"] = teacher.Fullname
		timedReservationJson["teacher_mobile"] = teacher.Mobile
	}
	result["timed_reservation"] = timedReservationJson

	return result
}

func (s *Service) EditTimedReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	timedReservationId := r.PostFormValue("timed_reservation_id")
	weekday := r.PostFormValue("weekday")
	startTime := r.PostFormValue("start_clock")
	endTime := r.PostFormValue("end_clock")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")
	force := strings.EqualFold(r.PostFormValue("force"), "FORCE")

	var result = map[string]interface{}{"state": "SUCCESS"}

	var timedReservationJson = make(map[string]interface{})
	timedReservation, err := s.w.EditTimetableByAdmin(timedReservationId, weekday, startTime, endTime, teacherUsername,
		teacherFullname, teacherMobile, force, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	timedReservationJson["timed_reservation_id"] = timedReservation.Id.Hex()
	timedReservationJson["weekday"] = timedReservation.Weekday
	timedReservationJson["start_clock"] = timedReservation.StartTime.Format("15:04")
	timedReservationJson["end_clock"] = timedReservation.EndTime.Format("15:04")
	timedReservationJson["status"] = timedReservation.Status.String()
	timedReservationJson["teacher_id"] = timedReservation.TeacherId
	if teacher, err := s.w.GetTeacherById(timedReservation.TeacherId); err == nil {
		timedReservationJson["teacher_username"] = teacher.Username
		timedReservationJson["teacher_fullname"] = teacher.Fullname
		timedReservationJson["teacher_mobile"] = teacher.Mobile
	}
	result["timed_reservation"] = timedReservationJson

	return result
}

func (s *Service) RemoveTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	r.ParseForm()
	timedReservationIds := []string(r.Form["timed_reservation_ids"])

	var result = map[string]interface{}{"state": "SUCCESS"}

	removed, err := s.w.RemoveTimetablesByAdmin(timedReservationIds, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	result["removed_count"] = removed

	return result
}

func (s *Service) OpenTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	r.ParseForm()
	timedReservationIds := []string(r.Form["timed_reservation_ids"])

	var result = map[string]interface{}{"state": "SUCCESS"}

	opened, err := s.w.OpenTimetablesByAdmin(timedReservationIds, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	result["opened_count"] = opened

	return result
}

func (s *Service) CloseTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	r.ParseForm()
	timedReservationIds := []string(r.Form["timed_reservation_ids"])

	var result = map[string]interface{}{"state": "SUCCESS"}

	closed, err := s.w.CloseTimetablesByAdmin(timedReservationIds, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	result["closed_count"] = closed

	return result
}
