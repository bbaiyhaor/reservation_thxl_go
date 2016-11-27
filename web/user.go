package web

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/service"
	"github.com/mijia/sweb/form"
	"github.com/mijia/sweb/render"
	"golang.org/x/net/context"
	"net/http"
)

type UserController struct {
	BaseMuxController
}

const (
	kUserApiBaseUrl = "/api/user"
)

func (uc *UserController) MuxHandlers(m JsonMuxer) {
	m.Get("/m", "EntryPage", uc.getEntryPage)
	m.Get("/m/student", "StudentPage", uc.getStudentPage)
	m.Get("/m/teacher", "TeacherPage", uc.getTeacherPage)
	// legacy
	m.Get("/reservation/admin/login", "AdminLoginPage", uc.getAdminLoginPageLegacy)
	m.Get("/reservation/admin", "AdminPage", LegacyAdminPageInjection(uc.getAdminPageLegacy))
	m.Get("/reservation/admin/timetable", "AdminTimetablePage", LegacyAdminPageInjection(uc.getAdminTimetablePageLegacy))

	m.PostJson(kUserApiBaseUrl+"/student/login", "StudentLogin", uc.studentLogin)
	m.PostJson(kUserApiBaseUrl+"/student/register", "StudentRegister", uc.studentRegister)
	m.PostJson(kUserApiBaseUrl+"/teacher/login", "TeacherLogin", uc.teacherLogin)
	m.PostJson(kUserApiBaseUrl+"/teacher/password/change", "TeacherChangePassword", RoleCookieInjection(uc.teacherChangePassword))
	m.PostJson(kUserApiBaseUrl+"/teacher/password/reset/sms", "TeacherResetPasswordSms", uc.teacherResetPasswordSms)
	m.PostJson(kUserApiBaseUrl+"/teacher/password/reset/verify", "TeacherResetPasswordVerify", uc.teacherResetPasswordVerify)
	m.PostJson(kUserApiBaseUrl+"/admin/login", "AdminLogin", uc.adminLogin)
	m.GetJson(kUserApiBaseUrl+"/logout", "Logout", RoleCookieInjection(uc.logout))
	m.GetJson(kUserApiBaseUrl+"/session", "UpdateSession", RoleCookieInjection(uc.updateSession))
}

func (uc *UserController) GetTemplates() []*render.TemplateSet {
	return []*render.TemplateSet{
		render.NewTemplateSet("entry", "desktop.html", "reservation/entry.html", "layout/desktop.html"),
		render.NewTemplateSet("student", "desktop.html", "reservation/student.html", "layout/desktop.html"),
		render.NewTemplateSet("teacher", "desktop.html", "reservation/teacher.html", "layout/desktop.html"),
		render.NewTemplateSet("admin_login", "desktop.html", "legacy/admin_login.html", "layout/desktop.html"),
		render.NewTemplateSet("admin", "desktop.html", "legacy/admin.html", "layout/desktop.html"),
		render.NewTemplateSet("admin_timetable", "desktop.html", "legacy/admin_timetable.html", "layout/desktop.html"),
	}
}

func (uc *UserController) getEntryPage(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "entry", params)
	return ctx
}

func (uc *UserController) getStudentPage(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "student", params)
	return ctx
}

func (uc *UserController) getTeacherPage(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "teacher", params)
	return ctx
}

func (uc *UserController) getAdminLoginPageLegacy(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "admin_login", params)
	return ctx
}

func (uc *UserController) getAdminPageLegacy(ctx context.Context, w http.ResponseWriter, r *http.Request, userId string, userType int) context.Context {
	if userType != model.USER_TYPE_ADMIN {
		http.Redirect(w, r, "/reservation/admin/login", http.StatusFound)
		return ctx
	} else if _, err := service.MongoClient().GetAdminById(userId); err != nil {
		http.Redirect(w, r, "/reservation/admin/login", http.StatusFound)
		return ctx
	}
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "admin", params)
	return ctx
}

func (uc *UserController) getAdminTimetablePageLegacy(ctx context.Context, w http.ResponseWriter, r *http.Request, userId string, userType int) context.Context {
	if userType != model.USER_TYPE_ADMIN {
		http.Redirect(w, r, "/reservation/admin/login", http.StatusFound)
		return ctx
	} else if _, err := service.MongoClient().GetAdminById(userId); err != nil {
		http.Redirect(w, r, "/reservation/admin/login", http.StatusFound)
		return ctx
	}
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "admin_timetable", params)
	return ctx
}

func (uc *UserController) studentRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	username := form.ParamString(r, "username", "")
	password := form.ParamString(r, "password", "")

	var result = make(map[string]interface{})

	student, err := service.Workflow().StudentRegister(username, password)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	if err = setUserCookie(w, student.Id.Hex(), student.Username, student.UserType); err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["user_id"] = student.Id.Hex()
	result["username"] = student.Username
	result["user_type"] = student.UserType
	result["fullname"] = student.Fullname

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) studentLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	username := form.ParamString(r, "username", "")
	password := form.ParamString(r, "password", "")

	var result = make(map[string]interface{})

	student, err := service.Workflow().StudentLogin(username, password)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	if err = setUserCookie(w, student.Id.Hex(), student.Username, student.UserType); err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["user_id"] = student.Id.Hex()
	result["username"] = student.Username
	result["user_type"] = student.UserType
	result["fullname"] = student.Fullname

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) teacherLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	username := form.ParamString(r, "username", "")
	password := form.ParamString(r, "password", "")

	var result = make(map[string]interface{})

	teacher, err := service.Workflow().TeacherLogin(username, password)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	if err = setUserCookie(w, teacher.Id.Hex(), teacher.Username, teacher.UserType); err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["user_id"] = teacher.Id.Hex()
	result["username"] = teacher.Username
	result["user_type"] = teacher.UserType
	result["fullname"] = teacher.Fullname

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) teacherChangePassword(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	username := form.ParamString(r, "username", "")
	oldPassword := form.ParamString(r, "old_password", "")
	newPassword := form.ParamString(r, "new_password", "")

	var result = make(map[string]interface{})

	teacher, err := service.Workflow().TeacherChangePassword(username, oldPassword, newPassword, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["teacher"] = service.Workflow().WrapTeacher(teacher)

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) teacherResetPasswordSms(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	username := form.ParamString(r, "username", "")
	mobile := form.ParamString(r, "mobile", "")

	var result = make(map[string]interface{})

	err := service.Workflow().TeacherResetPasswordSms(username, mobile)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) teacherResetPasswordVerify(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	username := form.ParamString(r, "username", "")
	newPassword := form.ParamString(r, "new_password", "")
	verifyCode := form.ParamString(r, "verify_code", "")

	var result = make(map[string]interface{})

	err := service.Workflow().TeacherRestPasswordVerify(username, newPassword, verifyCode)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) adminLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	username := form.ParamString(r, "username", "")
	password := form.ParamString(r, "password", "")

	var result = make(map[string]interface{})

	admin, err := service.Workflow().AdminLogin(username, password)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	if err = setUserCookie(w, admin.Id.Hex(), admin.Username, admin.UserType); err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["user_id"] = admin.Id.Hex()
	result["username"] = admin.Username
	result["user_type"] = admin.UserType
	result["redirect_url"] = "/reservation/admin"

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) logout(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	var result = make(map[string]interface{})

	switch userType {
	case model.USER_TYPE_ADMIN:
		result["redirect_url"] = "/reservation/admin"
	case model.USER_TYPE_TEACHER:
		result["redirect_url"] = "/m/teacher#/login"
	case model.USER_TYPE_STUDENT:
		result["redirect_url"] = "/m/student#/login"
	default:
		result["redirect_url"] = "/m"
	}
	clearUserCookie(w)

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) updateSession(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	result, err := service.Workflow().UpdateSession(userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	return http.StatusOK, wrapJsonOk(result)
}
