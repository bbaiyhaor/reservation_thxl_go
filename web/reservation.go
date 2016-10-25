package web

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/service"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
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
	m.GetJson(kCategoryApiBaseUrl+"/feedback", "GetFeedbackCategories", rc.GetFeedbackCategories)

	m.GetJson(kStudentApiBaseUrl+"/reservation/view", "ViewReservationsByStudent", RoleCookieInjection(rc.ViewReservationsByStudent))
	m.PostJson(kStudentApiBaseUrl+"/reservation/make", "MakeReservationByStudent", RoleCookieInjection(rc.MakeReservationByStudent))
	m.PostJson(kStudentApiBaseUrl+"/reservation/feedback/get", "GetFeedbackByStudent", RoleCookieInjection(rc.GetFeedbackByStudent))
	m.PostJson(kStudentApiBaseUrl+"/reservation/feedback/submit", "SubmitFeedbackByStudent", RoleCookieInjection(rc.SubmitFeedbackByStudent))

	m.GetJson(kTeacherApiBaseUrl+"/reservation/view", "ViewReservationsByTeacher", RoleCookieInjection(rc.ViewReservationsByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/reservation/feedback/get", "GetFeedbackByTeacher", RoleCookieInjection(rc.GetFeedbackByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/reservation/feedback/submit", "SubmitFeedbackByTeacher", RoleCookieInjection(rc.SubmitFeedbackByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/student/get", "GetStudentInfoByTeacher", RoleCookieInjection(rc.GetStudentInfoByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/student/query", "QueryStudentInfoByTeacher", RoleCookieInjection(rc.QueryStudentInfoByTeacher))

	m.GetJson(kAdminApiBaseUrl+"/timetable/view", "ViewTimedReservationsByAdmin", RoleCookieInjection(rc.ViewTimedReservationsByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/timetable/add", "AddTimedReservationByAdmin", RoleCookieInjection(rc.AddTimedReservationByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/timetable/edit", "EditTimedReservationByAdmin", RoleCookieInjection(rc.EditTimedReservationByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/timetable/remove", "RemoveTimedReservationsByAdmin", RoleCookieInjection(rc.RemoveTimedReservationsByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/timetable/open", "OpenTimedReservationsByAdmin", RoleCookieInjection(rc.OpenTimedReservationsByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/timetable/close", "CloseTimedReservationsByAdmin", RoleCookieInjection(rc.CloseTimedReservationsByAdmin))
	m.GetJson(kAdminApiBaseUrl+"/reservation/view", "ViewReservationsByAdmin", RoleCookieInjection(rc.ViewReservationsByAdmin))
	m.GetJson(kAdminApiBaseUrl+"/reservation/view/daily", "ViewDailyReservationsByAdmin", RoleCookieInjection(rc.ViewDailyReservationsByAdmin))
	m.GetJson(kAdminApiBaseUrl+"/reservation/export/today", "ExportTodayReservationsByAdmin", RoleCookieInjection(rc.ExportTodayReservationsByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/export/report", "ExportReportFormByAdmin", RoleCookieInjection(rc.ExportReportFormByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/export/report/monthly", "ExportReportMonthlyByAdmin", RoleCookieInjection(rc.ExportReportMonthlyByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/add", "AddReservationByAdmin", RoleCookieInjection(rc.AddReservationByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/edit", "EditReservationByAdmin", RoleCookieInjection(rc.EditReservationByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/remove", "RemoveReservationsByAdmin", RoleCookieInjection(rc.RemoveReservationsByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/cancel", "CancelReservationByAdmin", RoleCookieInjection(rc.CancelReservationByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/feedback/get", "GetFeedbackByAdmin", RoleCookieInjection(rc.GetFeedbackByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/feedback/submit", "SubmitFeedbackByAdmin", RoleCookieInjection(rc.SubmitFeedbackByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/student/set", "SetStudentByAdmin", RoleCookieInjection(rc.SetStudentByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/get", "GetStudentInfoByAdmin", RoleCookieInjection(rc.GetStudentInfoByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/search", "SearchStudentByAdmin", RoleCookieInjection(rc.SearchStudentByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/crisis/update", "UpdateStudentCrisisLevelByAdmin", RoleCookieInjection(rc.UpdateStudentCrisisLevelByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/archive/update", "UpdateStudentArchiveNumberByAdmin", RoleCookieInjection(rc.UpdateStudentArchiveNumberByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/password/reset", "ResetStudentPasswordByAdmin", RoleCookieInjection(rc.ResetStudentPasswordByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/account/delete", "DeleteStudentAccountByAdmin", RoleCookieInjection(rc.DeleteStudentAccountByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/export", "ExportStudentByAdmin", RoleCookieInjection(rc.ExportStudentByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/unbind", "UnbindStudentByAdmin", RoleCookieInjection(rc.UnbindStudentByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/bind", "BindStudentByAdmin", RoleCookieInjection(rc.BindStudentByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/query", "QueryStudentInfoByAdmin", RoleCookieInjection(rc.QueryStudentInfoByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/teacher/search", "SearchTeacherByAdmin", RoleCookieInjection(rc.SearchTeacherByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/teacher/workload", "GetTeacherWorkloadByAdmin", RoleCookieInjection(rc.GetTeacherWorkloadByAdmin))
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
	result["student"] = service.Workflow().WrapSimpleStudent(student)

	reservations, err := service.Workflow().GetReservationsByStudent(userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := service.Workflow().WrapSimpleReservation(res)
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
	reservationId := r.FormValue("reservation_id")
	sourceId := r.FormValue("source_id")
	startTime := r.FormValue("start_time")
	fullname := r.FormValue("student_fullname")
	gender := r.FormValue("student_gender")
	birthday := r.FormValue("student_birthday")
	school := r.FormValue("student_school")
	grade := r.FormValue("student_grade")
	currentAddress := r.FormValue("student_current_address")
	familyAddress := r.FormValue("student_family_address")
	mobile := r.FormValue("student_mobile")
	email := r.FormValue("student_email")
	experienceTime := r.FormValue("student_experience_time")
	experienceLocation := r.FormValue("student_experience_location")
	experienceTeacher := r.FormValue("student_experience_teacher")
	fatherAge := r.FormValue("student_father_age")
	fatherJob := r.FormValue("student_father_job")
	fatherEdu := r.FormValue("student_father_edu")
	motherAge := r.FormValue("student_mother_age")
	motherJob := r.FormValue("student_mother_job")
	motherEdu := r.FormValue("student_mother_edu")
	parentMarriage := r.FormValue("student_parent_marriage")
	siginificant := r.FormValue("student_significant")
	problem := r.FormValue("student_problem")

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().MakeReservationByStudent(reservationId, sourceId, startTime, fullname, gender, birthday,
		school, grade, currentAddress, familyAddress, mobile, email, experienceTime, experienceLocation, experienceTeacher,
		fatherAge, fatherJob, fatherEdu, motherAge, motherJob, motherEdu, parentMarriage, siginificant, problem,
		userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["reservation"] = service.Workflow().WrapSimpleReservation(reservation)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetFeedbackByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := r.FormValue("reservation_id")
	sourceId := r.FormValue("source_id")

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
	reservationId := r.FormValue("reservation_id")
	sourceId := r.FormValue("source_id")
	r.ParseForm()
	scores := []string(r.Form["scores"])
	scoresInt := []int{}
	for _, p := range scores {
		if pi, err := strconv.Atoi(p); err == nil {
			scoresInt = append(scoresInt, pi)
		}
	}

	var result = make(map[string]interface{})

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
	result["teacher"] = service.Workflow().WrapTeacher(teacher)

	reservations, err := service.Workflow().GetReservationsByTeacher(userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := service.Workflow().WrapReservation(res)
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
	reservationId := r.FormValue("reservation_id")
	sourceId := r.FormValue("source_id")

	var result = make(map[string]interface{})

	student, reservation, err := service.Workflow().GetFeedbackByTeacher(reservationId, sourceId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	feedback := reservation.TeacherFeedback.ToJson()
	feedback["crisis_level"] = student.CrisisLevel
	result["feedback"] = feedback

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SubmitFeedbackByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := r.FormValue("reservation_id")
	sourceId := r.FormValue("source_id")
	category := r.FormValue("category")
	r.ParseForm()
	participants := []string(r.Form["participants"])
	participantsInt := make([]int, 0)
	for _, p := range participants {
		if pi, err := strconv.Atoi(p); err == nil {
			participantsInt = append(participantsInt, pi)
		}
	}
	emphasis := r.FormValue("emphasis")
	severity := []string(r.Form["severity"])
	severityInt := make([]int, 0)
	for _, s := range severity {
		if si, err := strconv.Atoi(s); err == nil {
			severityInt = append(severityInt, si)
		}
	}
	medicalDiagnosis := []string(r.Form["medical_diagnosis"])
	medicalDiagnosisInt := make([]int, 0)
	for _, m := range medicalDiagnosis {
		if mi, err := strconv.Atoi(m); err == nil {
			medicalDiagnosisInt = append(medicalDiagnosisInt, mi)
		}
	}
	crisis := []string(r.Form["crisis"])
	crisisInt := make([]int, 0)
	for _, c := range crisis {
		if ci, err := strconv.Atoi(c); err == nil {
			crisisInt = append(crisisInt, ci)
		}
	}
	record := r.FormValue("record")
	crisisLevel := r.FormValue("crisis_level")

	var result = make(map[string]interface{})

	_, err := service.Workflow().SubmitFeedbackByTeacher(reservationId, sourceId, category, participantsInt, emphasis, severityInt,
		medicalDiagnosisInt, crisisInt, record, crisisLevel, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetStudentInfoByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.FormValue("student_id")

	var result = make(map[string]interface{})

	student, reservations, err := service.Workflow().GetStudentInfoByTeacher(studentId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["student"] = service.Workflow().WrapStudent(student)

	var reservationJson = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := service.Workflow().WrapReservation(res)
		reservationJson = append(reservationJson, resJson)
	}
	result["reservations"] = reservationJson

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) QueryStudentInfoByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentUsername := r.FormValue("student_username")

	var result = make(map[string]interface{})

	student, reservations, err := service.Workflow().QueryStudentInfoByTeacher(studentUsername, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["student"] = service.Workflow().WrapStudent(student)

	var reservationJson = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := service.Workflow().WrapReservation(res)
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
		resJson := service.Workflow().WrapReservation(res)
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
	fromDate := r.FormValue("from_date")

	var result = make(map[string]interface{})

	reservations, err := service.Workflow().GetReservationsDailyByAdmin(fromDate, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := service.Workflow().WrapReservation(res)
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
	startTime := r.FormValue("start_time")
	endTime := r.FormValue("end_time")
	teacherUsername := r.FormValue("teacher_username")
	teacherFullname := r.FormValue("teacher_fullname")
	teacherMobile := r.FormValue("teacher_mobile")
	force, err := strconv.ParseBool(r.FormValue("force"))
	if err != nil {
		return http.StatusOK, wrapJsonError("参数错误，请联系管理员")
	}

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().AddReservationByAdmin(startTime, endTime, teacherUsername, teacherFullname,
		teacherMobile, force, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["reservation"] = service.Workflow().WrapReservation(reservation)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) EditReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := r.FormValue("reservation_id")
	sourceId := r.FormValue("source_id")
	originalStartTime := r.FormValue("original_start_time")
	startTime := r.FormValue("start_time")
	endTime := r.FormValue("end_time")
	teacherUsername := r.FormValue("teacher_username")
	teacherFullname := r.FormValue("teacher_fullname")
	teacherMobile := r.FormValue("teacher_mobile")
	force, err := strconv.ParseBool(r.FormValue("force"))
	if err != nil {
		return http.StatusOK, wrapJsonError("参数错误，请联系管理员")
	}

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().EditReservationByAdmin(reservationId, sourceId, originalStartTime,
		startTime, endTime, teacherUsername, teacherFullname, teacherMobile, force, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["reservation"] = service.Workflow().WrapReservation(reservation)

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
	result["canceled_count"] = removed

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := r.FormValue("reservation_id")
	sourceId := r.FormValue("source_id")

	var result = make(map[string]interface{})

	student, reservation, err := service.Workflow().GetFeedbackByAdmin(reservationId, sourceId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	feedback := reservation.TeacherFeedback.ToJson()
	feedback["crisis_level"] = student.CrisisLevel
	result["feedback"] = feedback

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SubmitFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := r.FormValue("reservation_id")
	sourceId := r.FormValue("source_id")
	category := r.FormValue("category")
	r.ParseForm()
	participants := []string(r.Form["participants"])
	participantsInt := make([]int, 0)
	for _, p := range participants {
		if pi, err := strconv.Atoi(p); err == nil {
			participantsInt = append(participantsInt, pi)
		}
	}
	emphasis := r.FormValue("emphasis")
	severity := []string(r.Form["severity"])
	severityInt := make([]int, 0)
	for _, s := range severity {
		if si, err := strconv.Atoi(s); err == nil {
			severityInt = append(severityInt, si)
		}
	}
	medicalDiagnosis := []string(r.Form["medical_diagnosis"])
	medicalDiagnosisInt := make([]int, 0)
	for _, m := range medicalDiagnosis {
		if mi, err := strconv.Atoi(m); err == nil {
			medicalDiagnosisInt = append(medicalDiagnosisInt, mi)
		}
	}
	crisis := []string(r.Form["crisis"])
	crisisInt := make([]int, 0)
	for _, c := range crisis {
		if ci, err := strconv.Atoi(c); err == nil {
			crisisInt = append(crisisInt, ci)
		}
	}
	record := r.FormValue("record")
	crisisLevel := r.FormValue("crisis_level")

	var result = make(map[string]interface{})

	_, err := service.Workflow().SubmitFeedbackByAdmin(reservationId, sourceId, category, participantsInt, emphasis, severityInt,
		medicalDiagnosisInt, crisisInt, record, crisisLevel, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SetStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := r.FormValue("reservation_id")
	sourceId := r.FormValue("source_id")
	startTime := r.FormValue("start_time")
	studentUsername := r.FormValue("student_username")
	fullname := r.FormValue("student_fullname")
	gender := r.FormValue("student_gender")
	birthday := r.FormValue("student_birthday")
	school := r.FormValue("student_school")
	grade := r.FormValue("student_grade")
	currentAddress := r.FormValue("student_current_address")
	familyAddress := r.FormValue("student_family_address")
	mobile := r.FormValue("student_mobile")
	email := r.FormValue("student_email")
	experienceTime := r.FormValue("student_experience_time")
	experienceLocation := r.FormValue("student_experience_location")
	experienceTeacher := r.FormValue("student_experience_teacher")
	fatherAge := r.FormValue("student_father_age")
	fatherJob := r.FormValue("student_father_job")
	fatherEdu := r.FormValue("student_father_edu")
	motherAge := r.FormValue("student_mother_age")
	motherJob := r.FormValue("student_mother_job")
	motherEdu := r.FormValue("student_mother_edu")
	parentMarriage := r.FormValue("student_parent_marriage")
	siginificant := r.FormValue("student_significant")
	problem := r.FormValue("student_problem")
	sendSms, err := strconv.ParseBool(r.FormValue("student_sms"))
	if err != nil {
		return http.StatusOK, wrapJsonError("参数错误，请联系管理员")
	}

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().SetStudentByAdmin(reservationId, sourceId, startTime, studentUsername, fullname,
		gender, birthday, school, grade, currentAddress, familyAddress, mobile, email, experienceTime,
		experienceLocation, experienceTeacher, fatherAge, fatherJob, fatherEdu, motherAge, motherJob, motherEdu,
		parentMarriage, siginificant, problem, sendSms, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["reservation"] = service.Workflow().WrapReservation(reservation)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetStudentInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.FormValue("student_id")

	var result = make(map[string]interface{})

	student, reservations, err := service.Workflow().GetStudentInfoByAdmin(studentId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["student"] = service.Workflow().WrapStudent(student)

	var reservationJson = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := service.Workflow().WrapReservation(res)
		reservationJson = append(reservationJson, resJson)
	}
	result["reservations"] = reservationJson

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SearchStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentUsername := r.FormValue("student_username")

	var result = make(map[string]interface{})

	student, err := service.Workflow().GetStudentByUsername(studentUsername)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["student"] = service.Workflow().WrapStudent(student)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) UpdateStudentCrisisLevelByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.FormValue("student_id")
	crisisLevel := r.FormValue("crisis_level")

	var result = make(map[string]interface{})

	_, err := service.Workflow().UpdateStudentCrisisLevelByAdmin(studentId, crisisLevel, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) UpdateStudentArchiveNumberByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.FormValue("student_id")
	archiveCategory := r.FormValue("archive_category")
	archiveNumber := r.FormValue("archive_number")

	var result = make(map[string]interface{})

	_, err := service.Workflow().UpdateStudentArchiveNumberByAdmin(studentId, archiveCategory, archiveNumber, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) ResetStudentPasswordByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.FormValue("student_id")
	password := r.FormValue("password")

	var result = make(map[string]interface{})

	_, err := service.Workflow().ResetStudentPasswordByAdmin(studentId, password, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) DeleteStudentAccountByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.FormValue("student_id")

	var result = make(map[string]interface{})

	err := service.Workflow().DeleteStudentAccountByAdmin(studentId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) ExportStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.FormValue("student_id")

	var result = make(map[string]interface{})

	url, err := service.Workflow().ExportStudentByAdmin(studentId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["url"] = url

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) UnbindStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.FormValue("student_id")

	var result = make(map[string]interface{})

	student, err := service.Workflow().UnbindStudentByAdmin(studentId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["student"] = service.Workflow().WrapStudent(student)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) BindStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := r.FormValue("student_id")
	teacherUsername := r.FormValue("teacher_username")

	var result = make(map[string]interface{})

	student, err := service.Workflow().BindStudentByAdmin(studentId, teacherUsername, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["student"] = service.Workflow().WrapStudent(student)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) QueryStudentInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentUsername := r.FormValue("student_username")

	var result = make(map[string]interface{})

	student, reservations, err := service.Workflow().QueryStudentInfoByAdmin(studentUsername, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["student"] = service.Workflow().WrapStudent(student)

	var reservationJson = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := service.Workflow().WrapReservation(res)
		reservationJson = append(reservationJson, resJson)
	}
	result["reservations"] = reservationJson

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SearchTeacherByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	teacherUsername := r.FormValue("teacher_username")
	teacherFullname := r.FormValue("teacher_fullname")
	teacherMoble := r.FormValue("teacher_mobile")

	var result = make(map[string]interface{})

	teacher, err := service.Workflow().SearchTeacherByAdmin(teacherFullname, teacherUsername, teacherMoble,
		userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["teacher"] = service.Workflow().WrapTeacher(teacher)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetTeacherWorkloadByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	fromDate := r.FormValue("from_date")
	toDate := r.FormValue("to_date")

	var result = make(map[string]interface{})

	workload, err := service.Workflow().GetTeacherWorkloadByAdmin(fromDate, toDate, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["workload"] = workload

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) ExportReportFormByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	fromDate := r.FormValue("from_date")
	toDate := r.FormValue("to_date")

	var result = make(map[string]interface{})

	url, err := service.Workflow().ExportReportFormByAdmin(fromDate, toDate, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["url"] = url

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) ExportReportMonthlyByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	monthlyDate := r.FormValue("monthly_date")

	var result = make(map[string]interface{})

	reportUrl, keyCaseUrl, err := service.Workflow().ExportReportMonthlyByAdmin(monthlyDate, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["report_url"] = reportUrl
	result["key_case_url"] = keyCaseUrl

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
			trJson := service.Workflow().WrapTimedReservation(tr)
			array = append(array, trJson)
		}
		timetable[weekday.String()] = array
	}
	result["timed_reservations"] = timetable

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) AddTimedReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	weekday := r.FormValue("weekday")
	startTime := r.FormValue("start_clock")
	endTime := r.FormValue("end_clock")
	teacherUsername := r.FormValue("teacher_username")
	teacherFullname := r.FormValue("teacher_fullname")
	teacherMobile := r.FormValue("teacher_mobile")
	force, err := strconv.ParseBool(r.FormValue("force"))
	if err != nil {
		return http.StatusOK, wrapJsonError("参数错误，请联系管理员")
	}

	var result = make(map[string]interface{})

	timedReservation, err := service.Workflow().AddTimetableByAdmin(weekday, startTime, endTime, teacherUsername, teacherFullname,
		teacherMobile, force, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["timed_reservation"] = service.Workflow().WrapTimedReservation(timedReservation)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) EditTimedReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	timedReservationId := r.FormValue("timed_reservation_id")
	weekday := r.FormValue("weekday")
	startTime := r.FormValue("start_clock")
	endTime := r.FormValue("end_clock")
	teacherUsername := r.FormValue("teacher_username")
	teacherFullname := r.FormValue("teacher_fullname")
	teacherMobile := r.FormValue("teacher_mobile")
	force, err := strconv.ParseBool(r.FormValue("force"))
	if err != nil {
		return http.StatusOK, wrapJsonError("参数错误，请联系管理员")
	}

	var result = make(map[string]interface{})

	timedReservation, err := service.Workflow().EditTimetableByAdmin(timedReservationId, weekday, startTime, endTime, teacherUsername,
		teacherFullname, teacherMobile, force, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	result["timed_reservation"] = service.Workflow().WrapTimedReservation(timedReservation)

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
