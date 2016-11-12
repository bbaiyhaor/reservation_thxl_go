package web

import (
	"github.com/mijia/sweb/server"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
)

func RoleCookieInjection(handle func(http.ResponseWriter, *http.Request, string, int) (int, interface{})) JsonHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
		var userId string
		var userType int
		if cookie, err := r.Cookie("user_id"); err != nil {
			return http.StatusOK, wrapJsonError("请先登录")
		} else {
			userId = cookie.Value
		}
		if _, err := r.Cookie("username"); err != nil {
			return http.StatusOK, wrapJsonError("请先登录")
		}
		if cookie, err := r.Cookie("user_type"); err != nil {
			return http.StatusOK, wrapJsonError("请先登录")
		} else {
			if ut, err := strconv.Atoi(cookie.Value); err != nil {
				return http.StatusOK, wrapJsonError("请先登录")
			} else {
				userType = ut
			}
		}
		return handle(w, r, userId, userType)
	}
}

func LegacyAdminPageInjection(handle func(context.Context, http.ResponseWriter, *http.Request, string, int) context.Context) server.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		redirectUrl := "/reservation/admin/login"
		var userId string
		var userType int
		if cookie, err := r.Cookie("user_id"); err != nil {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
		} else {
			userId = cookie.Value
		}
		if _, err := r.Cookie("username"); err != nil {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
		}
		if cookie, err := r.Cookie("user_type"); err != nil {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
		} else if ut, err := strconv.Atoi(cookie.Value); err != nil {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
		} else {
			userType = ut
		}
		return handle(ctx, w, r, userId, userType)
	}
}
