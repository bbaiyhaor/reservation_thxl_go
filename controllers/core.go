package controllers

import (
	"encoding/json"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"html/template"
	"net/http"
	"strings"
)

func EntryPage(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	t := template.Must(template.ParseFiles("templates/entry.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func StudentLoginPage(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	t := template.Must(template.ParseFiles("templates/student_login.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func StudentRegisterPage(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	t := template.Must(template.ParseFiles("templates/student_register.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func StudentPage(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	if userType == models.TEACHER {
		http.Redirect(w, r, "/reservation/teacher", http.StatusFound)
		return nil
	} else if userType == models.ADMIN {
		http.Redirect(w, r, "/reservation/admin", http.StatusFound)
		return nil
	}
	t := template.Must(template.ParseFiles("templates/student.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func TeacherLoginPage(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	t := template.Must(template.ParseFiles("templates/teacher_login.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func TeacherPage(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	if userType == models.STUDENT {
		http.Redirect(w, r, "/reservation/student", http.StatusFound)
		return nil
	} else if userType == models.ADMIN {
		http.Redirect(w, r, "/reservation/admin", http.StatusFound)
		return nil
	}
	t := template.Must(template.ParseFiles("templates/teacher.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func AdminLoginPage(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	t := template.Must(template.ParseFiles("templates/admin_login.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func AdminPage(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	if userType == models.STUDENT {
		http.Redirect(w, r, "/reservation/student", http.StatusFound)
		return nil
	} else if userType == models.TEACHER {
		http.Redirect(w, r, "/reservation/teacher", http.StatusFound)
		return nil
	}
	t := template.Must(template.ParseFiles("templates/admin.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func AdminTimetablePage(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	if userType == models.STUDENT {
		http.Redirect(w, r, "/reservation/student", http.StatusFound)
		return nil
	} else if userType == models.TEACHER {
		http.Redirect(w, r, "/reservation/teacher", http.StatusFound)
		return nil
	}
	t := template.Must(template.ParseFiles("templates/admin_timetable.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

type ErrorMsg struct {
	State   string `json:"state"`
	Message string `json:"message"`
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	state := "FAILED"
	if strings.EqualFold(err.Error(), models.CHECK_MESSAGE) {
		state = models.CHECK_MESSAGE
	}
	if data, err := json.Marshal(ErrorMsg{
		State:   state,
		Message: err.Error(),
	}); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}