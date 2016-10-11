package main

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/config"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/models"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/utils"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/workflow"
	"flag"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

func main() {
	conf := flag.String("conf", "../deploy/thxl.conf", "conf file path")
	isSmock := flag.Bool("smock", true, "is smock server")
	flag.Parse()
	config.InitWithParams(*conf, *isSmock)
	log.Printf("config loaded: %+v", *config.Instance())
	// 数据库连接
	var session *mgo.Session
	var err error
	if config.Instance().IsSmockServer() {
		session, err = mgo.Dial("127.0.0.1:27017")
	} else {
		mongoDbDialInfo := mgo.DialInfo{
			Addrs:    []string{config.Instance().MongoHost},
			Timeout:  60 * time.Second,
			Database: config.Instance().MongoDatabase,
			Username: config.Instance().MongoUser,
			Password: config.Instance().MongoPassword,
		}
		session, err = mgo.DialWithInfo(&mongoDbDialInfo)
	}
	if err != nil {
		log.Fatalf("连接数据库失败：%v", err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	models.Mongo = session.DB("reservation_thxl")
	// 时区
	if time.Local, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		log.Fatalf("初始化时区失败：%v", err)
	}
	// Reminder
	today := utils.BeginOfDay(time.Now())
	from := today.AddDate(0, 0, 1)
	to := today.AddDate(0, 0, 2)
	reservations, err := models.GetReservationsBetweenTime(from, to)
	if err != nil {
		log.Printf("获取咨询列表失败：%v", err)
		return
	}
	succCnt, failCnt := 0, 0
	for _, reservation := range reservations {
		if reservation.Status == models.RESERVATED {
			if err = workflow.SendReminderSMS(reservation); err == nil {
				succCnt++
			} else {
				log.Printf("发送短信失败：%+v %+v", reservation, err)
				failCnt++
			}
		}
	}
	log.Printf("发送%d个预约记录的提醒短信，成功%d个，失败%d个", succCnt+failCnt, succCnt, failCnt)
}
