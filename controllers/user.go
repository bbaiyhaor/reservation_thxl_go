package controllers

import (
	"github.com/shudiwsh2009/reservation_thxl_go/buslogic"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"net/http"
	"time"
)

func StudentRegister(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var ul = buslogic.UserLogic{}

	student, err := ul.StudentRegister(username, password)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Value:   student.Id.Hex(),
		Path:    "/",
		Expires: time.Now().Local().AddDate(1, 0, 0),
		MaxAge:  365 * 24 * 60 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "username",
		Value:   student.Username,
		Path:    "/",
		Expires: time.Now().Local().AddDate(1, 0, 0),
		MaxAge:  365 * 24 * 60 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "user_type",
		Value:   student.UserType.IntStr(),
		Path:    "/",
		Expires: time.Now().Local().AddDate(1, 0, 0),
		MaxAge:  365 * 24 * 60 * 60,
	})
	result["url"] = "/reservation/student"

	return result
}

func StudentLogin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var ul = buslogic.UserLogic{}

	student, err := ul.StudentLogin(username, password)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Value:   student.Id.Hex(),
		Path:    "/",
		Expires: time.Now().Local().AddDate(1, 0, 0),
		MaxAge:  365 * 24 * 60 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "username",
		Value:   student.Username,
		Path:    "/",
		Expires: time.Now().Local().AddDate(1, 0, 0),
		MaxAge:  365 * 24 * 60 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "user_type",
		Value:   student.UserType.IntStr(),
		Path:    "/",
		Expires: time.Now().Local().AddDate(1, 0, 0),
		MaxAge:  365 * 24 * 60 * 60,
	})
	result["url"] = "/reservation/student"

	return result
}

func TeacherLogin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var ul = buslogic.UserLogic{}

	teacher, err := ul.TeacherLogin(username, password)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Value:   teacher.Id.Hex(),
		Path:    "/",
		Expires: time.Now().Local().AddDate(1, 0, 0),
		MaxAge:  365 * 24 * 60 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "username",
		Value:   teacher.Username,
		Path:    "/",
		Expires: time.Now().Local().AddDate(1, 0, 0),
		MaxAge:  365 * 24 * 60 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "user_type",
		Value:   teacher.UserType.IntStr(),
		Path:    "/",
		Expires: time.Now().Local().AddDate(1, 0, 0),
		MaxAge:  365 * 24 * 60 * 60,
	})
	switch teacher.UserType {
	case models.TEACHER:
		result["url"] = "/reservation/teacher"
	default:
		result["url"] = "/reservation/entry"
	}

	return result
}

func AdminLogin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var ul = buslogic.UserLogic{}

	admin, err := ul.AdminLogin(username, password)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Value:   admin.Id.Hex(),
		Path:    "/",
		Expires: time.Now().Local().AddDate(1, 0, 0),
		MaxAge:  365 * 24 * 60 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "username",
		Value:   admin.Username,
		Path:    "/",
		Expires: time.Now().Local().AddDate(1, 0, 0),
		MaxAge:  365 * 24 * 60 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "user_type",
		Value:   admin.UserType.IntStr(),
		Path:    "/",
		Expires: time.Now().Local().AddDate(1, 0, 0),
		MaxAge:  365 * 24 * 60 * 60,
	})
	switch admin.UserType {
	case models.ADMIN:
		result["url"] = "/reservation/admin"
	default:
		result["url"] = "/reservation/entry"
	}

	return result
}

func Logout(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	var result = map[string]interface{}{"state": "SUCCESS"}

	switch userType {
	case models.ADMIN:
		result["url"] = "/reservation/admin"
	case models.TEACHER:
		result["url"] = "/reservation/teacher"
	case models.STUDENT:
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
