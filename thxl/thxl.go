package main

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/config"
	tsvc "bitbucket.org/shudiwsh2009/reservation_thxl_go/service"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var needUserPath = regexp.MustCompile("^(/reservation/(student|teacher|admin)$|/(user/logout|(student|teacher|admin)/))")
var redirectStudentPath = regexp.MustCompile("^(/reservation/student$|/student/)")
var redirectTeacherPath = regexp.MustCompile("^(/reservation/teacher$|/teacher/)")
var redirectAdminPath = regexp.MustCompile("^(/reservation/admin|/admin/)")

func handleWithCookie(fn func(http.ResponseWriter, *http.Request, string, int) interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		// request log
		case http.MethodGet:
			if r.URL.RawQuery != "" {
				log.Printf("%s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
			} else {
				log.Printf("%s %s", r.Method, r.URL.Path)
			}
		case http.MethodPost:
			r.ParseForm()
			log.Printf("%s %s %+v", r.Method, r.URL.Path, r.PostForm)
		}
		if !needUserPath.MatchString(r.URL.Path) {
			if result := fn(w, r, "", 0); result != nil {
				// non-get response
				if r.Method != http.MethodGet {
					log.Printf("response: %s %+v", r.URL.Path, result)
				}
				if data, err := json.Marshal(result); err == nil {
					w.Header().Set("Content-Type", "application/json;charset=UTF-8")
					w.Write(data)
				} else {
					log.Printf("%v", err)
				}
			}
			return
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
			return
		} else {
			userId = cookie.Value
		}
		if _, err := r.Cookie("username"); err != nil {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
			return
		}
		if cookie, err := r.Cookie("user_type"); err != nil {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
			return
		} else {
			if ut, err := strconv.Atoi(cookie.Value); err != nil {
				http.Redirect(w, r, redirectUrl, http.StatusFound)
				return
			} else {
				userType = ut
			}
		}
		if result := fn(w, r, userId, userType); result != nil {
			// non-get authorized response
			if r.Method != http.MethodGet {
				log.Printf("response userId %s userType %d: %s %+v", userId, userType, r.URL.Path, result)
			}
			if data, err := json.Marshal(result); err == nil {
				w.Header().Set("Content-Type", "application/json;charset=UTF-8")
				w.Write(data)
			} else {
				log.Printf("%v", err)
			}
		}
	}
}

func main() {
	conf := flag.String("conf", "../deploy/thxl.conf", "conf file path")
	isSmock := flag.Bool("smock", true, "is smock server")
	flag.Parse()

	config.InitWithParams(*conf, *isSmock)
	log.Printf("config loaded: %+v", *config.Instance())
	service := tsvc.NewService()

	// TODO: Remove the following test codes

	// mux
	router := mux.NewRouter()
	// 加载页面处理器
	pageRouter := router.PathPrefix("/reservation").Methods("GET").Subrouter()
	pageRouter.HandleFunc("/", handleWithCookie(service.EntryPage))
	pageRouter.HandleFunc("/entry", handleWithCookie(service.EntryPage))
	pageRouter.HandleFunc("/student", handleWithCookie(service.StudentPage))
	pageRouter.HandleFunc("/student/login", handleWithCookie(service.StudentLoginPage))
	pageRouter.HandleFunc("/student/register", handleWithCookie(service.StudentRegisterPage))
	pageRouter.HandleFunc("/teacher", handleWithCookie(service.TeacherPage))
	pageRouter.HandleFunc("/teacher/login", handleWithCookie(service.TeacherLoginPage))
	pageRouter.HandleFunc("/admin", handleWithCookie(service.AdminPage))
	pageRouter.HandleFunc("/admin/login", handleWithCookie(service.AdminLoginPage))
	pageRouter.HandleFunc("/admin/timetable", handleWithCookie(service.AdminTimetablePage))
	pageRouter.HandleFunc("/protocol", handleWithCookie(service.ProtocolPage))
	// http加载处理器
	http.Handle("/", router)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets/"))))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
