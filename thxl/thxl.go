package main

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/config"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
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

func handleWithCookie(fn func(http.ResponseWriter, *http.Request, string, model.UserType) interface{}) http.HandlerFunc {
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
		var userType model.UserType
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
				userType = model.UserType(ut)
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
	// 加载动态处理器
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/student/login", handleWithCookie(service.StudentLogin)).Methods("POST")
	userRouter.HandleFunc("/student/register", handleWithCookie(service.StudentRegister)).Methods("POST")
	userRouter.HandleFunc("/teacher/login", handleWithCookie(service.TeacherLogin)).Methods("POST")
	userRouter.HandleFunc("/admin/login", handleWithCookie(service.AdminLogin)).Methods("POST")
	userRouter.HandleFunc("/logout", handleWithCookie(service.Logout)).Methods("GET")
	studentRouter := router.PathPrefix("/student").Subrouter()
	studentRouter.HandleFunc("/reservation/view", handleWithCookie(service.ViewReservationsByStudent)).Methods("GET")
	studentRouter.HandleFunc("/reservation/make", handleWithCookie(service.MakeReservationByStudent)).Methods("POST")
	studentRouter.HandleFunc("/reservation/feedback/get", handleWithCookie(service.GetFeedbackByStudent)).Methods("POST")
	studentRouter.HandleFunc("/reservation/feedback/submit", handleWithCookie(service.SubmitFeedbackByStudent)).Methods("POST")
	teacherRouter := router.PathPrefix("/teacher").Subrouter()
	teacherRouter.HandleFunc("/reservation/view", handleWithCookie(service.ViewReservationsByTeacher)).Methods("GET")
	teacherRouter.HandleFunc("/reservation/feedback/get", handleWithCookie(service.GetFeedbackByTeacher)).Methods("POST")
	teacherRouter.HandleFunc("/reservation/feedback/submit", handleWithCookie(service.SubmitFeedbackByTeacher)).Methods("POST")
	teacherRouter.HandleFunc("/student/get", handleWithCookie(service.GetStudentInfoByTeacher)).Methods("POST")
	teacherRouter.HandleFunc("/student/query", handleWithCookie(service.QueryStudentInfoByTeacher)).Methods("POST")
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/timetable/view", handleWithCookie(service.ViewTimedReservationsByAdmin)).Methods("GET")
	adminRouter.HandleFunc("/timetable/add", handleWithCookie(service.AddTimedReservationByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/timetable/edit", handleWithCookie(service.EditTimedReservationByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/timetable/remove", handleWithCookie(service.RemoveTimedReservationsByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/timetable/open", handleWithCookie(service.OpenTimedReservationsByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/timetable/close", handleWithCookie(service.CloseTimedReservationsByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/view", handleWithCookie(service.ViewReservationsByAdmin)).Methods("GET")
	adminRouter.HandleFunc("/reservation/view/daily", handleWithCookie(service.ViewDailyReservationsByAdmin)).Methods("GET")
	adminRouter.HandleFunc("/reservation/export/today", handleWithCookie(service.ExportTodayReservationsByAdmin)).Methods("GET")
	adminRouter.HandleFunc("/reservation/export/report", handleWithCookie(service.ExportReportFormByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/export/report/monthly", handleWithCookie(service.ExportReportMonthlyByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/add", handleWithCookie(service.AddReservationByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/edit", handleWithCookie(service.EditReservationByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/remove", handleWithCookie(service.RemoveReservationsByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/cancel", handleWithCookie(service.CancelReservationByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/feedback/get", handleWithCookie(service.GetFeedbackByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/feedback/submit", handleWithCookie(service.SubmitFeedbackByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/student/set", handleWithCookie(service.SetStudentByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/get", handleWithCookie(service.GetStudentInfoByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/search", handleWithCookie(service.SearchStudentByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/crisis/update", handleWithCookie(service.UpdateStudentCrisisLevelByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/archive/update", handleWithCookie(service.UpdateStudentArchiveNumberByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/password/reset", handleWithCookie(service.ResetStudentPasswordByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/account/delete", handleWithCookie(service.DeleteStudentAccountByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/export", handleWithCookie(service.ExportStudentByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/unbind", handleWithCookie(service.UnbindStudentByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/bind", handleWithCookie(service.BindStudentByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/query", handleWithCookie(service.QueryStudentInfoByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/teacher/search", handleWithCookie(service.SearchTeacherByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/teacher/workload", handleWithCookie(service.GetTeacherWorkloadByAdmin)).Methods("POST")
	categoryRouter := router.PathPrefix("/category").Subrouter()
	categoryRouter.HandleFunc("/feedback", handleWithCookie(service.GetFeedbackCategories)).Methods("GET")
	// http加载处理器
	http.Handle("/", router)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets/"))))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
