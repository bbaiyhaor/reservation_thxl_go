package controllers

import (
	"github.com/shudiwsh2009/reservation_thxl_go/buslogic"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"net/http"
	"time"
)

func ViewReservationsByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	var result = map[string]interface{}{"state": "SUCCESS"}
	var ul = buslogic.UserLogic{}
	var rl = buslogic.ReservationLogic{}

	teacher, err := ul.GetTeacherById(userId)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	var teacherJson = make(map[string]interface{})
	teacherJson["teacher_fullname"] = teacher.Fullname
	teacherJson["teacher_mobile"] = teacher.Mobile
	result["teacher_info"] = teacherJson

	reservations, err := rl.GetReservationsByTeacher(userId, userType)
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

func GetFeedbackByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var tl = buslogic.TeacherLogic{}

	var feedback = make(map[string]interface{})
	reservation, err := tl.GetFeedbackByTeacher(reservationId, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	feedback["problem"] = reservation.TeacherFeedback.Problem
	feedback["solution"] = reservation.TeacherFeedback.Record
	result["feedback"] = feedback

	return result
}

func SubmitFeedbackByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	problem := r.PostFormValue("problem")
	record := r.PostFormValue("record")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var tl = buslogic.TeacherLogic{}

	_, err := tl.SubmitFeedbackByTeacher(reservationId, problem, record, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}

	return result
}

func GetStudentInfoByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var tl = buslogic.TeacherLogic{}
	var ul = buslogic.UserLogic{}

	var studentJson = make(map[string]interface{})
	student, err := tl.GetStudentInfoByTeacher(reservationId, userId, userType)
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
