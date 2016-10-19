package web

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/ifs"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/service"
	"github.com/mijia/sweb/render"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"time"
)

type UserController struct {
	ifs.BaseMux
}

const (
	kUserApiBaseUrl = "/api/user"
)

func (uc *UserController) MuxHandlers(s ifs.Muxer) {
	s.Get("/m", "HomePage", uc.GetHomePage)

	baseUrl := kUserApiBaseUrl
	s.PostJson(baseUrl+"/student/login", "StudentLogin", uc.StudentLogin)
	s.PostJson(baseUrl+"/student/register", "StudentRegister", uc.StudentRegister)
	s.PostJson(baseUrl+"/teacher/login", "TeacherLogin", uc.TeacherLogin)
	s.PostJson(baseUrl+"/admin/login", "AdminLogin", uc.AdminLogin)
	s.GetJson(baseUrl+"/logout", "Logout", RoleCookieInjection(uc.Logout))
}

func (uc *UserController) GetTemplates() []*render.TemplateSet {
	return []*render.TemplateSet{
		render.NewTemplateSet("user_mobile", "desktop.html", "mobile/user_mobile.html", "layout/desktop.html"),
	}
}

func (uc *UserController) GetHomePage(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "user_mobile", params)
	return ctx
}

func (uc *UserController) StudentRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) interface{} {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	var result = map[string]interface{}{"status": "SUCCESS"}

	student, err := service.Workflow().StudentRegister(username, password)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    student.Id.Hex(),
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    student.Username,
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "user_type",
		Value:    strconv.Itoa(student.UserType),
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	result["url"] = "/reservation/student"

	return result
}

func (uc *UserController) StudentLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) interface{} {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	var result = map[string]interface{}{"status": "SUCCESS"}

	student, err := service.Workflow().StudentLogin(username, password)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    student.Id.Hex(),
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    student.Username,
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "user_type",
		Value:    strconv.Itoa(student.UserType),
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	result["url"] = "/reservation/student"

	return result
}

func (uc *UserController) TeacherLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) interface{} {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	var result = map[string]interface{}{"status": "SUCCESS"}

	teacher, err := service.Workflow().TeacherLogin(username, password)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    teacher.Id.Hex(),
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    teacher.Username,
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "user_type",
		Value:    strconv.Itoa(teacher.UserType),
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	switch teacher.UserType {
	case model.USER_TYPE_TEACHER:
		result["url"] = "/reservation/teacher"
	default:
		result["url"] = "/reservation/entry"
	}

	return result
}

func (uc *UserController) AdminLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) interface{} {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	var result = map[string]interface{}{"status": "SUCCESS"}

	admin, err := service.Workflow().AdminLogin(username, password)
	if err != nil {
		result["status"] = "FAIL"
		result["error"] = err.Error()
		return result
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    admin.Id.Hex(),
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    admin.Username,
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "user_type",
		Value:    strconv.Itoa(admin.UserType),
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	switch admin.UserType {
	case model.USER_TYPE_ADMIN:
		result["url"] = "/reservation/admin"
	default:
		result["url"] = "/reservation/entry"
	}

	return result
}

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request, userId string, userType int) interface{} {
	var result = map[string]interface{}{"status": "SUCCESS"}

	switch userType {
	case model.USER_TYPE_ADMIN:
		result["url"] = "/reservation/admin"
	case model.USER_TYPE_TEACHER:
		result["url"] = "/reservation/teacher"
	case model.USER_TYPE_STUDENT:
		result["url"] = "/reservation/student"
	default:
		result["url"] = "/reservation/entry"
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "user_id",
		Path:   "/",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "username",
		Path:   "/",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "user_type",
		Path:   "/",
		MaxAge: -1,
	})

	return result
}
