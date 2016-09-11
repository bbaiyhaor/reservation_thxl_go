package main

import (
	"fmt"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"github.com/shudiwsh2009/reservation_thxl_go/workflow"
	"gopkg.in/mgo.v2"
	"time"
	"flag"
)

func main() {
	appEnv := flag.String("app-env", "STAGING", "app environment")
	smsUid := flag.String("sms-uid", "", "sms uid")
	smsKey := flag.String("sms-key", "", "sms key")
	flag.Parse()
	utils.APP_ENV = *appEnv
	utils.SMS_UID = *smsUid
	utils.SMS_KEY = *smsKey
	// 数据库连接
	mongoDbDialInfo := mgo.DialInfo{
		Addrs:		[]string{"127.0.0.1:27017"},
		Timeout:	60 * time.Second,
		Database:	"admin",
		Username:	"admin",
		Password:	"THXLFZZX",
	}
	var session *mgo.Session
	var err error
	if utils.APP_ENV == "ONLINE" {
		session, err = mgo.DialWithInfo(&mongoDbDialInfo)
	} else {
		session, err = mgo.Dial("127.0.0.1:27017")
	}
	if err != nil {
		fmt.Printf("连接数据库失败：%v\n", err)
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	models.Mongo = session.DB("reservation_thxl")
	// 时区
	if utils.Location, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		fmt.Errorf("初始化时区失败：%v", err)
		return
	}
	// Reminder
	today := utils.GetToday()
	from := today.AddDate(0, 0, 1)
	to := today.AddDate(0, 0, 2)
	reservations, err := models.GetReservationsBetweenTime(from, to)
	if err != nil {
		fmt.Errorf("获取咨询列表失败：%v", err)
		return
	}
	for _, reservation := range reservations {
		if reservation.Status == models.RESERVATED {
			workflow.SendReminderSMS(reservation)
		}
	}
}
