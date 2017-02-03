package main

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/buslogic"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/config"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"flag"
	"github.com/mijia/sweb/log"
)

func main() {
	conf := flag.String("conf", "deploy/thxl.conf", "conf file path")
	isSmock := flag.Bool("smock", true, "is smock server")
	method := flag.String("method", "", "method")
	mailTo := flag.String("mail-to", "shudiwsh2009@gmail.com", "mail to list")
	username := flag.String("username", "", "username")
	password := flag.String("password", "", "password")
	userType := flag.Int("usertype", model.USER_TYPE_UNKNOWN, "usertype")
	days := flag.Int("days", 0, "shift reservation time in days")
	flag.Parse()

	config.InitWithParams(*conf, *isSmock)
	log.Infof("config loaded: %+v", *config.Instance())
	workflow := buslogic.NewWorkflow()

	if *method == "reminder" {
		// 每晚发送第二天咨询的提醒短信
		workflow.SendTomorrowReservationReminderSMS()
	} else if *method == "timetable" {
		// 每天早上发送当天咨询安排表邮件
		workflow.SendTodayTimetableMail(*mailTo)
	} else if *method == "import-archive-file" {
		// 导入档案列表，仅需执行一次
		err := workflow.ImportArchiveFromCSVFile()
		if err != nil {
			log.Errorf("import archive file failed, err: %+v", err)
			return
		}
	} else if *method == "reset-user-password" {
		// 重置用户密码
		if err := workflow.ResetUserPassword(*username, *userType, *password); err != nil {
			log.Errorf("reset user password failed, err: %+v", err)
			return
		}
	} else if *method == "shift-reservation-time" {
		// 将所有咨询的开始时间和结束时间变更数天
		if err := workflow.ShiftReservationTimeInDays(*days); err != nil {
			log.Errorf("fail to shift reservation: %+v", err)
			return
		}
	} else if *method == "add-new-admin" {
		// 添加新管理员
		if _, err := workflow.AddNewAdmin(*username, *password); err != nil {
			log.Errorf("fail to add new admin: %+v", err)
			return
		}
	}
	log.Info("Success")
}
