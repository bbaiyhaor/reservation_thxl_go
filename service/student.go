package service

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"net/http"
	"strconv"
	"time"
)

func (s *Service) ViewReservationsByStudent(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	var result = map[string]interface{}{"state": "SUCCESS"}

	student, err := s.w.GetStudentById(userId)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	var studentJson = make(map[string]interface{})
	studentJson["student_fullname"] = student.Fullname
	studentJson["student_gender"] = student.Gender
	studentJson["student_email"] = student.Email
	studentJson["student_school"] = student.School
	studentJson["student_grade"] = student.Grade
	studentJson["student_current_address"] = student.CurrentAddress
	studentJson["student_mobile"] = student.Mobile
	studentJson["student_birthday"] = student.Birthday
	studentJson["student_family_address"] = student.FamilyAddress
	studentJson["student_experience_time"] = student.Experience.Time
	studentJson["student_experience_location"] = student.Experience.Location
	studentJson["student_experience_teacher"] = student.Experience.Teacher
	studentJson["student_father_age"] = student.FatherAge
	studentJson["student_father_job"] = student.FatherJob
	studentJson["student_father_edu"] = student.FatherEdu
	studentJson["student_mother_age"] = student.MotherAge
	studentJson["student_mother_job"] = student.MotherJob
	studentJson["student_mother_edu"] = student.MotherEdu
	studentJson["student_parent_marriage"] = student.ParentMarriage
	studentJson["student_significant"] = student.Significant
	studentJson["student_problem"] = student.Problem
	result["student_info"] = studentJson

	reservations, err := s.w.GetReservationsByStudent(userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := make(map[string]interface{})
		resJson["reservation_id"] = res.Id.Hex()
		resJson["start_time"] = res.StartTime.Format("2006-01-02 15:04")
		resJson["end_time"] = res.EndTime.Format("2006-01-02 15:04")
		resJson["source"] = res.Source.String()
		resJson["source_id"] = res.SourceId
		if teacher, err := s.w.GetTeacherById(res.TeacherId); err == nil {
			resJson["teacher_fullname"] = teacher.Fullname
		}
		if res.Status == model.AVAILABLE {
			resJson["status"] = model.AVAILABLE.String()
		} else if res.Status == model.RESERVATED && res.StartTime.Before(time.Now()) && res.StudentId == student.Id.Hex() {
			resJson["status"] = model.FEEDBACK.String()
		} else {
			resJson["status"] = model.RESERVATED.String()
		}
		array = append(array, resJson)
	}
	result["reservations"] = array

	return result
}

func (s *Service) MakeReservationByStudent(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")
	startTime := r.PostFormValue("start_time")
	fullname := r.PostFormValue("student_fullname")
	gender := r.PostFormValue("student_gender")
	birthday := r.PostFormValue("student_birthday")
	school := r.PostFormValue("student_school")
	grade := r.PostFormValue("student_grade")
	currentAddress := r.PostFormValue("student_current_address")
	familyAddress := r.PostFormValue("student_family_address")
	mobile := r.PostFormValue("student_mobile")
	email := r.PostFormValue("student_email")
	experienceTime := r.PostFormValue("student_experience_time")
	experienceLocation := r.PostFormValue("student_experience_location")
	experienceTeacher := r.PostFormValue("student_experience_teacher")
	fatherAge := r.PostFormValue("student_father_age")
	fatherJob := r.PostFormValue("student_father_job")
	fatherEdu := r.PostFormValue("student_father_edu")
	motherAge := r.PostFormValue("student_mother_age")
	motherJob := r.PostFormValue("student_mother_job")
	motherEdu := r.PostFormValue("student_mother_edu")
	parentMarriage := r.PostFormValue("student_parent_marriage")
	siginificant := r.PostFormValue("student_significant")
	problem := r.PostFormValue("student_problem")

	var result = map[string]interface{}{"state": "SUCCESS"}

	var reservationJson = make(map[string]interface{})
	reservation, err := s.w.MakeReservationByStudent(reservationId, sourceId, startTime, fullname, gender, birthday,
		school, grade, currentAddress, familyAddress, mobile, email, experienceTime, experienceLocation, experienceTeacher,
		fatherAge, fatherJob, fatherEdu, motherAge, motherJob, motherEdu, parentMarriage, siginificant, problem,
		userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	reservationJson["start_time"] = reservation.StartTime.Format("2006-01-02 15:04")
	reservationJson["end_time"] = reservation.EndTime.Format("2006-01-02 15:04")
	if teacher, err := s.w.GetTeacherById(reservation.TeacherId); err == nil {
		reservationJson["teacher_fullname"] = teacher.Fullname
	}
	result["reservation"] = reservationJson

	return result
}

func (s *Service) GetFeedbackByStudent(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")

	var result = map[string]interface{}{"state": "SUCCESS"}

	var feedbackJson = make(map[string]interface{})
	reservation, err := s.w.GetFeedbackByStudent(reservationId, sourceId, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	feedbackJson["scores"] = reservation.StudentFeedback.Scores
	result["feedback"] = feedbackJson

	return result
}

func (s *Service) SubmitFeedbackByStudent(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")
	r.ParseForm()
	scores := []string(r.Form["scores"])

	var result = map[string]interface{}{"state": "SUCCESS"}

	scoresInt := []int{}
	for _, p := range scores {
		if pi, err := strconv.Atoi(p); err == nil {
			scoresInt = append(scoresInt, pi)
		}
	}
	_, err := s.w.SubmitFeedbackByStudent(reservationId, sourceId, scoresInt, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}

	return result
}
