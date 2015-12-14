package controllers

import (
	"github.com/shudiwsh2009/reservation_thxl_go/buslogic"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"net/http"
)

func ViewTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}
	var ul = buslogic.UserLogic{}

	timedReservations, err := al.ViewTimetableByAdmin(userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	var timetable = make(map[string]interface{})
	for weekday, trs := range timedReservations {
		var array = make([]interface{}, 0)
		for _, tr := range trs {
			trJson := make(map[string]interface{})
			trJson["timed_reservation_id"] = tr.Id.Hex()
			trJson["weekday"] = tr.Weekday
			trJson["start_clock"] = tr.StartTime.In(utils.Location).Format(utils.CLOCK_PATTERN)
			trJson["end_clock"] = tr.EndTime.In(utils.Location).Format(utils.CLOCK_PATTERN)
			trJson["status"] = tr.Status.String()
			trJson["teacher_id"] = tr.TeacherId
			if teacher, err := ul.GetTeacherById(tr.TeacherId); err == nil {
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

func AddTimedReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	weekday := r.PostFormValue("weekday")
	startTime := r.PostFormValue("start_clock")
	endTime := r.PostFormValue("end_clock")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}
	var ul = buslogic.UserLogic{}

	var timedReservationJson = make(map[string]interface{})
	timedReservation, err := al.AddTimetableByAdmin(weekday, startTime, endTime, teacherUsername, teacherFullname,
		teacherMobile, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	timedReservationJson["timed_reservation_id"] = timedReservation.Id.Hex()
	timedReservationJson["weekday"] = timedReservation.Weekday
	timedReservationJson["start_clock"] = timedReservation.StartTime.In(utils.Location).Format(utils.CLOCK_PATTERN)
	timedReservationJson["end_clock"] = timedReservation.EndTime.In(utils.Location).Format(utils.CLOCK_PATTERN)
	timedReservationJson["status"] = timedReservation.Status.String()
	timedReservationJson["teacher_id"] = timedReservation.TeacherId
	if teacher, err := ul.GetTeacherById(timedReservation.TeacherId); err == nil {
		timedReservationJson["teacher_username"] = teacher.Username
		timedReservationJson["teacher_fullname"] = teacher.Fullname
		timedReservationJson["teacher_mobile"] = teacher.Mobile
	}
	result["timed_reservation"] = timedReservationJson

	return result
}

func EditTimedReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	timedReservationId := r.PostFormValue("timed_reservation_id")
	weekday := r.PostFormValue("weekday")
	startTime := r.PostFormValue("start_clock")
	endTime := r.PostFormValue("end_clock")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}
	var ul = buslogic.UserLogic{}

	var timedReservationJson = make(map[string]interface{})
	timedReservation, err := al.EditTimetableByAdmin(timedReservationId, weekday, startTime, endTime, teacherUsername,
		teacherFullname, teacherMobile, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	timedReservationJson["timed_reservation_id"] = timedReservation.Id.Hex()
	timedReservationJson["weekday"] = timedReservation.Weekday
	timedReservationJson["start_clock"] = timedReservation.StartTime.In(utils.Location).Format(utils.CLOCK_PATTERN)
	timedReservationJson["end_clock"] = timedReservation.EndTime.In(utils.Location).Format(utils.CLOCK_PATTERN)
	timedReservationJson["status"] = timedReservation.Status.String()
	timedReservationJson["teacher_id"] = timedReservation.TeacherId
	if teacher, err := ul.GetTeacherById(timedReservation.TeacherId); err == nil {
		timedReservationJson["teacher_username"] = teacher.Username
		timedReservationJson["teacher_fullname"] = teacher.Fullname
		timedReservationJson["teacher_mobile"] = teacher.Mobile
	}
	result["timed_reservation"] = timedReservationJson

	return result
}

func RemoveTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	r.ParseForm()
	timedReservationIds := []string(r.Form["timed_reservation_ids"])

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	removed, err := al.RemoveTimetablesByAdmin(timedReservationIds, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	result["removed_count"] = removed

	return result
}
