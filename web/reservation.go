package web

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/service"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ReservationController struct {
	BaseMuxController
}

const (
	kCategoryApiBaseUrl = "/api/category"
	kStudentApiBaseUrl  = "/api/student"
	kTeacherApiBaseUrl  = "/api/teacher"
	kAdminApiBaseUrl    = "/api/admin"
)

func (rc *ReservationController) MuxHandlers(m JsonMuxer) {
	categoryBaseUrl := kCategoryApiBaseUrl
	m.GetJson(categoryBaseUrl+"/feedback", "GetFeedbackCategories", rc.GetFeedbackCategories)

	studentBaseUrl := kStudentApiBaseUrl
	m.GetJson(studentBaseUrl+"/reservation/view", "ViewReservationsByStudent", RoleCookieInjection(rc.ViewReservationsByStudent))
	m.PostJson(studentBaseUrl+"/reservation/make", "MakeReservationByStudent", RoleCookieInjection(rc.MakeReservationByStudent))
	m.PostJson(studentBaseUrl+"/reservation/feedback/get", "GetFeedbackByStudent", RoleCookieInjection(rc.GetFeedbackByStudent))
	m.PostJson(studentBaseUrl+"/reservation/feedback/submit", "SubmitFeedbackByStudent", RoleCookieInjection(rc.SubmitFeedbackByStudent))

	teacherBaseUrl := kTeacherApiBaseUrl
	m.GetJson(teacherBaseUrl+"/reservation/view", "ViewReservationsByTeacher", RoleCookieInjection(rc.ViewReservationsByTeacher))
	m.PostJson(teacherBaseUrl+"/reservation/feedback/get", "GetFeedbackByTeacher", RoleCookieInjection(rc.GetFeedbackByTeacher))
	m.PostJson(teacherBaseUrl+"/reservation/feedback/submit", "SubmitFeedbackByTeacher", RoleCookieInjection(rc.SubmitFeedbackByTeacher))
	m.PostJson(teacherBaseUrl+"/student/get", "GetStudentInfoByTeacher", RoleCookieInjection(rc.GetStudentInfoByTeacher))
	m.PostJson(teacherBaseUrl+"/student/query", "QueryStudentInfoByTeacher", RoleCookieInjection(rc.QueryStudentInfoByTeacher))

	adminBaseUrl := kAdminApiBaseUrl
	m.GetJson(adminBaseUrl+"/timetable/view", "ViewTimedReservationsByAdmin", RoleCookieInjection(rc.ViewTimedReservationsByAdmin))
	m.PostJson(adminBaseUrl+"/timetable/add", "AddTimedReservationByAdmin", RoleCookieInjection(rc.AddTimedReservationByAdmin))
	m.PostJson(adminBaseUrl+"/timetable/edit", "EditTimedReservationByAdmin", RoleCookieInjection(rc.EditTimedReservationByAdmin))
	m.PostJson(adminBaseUrl+"/timetable/remove", "RemoveTimedReservationsByAdmin", RoleCookieInjection(rc.RemoveTimedReservationsByAdmin))
	m.PostJson(adminBaseUrl+"/timetable/open", "OpenTimedReservationsByAdmin", RoleCookieInjection(rc.OpenTimedReservationsByAdmin))
	m.PostJson(adminBaseUrl+"/timetable/close", "CloseTimedReservationsByAdmin", RoleCookieInjection(rc.CloseTimedReservationsByAdmin))
	m.GetJson(adminBaseUrl+"/reservation/view", "ViewReservationsByAdmin", RoleCookieInjection(rc.ViewReservationsByAdmin))
	m.GetJson(adminBaseUrl+"/reservation/view/daily", "ViewDailyReservationsByAdmin", RoleCookieInjection(rc.ViewDailyReservationsByAdmin))
	m.GetJson(adminBaseUrl+"/reservation/export/today", "ExportTodayReservationsByAdmin", RoleCookieInjection(rc.ExportTodayReservationsByAdmin))
	m.PostJson(adminBaseUrl+"/reservation/export/report", "ExportReportFormByAdmin", RoleCookieInjection(rc.ExportReportFormByAdmin))
	m.PostJson(adminBaseUrl+"/reservation/export/report/monthly", "ExportReportMonthlyByAdmin", RoleCookieInjection(rc.ExportReportMonthlyByAdmin))
	m.PostJson(adminBaseUrl+"/reservation/add", "AddReservationByAdmin", RoleCookieInjection(rc.AddReservationByAdmin))
	m.PostJson(adminBaseUrl+"/reservation/edit", "EditReservationByAdmin", RoleCookieInjection(rc.EditReservationByAdmin))
	m.PostJson(adminBaseUrl+"/reservation/remove", "RemoveReservationsByAdmin", RoleCookieInjection(rc.RemoveReservationsByAdmin))
	m.PostJson(adminBaseUrl+"/reservation/cancel", "CancelReservationByAdmin", RoleCookieInjection(rc.CancelReservationByAdmin))
	m.PostJson(adminBaseUrl+"/reservation/feedback/get", "GetFeedbackByAdmin", RoleCookieInjection(rc.GetFeedbackByAdmin))
	m.PostJson(adminBaseUrl+"/reservation/feedback/submit", "SubmitFeedbackByAdmin", RoleCookieInjection(rc.SubmitFeedbackByAdmin))
	m.PostJson(adminBaseUrl+"/reservation/student/set", "SetStudentByAdmin", RoleCookieInjection(rc.SetStudentByAdmin))
	m.PostJson(adminBaseUrl+"/student/get", "GetStudentInfoByAdmin", RoleCookieInjection(rc.GetStudentInfoByAdmin))
	m.PostJson(adminBaseUrl+"/student/search", "SearchStudentByAdmin", RoleCookieInjection(rc.SearchStudentByAdmin))
	m.PostJson(adminBaseUrl+"/student/crisis/update", "UpdateStudentCrisisLevelByAdmin", RoleCookieInjection(rc.UpdateStudentCrisisLevelByAdmin))
	m.PostJson(adminBaseUrl+"/student/archive/update", "UpdateStudentArchiveNumberByAdmin", RoleCookieInjection(rc.UpdateStudentArchiveNumberByAdmin))
	m.PostJson(adminBaseUrl+"/student/password/reset", "ResetStudentPasswordByAdmin", RoleCookieInjection(rc.ResetStudentPasswordByAdmin))
	m.PostJson(adminBaseUrl+"/student/account/delete", "DeleteStudentAccountByAdmin", RoleCookieInjection(rc.DeleteStudentAccountByAdmin))
	m.PostJson(adminBaseUrl+"/student/export", "ExportStudentByAdmin", RoleCookieInjection(rc.ExportStudentByAdmin))
	m.PostJson(adminBaseUrl+"/student/unbind", "UnbindStudentByAdmin", RoleCookieInjection(rc.UnbindStudentByAdmin))
	m.PostJson(adminBaseUrl+"/student/bind", "BindStudentByAdmin", RoleCookieInjection(rc.BindStudentByAdmin))
	m.PostJson(adminBaseUrl+"/student/query", "QueryStudentInfoByAdmin", RoleCookieInjection(rc.QueryStudentInfoByAdmin))
	m.PostJson(adminBaseUrl+"/teacher/search", "SearchTeacherByAdmin", RoleCookieInjection(rc.SearchTeacherByAdmin))
	m.PostJson(adminBaseUrl+"/teacher/workload", "GetTeacherWorkloadByAdmin", RoleCookieInjection(rc.GetTeacherWorkloadByAdmin))
}

//==================== category ====================
func (rc *ReservationController) GetFeedbackCategories(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	var result = make(map[string]interface{})

	result["first_category"] = model.FeedbackFirstCategory
	result["second_category"] = model.FeedbackSecondCategory

	return http.StatusOK, wrapJsonOk(result)
}

//==================== student ====================
func (rc *ReservationController) ViewReservationsByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	var result = make(map[string]interface{})

	student, err := service.Workflow().GetStudentById(userId)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	var studentJson = make(map[string]interface{})
	studentJson["fullname"] = student.Fullname
	studentJson["gender"] = student.Gender
	studentJson["email"] = student.Email
	studentJson["school"] = student.School
	studentJson["grade"] = student.Grade
	studentJson["current_address"] = student.CurrentAddress
	studentJson["mobile"] = student.Mobile
	studentJson["birthday"] = student.Birthday
	studentJson["family_address"] = student.FamilyAddress
	studentJson["experience_time"] = student.Experience.Time
	studentJson["experience_location"] = student.Experience.Location
	studentJson["experience_teacher"] = student.Experience.Teacher
	studentJson["father_age"] = student.FatherAge
	studentJson["father_job"] = student.FatherJob
	studentJson["father_edu"] = student.FatherEdu
	studentJson["mother_age"] = student.MotherAge
	studentJson["mother_job"] = student.MotherJob
	studentJson["mother_edu"] = student.MotherEdu
	studentJson["parent_marriage"] = student.ParentMarriage
	studentJson["significant"] = student.Significant
	studentJson["problem"] = student.Problem
	result["student"] = studentJson

	reservations, err := service.Workflow().GetReservationsByStudent(userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) MakeReservationByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
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

	var result = make(map[string]interface{})

	var reservationJson = make(map[string]interface{})
	reservation, err := service.Workflow().MakeReservationByStudent(reservationId, sourceId, startTime, fullname, gender, birthday,
		school, grade, currentAddress, familyAddress, mobile, email, experienceTime, experienceLocation, experienceTeacher,
		fatherAge, fatherJob, fatherEdu, motherAge, motherJob, motherEdu, parentMarriage, siginificant, problem,
		userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	reservationJson["start_time"] = reservation.StartTime.Format("2006-01-02 15:04")
	reservationJson["end_time"] = reservation.EndTime.Format("2006-01-02 15:04")
	if teacher, err := service.Workflow().GetTeacherById(reservation.TeacherId); err == nil {
		reservationJson["teacher_fullname"] = teacher.Fullname
	}
	result["reservation"] = reservationJson

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetFeedbackByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")

	var result = make(map[string]interface{})

	var feedbackJson = make(map[string]interface{})
	reservation, err := service.Workflow().GetFeedbackByStudent(reservationId, sourceId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	feedbackJson["scores"] = reservation.StudentFeedback.Scores
	result["feedback"] = feedbackJson

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SubmitFeedbackByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")
	r.ParseForm()
	scores := []string(r.Form["scores"])

	var result = make(map[string]interface{})

	scoresInt := []int{}
	for _, p := range scores {
		if pi, err := strconv.Atoi(p); err == nil {
			scoresInt = append(scoresInt, pi)
		}
	}
	_, err := service.Workflow().SubmitFeedbackByStudent(reservationId, sourceId, scoresInt, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}

	return http.StatusOK, wrapJsonOk(result)
}

//==================== teacher ====================
func (rc *ReservationController) ViewReservationsByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	var result = make(map[string]interface{})

	teacher, err := service.Workflow().GetTeacherById(userId)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	var teacherJson = make(map[string]interface{})
	teacherJson["teacher_fullname"] = teacher.Fullname
	teacherJson["teacher_mobile"] = teacher.Mobile
	result["teacher_info"] = teacherJson

	reservations, err := service.Workflow().GetReservationsByTeacher(userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetFeedbackByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")

	var result = make(map[string]interface{})

	var feedback = make(map[string]interface{})
	student, reservation, err := service.Workflow().GetFeedbackByTeacher(reservationId, sourceId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SubmitFeedbackByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
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

	var result = make(map[string]interface{})

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
		return http.StatusOK, wrapJsonError(err.Error())
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetStudentInfoByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.PostFormValue("student_id")

	var result = make(map[string]interface{})

	var studentJson = make(map[string]interface{})
	student, reservations, err := service.Workflow().GetStudentInfoByTeacher(studentId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) QueryStudentInfoByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentUsername := r.PostFormValue("student_username")

	var result = make(map[string]interface{})

	var studentJson = make(map[string]interface{})
	student, reservations, err := service.Workflow().QueryStudentInfoByTeacher(studentUsername, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

//==================== admin ====================
func (rc *ReservationController) ViewReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	var result = make(map[string]interface{})

	reservations, err := service.Workflow().GetReservationsByAdmin(userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) ViewDailyReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["from_date"]) == 0 {
		return http.StatusOK, wrapJsonError("参数错误")
	}
	fromDate := queryForm["from_date"][0]

	var result = make(map[string]interface{})

	reservations, err := service.Workflow().GetReservationsDailyByAdmin(fromDate, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) ExportTodayReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	var result = make(map[string]interface{})

	url, err := service.Workflow().ExportTodayReservationTimetableByAdmin(userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["url"] = url

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) AddReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	startTime := r.PostFormValue("start_time")
	endTime := r.PostFormValue("end_time")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")
	force := strings.EqualFold(r.PostFormValue("force"), "FORCE")

	var result = make(map[string]interface{})

	var reservationJson = make(map[string]interface{})
	reservation, err := service.Workflow().AddReservationByAdmin(startTime, endTime, teacherUsername, teacherFullname,
		teacherMobile, force, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) EditReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")
	originalStartTime := r.PostFormValue("original_start_time")
	startTime := r.PostFormValue("start_time")
	endTime := r.PostFormValue("end_time")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")
	force := strings.EqualFold(r.PostFormValue("force"), "FORCE")

	var result = make(map[string]interface{})

	var reservationJson = make(map[string]interface{})
	reservation, err := service.Workflow().EditReservationByAdmin(reservationId, sourceId, originalStartTime,
		startTime, endTime, teacherUsername, teacherFullname, teacherMobile, force, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) RemoveReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	r.ParseForm()
	reservationIds := []string(r.Form["reservation_ids"])
	sourceIds := []string(r.Form["source_ids"])
	startTimes := []string(r.Form["start_times"])

	var result = make(map[string]interface{})

	removed, err := service.Workflow().RemoveReservationsByAdmin(reservationIds, sourceIds, startTimes, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["removed_count"] = removed

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) CancelReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	r.ParseForm()
	reservationIds := []string(r.Form["reservation_ids"])
	sourceIds := []string(r.Form["source_ids"])

	var result = make(map[string]interface{})

	removed, err := service.Workflow().CancelReservationsByAdmin(reservationIds, sourceIds, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["removed_count"] = removed

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")

	var result = make(map[string]interface{})

	var feedback = make(map[string]interface{})
	student, reservation, err := service.Workflow().GetFeedbackByAdmin(reservationId, sourceId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SubmitFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
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

	var result = make(map[string]interface{})

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
		return http.StatusOK, wrapJsonError(err.Error())
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SetStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
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
		return http.StatusOK, wrapJsonError("参数错误，请联系管理员")
	}

	var result = make(map[string]interface{})

	var reservationJson = make(map[string]interface{})
	reservation, err := service.Workflow().SetStudentByAdmin(reservationId, sourceId, startTime, studentUsername, fullname,
		gender, birthday, school, grade, currentAddress, familyAddress, mobile, email, experienceTime,
		experienceLocation, experienceTeacher, fatherAge, fatherJob, fatherEdu, motherAge, motherJob, motherEdu,
		parentMarriage, siginificant, problem, sendSms, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetStudentInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.PostFormValue("student_id")

	var result = make(map[string]interface{})

	var studentJson = make(map[string]interface{})
	student, reservations, err := service.Workflow().GetStudentInfoByAdmin(studentId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SearchStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentUsername := r.PostFormValue("student_username")

	var result = make(map[string]interface{})

	var studentJson = make(map[string]interface{})
	student, err := service.Workflow().GetStudentByUsername(studentUsername)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) UpdateStudentCrisisLevelByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.PostFormValue("student_id")
	crisisLevel := r.PostFormValue("crisis_level")

	var result = make(map[string]interface{})

	_, err := service.Workflow().UpdateStudentCrisisLevelByAdmin(studentId, crisisLevel, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) UpdateStudentArchiveNumberByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.PostFormValue("student_id")
	archiveCategory := r.PostFormValue("archive_category")
	archiveNumber := r.PostFormValue("archive_number")

	var result = make(map[string]interface{})

	_, err := service.Workflow().UpdateStudentArchiveNumberByAdmin(studentId, archiveCategory, archiveNumber, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) ResetStudentPasswordByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.PostFormValue("student_id")
	password := r.PostFormValue("password")

	var result = make(map[string]interface{})

	_, err := service.Workflow().ResetStudentPasswordByAdmin(studentId, password, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) DeleteStudentAccountByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.PostFormValue("student_id")

	var result = make(map[string]interface{})

	err := service.Workflow().DeleteStudentAccountByAdmin(studentId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) ExportStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.PostFormValue("student_id")

	var result = make(map[string]interface{})

	url, err := service.Workflow().ExportStudentByAdmin(studentId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["url"] = url

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) UnbindStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.PostFormValue("student_id")

	var result = make(map[string]interface{})

	var studentJson = make(map[string]interface{})
	student, err := service.Workflow().UnbindStudentByAdmin(studentId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) BindStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.PostFormValue("student_id")
	teacherUsername := r.PostFormValue("teacher_username")

	var result = make(map[string]interface{})

	var studentJson = make(map[string]interface{})
	student, err := service.Workflow().BindStudentByAdmin(studentId, teacherUsername, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) QueryStudentInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentUsername := r.PostFormValue("student_username")

	var result = make(map[string]interface{})

	var studentJson = make(map[string]interface{})
	student, reservations, err := service.Workflow().QueryStudentInfoByAdmin(studentUsername, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SearchTeacherByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMoble := r.PostFormValue("teacher_mobile")

	var result = make(map[string]interface{})

	var teacherJson = make(map[string]interface{})
	teacher, err := service.Workflow().SearchTeacherByAdmin(teacherFullname, teacherUsername, teacherMoble,
		userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	teacherJson["teacher_id"] = teacher.Id.Hex()
	teacherJson["teacher_username"] = teacher.Username
	teacherJson["teacher_fullname"] = teacher.Fullname
	teacherJson["teacher_mobile"] = teacher.Mobile
	result["teacher"] = teacherJson

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetTeacherWorkloadByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	fromDate := r.PostFormValue("from_date")
	toDate := r.PostFormValue("to_date")

	var result = make(map[string]interface{})

	workload, err := service.Workflow().GetTeacherWorkloadByAdmin(fromDate, toDate, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["workload"] = workload

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) ExportReportFormByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	fromDate := r.PostFormValue("from_date")
	toDate := r.PostFormValue("to_date")

	var result = make(map[string]interface{})

	url, err := service.Workflow().ExportReportFormByAdmin(fromDate, toDate, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["url"] = url

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) ExportReportMonthlyByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	monthlyDate := r.PostFormValue("monthly_date")

	var result = make(map[string]interface{})

	reportUrl, keyCaseUrl, err := service.Workflow().ExportReportMonthlyByAdmin(monthlyDate, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["report"] = reportUrl
	result["key_case"] = keyCaseUrl

	return http.StatusOK, wrapJsonOk(result)
}

//==================== timetable ====================
func (rc *ReservationController) ViewTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	var result = make(map[string]interface{})

	timedReservations, err := service.Workflow().ViewTimetableByAdmin(userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) AddTimedReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	weekday := r.PostFormValue("weekday")
	startTime := r.PostFormValue("start_clock")
	endTime := r.PostFormValue("end_clock")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")
	force := strings.EqualFold(r.PostFormValue("force"), "FORCE")

	var result = make(map[string]interface{})

	var timedReservationJson = make(map[string]interface{})
	timedReservation, err := service.Workflow().AddTimetableByAdmin(weekday, startTime, endTime, teacherUsername, teacherFullname,
		teacherMobile, force, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) EditTimedReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	timedReservationId := r.PostFormValue("timed_reservation_id")
	weekday := r.PostFormValue("weekday")
	startTime := r.PostFormValue("start_clock")
	endTime := r.PostFormValue("end_clock")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")
	force := strings.EqualFold(r.PostFormValue("force"), "FORCE")

	var result = make(map[string]interface{})

	var timedReservationJson = make(map[string]interface{})
	timedReservation, err := service.Workflow().EditTimetableByAdmin(timedReservationId, weekday, startTime, endTime, teacherUsername,
		teacherFullname, teacherMobile, force, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
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

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) RemoveTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	r.ParseForm()
	timedReservationIds := []string(r.Form["timed_reservation_ids"])

	var result = make(map[string]interface{})

	removed, err := service.Workflow().RemoveTimetablesByAdmin(timedReservationIds, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["removed_count"] = removed

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) OpenTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	r.ParseForm()
	timedReservationIds := []string(r.Form["timed_reservation_ids"])

	var result = make(map[string]interface{})

	opened, err := service.Workflow().OpenTimetablesByAdmin(timedReservationIds, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["opened_count"] = opened

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) CloseTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	r.ParseForm()
	timedReservationIds := []string(r.Form["timed_reservation_ids"])

	var result = make(map[string]interface{})

	closed, err := service.Workflow().CloseTimetablesByAdmin(timedReservationIds, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["closed_count"] = closed

	return http.StatusOK, wrapJsonOk(result)
}
