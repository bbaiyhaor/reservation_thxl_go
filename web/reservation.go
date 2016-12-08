package web

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/service"
	"github.com/mijia/sweb/form"
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
	m.GetJson(kCategoryApiBaseUrl+"/feedback", "GetFeedbackCategories", rc.getFeedbackCategories)

	m.GetJson(kStudentApiBaseUrl+"/reservation/view", "ViewReservationsByStudent", RoleCookieInjection(rc.viewReservationsByStudent))
	m.PostJson(kStudentApiBaseUrl+"/reservation/make", "MakeReservationByStudent", RoleCookieInjection(rc.makeReservationByStudent))
	m.PostJson(kStudentApiBaseUrl+"/reservation/feedback/get", "GetFeedbackByStudent", RoleCookieInjection(rc.getFeedbackByStudent))
	m.PostJson(kStudentApiBaseUrl+"/reservation/feedback/submit", "SubmitFeedbackByStudent", RoleCookieInjection(rc.submitFeedbackByStudent))

	m.GetJson(kTeacherApiBaseUrl+"/reservation/view", "ViewReservationsByTeacher", RoleCookieInjection(rc.viewReservationsByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/reservation/feedback/get", "GetFeedbackByTeacher", RoleCookieInjection(rc.getFeedbackByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/reservation/feedback/submit", "SubmitFeedbackByTeacher", RoleCookieInjection(rc.submitFeedbackByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/student/get", "GetStudentInfoByTeacher", RoleCookieInjection(rc.getStudentInfoByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/student/query", "QueryStudentInfoByTeacher", RoleCookieInjection(rc.queryStudentInfoByTeacher))

	m.GetJson(kAdminApiBaseUrl+"/timetable/view", "ViewTimedReservationsByAdmin", RoleCookieInjection(rc.viewTimedReservationsByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/timetable/add", "AddTimedReservationByAdmin", RoleCookieInjection(rc.addTimedReservationByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/timetable/edit", "EditTimedReservationByAdmin", RoleCookieInjection(rc.editTimedReservationByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/timetable/remove", "RemoveTimedReservationsByAdmin", RoleCookieInjection(rc.removeTimedReservationsByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/timetable/open", "OpenTimedReservationsByAdmin", RoleCookieInjection(rc.openTimedReservationsByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/timetable/close", "CloseTimedReservationsByAdmin", RoleCookieInjection(rc.closeTimedReservationsByAdmin))
	m.GetJson(kAdminApiBaseUrl+"/reservation/view", "ViewReservationsByAdmin", RoleCookieInjection(rc.viewReservationsByAdmin))
	m.GetJson(kAdminApiBaseUrl+"/reservation/view/daily", "ViewDailyReservationsByAdmin", RoleCookieInjection(rc.viewDailyReservationsByAdmin))
	m.GetJson(kAdminApiBaseUrl+"/reservation/view/teacher/username", "ViewReservationsWithTeacherUsernameByAdmin", RoleCookieInjection(rc.viewReservationsWithTeacherUsernameByAdmin))
	m.GetJson(kAdminApiBaseUrl+"/reservation/export/today", "ExportTodayReservationsByAdmin", RoleCookieInjection(rc.exportTodayReservationsByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/export/report", "ExportReportFormByAdmin", RoleCookieInjection(rc.exportReportFormByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/export/report/monthly", "ExportReportMonthlyByAdmin", RoleCookieInjection(rc.exportReportMonthlyByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/add", "AddReservationByAdmin", RoleCookieInjection(rc.addReservationByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/edit", "EditReservationByAdmin", RoleCookieInjection(rc.editReservationByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/remove", "RemoveReservationsByAdmin", RoleCookieInjection(rc.removeReservationsByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/cancel", "CancelReservationByAdmin", RoleCookieInjection(rc.cancelReservationByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/feedback/get", "GetFeedbackByAdmin", RoleCookieInjection(rc.getFeedbackByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/feedback/submit", "SubmitFeedbackByAdmin", RoleCookieInjection(rc.submitFeedbackByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/student/set", "SetStudentByAdmin", RoleCookieInjection(rc.setStudentByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/get", "GetStudentInfoByAdmin", RoleCookieInjection(rc.getStudentInfoByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/search", "SearchStudentByAdmin", RoleCookieInjection(rc.searchStudentByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/crisis/update", "UpdateStudentCrisisLevelByAdmin", RoleCookieInjection(rc.updateStudentCrisisLevelByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/archive/update", "UpdateStudentArchiveNumberByAdmin", RoleCookieInjection(rc.updateStudentArchiveNumberByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/password/reset", "ResetStudentPasswordByAdmin", RoleCookieInjection(rc.resetStudentPasswordByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/account/delete", "DeleteStudentAccountByAdmin", RoleCookieInjection(rc.deleteStudentAccountByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/export", "ExportStudentByAdmin", RoleCookieInjection(rc.exportStudentByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/unbind", "UnbindStudentByAdmin", RoleCookieInjection(rc.unbindStudentByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/bind", "BindStudentByAdmin", RoleCookieInjection(rc.bindStudentByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/student/query", "QueryStudentInfoByAdmin", RoleCookieInjection(rc.queryStudentInfoByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/teacher/password/reset", "ResetTeacherPasswordByAdmin", RoleCookieInjection(rc.resetTeacherPasswordByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/teacher/search", "SearchTeacherByAdmin", RoleCookieInjection(rc.searchTeacherByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/teacher/workload", "GetTeacherWorkloadByAdmin", RoleCookieInjection(rc.getTeacherWorkloadByAdmin))
}

//==================== category ====================
func (rc *ReservationController) getFeedbackCategories(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	var result = make(map[string]interface{})

	result["first_category"] = model.FeedbackFirstCategory
	result["second_category"] = model.FeedbackSecondCategory

	return http.StatusOK, wrapJsonOk(result)
}

//==================== student ====================
func (rc *ReservationController) viewReservationsByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	var result = make(map[string]interface{})

	student, reservations, err := service.Workflow().GetReservationsByStudent(userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["student"] = service.Workflow().WrapSimpleStudent(student)
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := service.Workflow().WrapSimpleReservation(res)
		if res.StudentId == student.Id.Hex() {
			resJson["student_id"] = res.StudentId
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

func (rc *ReservationController) makeReservationByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	sourceId := form.ParamString(r, "source_id", "")
	startTime := form.ParamString(r, "start_time", "")
	fullname := form.ParamString(r, "student_fullname", "")
	gender := form.ParamString(r, "student_gender", "")
	birthday := form.ParamString(r, "student_birthday", "")
	school := form.ParamString(r, "student_school", "")
	grade := form.ParamString(r, "student_grade", "")
	currentAddress := form.ParamString(r, "student_current_address", "")
	familyAddress := form.ParamString(r, "student_family_address", "")
	mobile := form.ParamString(r, "student_mobile", "")
	email := form.ParamString(r, "student_email", "")
	experienceTime := form.ParamString(r, "student_experience_time", "")
	experienceLocation := form.ParamString(r, "student_experience_location", "")
	experienceTeacher := form.ParamString(r, "student_experience_teacher", "")
	fatherAge := form.ParamString(r, "student_father_age", "")
	fatherJob := form.ParamString(r, "student_father_job", "")
	fatherEdu := form.ParamString(r, "student_father_edu", "")
	motherAge := form.ParamString(r, "student_mother_age", "")
	motherJob := form.ParamString(r, "student_mother_job", "")
	motherEdu := form.ParamString(r, "student_mother_edu", "")
	parentMarriage := form.ParamString(r, "student_parent_marriage", "")
	siginificant := form.ParamString(r, "student_significant", "")
	problem := form.ParamString(r, "student_problem", "")

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().MakeReservationByStudent(reservationId, sourceId, startTime, fullname, gender, birthday,
		school, grade, currentAddress, familyAddress, mobile, email, experienceTime, experienceLocation, experienceTeacher,
		fatherAge, fatherJob, fatherEdu, motherAge, motherJob, motherEdu, parentMarriage, siginificant, problem,
		userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["reservation"] = service.Workflow().WrapSimpleReservation(reservation)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) getFeedbackByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	sourceId := form.ParamString(r, "source_id", "")

	var result = make(map[string]interface{})

	var feedbackJson = make(map[string]interface{})
	reservation, err := service.Workflow().GetFeedbackByStudent(reservationId, sourceId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	feedbackJson["scores"] = reservation.StudentFeedback.Scores
	result["feedback"] = feedbackJson

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) submitFeedbackByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	sourceId := form.ParamString(r, "source_id", "")
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
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

//==================== teacher ====================
func (rc *ReservationController) viewReservationsByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	var result = make(map[string]interface{})

	teacher, reservations, err := service.Workflow().GetReservationsByTeacher(userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["teacher"] = service.Workflow().WrapTeacher(teacher)
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

func (rc *ReservationController) getFeedbackByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	sourceId := form.ParamString(r, "source_id", "")

	var result = make(map[string]interface{})

	student, reservation, err := service.Workflow().GetFeedbackByTeacher(reservationId, sourceId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	feedback := reservation.TeacherFeedback.ToJson()
	feedback["crisis_level"] = student.CrisisLevel
	result["feedback"] = feedback

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) submitFeedbackByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	sourceId := form.ParamString(r, "source_id", "")
	category := form.ParamString(r, "category", "")
	participants := []string(r.Form["participants"])
	participantsInt := make([]int, 0)
	for _, p := range participants {
		if pi, err := strconv.Atoi(p); err == nil {
			participantsInt = append(participantsInt, pi)
		}
	}
	emphasis := form.ParamString(r, "emphasis", "")
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
	record := form.ParamString(r, "record", "")
	crisisLevel := form.ParamString(r, "crisis_level", "")

	var result = make(map[string]interface{})

	_, err := service.Workflow().SubmitFeedbackByTeacher(reservationId, sourceId, category, participantsInt, emphasis, severityInt,
		medicalDiagnosisInt, crisisInt, record, crisisLevel, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) getStudentInfoByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := form.ParamString(r, "student_id", "")

	var result = make(map[string]interface{})

	student, reservations, err := service.Workflow().GetStudentInfoByTeacher(studentId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
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

func (rc *ReservationController) queryStudentInfoByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentUsername := form.ParamString(r, "student_username", "")

	var result = make(map[string]interface{})

	student, reservations, err := service.Workflow().QueryStudentInfoByTeacher(studentUsername, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
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
func (rc *ReservationController) viewReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	var result = make(map[string]interface{})

	admin, reservations, err := service.Workflow().GetReservationsByAdmin(userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["admin"] = service.Workflow().WrapAdmin(admin)
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

func (rc *ReservationController) viewDailyReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	fromDate := form.ParamString(r, "from_date", "")

	var result = make(map[string]interface{})

	admin, reservations, err := service.Workflow().GetReservationsDailyByAdmin(fromDate, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["admin"] = service.Workflow().WrapAdmin(admin)
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

func (rc *ReservationController) viewReservationsWithTeacherUsernameByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	teacherUsername := form.ParamString(r, "teacher_username", "")

	var result = make(map[string]interface{})

	reservations, err := service.Workflow().GetReservationsWithTeacherUsernameByAdmin(teacherUsername, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
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

func (rc *ReservationController) exportTodayReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	var result = make(map[string]interface{})

	url, err := service.Workflow().ExportTodayReservationTimetableByAdmin(userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["url"] = "/" + url

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) addReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	startTime := form.ParamString(r, "start_time", "")
	endTime := form.ParamString(r, "end_time", "")
	teacherUsername := form.ParamString(r, "teacher_username", "")
	teacherFullname := form.ParamString(r, "teacher_fullname", "")
	teacherMobile := form.ParamString(r, "teacher_mobile", "")
	var forceBool bool
	if force := form.ParamString(r, "force", ""); force != "" {
		forceBool, _ = strconv.ParseBool(r.FormValue("force"))
	}

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().AddReservationByAdmin(startTime, endTime, teacherUsername, teacherFullname,
		teacherMobile, forceBool, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["reservation"] = service.Workflow().WrapReservation(reservation)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) editReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	sourceId := form.ParamString(r, "source_id", "")
	originalStartTime := form.ParamString(r, "original_start_time", "")
	startTime := form.ParamString(r, "start_time", "")
	endTime := form.ParamString(r, "end_time", "")
	teacherUsername := form.ParamString(r, "teacher_username", "")
	teacherFullname := form.ParamString(r, "teacher_fullname", "")
	teacherMobile := form.ParamString(r, "teacher_mobile", "")
	var forceBool bool
	if force := form.ParamString(r, "force", ""); force != "" {
		forceBool, _ = strconv.ParseBool(r.FormValue("force"))
	}

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().EditReservationByAdmin(reservationId, sourceId, originalStartTime,
		startTime, endTime, teacherUsername, teacherFullname, teacherMobile, forceBool, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["reservation"] = service.Workflow().WrapReservation(reservation)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) removeReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationIds := []string(r.Form["reservation_ids"])
	sourceIds := []string(r.Form["source_ids"])
	startTimes := []string(r.Form["start_times"])

	var result = make(map[string]interface{})

	removed, err := service.Workflow().RemoveReservationsByAdmin(reservationIds, sourceIds, startTimes, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["removed_count"] = removed

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) cancelReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationIds := []string(r.Form["reservation_ids"])
	sourceIds := []string(r.Form["source_ids"])

	var result = make(map[string]interface{})

	removed, err := service.Workflow().CancelReservationsByAdmin(reservationIds, sourceIds, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["canceled_count"] = removed

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) getFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	sourceId := form.ParamString(r, "source_id", "")

	var result = make(map[string]interface{})

	student, reservation, err := service.Workflow().GetFeedbackByAdmin(reservationId, sourceId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	feedback := reservation.TeacherFeedback.ToJson()
	feedback["crisis_level"] = student.CrisisLevel
	result["feedback"] = feedback

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) submitFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	sourceId := form.ParamString(r, "source_id", "")
	category := form.ParamString(r, "category", "")
	participants := []string(r.Form["participants"])
	participantsInt := make([]int, 0)
	for _, p := range participants {
		if pi, err := strconv.Atoi(p); err == nil {
			participantsInt = append(participantsInt, pi)
		}
	}
	emphasis := form.ParamString(r, "emphasis", "")
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
	record := form.ParamString(r, "record", "")
	crisisLevel := form.ParamString(r, "crisis_level", "")

	var result = make(map[string]interface{})

	_, err := service.Workflow().SubmitFeedbackByAdmin(reservationId, sourceId, category, participantsInt, emphasis, severityInt,
		medicalDiagnosisInt, crisisInt, record, crisisLevel, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) setStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	sourceId := form.ParamString(r, "source_id", "")
	startTime := form.ParamString(r, "start_time", "")
	studentUsername := form.ParamString(r, "student_username", "")
	fullname := form.ParamString(r, "student_fullname", "")
	gender := form.ParamString(r, "student_gender", "")
	birthday := form.ParamString(r, "student_birthday", "")
	school := form.ParamString(r, "student_school", "")
	grade := form.ParamString(r, "student_grade", "")
	currentAddress := form.ParamString(r, "student_current_address", "")
	familyAddress := form.ParamString(r, "student_family_address", "")
	mobile := form.ParamString(r, "student_mobile", "")
	email := form.ParamString(r, "student_email", "")
	experienceTime := form.ParamString(r, "student_experience_time", "")
	experienceLocation := form.ParamString(r, "student_experience_location", "")
	experienceTeacher := form.ParamString(r, "student_experience_teacher", "")
	fatherAge := form.ParamString(r, "student_father_age", "")
	fatherJob := form.ParamString(r, "student_father_job", "")
	fatherEdu := form.ParamString(r, "student_father_edu", "")
	motherAge := form.ParamString(r, "student_mother_age", "")
	motherJob := form.ParamString(r, "student_mother_job", "")
	motherEdu := form.ParamString(r, "student_mother_edu", "")
	parentMarriage := form.ParamString(r, "student_parent_marriage", "")
	siginificant := form.ParamString(r, "student_significant", "")
	problem := form.ParamString(r, "student_problem", "")
	sendSms, err := strconv.ParseBool(r.FormValue("student_sms"))
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().SetStudentByAdmin(reservationId, sourceId, startTime, studentUsername, fullname,
		gender, birthday, school, grade, currentAddress, familyAddress, mobile, email, experienceTime,
		experienceLocation, experienceTeacher, fatherAge, fatherJob, fatherEdu, motherAge, motherJob, motherEdu,
		parentMarriage, siginificant, problem, sendSms, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["reservation"] = service.Workflow().WrapReservation(reservation)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) getStudentInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := form.ParamString(r, "student_id", "")

	var result = make(map[string]interface{})

	student, reservations, err := service.Workflow().GetStudentInfoByAdmin(studentId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
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

func (rc *ReservationController) searchStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentUsername := form.ParamString(r, "student_username", "")

	var result = make(map[string]interface{})

	student, _, err := service.Workflow().QueryStudentInfoByAdmin(studentUsername, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["student"] = service.Workflow().WrapStudent(student)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) updateStudentCrisisLevelByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := form.ParamString(r, "student_id", "")
	crisisLevel := form.ParamString(r, "crisis_level", "")

	var result = make(map[string]interface{})

	_, err := service.Workflow().UpdateStudentCrisisLevelByAdmin(studentId, crisisLevel, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) updateStudentArchiveNumberByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := form.ParamString(r, "student_id", "")
	archiveCategory := form.ParamString(r, "archive_category", "")
	archiveNumber := form.ParamString(r, "archive_number", "")

	var result = make(map[string]interface{})

	_, err := service.Workflow().UpdateStudentArchiveNumberByAdmin(studentId, archiveCategory, archiveNumber, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) resetStudentPasswordByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := form.ParamString(r, "student_id", "")
	password := form.ParamString(r, "password", "")

	var result = make(map[string]interface{})

	_, err := service.Workflow().ResetStudentPasswordByAdmin(studentId, password, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) deleteStudentAccountByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := form.ParamString(r, "student_id", "")

	var result = make(map[string]interface{})

	err := service.Workflow().DeleteStudentAccountByAdmin(studentId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) exportStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := form.ParamString(r, "student_id", "")

	var result = make(map[string]interface{})

	url, err := service.Workflow().ExportStudentByAdmin(studentId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["url"] = "/" + url

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) unbindStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := form.ParamString(r, "student_id", "")

	var result = make(map[string]interface{})

	student, err := service.Workflow().UnbindStudentByAdmin(studentId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["student"] = service.Workflow().WrapStudent(student)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) bindStudentByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentId := form.ParamString(r, "student_id", "")
	teacherUsername := form.ParamString(r, "teacher_username", "")

	var result = make(map[string]interface{})

	student, err := service.Workflow().BindStudentByAdmin(studentId, teacherUsername, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["student"] = service.Workflow().WrapStudent(student)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) queryStudentInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	studentUsername := form.ParamString(r, "student_username", "")

	var result = make(map[string]interface{})

	student, reservations, err := service.Workflow().QueryStudentInfoByAdmin(studentUsername, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
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

func (rc *ReservationController) resetTeacherPasswordByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	teacherUsername := form.ParamString(r, "teacher_username", "")
	teacherFullname := form.ParamString(r, "teacher_fullname", "")
	teacherMobile := form.ParamString(r, "teacher_mobile", "")
	password := form.ParamString(r, "password", "")

	var result = make(map[string]interface{})

	teacher, err := service.Workflow().ResetTeacherPasswordByAdmin(teacherUsername, teacherFullname, teacherMobile, password, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["teacher"] = service.Workflow().WrapTeacher(teacher)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) searchTeacherByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	teacherUsername := form.ParamString(r, "teacher_username", "")
	teacherFullname := form.ParamString(r, "teacher_fullname", "")
	teacherMoble := form.ParamString(r, "teacher_mobile", "")

	var result = make(map[string]interface{})

	teacher, err := service.Workflow().SearchTeacherByAdmin(teacherFullname, teacherUsername, teacherMoble,
		userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["teacher"] = service.Workflow().WrapTeacher(teacher)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) getTeacherWorkloadByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	fromDate := form.ParamString(r, "from_date", "")
	toDate := form.ParamString(r, "to_date", "")

	var result = make(map[string]interface{})

	workload, err := service.Workflow().GetTeacherWorkloadByAdmin(fromDate, toDate, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["workload"] = workload

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) exportReportFormByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	fromDate := form.ParamString(r, "from_date", "")
	toDate := form.ParamString(r, "to_date", "")

	var result = make(map[string]interface{})

	url, err := service.Workflow().ExportReportFormByAdmin(fromDate, toDate, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["url"] = "/" + url

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) exportReportMonthlyByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	monthlyDate := form.ParamString(r, "monthly_date", "")

	var result = make(map[string]interface{})

	reportUrl, keyCaseUrl, err := service.Workflow().ExportReportMonthlyByAdmin(monthlyDate, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["report_url"] = "/" + reportUrl
	result["key_case_url"] = "/" + keyCaseUrl

	return http.StatusOK, wrapJsonOk(result)
}

//==================== timetable ====================
func (rc *ReservationController) viewTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	var result = make(map[string]interface{})

	timedReservations, err := service.Workflow().ViewTimetableByAdmin(userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
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

func (rc *ReservationController) addTimedReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	weekday := form.ParamString(r, "weekday", "")
	startTime := form.ParamString(r, "start_clock", "")
	endTime := form.ParamString(r, "end_clock", "")
	teacherUsername := form.ParamString(r, "teacher_username", "")
	teacherFullname := form.ParamString(r, "teacher_fullname", "")
	teacherMobile := form.ParamString(r, "teacher_mobile", "")
	var forceBool bool
	if force := form.ParamString(r, "force", ""); force != "" {
		forceBool, _ = strconv.ParseBool(r.FormValue("force"))
	}

	var result = make(map[string]interface{})

	timedReservation, err := service.Workflow().AddTimetableByAdmin(weekday, startTime, endTime, teacherUsername, teacherFullname,
		teacherMobile, forceBool, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["timed_reservation"] = service.Workflow().WrapTimedReservation(timedReservation)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) editTimedReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	timedReservationId := form.ParamString(r, "timed_reservation_id", "")
	weekday := form.ParamString(r, "weekday", "")
	startTime := form.ParamString(r, "start_clock", "")
	endTime := form.ParamString(r, "end_clock", "")
	teacherUsername := form.ParamString(r, "teacher_username", "")
	teacherFullname := form.ParamString(r, "teacher_fullname", "")
	teacherMobile := form.ParamString(r, "teacher_mobile", "")
	var forceBool bool
	if force := form.ParamString(r, "force", ""); force != "" {
		forceBool, _ = strconv.ParseBool(r.FormValue("force"))
	}

	var result = make(map[string]interface{})

	timedReservation, err := service.Workflow().EditTimetableByAdmin(timedReservationId, weekday, startTime, endTime, teacherUsername,
		teacherFullname, teacherMobile, forceBool, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["timed_reservation"] = service.Workflow().WrapTimedReservation(timedReservation)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) removeTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	timedReservationIds := []string(r.Form["timed_reservation_ids"])

	var result = make(map[string]interface{})

	removed, err := service.Workflow().RemoveTimetablesByAdmin(timedReservationIds, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["removed_count"] = removed

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) openTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	timedReservationIds := []string(r.Form["timed_reservation_ids"])

	var result = make(map[string]interface{})

	opened, err := service.Workflow().OpenTimetablesByAdmin(timedReservationIds, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["opened_count"] = opened

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) closeTimedReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	timedReservationIds := []string(r.Form["timed_reservation_ids"])

	var result = make(map[string]interface{})

	closed, err := service.Workflow().CloseTimetablesByAdmin(timedReservationIds, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["closed_count"] = closed

	return http.StatusOK, wrapJsonOk(result)
}
