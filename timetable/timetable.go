package main

import (
	"flag"
	"fmt"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"github.com/shudiwsh2009/reservation_thxl_go/workflow"
	"gopkg.in/mgo.v2"
	"log"
	"sort"
	"time"
)

func main() {
	appEnv := flag.String("app-env", "STAGING", "app environment")
	smsUid := flag.String("sms-uid", "", "sms uid")
	smsKey := flag.String("sms-key", "", "sms key")
	mailSmtp := flag.String("mail-smtp", "", "mail smtp")
	mailUsername := flag.String("mail-username", "", "mail username")
	mailPassword := flag.String("mail-password", "", "mail password")
	flag.Parse()
	utils.APP_ENV = *appEnv
	utils.SMS_UID = *smsUid
	utils.SMS_KEY = *smsKey
	utils.MAIL_SMTP = *mailSmtp
	utils.MAIL_USERNAME = *mailUsername
	utils.MAIL_PASSWORD = *mailPassword
	log.Printf("loading config: %s %s %s %s %s %s", utils.APP_ENV, utils.SMS_UID, utils.SMS_KEY, utils.MAIL_SMTP, utils.MAIL_USERNAME, utils.MAIL_PASSWORD)
	// 数据库连接
	mongoDbDialInfo := mgo.DialInfo{
		Addrs:    []string{"127.0.0.1:27017"},
		Timeout:  60 * time.Second,
		Database: "admin",
		Username: "admin",
		Password: "THXLFZZX",
	}
	var session *mgo.Session
	var err error
	if utils.APP_ENV == "ONLINE" {
		session, err = mgo.DialWithInfo(&mongoDbDialInfo)
	} else {
		session, err = mgo.Dial("127.0.0.1:27017")
	}
	if err != nil {
		log.Printf("连接数据库失败：%v", err)
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	models.Mongo = session.DB("reservation_thxl")
	// 时区
	if utils.Location, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		log.Printf("初始化时区失败：%v", err)
		return
	}
	// timetable
	today := utils.GetToday()
	tomorrow := today.AddDate(0, 0, 1)
	reservations, err := models.GetReservationsBetweenTime(today, tomorrow)
	if err != nil {
		log.Printf(err)
		return
	}
	todayDate := today.Format(utils.DATE_PATTERN)
	if timedReservations, err := models.GetTimedReservationsByWeekday(today.Weekday()); err == nil {
		for _, tr := range timedReservations {
			if !tr.Exceptions[todayDate] && !tr.Timed[todayDate] {
				reservations = append(reservations, tr.ToReservation(today))
			}
		}
	}
	sort.Sort(models.ReservationSlice(reservations))
	filename := "timetable_" + todayDate + utils.CsvSuffix
	if err = workflow.ExportTodayReservationTimetable(reservations, filename); err != nil {
		log.Printf(err)
		return
	}
	// email
	title := fmt.Sprintf("【心理发展中心】%s咨询安排表", todayDate)
	if err := workflow.SendEmail(title, title, []string{fmt.Sprintf("%s%s", utils.ExportFolder, filename)},
		workflow.EMAIL_TO_DEVELOPER); err != nil {
		log.Printf("发送邮件失败：%v", err)
	}
}
