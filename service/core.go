package service

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

func (s *Service) EntryPage(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	t := template.Must(template.ParseFiles("../templates/entry.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (s *Service) StudentLoginPage(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	t := template.Must(template.ParseFiles("../templates/student_login.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (s *Service) StudentRegisterPage(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	t := template.Must(template.ParseFiles("../templates/student_register.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (s *Service) StudentPage(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	if userType == model.TEACHER {
		http.Redirect(w, r, "/reservation/teacher", http.StatusFound)
		return nil
	} else if userType == model.ADMIN {
		http.Redirect(w, r, "/reservation/admin", http.StatusFound)
		return nil
	}
	t := template.Must(template.ParseFiles("../templates/student.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (s *Service) TeacherLoginPage(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	t := template.Must(template.ParseFiles("../templates/teacher_login.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (s *Service) TeacherPage(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	if userType == model.STUDENT {
		http.Redirect(w, r, "/reservation/student", http.StatusFound)
		return nil
	} else if userType == model.ADMIN {
		http.Redirect(w, r, "/reservation/admin", http.StatusFound)
		return nil
	}
	t := template.Must(template.ParseFiles("../templates/teacher.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (s *Service) AdminLoginPage(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	t := template.Must(template.ParseFiles("../templates/admin_login.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (s *Service) AdminPage(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	if userType == model.STUDENT {
		http.Redirect(w, r, "/reservation/student", http.StatusFound)
		return nil
	} else if userType == model.TEACHER {
		http.Redirect(w, r, "/reservation/teacher", http.StatusFound)
		return nil
	}
	t := template.Must(template.ParseFiles("../templates/admin.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (s *Service) AdminTimetablePage(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	if userType == model.STUDENT {
		http.Redirect(w, r, "/reservation/student", http.StatusFound)
		return nil
	} else if userType == model.TEACHER {
		http.Redirect(w, r, "/reservation/teacher", http.StatusFound)
		return nil
	}
	t := template.Must(template.ParseFiles("../templates/admin_timetable.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (s *Service) ProtocolPage(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	t := template.Must(template.ParseFiles("../templates/protocol.html"))
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

func (s *Service) ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	state := "FAILED"
	if err.Error() == model.CHECK_MESSAGE {
		state = model.CHECK_MESSAGE
	} else {
		log.Printf("error %s %v", r.URL.Path, err)
	}
	if data, err := json.Marshal(ErrorMsg{
		State:   state,
		Message: err.Error(),
	}); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}
