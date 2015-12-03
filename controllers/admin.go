package controllers

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxl_go/buslogic"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"net/http"
	"net/url"
	"time"
)

func ViewReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	var result = map[string]interface{}{"state": "SUCCESS"}
	var rl = buslogic.ReservationLogic{}

	reservations, err := rl.GetReservationsByAdmin(userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := make(map[string]interface{})
		resJson["reservation_id"] = res.Id
		resJson["start_time"] = res.StartTime.In(utils.Location).Format(utils.TIME_PATTERN)
		resJson["end_time"] = res.EndTime.In(utils.Location).Format(utils.TIME_PATTERN)
		resJson["teacher_username"] = res.TeacherUsername
		resJson["teacher_fullname"] = res.TeacherFullname
		resJson["teacher_mobile"] = res.TeacherMobile
		if res.Status == models.AVAILABLE {
			resJson["status"] = models.AVAILABLE.String()
		} else if res.Status == models.RESERVATED && res.StartTime.Before(time.Now().In(utils.Location)) {
			resJson["status"] = models.FEEDBACK.String()
		} else {
			resJson["status"] = models.RESERVATED.String()
		}
		array = append(array, resJson)
	}
	result["reservations"] = array

	return result
}

func ViewMonthlyReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["from_time"]) == 0 {
		ErrorHandler(w, r, errors.New("参数错误"))
		return nil
	}
	fromTime := queryForm["from_time"][0]

	var result = map[string]interface{}{"state": "SUCCESS"}
	var rl = buslogic.ReservationLogic{}

	reservations, err := rl.GetReservationsMonthlyByAdmin(fromTime, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := make(map[string]interface{})
		resJson["reservation_id"] = res.Id
		resJson["start_time"] = res.StartTime.In(utils.Location).Format(utils.TIME_PATTERN)
		resJson["end_time"] = res.EndTime.In(utils.Location).Format(utils.TIME_PATTERN)
		resJson["teacher_username"] = res.TeacherUsername
		resJson["teacher_fullname"] = res.TeacherFullname
		resJson["teacher_mobile"] = res.TeacherMobile
		if res.Status == models.AVAILABLE {
			resJson["status"] = models.AVAILABLE.String()
		} else if res.Status == models.RESERVATED && res.StartTime.Before(time.Now().In(utils.Location)) {
			resJson["status"] = models.FEEDBACK.String()
		} else {
			resJson["status"] = models.RESERVATED.String()
		}
		array = append(array, resJson)
	}
	result["reservations"] = array

	return result
}

func ExportReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["from_time"]) == 0 {
		ErrorHandler(w, r, errors.New("参数错误"))
		return nil
	}
	fromTime := queryForm["from_time"][0]

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	url, err := al.ExportReservationTimetable(fromTime, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	result["url"] = url

	return result
}

func AddReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	startTime := r.PostFormValue("start_time")
	endTime := r.PostFormValue("end_time")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	var reservationJson = make(map[string]interface{})
	reservation, err := al.AddReservationByAdmin(startTime, endTime, teacherUsername, teacherFullname,
		teacherMobile, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	reservationJson["reservation_id"] = reservation.Id
	reservationJson["start_time"] = reservation.StartTime.In(utils.Location).Format(utils.TIME_PATTERN)
	reservationJson["end_time"] = reservation.EndTime.In(utils.Location).Format(utils.TIME_PATTERN)
	reservationJson["teacher_username"] = reservation.TeacherUsername
	reservationJson["teacher_fullname"] = reservation.TeacherFullname
	reservationJson["teacher_mobile"] = reservation.TeacherMobile
	result["reservation"] = reservationJson

	return result
}

func EditReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	startTime := r.PostFormValue("start_time")
	endTime := r.PostFormValue("end_time")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	var reservationJson = make(map[string]interface{})
	reservation, err := al.EditReservationByAdmin(reservationId, startTime, endTime, teacherUsername,
		teacherFullname, teacherMobile, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	reservationJson["reservation_id"] = reservation.Id
	reservationJson["start_time"] = reservation.StartTime.In(utils.Location).Format(utils.TIME_PATTERN)
	reservationJson["end_time"] = reservation.EndTime.In(utils.Location).Format(utils.TIME_PATTERN)
	reservationJson["teacher_username"] = reservation.TeacherUsername
	reservationJson["teacher_fullname"] = reservation.TeacherFullname
	reservationJson["teacher_mobile"] = reservation.TeacherMobile
	result["reservation"] = reservationJson

	return result
}

func RemoveReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	r.ParseForm()
	reservationIds := []string(r.Form["reservation_ids"])

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	removed, err := al.RemoveReservationsByAdmin(reservationIds, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	result["removed_count"] = removed

	return result
}

func CancelReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	r.ParseForm()
	reservationIds := []string(r.Form["reservation_ids"])

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	removed, err := al.CancelReservationsByAdmin(reservationIds, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	result["removed_count"] = removed

	return result
}

func GetFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	var feedback = make(map[string]interface{})
	reservation, err := al.GetFeedbackByAdmin(reservationId, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	feedback["problem"] = reservation.TeacherFeedback.Problem
	feedback["record"] = reservation.TeacherFeedback.Record
	result["feedback"] = feedback

	return result
}

func SubmitFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	problem := r.PostFormValue("problem")
	record := r.PostFormValue("record")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	_, err := al.SubmitFeedbackByAdmin(reservationId, problem, record, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}

	return result
}

func GetStudentInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}
	var ul = buslogic.UserLogic{}

	var studentJson = make(map[string]interface{})
	student, err := al.GetStudentInfoByAdmin(reservationId, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	studentJson["student_username"] = student.Username
	studentJson["student_fullname"] = student.Fullname
	studentJson["student_gender"] = student.Gender
	studentJson["student_email"] = student.Email
	studentJson["student_school"] = student.School
	studentJson["student_grade"] = student.Grade
	studentJson["student_current_address"] = student.CurrentAddress
	studentJson["student_mobile"] = student.Mobile
	studentJson["student_birthday"] = student.Birthday
	studentJson["student_family_address"] = student.FamilyAddress
	if !student.Experience.IsEmpty() {
		studentJson["student_experience_time"] = student.Experience.Time
		studentJson["student_experience_location"] = student.Experience.Location
		studentJson["student_experience_teacher"] = student.Experience.Teacher
	}
	studentJson["student_father_age"] = student.FatherAge
	studentJson["student_father_job"] = student.FatherJob
	studentJson["student_father_edu"] = student.FatherEdu
	studentJson["student_mother_age"] = student.MotherAge
	studentJson["student_mother_job"] = student.MotherJob
	studentJson["student_mother_edu"] = student.MotherEdu
	studentJson["student_parent_marriage"] = student.ParentMarriage
	studentJson["student_significant"] = student.Significant
	studentJson["student_problem"] = student.Problem
	if len(student.BindedTeacher) != 0 {
		teacher, err := ul.GetTeacherByUsername(student.BindedTeacher)
		if err != nil {
			studentJson["student_binded_teacher_username"] = "无"
			studentJson["student_binded_teacher_fullname"] = ""
		}
		studentJson["student_binded_teacher_username"] = teacher.Username
		studentJson["student_binded_teacher_fullname"] = teacher.Fullname
	} else {
		studentJson["student_binded_teacher_username"] = "无"
		studentJson["student_binded_teacher_fullname"] = ""
	}
	result["student_info"] = studentJson

	return result
}

func ExportStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	studentUsername := r.PostFormValue("student_username")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	url, err := al.ExportStudentByAdmin(studentUsername, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	result["url"] = url

	return result
}

func UnbindStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	studentUsername := r.PostFormValue("student_username")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}
	var ul = buslogic.UserLogic{}

	var studentJson = make(map[string]interface{})
	student, err := al.UnbindStudentByAdmin(studentUsername, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	if len(student.BindedTeacher) != 0 {
		teacher, err := ul.GetTeacherByUsername(student.BindedTeacher)
		if err != nil {
			studentJson["student_binded_teacher_username"] = "无"
			studentJson["student_binded_teacher_fullname"] = ""
		}
		studentJson["student_binded_teacher_username"] = teacher.Username
		studentJson["student_binded_teacher_fullname"] = teacher.Fullname
	} else {
		studentJson["student_binded_teacher_username"] = "无"
		studentJson["student_binded_teacher_fullname"] = ""
	}
	result["student_info"] = studentJson

	return result
}

func BindStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	studentUsername := r.PostFormValue("student_username")
	teacherUsername := r.PostFormValue("teacher_username")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}
	var ul = buslogic.UserLogic{}

	var studentJson = make(map[string]interface{})
	student, err := al.BindStudentByAdmin(studentUsername, teacherUsername, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	} else if len(student.BindedTeacher) != 0 {
		teacher, err := ul.GetTeacherByUsername(student.BindedTeacher)
		if err != nil {
			studentJson["student_binded_teacher_username"] = "无"
			studentJson["student_binded_teacher_fullname"] = ""
		}
		studentJson["student_binded_teacher_username"] = teacher.Username
		studentJson["student_binded_teacher_fullname"] = teacher.Fullname
	} else {
		studentJson["student_binded_teacher_username"] = "无"
		studentJson["student_binded_teacher_fullname"] = ""
	}
	result["student_info"] = studentJson

	return result
}

func QueryStudentInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	studentUsername := r.PostFormValue("student_username")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}
	var ul = buslogic.UserLogic{}

	var studentJson = make(map[string]interface{})
	student, err := al.QueryStudentInfoByAdmin(studentUsername, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	studentJson["student_username"] = student.Username
	studentJson["student_fullname"] = student.Fullname
	studentJson["student_gender"] = student.Gender
	studentJson["student_email"] = student.Email
	studentJson["student_school"] = student.School
	studentJson["student_grade"] = student.Grade
	studentJson["student_current_address"] = student.CurrentAddress
	studentJson["student_mobile"] = student.Mobile
	studentJson["student_birthday"] = student.Birthday
	studentJson["student_family_address"] = student.FamilyAddress
	if !student.Experience.IsEmpty() {
		studentJson["student_experience_time"] = student.Experience.Time
		studentJson["student_experience_location"] = student.Experience.Location
		studentJson["student_experience_teacher"] = student.Experience.Teacher
	}
	studentJson["student_father_age"] = student.FatherAge
	studentJson["student_father_job"] = student.FatherJob
	studentJson["student_father_edu"] = student.FatherEdu
	studentJson["student_mother_age"] = student.MotherAge
	studentJson["student_mother_job"] = student.MotherJob
	studentJson["student_mother_edu"] = student.MotherEdu
	studentJson["student_parent_marriage"] = student.ParentMarriage
	studentJson["student_significant"] = student.Significant
	studentJson["student_problem"] = student.Problem
	if len(student.BindedTeacher) != 0 {
		teacher, err := ul.GetTeacherByUsername(student.BindedTeacher)
		if err != nil {
			studentJson["student_binded_teacher_username"] = "无"
			studentJson["student_binded_teacher_fullname"] = ""
		}
		studentJson["student_binded_teacher_username"] = teacher.Username
		studentJson["student_binded_teacher_fullname"] = teacher.Fullname
	} else {
		studentJson["student_binded_teacher_username"] = "无"
		studentJson["student_binded_teacher_fullname"] = ""
	}
	result["student_info"] = studentJson

	return result
}

func SearchTeacherByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMoble := r.PostFormValue("teacher_mobile")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	var teacherJson = make(map[string]interface{})
	teacher, err := al.SearchTeacherByAdmin(teacherFullname, teacherUsername, teacherMoble,
		userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	teacherJson["teacher_username"] = teacher.Username
	teacherJson["teacher_fullname"] = teacher.Fullname
	teacherJson["teacher_mobile"] = teacher.Mobile
	result["teacher"] = teacherJson

	return result
}
