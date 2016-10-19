package web

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/ifs"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/service"
	"github.com/mijia/sweb/render"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ReservationController struct {
	ifs.BaseMux
}

const (
	kCategoryApiBaseUrl = "/api/category"
	kStudentApiBaseUrl  = "/api/student"
	kTeacherApiBaseUrl  = "/api/teacher"
	kAdminApiBaseUrl    = "/api/admin"
)

func (rc *ReservationController) MuxHandlers(s ifs.Muxer) {
	categoryBaseUrl := kCategoryApiBaseUrl
	s.GetJson(categoryBaseUrl+"/feedback", "GetFeedbackCategories", rc.GetFeedbackCategories)

	studentBaseUrl := kStudentApiBaseUrl
	s.GetJson(studentBaseUrl+"/reservation/view", "ViewReservationsByStudent", RoleCookieInjection(rc.ViewReservationsByStudent))
	s.PostJson(studentBaseUrl+"/reservation/make", "MakeReservationByStudent", RoleCookieInjection(rc.MakeReservationByStudent))
	s.PostJson(studentBaseUrl+"/reservation/feedback/get", "GetFeedbackByStudent", RoleCookieInjection(rc.GetFeedbackByStudent))
	s.PostJson(studentBaseUrl+"/reservation/feedback/submit", "SubmitFeedbackByStudent", RoleCookieInjection(rc.SubmitFeedbackByStudent))

	teacherBaseUrl := kTeacherApiBaseUrl
	s.GetJson(teacherBaseUrl+"/reservation/view", "ViewReservationsByTeacher", RoleCookieInjection(rc.ViewReservationsByTeacher))
	s.PostJson(teacherBaseUrl+"/reservation/feedback/get", "GetFeedbackByTeacher", RoleCookieInjection(rc.GetFeedbackByTeacher))
	s.PostJson(teacherBaseUrl+"/reservation/feedback/submit", "SubmitFeedbackByTeacher", RoleCookieInjection(rc.SubmitFeedbackByTeacher))
	s.PostJson(teacherBaseUrl+"/student/get", "GetStudentInfoByTeacher", RoleCookieInjection(rc.GetStudentInfoByTeacher))
	s.PostJson(teacherBaseUrl+"/student/query", "QueryStudentInfoByTeacher", RoleCookieInjection(rc.QueryStudentInfoByTeacher))

	adminBaseUrl := kAdminApiBaseUrl
	s.GetJson(adminBaseUrl+"/timetable/view", "ViewTimedReservationsByAdmin", RoleCookieInjection(rc.ViewTimedReservationsByAdmin))
	s.PostJson(adminBaseUrl+"/timetable/add", "AddTimedReservationByAdmin", RoleCookieInjection(rc.AddTimedReservationByAdmin))
	s.PostJson(adminBaseUrl+"/timetable/edit", "EditTimedReservationByAdmin", RoleCookieInjection(rc.EditTimedReservationByAdmin))
	s.PostJson(adminBaseUrl+"/timetable/remove", "RemoveTimedReservationsByAdmin", RoleCookieInjection(rc.RemoveTimedReservationsByAdmin))
	s.PostJson(adminBaseUrl+"/timetable/open", "OpenTimedReservationsByAdmin", RoleCookieInjection(rc.OpenTimedReservationsByAdmin))
	s.PostJson(adminBaseUrl+"/timetable/close", "CloseTimedReservationsByAdmin", RoleCookieInjection(rc.CloseTimedReservationsByAdmin))
	s.GetJson(adminBaseUrl+"/reservation/view", "ViewReservationsByAdmin", RoleCookieInjection(rc.ViewReservationsByAdmin))
	s.GetJson(adminBaseUrl+"/reservation/view/daily", "ViewDailyReservationsByAdmin", RoleCookieInjection(rc.ViewDailyReservationsByAdmin))
	s.GetJson(adminBaseUrl+"/reservation/export/today", "ExportTodayReservationsByAdmin", RoleCookieInjection(rc.ExportTodayReservationsByAdmin))
	s.PostJson(adminBaseUrl+"/reservation/export/report", "ExportReportFormByAdmin", RoleCookieInjection(rc.ExportReportFormByAdmin))
	s.PostJson(adminBaseUrl+"/reservation/export/report/monthly", "ExportReportMonthlyByAdmin", RoleCookieInjection(rc.ExportReportMonthlyByAdmin))
	s.PostJson(adminBaseUrl+"/reservation/add", "AddReservationByAdmin", RoleCookieInjection(rc.AddReservationByAdmin))
	s.PostJson(adminBaseUrl+"/reservation/edit", "EditReservationByAdmin", RoleCookieInjection(rc.EditReservationByAdmin))
	s.PostJson(adminBaseUrl+"/reservation/remove", "RemoveReservationsByAdmin", RoleCookieInjection(rc.RemoveReservationsByAdmin))
	s.PostJson(adminBaseUrl+"/reservation/cancel", "CancelReservationByAdmin", RoleCookieInjection(rc.CancelReservationByAdmin))
	s.PostJson(adminBaseUrl+"/reservation/feedback/get", "GetFeedbackByAdmin", RoleCookieInjection(rc.GetFeedbackByAdmin))
	s.PostJson(adminBaseUrl+"/reservation/feedback/submit", "SubmitFeedbackByAdmin", RoleCookieInjection(rc.SubmitFeedbackByAdmin))
	s.PostJson(adminBaseUrl+"/reservation/student/set", "SetStudentByAdmin", RoleCookieInjection(rc.SetStudentByAdmin))
	s.PostJson(adminBaseUrl+"/student/get", "GetStudentInfoByAdmin", RoleCookieInjection(rc.GetStudentInfoByAdmin))
	s.PostJson(adminBaseUrl+"/student/search", "SearchStudentByAdmin", RoleCookieInjection(rc.SearchStudentByAdmin))
	s.PostJson(adminBaseUrl+"/student/crisis/update", "UpdateStudentCrisisLevelByAdmin", RoleCookieInjection(rc.UpdateStudentCrisisLevelByAdmin))
	s.PostJson(adminBaseUrl+"/student/archive/update", "UpdateStudentArchiveNumberByAdmin", RoleCookieInjection(rc.UpdateStudentArchiveNumberByAdmin))
	s.PostJson(adminBaseUrl+"/student/password/reset", "ResetStudentPasswordByAdmin", RoleCookieInjection(rc.ResetStudentPasswordByAdmin))
	s.PostJson(adminBaseUrl+"/student/account/delete", "DeleteStudentAccountByAdmin", RoleCookieInjection(rc.DeleteStudentAccountByAdmin))
	s.PostJson(adminBaseUrl+"/student/export", "ExportStudentByAdmin", RoleCookieInjection(rc.ExportStudentByAdmin))
	s.PostJson(adminBaseUrl+"/student/unbind", "UnbindStudentByAdmin", RoleCookieInjection(rc.UnbindStudentByAdmin))
	s.PostJson(adminBaseUrl+"/student/bind", "BindStudentByAdmin", RoleCookieInjection(rc.BindStudentByAdmin))
	s.PostJson(adminBaseUrl+"/student/query", "QueryStudentInfoByAdmin", RoleCookieInjection(rc.QueryStudentInfoByAdmin))
	s.PostJson(adminBaseUrl+"/teacher/search", "SearchTeacherByAdmin", RoleCookieInjection(rc.SearchTeacherByAdmin))
	s.PostJson(adminBaseUrl+"/teacher/workload", "GetTeacherWorkloadByAdmin", RoleCookieInjection(rc.GetTeacherWorkloadByAdmin))
}

func (rc *ReservationController) GetTemplates() []*render.TemplateSet {
	return []*render.TemplateSet{}
}

//==================== category ====================
func (rc *ReservationController) GetFeedbackCategories(ctx context.Context, w http.ResponseWriter, r *http.Request) interface{} {
	var result = map[string]interface{}{"status": "SUCCESS"}

	result["first_category"] = model.FeedbackFirstCategory
	result["second_category"] = model.FeedbackSecondCategory

	return result
}

//==================== student ====================
func (rc *ReservationController) ViewReservationsByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	var result = map[string]interface{}{"status": "SUCCESS"}

	student, err := service.Workflow().GetStudentById(userId)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
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

	reservations, err := service.Workflow().GetReservationsByStudent(userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := make(map[string]interface{})
		resJson["reservation_id"] = res.Id.Hex()
		resJson["start_time"] = res.StartTime.Format("2006-01-02 15:04")
		resJson["end_time"] = res.EndTime.Format("2006-01-02 15:04")
		resJson["source"] = res.Source
		resJson["source_id"] = res.SourceId
		if teacher, err := service.Workflow().GetTeacherById(res.TeacherId); err == nil {
			resJson["teacher_fullname"] = teacher.Fullname
		}
		if res.Status == model.RESERVATION_STATUS_AVAILABLE {
			resJson["status"] = model.RESERVATION_STATUS_AVAILABLE
		} else if res.Status == model.RESERVATION_STATUS_RESERVATED && res.StartTime.Before(time.Now()) && res.StudentId == student.Id.Hex() {
			resJson["status"] = model.RESERVATION_STATUS_FEEDBACK
		} else {
			resJson["status"] = model.RESERVATION_STATUS_RESERVATED
		}
		array = append(array, resJson)
	}
	result["reservations"] = array

	return result
}

func (rc *ReservationController) MakeReservationByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
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

	var result = map[string]interface{}{"status": "SUCCESS"}

	var reservationJson = make(map[string]interface{})
	reservation, err := service.Workflow().MakeReservationByStudent(reservationId, sourceId, startTime, fullname, gender, birthday,
		school, grade, currentAddress, familyAddress, mobile, email, experienceTime, experienceLocation, experienceTeacher,
		fatherAge, fatherJob, fatherEdu, motherAge, motherJob, motherEdu, parentMarriage, siginificant, problem,
		userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	reservationJson["start_time"] = reservation.StartTime.Format("2006-01-02 15:04")
	reservationJson["end_time"] = reservation.EndTime.Format("2006-01-02 15:04")
	if teacher, err := service.Workflow().GetTeacherById(reservation.TeacherId); err == nil {
		reservationJson["teacher_fullname"] = teacher.Fullname
	}
	result["reservation"] = reservationJson

	return result
}

func (rc *ReservationController) GetFeedbackByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")

	var result = map[string]interface{}{"status": "SUCCESS"}

	var feedbackJson = make(map[string]interface{})
	reservation, err := service.Workflow().GetFeedbackByStudent(reservationId, sourceId, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	feedbackJson["scores"] = reservation.StudentFeedback.Scores
	result["feedback"] = feedbackJson

	return result
}

func (rc *ReservationController) SubmitFeedbackByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")
	r.ParseForm()
	scores := []string(r.Form["scores"])

	var result = map[string]interface{}{"status": "SUCCESS"}

	scoresInt := []int{}
	for _, p := range scores {
		if pi, err := strconv.Atoi(p); err == nil {
			scoresInt = append(scoresInt, pi)
		}
	}
	_, err := service.Workflow().SubmitFeedbackByStudent(reservationId, sourceId, scoresInt, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}

	return result
}

//==================== teacher ====================
func (rc *ReservationController) ViewReservationsByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	var result = map[string]interface{}{"status": "SUCCESS"}

	teacher, err := service.Workflow().GetTeacherById(userId)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	var teacherJson = make(map[string]interface{})
	teacherJson["teacher_fullname"] = teacher.Fullname
	teacherJson["teacher_mobile"] = teacher.Mobile
	result["teacher_info"] = teacherJson

	reservations, err := service.Workflow().GetReservationsByTeacher(userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := make(map[string]interface{})
		resJson["reservation_id"] = res.Id.Hex()
		resJson["start_time"] = res.StartTime.Format("2006-01-02 15:04")
		resJson["end_time"] = res.EndTime.Format("2006-01-02 15:04")
		resJson["source"] = res.Source
		resJson["source_id"] = res.SourceId
		resJson["student_id"] = res.StudentId
		if student, err := service.Workflow().GetStudentById(res.StudentId); err == nil {
			resJson["student_crisis_level"] = student.CrisisLevel
		}
		resJson["teacher_id"] = res.TeacherId
		if teacher, err := service.Workflow().GetTeacherById(res.TeacherId); err == nil {
			resJson["teacher_fullname"] = teacher.Fullname
			resJson["teacher_mobile"] = teacher.Mobile
		}
		if res.Status == model.RESERVATION_STATUS_AVAILABLE {
			resJson["status"] = model.RESERVATION_STATUS_AVAILABLE
		} else if res.Status == model.RESERVATION_STATUS_RESERVATED && res.StartTime.Before(time.Now()) {
			resJson["status"] = model.RESERVATION_STATUS_FEEDBACK
		} else {
			resJson["status"] = model.RESERVATION_STATUS_RESERVATED
		}
		array = append(array, resJson)
	}
	result["reservations"] = array

	return result
}

func (rc *ReservationController) GetFeedbackByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")

	var result = map[string]interface{}{"status": "SUCCESS"}

	var feedback = make(map[string]interface{})
	student, reservation, err := service.Workflow().GetFeedbackByTeacher(reservationId, sourceId, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
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

func (rc *ReservationController) SubmitFeedbackByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
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

	var result = map[string]interface{}{"status": "SUCCESS"}

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
	_, err := service.Workflow().SubmitFeedbackByTeacher(reservationId, sourceId, category, participantsInt, emphasis, severityInt,
		medicalDiagnosisInt, crisisInt, record, crisisLevel, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}

	return result
}

func (rc *ReservationController) GetStudentInfoByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	studentId := r.PostFormValue("student_id")

	var result = map[string]interface{}{"status": "SUCCESS"}

	var studentJson = make(map[string]interface{})
	student, reservations, err := service.Workflow().GetStudentInfoByTeacher(studentId, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
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
		teacher, err := service.Workflow().GetTeacherById(student.BindedTeacherId)
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
		if res.Status == model.RESERVATION_STATUS_AVAILABLE {
			resJson["status"] = model.RESERVATION_STATUS_AVAILABLE
		} else if res.Status == model.RESERVATION_STATUS_RESERVATED && res.StartTime.Before(time.Now()) {
			resJson["status"] = model.RESERVATION_STATUS_FEEDBACK
		} else {
			resJson["status"] = model.RESERVATION_STATUS_RESERVATED
		}
		resJson["student_id"] = res.StudentId
		resJson["teacher_id"] = res.TeacherId
		if teacher, err := service.Workflow().GetTeacherById(res.TeacherId); err == nil {
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

func (rc *ReservationController) QueryStudentInfoByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	studentUsername := r.PostFormValue("student_username")

	var result = map[string]interface{}{"status": "SUCCESS"}

	var studentJson = make(map[string]interface{})
	student, reservations, err := service.Workflow().QueryStudentInfoByTeacher(studentUsername, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
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
		teacher, err := service.Workflow().GetTeacherById(student.BindedTeacherId)
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
		if res.Status == model.RESERVATION_STATUS_AVAILABLE {
			resJson["status"] = model.RESERVATION_STATUS_AVAILABLE
		} else if res.Status == model.RESERVATION_STATUS_RESERVATED && res.StartTime.Before(time.Now()) {
			resJson["status"] = model.RESERVATION_STATUS_FEEDBACK
		} else {
			resJson["status"] = model.RESERVATION_STATUS_RESERVATED
		}
		resJson["student_id"] = res.StudentId
		resJson["teacher_id"] = res.TeacherId
		if teacher, err := service.Workflow().GetTeacherById(res.TeacherId); err == nil {
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

//==================== admin ====================
func (rc *ReservationController) ViewReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	var result = map[string]interface{}{"status": "SUCCESS"}

	reservations, err := service.Workflow().GetReservationsByAdmin(userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := make(map[string]interface{})
		resJson["reservation_id"] = res.Id
		resJson["start_time"] = res.StartTime.Format("2006-01-02 15:04")
		resJson["end_time"] = res.EndTime.Format("2006-01-02 15:04")
		resJson["source"] = res.Source
		resJson["source_id"] = res.SourceId
		resJson["student_id"] = res.StudentId
		if student, err := service.Workflow().GetStudentById(res.StudentId); err == nil {
			resJson["student_crisis_level"] = student.CrisisLevel
		}
		resJson["teacher_id"] = res.TeacherId
		if teacher, err := service.Workflow().GetTeacherById(res.TeacherId); err == nil {
			resJson["teacher_username"] = teacher.Username
			resJson["teacher_fullname"] = teacher.Fullname
			resJson["teacher_mobile"] = teacher.Mobile
		}
		if res.Status == model.RESERVATION_STATUS_AVAILABLE {
			resJson["status"] = model.RESERVATION_STATUS_AVAILABLE
		} else if res.Status == model.RESERVATION_STATUS_RESERVATED && res.StartTime.Before(time.Now()) {
			resJson["status"] = model.RESERVATION_STATUS_FEEDBACK
		} else {
			resJson["status"] = model.RESERVATION_STATUS_RESERVATED
		}
		array = append(array, resJson)
	}
	result["reservations"] = array

	return result
}

func (rc *ReservationController) ViewDailyReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["from_date"]) == 0 {
		return map[string]string{
			"status": "FAIL",
			"error":  "参数错误",
		}
		return nil
	}
	fromDate := queryForm["from_date"][0]

	var result = map[string]interface{}{"status": "SUCCESS"}

	reservations, err := service.Workflow().GetReservationsDailyByAdmin(fromDate, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := make(map[string]interface{})
		resJson["reservation_id"] = res.Id
		resJson["start_time"] = res.StartTime.Format("2006-01-02 15:04")
		resJson["end_time"] = res.EndTime.Format("2006-01-02 15:04")
		resJson["source"] = res.Source
		resJson["source_id"] = res.SourceId
		resJson["student_id"] = res.StudentId
		if student, err := service.Workflow().GetStudentById(res.StudentId); err == nil {
			resJson["student_crisis_level"] = student.CrisisLevel
		}
		resJson["teacher_id"] = res.TeacherId
		if teacher, err := service.Workflow().GetTeacherById(res.TeacherId); err == nil {
			resJson["teacher_username"] = teacher.Username
			resJson["teacher_fullname"] = teacher.Fullname
			resJson["teacher_mobile"] = teacher.Mobile
		}
		if res.Status == model.RESERVATION_STATUS_AVAILABLE {
			resJson["status"] = model.RESERVATION_STATUS_AVAILABLE
		} else if res.Status == model.RESERVATION_STATUS_RESERVATED && res.StartTime.Before(time.Now()) {
			resJson["status"] = model.RESERVATION_STATUS_FEEDBACK
		} else {
			resJson["status"] = model.RESERVATION_STATUS_RESERVATED
		}
		array = append(array, resJson)
	}
	result["reservations"] = array

	return result
}

func (rc *ReservationController) ExportTodayReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	var result = map[string]interface{}{"status": "SUCCESS"}

	url, err := service.Workflow().ExportTodayReservationTimetableByAdmin(userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	result["url"] = url

	return result
}

func (rc *ReservationController) AddReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	startTime := r.PostFormValue("start_time")
	endTime := r.PostFormValue("end_time")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")
	force := strings.EqualFold(r.PostFormValue("force"), "FORCE")

	var result = map[string]interface{}{"status": "SUCCESS"}

	var reservationJson = make(map[string]interface{})
	reservation, err := service.Workflow().AddReservationByAdmin(startTime, endTime, teacherUsername, teacherFullname,
		teacherMobile, force, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	reservationJson["reservation_id"] = reservation.Id
	reservationJson["start_time"] = reservation.StartTime.Format("2006-01-02 15:04")
	reservationJson["end_time"] = reservation.EndTime.Format("2006-01-02 15:04")
	reservationJson["source"] = reservation.Source
	reservationJson["source_id"] = reservation.SourceId
	reservationJson["student_id"] = reservation.StudentId
	reservationJson["teacher_id"] = reservation.TeacherId
	if teacher, err := service.Workflow().GetTeacherById(reservation.TeacherId); err == nil {
		reservationJson["teacher_username"] = teacher.Username
		reservationJson["teacher_fullname"] = teacher.Fullname
		reservationJson["teacher_mobile"] = teacher.Mobile
	}
	result["reservation"] = reservationJson

	return result
}

func (rc *ReservationController) EditReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")
	originalStartTime := r.PostFormValue("original_start_time")
	startTime := r.PostFormValue("start_time")
	endTime := r.PostFormValue("end_time")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")
	force := strings.EqualFold(r.PostFormValue("force"), "FORCE")

	var result = map[string]interface{}{"status": "SUCCESS"}

	var reservationJson = make(map[string]interface{})
	reservation, err := service.Workflow().EditReservationByAdmin(reservationId, sourceId, originalStartTime,
		startTime, endTime, teacherUsername, teacherFullname, teacherMobile, force, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	reservationJson["reservation_id"] = reservation.Id
	reservationJson["start_time"] = reservation.StartTime.Format("2006-01-02 15:04")
	reservationJson["end_time"] = reservation.EndTime.Format("2006-01-02 15:04")
	reservationJson["source"] = reservation.Source
	reservationJson["source_id"] = reservation.SourceId
	reservationJson["student_id"] = reservation.StudentId
	reservationJson["teacher_id"] = reservation.TeacherId
	if teacher, err := service.Workflow().GetTeacherById(reservation.TeacherId); err == nil {
		reservationJson["teacher_username"] = teacher.Username
		reservationJson["teacher_fullname"] = teacher.Fullname
		reservationJson["teacher_mobile"] = teacher.Mobile
	}
	result["reservation"] = reservationJson

	return result
}

func (rc *ReservationController) RemoveReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	r.ParseForm()
	reservationIds := []string(r.Form["reservation_ids"])
	sourceIds := []string(r.Form["source_ids"])
	startTimes := []string(r.Form["start_times"])

	var result = map[string]interface{}{"status": "SUCCESS"}

	removed, err := service.Workflow().RemoveReservationsByAdmin(reservationIds, sourceIds, startTimes, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	result["removed_count"] = removed

	return result
}

func (rc *ReservationController) CancelReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	r.ParseForm()
	reservationIds := []string(r.Form["reservation_ids"])
	sourceIds := []string(r.Form["source_ids"])

	var result = map[string]interface{}{"status": "SUCCESS"}

	removed, err := service.Workflow().CancelReservationsByAdmin(reservationIds, sourceIds, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	result["removed_count"] = removed

	return result
}

func (rc *ReservationController) GetFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")

	var result = map[string]interface{}{"status": "SUCCESS"}

	var feedback = make(map[string]interface{})
	student, reservation, err := service.Workflow().GetFeedbackByAdmin(reservationId, sourceId, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
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

func (rc *ReservationController) SubmitFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
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

	var result = map[string]interface{}{"status": "SUCCESS"}

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
	_, err := service.Workflow().SubmitFeedbackByAdmin(reservationId, sourceId, category, participantsInt, emphasis, severityInt,
		medicalDiagnosisInt, crisisInt, record, crisisLevel, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}

	return result
}

func (rc *ReservationController) SetStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
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
		return map[string]string{
			"status": "FAIL",
			"error":  "参数错误，请联系管理员",
		}
	}

	var result = map[string]interface{}{"status": "SUCCESS"}

	var reservationJson = make(map[string]interface{})
	reservation, err := service.Workflow().SetStudentByAdmin(reservationId, sourceId, startTime, studentUsername, fullname, gender,
		birthday, school, grade, currentAddress, familyAddress, mobile, email, experienceTime,
		experienceLocation, experienceTeacher, fatherAge, fatherJob, fatherEdu, motherAge, motherJob, motherEdu,
		parentMarriage, siginificant, problem, sendSms, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	reservationJson["reservation_id"] = reservation.Id
	reservationJson["start_time"] = reservation.StartTime.Format("2006-01-02 15:04")
	reservationJson["end_time"] = reservation.EndTime.Format("2006-01-02 15:04")
	reservationJson["source"] = reservation.Source
	reservationJson["source_id"] = reservation.SourceId
	reservationJson["student_id"] = reservation.StudentId
	reservationJson["teacher_id"] = reservation.TeacherId
	if teacher, err := service.Workflow().GetTeacherById(reservation.TeacherId); err == nil {
		reservationJson["teacher_username"] = teacher.Username
		reservationJson["teacher_fullname"] = teacher.Fullname
		reservationJson["teacher_mobile"] = teacher.Mobile
	}
	result["reservation"] = reservationJson

	return result
}

func (rc *ReservationController) GetStudentInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	studentId := r.PostFormValue("student_id")

	var result = map[string]interface{}{"status": "SUCCESS"}

	var studentJson = make(map[string]interface{})
	student, reservations, err := service.Workflow().GetStudentInfoByAdmin(studentId, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
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
		teacher, err := service.Workflow().GetTeacherById(student.BindedTeacherId)
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
		if res.Status == model.RESERVATION_STATUS_AVAILABLE {
			resJson["status"] = model.RESERVATION_STATUS_AVAILABLE
		} else if res.Status == model.RESERVATION_STATUS_RESERVATED && res.StartTime.Before(time.Now()) {
			resJson["status"] = model.RESERVATION_STATUS_FEEDBACK
		} else {
			resJson["status"] = model.RESERVATION_STATUS_RESERVATED
		}
		resJson["student_id"] = res.StudentId
		resJson["teacher_id"] = res.TeacherId
		if teacher, err := service.Workflow().GetTeacherById(res.TeacherId); err == nil {
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

func (rc *ReservationController) SearchStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	studentUsername := r.PostFormValue("student_username")

	var result = map[string]interface{}{"status": "SUCCESS"}

	var studentJson = make(map[string]interface{})
	student, err := service.Workflow().GetStudentByUsername(studentUsername)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
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
		teacher, err := service.Workflow().GetTeacherById(student.BindedTeacherId)
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

func (rc *ReservationController) UpdateStudentCrisisLevelByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	studentId := r.PostFormValue("student_id")
	crisisLevel := r.PostFormValue("crisis_level")

	var result = map[string]interface{}{"status": "SUCCESS"}

	_, err := service.Workflow().UpdateStudentCrisisLevelByAdmin(studentId, crisisLevel, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}

	return result
}

func (rc *ReservationController) UpdateStudentArchiveNumberByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	studentId := r.PostFormValue("student_id")
	archiveCategory := r.PostFormValue("archive_category")
	archiveNumber := r.PostFormValue("archive_number")

	var result = map[string]interface{}{"status": "SUCCESS"}

	_, err := service.Workflow().UpdateStudentArchiveNumberByAdmin(studentId, archiveCategory, archiveNumber, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}

	return result
}

func (rc *ReservationController) ResetStudentPasswordByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	studentId := r.PostFormValue("student_id")
	password := r.PostFormValue("password")

	var result = map[string]interface{}{"status": "SUCCESS"}

	_, err := service.Workflow().ResetStudentPasswordByAdmin(studentId, password, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}

	return result
}

func (rc *ReservationController) DeleteStudentAccountByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	studentId := r.PostFormValue("student_id")

	var result = map[string]interface{}{"status": "SUCCESS"}

	err := service.Workflow().DeleteStudentAccountByAdmin(studentId, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}

	return result
}

func (rc *ReservationController) ExportStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	studentId := r.PostFormValue("student_id")

	var result = map[string]interface{}{"status": "SUCCESS"}

	url, err := service.Workflow().ExportStudentByAdmin(studentId, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	result["url"] = url

	return result
}

func (rc *ReservationController) UnbindStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	studentId := r.PostFormValue("student_id")

	var result = map[string]interface{}{"status": "SUCCESS"}

	var studentJson = make(map[string]interface{})
	student, err := service.Workflow().UnbindStudentByAdmin(studentId, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	if len(student.BindedTeacherId) != 0 {
		teacher, err := service.Workflow().GetTeacherById(student.BindedTeacherId)
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

func (rc *ReservationController) BindStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	studentId := r.PostFormValue("student_id")
	teacherUsername := r.PostFormValue("teacher_username")

	var result = map[string]interface{}{"status": "SUCCESS"}

	var studentJson = make(map[string]interface{})
	student, err := service.Workflow().BindStudentByAdmin(studentId, teacherUsername, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	} else if len(student.BindedTeacherId) != 0 {
		teacher, err := service.Workflow().GetTeacherById(student.BindedTeacherId)
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

func (rc *ReservationController) QueryStudentInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	studentUsername := r.PostFormValue("student_username")

	var result = map[string]interface{}{"status": "SUCCESS"}

	var studentJson = make(map[string]interface{})
	student, reservations, err := service.Workflow().QueryStudentInfoByAdmin(studentUsername, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
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
		teacher, err := service.Workflow().GetTeacherById(student.BindedTeacherId)
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
		if res.Status == model.RESERVATION_STATUS_AVAILABLE {
			resJson["status"] = model.RESERVATION_STATUS_AVAILABLE
		} else if res.Status == model.RESERVATION_STATUS_RESERVATED && res.StartTime.Before(time.Now()) {
			resJson["status"] = model.RESERVATION_STATUS_FEEDBACK
		} else {
			resJson["status"] = model.RESERVATION_STATUS_RESERVATED
		}
		resJson["student_id"] = res.StudentId
		resJson["teacher_id"] = res.TeacherId
		if teacher, err := service.Workflow().GetTeacherById(res.TeacherId); err == nil {
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

func (rc *ReservationController) SearchTeacherByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMoble := r.PostFormValue("teacher_mobile")

	var result = map[string]interface{}{"status": "SUCCESS"}

	var teacherJson = make(map[string]interface{})
	teacher, err := service.Workflow().SearchTeacherByAdmin(teacherFullname, teacherUsername, teacherMoble,
		userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	teacherJson["teacher_id"] = teacher.Id.Hex()
	teacherJson["teacher_username"] = teacher.Username
	teacherJson["teacher_fullname"] = teacher.Fullname
	teacherJson["teacher_mobile"] = teacher.Mobile
	result["teacher"] = teacherJson

	return result
}

func (rc *ReservationController) GetTeacherWorkloadByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	fromDate := r.PostFormValue("from_date")
	toDate := r.PostFormValue("to_date")

	var result = map[string]interface{}{"status": "SUCCESS"}

	workload, err := service.Workflow().GetTeacherWorkloadByAdmin(fromDate, toDate, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	result["workload"] = workload

	return result
}

func (rc *ReservationController) ExportReportFormByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	fromDate := r.PostFormValue("from_date")
	toDate := r.PostFormValue("to_date")

	var result = map[string]interface{}{"status": "SUCCESS"}

	url, err := service.Workflow().ExportReportFormByAdmin(fromDate, toDate, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	result["url"] = url

	return result
}

func (rc *ReservationController) ExportReportMonthlyByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	monthlyDate := r.PostFormValue("monthly_date")

	var result = map[string]interface{}{"status": "SUCCESS"}

	reportUrl, keyCaseUrl, err := service.Workflow().ExportReportMonthlyByAdmin(monthlyDate, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	result["report"] = reportUrl
	result["key_case"] = keyCaseUrl

	return result
}

//==================== timetable ====================
func (rc *ReservationController) ViewTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	var result = map[string]interface{}{"status": "SUCCESS"}

	timedReservations, err := service.Workflow().ViewTimetableByAdmin(userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
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
			trJson["status"] = tr.Status
			trJson["teacher_id"] = tr.TeacherId
			if teacher, err := service.Workflow().GetTeacherById(tr.TeacherId); err == nil {
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

func (rc *ReservationController) AddTimedReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	weekday := r.PostFormValue("weekday")
	startTime := r.PostFormValue("start_clock")
	endTime := r.PostFormValue("end_clock")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")
	force := strings.EqualFold(r.PostFormValue("force"), "FORCE")

	var result = map[string]interface{}{"status": "SUCCESS"}

	var timedReservationJson = make(map[string]interface{})
	timedReservation, err := service.Workflow().AddTimetableByAdmin(weekday, startTime, endTime, teacherUsername, teacherFullname,
		teacherMobile, force, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	timedReservationJson["timed_reservation_id"] = timedReservation.Id.Hex()
	timedReservationJson["weekday"] = timedReservation.Weekday
	timedReservationJson["start_clock"] = timedReservation.StartTime.Format("15:04")
	timedReservationJson["end_clock"] = timedReservation.EndTime.Format("15:04")
	timedReservationJson["status"] = timedReservation.Status
	timedReservationJson["teacher_id"] = timedReservation.TeacherId
	if teacher, err := service.Workflow().GetTeacherById(timedReservation.TeacherId); err == nil {
		timedReservationJson["teacher_username"] = teacher.Username
		timedReservationJson["teacher_fullname"] = teacher.Fullname
		timedReservationJson["teacher_mobile"] = teacher.Mobile
	}
	result["timed_reservation"] = timedReservationJson

	return result
}

func (rc *ReservationController) EditTimedReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	timedReservationId := r.PostFormValue("timed_reservation_id")
	weekday := r.PostFormValue("weekday")
	startTime := r.PostFormValue("start_clock")
	endTime := r.PostFormValue("end_clock")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")
	force := strings.EqualFold(r.PostFormValue("force"), "FORCE")

	var result = map[string]interface{}{"status": "SUCCESS"}

	var timedReservationJson = make(map[string]interface{})
	timedReservation, err := service.Workflow().EditTimetableByAdmin(timedReservationId, weekday, startTime, endTime, teacherUsername,
		teacherFullname, teacherMobile, force, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	timedReservationJson["timed_reservation_id"] = timedReservation.Id.Hex()
	timedReservationJson["weekday"] = timedReservation.Weekday
	timedReservationJson["start_clock"] = timedReservation.StartTime.Format("15:04")
	timedReservationJson["end_clock"] = timedReservation.EndTime.Format("15:04")
	timedReservationJson["status"] = timedReservation.Status
	timedReservationJson["teacher_id"] = timedReservation.TeacherId
	if teacher, err := service.Workflow().GetTeacherById(timedReservation.TeacherId); err == nil {
		timedReservationJson["teacher_username"] = teacher.Username
		timedReservationJson["teacher_fullname"] = teacher.Fullname
		timedReservationJson["teacher_mobile"] = teacher.Mobile
	}
	result["timed_reservation"] = timedReservationJson

	return result
}

func (rc *ReservationController) RemoveTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	r.ParseForm()
	timedReservationIds := []string(r.Form["timed_reservation_ids"])

	var result = map[string]interface{}{"status": "SUCCESS"}

	removed, err := service.Workflow().RemoveTimetablesByAdmin(timedReservationIds, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	result["removed_count"] = removed

	return result
}

func (rc *ReservationController) OpenTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	r.ParseForm()
	timedReservationIds := []string(r.Form["timed_reservation_ids"])

	var result = map[string]interface{}{"status": "SUCCESS"}

	opened, err := service.Workflow().OpenTimetablesByAdmin(timedReservationIds, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	result["opened_count"] = opened

	return result
}

func (rc *ReservationController) CloseTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	r.ParseForm()
	timedReservationIds := []string(r.Form["timed_reservation_ids"])

	var result = map[string]interface{}{"status": "SUCCESS"}

	closed, err := service.Workflow().CloseTimetablesByAdmin(timedReservationIds, userId, userType)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	result["closed_count"] = closed

	return result
}
