package service

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (s *Service) ViewReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	var result = map[string]interface{}{"state": "SUCCESS"}

	reservations, err := s.w.GetReservationsByAdmin(userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := make(map[string]interface{})
		resJson["reservation_id"] = res.Id
		resJson["start_time"] = res.StartTime.Format("2006-01-02 15:04")
		resJson["end_time"] = res.EndTime.Format("2006-01-02 15:04")
		resJson["source"] = res.Source.String()
		resJson["source_id"] = res.SourceId
		resJson["student_id"] = res.StudentId
		if student, err := s.w.GetStudentById(res.StudentId); err == nil {
			resJson["student_crisis_level"] = student.CrisisLevel
		}
		resJson["teacher_id"] = res.TeacherId
		if teacher, err := s.w.GetTeacherById(res.TeacherId); err == nil {
			resJson["teacher_username"] = teacher.Username
			resJson["teacher_fullname"] = teacher.Fullname
			resJson["teacher_mobile"] = teacher.Mobile
		}
		if res.Status == model.AVAILABLE {
			resJson["status"] = model.AVAILABLE.String()
		} else if res.Status == model.RESERVATED && res.StartTime.Before(time.Now()) {
			resJson["status"] = model.FEEDBACK.String()
		} else {
			resJson["status"] = model.RESERVATED.String()
		}
		array = append(array, resJson)
	}
	result["reservations"] = array

	return result
}

func (s *Service) ViewDailyReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["from_date"]) == 0 {
		s.ErrorHandler(w, r, errors.New("参数错误"))
		return nil
	}
	fromDate := queryForm["from_date"][0]

	var result = map[string]interface{}{"state": "SUCCESS"}

	reservations, err := s.w.GetReservationsDailyByAdmin(fromDate, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := make(map[string]interface{})
		resJson["reservation_id"] = res.Id
		resJson["start_time"] = res.StartTime.Format("2006-01-02 15:04")
		resJson["end_time"] = res.EndTime.Format("2006-01-02 15:04")
		resJson["source"] = res.Source.String()
		resJson["source_id"] = res.SourceId
		resJson["student_id"] = res.StudentId
		if student, err := s.w.GetStudentById(res.StudentId); err == nil {
			resJson["student_crisis_level"] = student.CrisisLevel
		}
		resJson["teacher_id"] = res.TeacherId
		if teacher, err := s.w.GetTeacherById(res.TeacherId); err == nil {
			resJson["teacher_username"] = teacher.Username
			resJson["teacher_fullname"] = teacher.Fullname
			resJson["teacher_mobile"] = teacher.Mobile
		}
		if res.Status == model.AVAILABLE {
			resJson["status"] = model.AVAILABLE.String()
		} else if res.Status == model.RESERVATED && res.StartTime.Before(time.Now()) {
			resJson["status"] = model.FEEDBACK.String()
		} else {
			resJson["status"] = model.RESERVATED.String()
		}
		array = append(array, resJson)
	}
	result["reservations"] = array

	return result
}

func (s *Service) ExportTodayReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	var result = map[string]interface{}{"state": "SUCCESS"}

	url, err := s.w.ExportTodayReservationTimetableByAdmin(userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	result["url"] = url

	return result
}

func (s *Service) AddReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	startTime := r.PostFormValue("start_time")
	endTime := r.PostFormValue("end_time")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")
	force := strings.EqualFold(r.PostFormValue("force"), "FORCE")

	var result = map[string]interface{}{"state": "SUCCESS"}

	var reservationJson = make(map[string]interface{})
	reservation, err := s.w.AddReservationByAdmin(startTime, endTime, teacherUsername, teacherFullname,
		teacherMobile, force, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	reservationJson["reservation_id"] = reservation.Id
	reservationJson["start_time"] = reservation.StartTime.Format("2006-01-02 15:04")
	reservationJson["end_time"] = reservation.EndTime.Format("2006-01-02 15:04")
	reservationJson["source"] = reservation.Source.String()
	reservationJson["source_id"] = reservation.SourceId
	reservationJson["student_id"] = reservation.StudentId
	reservationJson["teacher_id"] = reservation.TeacherId
	if teacher, err := s.w.GetTeacherById(reservation.TeacherId); err == nil {
		reservationJson["teacher_username"] = teacher.Username
		reservationJson["teacher_fullname"] = teacher.Fullname
		reservationJson["teacher_mobile"] = teacher.Mobile
	}
	result["reservation"] = reservationJson

	return result
}

func (s *Service) EditReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")
	originalStartTime := r.PostFormValue("original_start_time")
	startTime := r.PostFormValue("start_time")
	endTime := r.PostFormValue("end_time")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")
	force := strings.EqualFold(r.PostFormValue("force"), "FORCE")

	var result = map[string]interface{}{"state": "SUCCESS"}

	var reservationJson = make(map[string]interface{})
	reservation, err := s.w.EditReservationByAdmin(reservationId, sourceId, originalStartTime,
		startTime, endTime, teacherUsername, teacherFullname, teacherMobile, force, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	reservationJson["reservation_id"] = reservation.Id
	reservationJson["start_time"] = reservation.StartTime.Format("2006-01-02 15:04")
	reservationJson["end_time"] = reservation.EndTime.Format("2006-01-02 15:04")
	reservationJson["source"] = reservation.Source.String()
	reservationJson["source_id"] = reservation.SourceId
	reservationJson["student_id"] = reservation.StudentId
	reservationJson["teacher_id"] = reservation.TeacherId
	if teacher, err := s.w.GetTeacherById(reservation.TeacherId); err == nil {
		reservationJson["teacher_username"] = teacher.Username
		reservationJson["teacher_fullname"] = teacher.Fullname
		reservationJson["teacher_mobile"] = teacher.Mobile
	}
	result["reservation"] = reservationJson

	return result
}

func (s *Service) RemoveReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	r.ParseForm()
	reservationIds := []string(r.Form["reservation_ids"])
	sourceIds := []string(r.Form["source_ids"])
	startTimes := []string(r.Form["start_times"])

	var result = map[string]interface{}{"state": "SUCCESS"}

	removed, err := s.w.RemoveReservationsByAdmin(reservationIds, sourceIds, startTimes, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	result["removed_count"] = removed

	return result
}

func (s *Service) CancelReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	r.ParseForm()
	reservationIds := []string(r.Form["reservation_ids"])
	sourceIds := []string(r.Form["source_ids"])

	var result = map[string]interface{}{"state": "SUCCESS"}

	removed, err := s.w.CancelReservationsByAdmin(reservationIds, sourceIds, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	result["removed_count"] = removed

	return result
}

func (s *Service) GetFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")

	var result = map[string]interface{}{"state": "SUCCESS"}

	var feedback = make(map[string]interface{})
	student, reservation, err := s.w.GetFeedbackByAdmin(reservationId, sourceId, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	feedback["category"] = reservation.TeacherFeedback.Category
	if len(reservation.TeacherFeedback.Participants) != len(model.PARTICIPANTS) {
		feedback["participants"] = make([]int, len(model.PARTICIPANTS))
	} else {
		feedback["participants"] = reservation.TeacherFeedback.Participants
	}
	feedback["emphasis"] = reservation.TeacherFeedback.Emphasis
	if len(reservation.TeacherFeedback.Severity) != len(model.SEVERITY) {
		feedback["severity"] = make([]int, len(model.SEVERITY))
	} else {
		feedback["severity"] = reservation.TeacherFeedback.Severity
	}
	if len(reservation.TeacherFeedback.MedicalDiagnosis) != len(model.MEDICAL_DIAGNOSIS) {
		feedback["medical_diagnosis"] = make([]int, len(model.MEDICAL_DIAGNOSIS))
	} else {
		feedback["medical_diagnosis"] = reservation.TeacherFeedback.MedicalDiagnosis
	}
	if len(reservation.TeacherFeedback.Crisis) != len(model.CRISIS) {
		feedback["crisis"] = make([]int, len(model.CRISIS))
	} else {
		feedback["crisis"] = reservation.TeacherFeedback.Crisis
	}
	feedback["record"] = reservation.TeacherFeedback.Record
	feedback["crisis_level"] = student.CrisisLevel
	result["feedback"] = feedback

	return result
}

func (s *Service) SubmitFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")
	category := r.PostFormValue("category")
	r.ParseForm()
	participants := []string(r.Form["participants"])
	emphasis := r.PostFormValue("emphasis")
	severity := []string(r.Form["severity"])
	medicalDiagnosis := []string(r.Form["medical_diagnosis"])
	crisis := []string(r.Form["crisis"])
	record := r.PostFormValue("record")
	crisisLevel := r.PostFormValue("crisis_level")

	var result = map[string]interface{}{"state": "SUCCESS"}

	participantsInt := make([]int, 0)
	for _, p := range participants {
		if pi, err := strconv.Atoi(p); err == nil {
			participantsInt = append(participantsInt, pi)
		}
	}
	severityInt := make([]int, 0)
	for _, s := range severity {
		if si, err := strconv.Atoi(s); err == nil {
			severityInt = append(severityInt, si)
		}
	}
	medicalDiagnosisInt := make([]int, 0)
	for _, m := range medicalDiagnosis {
		if mi, err := strconv.Atoi(m); err == nil {
			medicalDiagnosisInt = append(medicalDiagnosisInt, mi)
		}
	}
	crisisInt := make([]int, 0)
	for _, c := range crisis {
		if ci, err := strconv.Atoi(c); err == nil {
			crisisInt = append(crisisInt, ci)
		}
	}
	_, err := s.w.SubmitFeedbackByAdmin(reservationId, sourceId, category, participantsInt, emphasis, severityInt,
		medicalDiagnosisInt, crisisInt, record, crisisLevel, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}

	return result
}

func (s *Service) SetStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")
	startTime := r.PostFormValue("start_time")
	studentUsername := r.PostFormValue("student_username")
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
	sendSms, err := strconv.ParseBool(r.PostFormValue("student_sms"))
	if err != nil {
		s.ErrorHandler(w, r, errors.New("参数错误，请联系管理员"))
		return nil
	}

	var result = map[string]interface{}{"state": "SUCCESS"}

	var reservationJson = make(map[string]interface{})
	reservation, err := s.w.SetStudentByAdmin(reservationId, sourceId, startTime, studentUsername, fullname, gender,
		birthday, school, grade, currentAddress, familyAddress, mobile, email, experienceTime,
		experienceLocation, experienceTeacher, fatherAge, fatherJob, fatherEdu, motherAge, motherJob, motherEdu,
		parentMarriage, siginificant, problem, sendSms, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	reservationJson["reservation_id"] = reservation.Id
	reservationJson["start_time"] = reservation.StartTime.Format("2006-01-02 15:04")
	reservationJson["end_time"] = reservation.EndTime.Format("2006-01-02 15:04")
	reservationJson["source"] = reservation.Source.String()
	reservationJson["source_id"] = reservation.SourceId
	reservationJson["student_id"] = reservation.StudentId
	reservationJson["teacher_id"] = reservation.TeacherId
	if teacher, err := s.w.GetTeacherById(reservation.TeacherId); err == nil {
		reservationJson["teacher_username"] = teacher.Username
		reservationJson["teacher_fullname"] = teacher.Fullname
		reservationJson["teacher_mobile"] = teacher.Mobile
	}
	result["reservation"] = reservationJson

	return result
}

func (s *Service) GetStudentInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	studentId := r.PostFormValue("student_id")

	var result = map[string]interface{}{"state": "SUCCESS"}

	var studentJson = make(map[string]interface{})
	student, reservations, err := s.w.GetStudentInfoByAdmin(studentId, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	studentJson["student_id"] = student.Id.Hex()
	studentJson["student_username"] = student.Username
	studentJson["student_fullname"] = student.Fullname
	studentJson["student_archive_category"] = student.ArchiveCategory
	studentJson["student_archive_number"] = student.ArchiveNumber
	studentJson["student_crisis_level"] = student.CrisisLevel
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
	if len(student.BindedTeacherId) != 0 {
		teacher, err := s.w.GetTeacherById(student.BindedTeacherId)
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

	var reservationJson = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := make(map[string]interface{})
		resJson["start_time"] = res.StartTime.Format("2006-01-02 15:04")
		resJson["end_time"] = res.EndTime.Format("2006-01-02 15:04")
		if res.Status == model.AVAILABLE {
			resJson["status"] = model.AVAILABLE.String()
		} else if res.Status == model.RESERVATED && res.StartTime.Before(time.Now()) {
			resJson["status"] = model.FEEDBACK.String()
		} else {
			resJson["status"] = model.RESERVATED.String()
		}
		resJson["student_id"] = res.StudentId
		resJson["teacher_id"] = res.TeacherId
		if teacher, err := s.w.GetTeacherById(res.TeacherId); err == nil {
			resJson["teacher_username"] = teacher.Username
			resJson["teacher_fullname"] = teacher.Fullname
			resJson["teacher_mobile"] = teacher.Mobile
		}
		resJson["student_feedback"] = res.StudentFeedback.ToJson()
		resJson["teacher_feedback"] = res.TeacherFeedback.ToJson()
		reservationJson = append(reservationJson, resJson)
	}
	result["reservations"] = reservationJson

	return result
}

func (s *Service) SearchStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	studentUsername := r.PostFormValue("student_username")

	var result = map[string]interface{}{"state": "SUCCESS"}

	var studentJson = make(map[string]interface{})
	student, err := s.w.GetStudentByUsername(studentUsername)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	studentJson["student_id"] = student.Id.Hex()
	studentJson["student_username"] = student.Username
	studentJson["student_fullname"] = student.Fullname
	studentJson["student_archive_category"] = student.ArchiveCategory
	studentJson["student_archive_number"] = student.ArchiveNumber
	studentJson["student_crisis_level"] = student.CrisisLevel
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
	if len(student.BindedTeacherId) != 0 {
		teacher, err := s.w.GetTeacherById(student.BindedTeacherId)
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

func (s *Service) UpdateStudentCrisisLevelByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	studentId := r.PostFormValue("student_id")
	crisisLevel := r.PostFormValue("crisis_level")

	var result = map[string]interface{}{"state": "SUCCESS"}

	_, err := s.w.UpdateStudentCrisisLevelByAdmin(studentId, crisisLevel, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}

	return result
}

func (s *Service) UpdateStudentArchiveNumberByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	studentId := r.PostFormValue("student_id")
	archiveCategory := r.PostFormValue("archive_category")
	archiveNumber := r.PostFormValue("archive_number")

	var result = map[string]interface{}{"state": "SUCCESS"}

	_, err := s.w.UpdateStudentArchiveNumberByAdmin(studentId, archiveCategory, archiveNumber, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}

	return result
}

func (s *Service) ResetStudentPasswordByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	studentId := r.PostFormValue("student_id")
	password := r.PostFormValue("password")

	var result = map[string]interface{}{"state": "SUCCESS"}

	_, err := s.w.ResetStudentPasswordByAdmin(studentId, password, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}

	return result
}

func (s *Service) DeleteStudentAccountByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	studentId := r.PostFormValue("student_id")

	var result = map[string]interface{}{"state": "SUCCESS"}

	err := s.w.DeleteStudentAccountByAdmin(studentId, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}

	return result
}

func (s *Service) ExportStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	studentId := r.PostFormValue("student_id")

	var result = map[string]interface{}{"state": "SUCCESS"}

	url, err := s.w.ExportStudentByAdmin(studentId, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	result["url"] = url

	return result
}

func (s *Service) UnbindStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	studentId := r.PostFormValue("student_id")

	var result = map[string]interface{}{"state": "SUCCESS"}

	var studentJson = make(map[string]interface{})
	student, err := s.w.UnbindStudentByAdmin(studentId, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	if len(student.BindedTeacherId) != 0 {
		teacher, err := s.w.GetTeacherById(student.BindedTeacherId)
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

func (s *Service) BindStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	studentId := r.PostFormValue("student_id")
	teacherUsername := r.PostFormValue("teacher_username")

	var result = map[string]interface{}{"state": "SUCCESS"}

	var studentJson = make(map[string]interface{})
	student, err := s.w.BindStudentByAdmin(studentId, teacherUsername, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	} else if len(student.BindedTeacherId) != 0 {
		teacher, err := s.w.GetTeacherById(student.BindedTeacherId)
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

func (s *Service) QueryStudentInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	studentUsername := r.PostFormValue("student_username")

	var result = map[string]interface{}{"state": "SUCCESS"}

	var studentJson = make(map[string]interface{})
	student, reservations, err := s.w.QueryStudentInfoByAdmin(studentUsername, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	studentJson["student_id"] = student.Id.Hex()
	studentJson["student_username"] = student.Username
	studentJson["student_fullname"] = student.Fullname
	studentJson["student_archive_category"] = student.ArchiveCategory
	studentJson["student_archive_number"] = student.ArchiveNumber
	studentJson["student_crisis_level"] = student.CrisisLevel
	studentJson["student_key_case"] = student.KeyCase
	studentJson["student_medical_diagnosis"] = student.MedicalDiagnosis
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
	if len(student.BindedTeacherId) != 0 {
		teacher, err := s.w.GetTeacherById(student.BindedTeacherId)
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

	var reservationJson = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := make(map[string]interface{})
		resJson["start_time"] = res.StartTime.Format("2006-01-02 15:04")
		resJson["end_time"] = res.EndTime.Format("2006-01-02 15:04")
		if res.Status == model.AVAILABLE {
			resJson["status"] = model.AVAILABLE.String()
		} else if res.Status == model.RESERVATED && res.StartTime.Before(time.Now()) {
			resJson["status"] = model.FEEDBACK.String()
		} else {
			resJson["status"] = model.RESERVATED.String()
		}
		resJson["student_id"] = res.StudentId
		resJson["teacher_id"] = res.TeacherId
		if teacher, err := s.w.GetTeacherById(res.TeacherId); err == nil {
			resJson["teacher_username"] = teacher.Username
			resJson["teacher_fullname"] = teacher.Fullname
			resJson["teacher_mobile"] = teacher.Mobile
		}
		resJson["student_feedback"] = res.StudentFeedback.ToJson()
		resJson["teacher_feedback"] = res.TeacherFeedback.ToJson()
		reservationJson = append(reservationJson, resJson)
	}
	result["reservations"] = reservationJson

	return result
}

func (s *Service) SearchTeacherByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMoble := r.PostFormValue("teacher_mobile")

	var result = map[string]interface{}{"state": "SUCCESS"}

	var teacherJson = make(map[string]interface{})
	teacher, err := s.w.SearchTeacherByAdmin(teacherFullname, teacherUsername, teacherMoble,
		userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	teacherJson["teacher_id"] = teacher.Id.Hex()
	teacherJson["teacher_username"] = teacher.Username
	teacherJson["teacher_fullname"] = teacher.Fullname
	teacherJson["teacher_mobile"] = teacher.Mobile
	result["teacher"] = teacherJson

	return result
}

func (s *Service) GetTeacherWorkloadByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	fromDate := r.PostFormValue("from_date")
	toDate := r.PostFormValue("to_date")

	var result = map[string]interface{}{"state": "SUCCESS"}

	workload, err := s.w.GetTeacherWorkloadByAdmin(fromDate, toDate, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	result["workload"] = workload

	return result
}

func (s *Service) ExportReportFormByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	fromDate := r.PostFormValue("from_date")
	toDate := r.PostFormValue("to_date")

	var result = map[string]interface{}{"state": "SUCCESS"}

	url, err := s.w.ExportReportFormByAdmin(fromDate, toDate, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	result["url"] = url

	return result
}

func (s *Service) ExportReportMonthlyByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	monthlyDate := r.PostFormValue("monthly_date")

	var result = map[string]interface{}{"state": "SUCCESS"}

	reportUrl, keyCaseUrl, err := s.w.ExportReportMonthlyByAdmin(monthlyDate, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	result["report"] = reportUrl
	result["key_case"] = keyCaseUrl

	return result
}
