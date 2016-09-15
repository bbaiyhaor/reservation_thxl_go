package main

import (
	"flag"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"github.com/shudiwsh2009/reservation_thxl_go/workflow"
	"gopkg.in/mgo.v2"
	"log"
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
	// Reminder
	today := utils.GetToday()
	from := today.AddDate(0, 0, 1)
	to := today.AddDate(0, 0, 2)
	reservations, err := models.GetReservationsBetweenTime(from, to)
	if err != nil {
		log.Printf("获取咨询列表失败：%v", err)
		return
	}
	for _, reservation := range reservations {
		if reservation.Status == models.RESERVATED {
			workflow.SendReminderSMS(reservation)
		}
	}
}
