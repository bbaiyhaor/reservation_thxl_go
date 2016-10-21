package web

import (
	"golang.org/x/net/context"
	"net/http"
	"regexp"
	"strconv"
)

var needUserPath = regexp.MustCompile("^(/reservation/(student|teacher|admin)$|/(user/logout|(student|teacher|admin)/))")
var redirectStudentPath = regexp.MustCompile("^(/reservation/student$|/student/)")
var redirectTeacherPath = regexp.MustCompile("^(/reservation/teacher$|/teacher/)")
var redirectAdminPath = regexp.MustCompile("^(/reservation/admin|/admin/)")

func RoleCookieInjection(handle func(http.ResponseWriter, *http.Request, string, int) (int, interface{})) JsonHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
		if !needUserPath.MatchString(r.URL.Path) {
			return handle(w, r, "", 0)
		}
		redirectUrl := "/reservation/entry"
		if redirectStudentPath.MatchString(r.URL.Path) {
			redirectUrl = "/reservation/student/login"
		} else if redirectTeacherPath.MatchString(r.URL.Path) {
			redirectUrl = "/reservation/teacher/login"
		} else if redirectAdminPath.MatchString(r.URL.Path) {
			redirectUrl = "/reservation/admin/login"
		}
		var userId string
		var userType int
		if cookie, err := r.Cookie("user_id"); err != nil {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
			return http.StatusOK, wrapJsonError(err.Error())
		} else {
			userId = cookie.Value
		}
		if _, err := r.Cookie("username"); err != nil {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
			return http.StatusOK, wrapJsonError(err.Error())
		}
		if cookie, err := r.Cookie("user_type"); err != nil {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
			return http.StatusOK, wrapJsonError(err.Error())
		} else {
			if ut, err := strconv.Atoi(cookie.Value); err != nil {
				http.Redirect(w, r, redirectUrl, http.StatusFound)
				return http.StatusOK, wrapJsonError(err.Error())
			} else {
				userType = ut
			}
		}
		return handle(w, r, userId, userType)
	}
}
