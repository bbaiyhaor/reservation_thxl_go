package web

import (
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
