package web

import (
	"golang.org/x/net/context"
	"net/http"
	"regexp"
	"strconv"
)

var redirectStudentPath = regexp.MustCompile("^(/api/student|/api/user/student)")
var redirectTeacherPath = regexp.MustCompile("^(/api/teacher|/api/user/teacher)")
var redirectAdminPath = regexp.MustCompile("^(/api/admin|/api/user/admin)")

func RoleCookieInjection(handle func(http.ResponseWriter, *http.Request, string, int) (int, interface{})) JsonHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
		redirectUrl := "/m"
		if redirectStudentPath.MatchString(r.URL.Path) {
			redirectUrl = "/m/student"
		} else if redirectTeacherPath.MatchString(r.URL.Path) {
			redirectUrl = "/m/teacher"
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
