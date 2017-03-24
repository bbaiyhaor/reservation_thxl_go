package web

import (
	"github.com/mijia/sweb/server"
	"golang.org/x/net/context"
	"net/http"
	"github.com/mijia/sweb/form"
	"github.com/shudiwsh2009/reservation_thxl_go/rerror"
	"github.com/shudiwsh2009/reservation_thxl_go/service"
	"github.com/shudiwsh2009/reservation_thxl_go/model"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
)

func RoleCookieInjection(handle func(http.ResponseWriter, *http.Request, string, int) (int, interface{})) JsonHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
		userId, _, userType, err := getSession(r)
		if err != nil {
			return http.StatusOK, wrapJsonError(err)
		}
		return handle(w, r, userId, userType)
	}
}

func LegacyAdminPageInjection(handle func(context.Context, http.ResponseWriter, *http.Request, string, int) context.Context) server.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		redirectUrl := "/reservation/admin/login"
		userId, _, userType, err := getSession(r)
		if err != nil {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
		}
		return handle(ctx, w, r, userId, userType)
	}
}

func RequestPasswordCheck(r *http.Request, userId string, userType int) error {
	password := form.ParamString(r, "password", "")
	if password == "" {
		return rerror.NewRErrorCodeContext("password is empty", nil, rerror.ERROR_MISSING_PARAM, "password")
	}
	switch userType {
	case model.USER_TYPE_STUDENT:
		student, err := service.MongoClient().GetStudentById(userId)
		if err == nil {
			if (student.Salt == "" && utils.ValidatePassword(password, student.EncryptedPassword)) ||
				(student.Salt != "" && student.Password == model.EncodePassword(student.Salt, password)) {
				return nil
			}
		}
	case model.USER_TYPE_TEACHER:
		teacher, err := service.MongoClient().GetTeacherById(userId)
		if err == nil {
			if (teacher.Salt == "" && utils.ValidatePassword(password, teacher.EncryptedPassword)) ||
				(teacher.Salt != "" && teacher.Password == model.EncodePassword(teacher.Salt, password)) {
				return nil
			}
		}
	case model.USER_TYPE_ADMIN:
		admin, err := service.MongoClient().GetAdminById(userId)
		if err == nil {
			if (admin.Salt == "" && utils.ValidatePassword(password, admin.EncryptedPassword)) ||
				(admin.Salt != "" && admin.Password == model.EncodePassword(admin.Salt, password)) {
				return nil
			}
		}
	}
	return rerror.NewRErrorCode("request password check failed", nil, rerror.ERROR_NOT_AUTHORIZED)
}