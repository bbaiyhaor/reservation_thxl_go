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
	m.Get("/m", "EntryPage", uc.GetEntryPage)
	m.Get("/m/student", "StudentPage", uc.GetStudentPage)
	m.Get("/m/teacher", "TeacherPage", uc.GetTeacherPage)
	// legacy
	m.Get("/reservation", "LegacyEntryPage", uc.GetEntryPage)
	m.Get("/reservation/student", "LegacyStudentPage", uc.GetStudentPage)
	m.Get("/reservation/teacher", "LegacyTeacherPage", uc.GetTeacherPage)
	m.Get("/reservation/admin/login", "AdminLoginPage", uc.GetAdminLoginPageLegacy)
	m.Get("/reservation/admin", "AdminPage", LegacyAdminPageInjection(uc.GetAdminPageLegacy))
	m.Get("/reservation/admin/timetable", "AdminTimetablePage", LegacyAdminPageInjection(uc.GetAdminTimetablePageLegacy))

	m.PostJson(kUserApiBaseUrl+"/student/login", "StudentLogin", uc.StudentLogin)
	m.PostJson(kUserApiBaseUrl+"/student/register", "StudentRegister", uc.StudentRegister)
	m.PostJson(kUserApiBaseUrl+"/teacher/login", "TeacherLogin", uc.TeacherLogin)
	m.PostJson(kUserApiBaseUrl+"/teacher/password/change", "TeacherChangePassword", RoleCookieInjection(uc.TeacherChangePassword))
	m.PostJson(kUserApiBaseUrl+"/teacher/password/reset/sms", "TeacherResetPasswordSms", uc.TeacherResetPasswordSms)
	m.PostJson(kUserApiBaseUrl+"/teacher/password/reset/verify", "TeacherResetPasswordVerify", uc.TeacherResetPasswordVerify)
	m.PostJson(kUserApiBaseUrl+"/admin/login", "AdminLogin", uc.AdminLogin)
	m.PostJson(kUserApiBaseUrl+"/admin/password/change", "AdminChangePassword", RoleCookieInjection(uc.AdminChangePassword))
	m.GetJson(kUserApiBaseUrl+"/logout", "Logout", RoleCookieInjection(uc.Logout))
	m.GetJson(kUserApiBaseUrl+"/session", "UpdateSession", RoleCookieInjection(uc.UpdateSession))
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

func (uc *UserController) GetEntryPage(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "entry", params)
	return ctx
}

func (uc *UserController) GetStudentPage(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "student", params)
	return ctx
}

func (uc *UserController) GetTeacherPage(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "teacher", params)
	return ctx
}

func (uc *UserController) GetAdminLoginPageLegacy(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "admin_login", params)
	return ctx
}

func (uc *UserController) GetAdminPageLegacy(ctx context.Context, w http.ResponseWriter, r *http.Request, userId string, userType int) context.Context {
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

func (uc *UserController) GetAdminTimetablePageLegacy(ctx context.Context, w http.ResponseWriter, r *http.Request, userId string, userType int) context.Context {
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

func (uc *UserController) StudentRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	username := form.ParamString(r, "username", "")
	password := form.ParamString(r, "password", "")

	var result = make(map[string]interface{})

	student, err := service.Workflow().StudentRegister(username, password)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["user"] = service.Workflow().WrapSimpleStudent(student)

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) StudentLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	username := form.ParamString(r, "username", "")
	password := form.ParamString(r, "password", "")

	var result = make(map[string]interface{})

	student, err := service.Workflow().StudentLogin(username, password)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	if err = setSession(w, student.Id.Hex(), student.Username, student.UserType); err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["user"] = service.Workflow().WrapSimpleStudent(student)

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) TeacherLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	username := form.ParamString(r, "username", "")
	password := form.ParamString(r, "password", "")

	var result = make(map[string]interface{})

	teacher, err := service.Workflow().TeacherLogin(username, password)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	if err = setSession(w, teacher.Id.Hex(), teacher.Username, teacher.UserType); err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["user"] = service.Workflow().WrapTeacher(teacher)

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) TeacherChangePassword(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
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

func (uc *UserController) TeacherResetPasswordSms(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	username := form.ParamString(r, "username", "")
	fullname := form.ParamString(r, "fullname", "")
	mobile := form.ParamString(r, "mobile", "")

	var result = make(map[string]interface{})

	err := service.Workflow().TeacherResetPasswordSms(username, fullname, mobile)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) TeacherResetPasswordVerify(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
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

func (uc *UserController) AdminLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	username := form.ParamString(r, "username", "")
	password := form.ParamString(r, "password", "")

	var result = make(map[string]interface{})

	admin, err := service.Workflow().AdminLogin(username, password)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	if err = setSession(w, admin.Id.Hex(), admin.Username, admin.UserType); err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["user"] = service.Workflow().WrapAdmin(admin)
	result["redirect_url"] = "/reservation/admin"

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) AdminChangePassword(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	username := form.ParamString(r, "username", "")
	oldPassword := form.ParamString(r, "old_password", "")
	newPassword := form.ParamString(r, "new_password", "")

	var result = make(map[string]interface{})

	_, err := service.Workflow().AdminChangePassword(username, oldPassword, newPassword, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
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
	clearSession(w, r)

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) UpdateSession(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	result, err := service.Workflow().UpdateSession(userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	return http.StatusOK, wrapJsonOk(result)
}
