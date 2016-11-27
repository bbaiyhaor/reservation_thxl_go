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
	flag.Parse()

	config.InitWithParams(*conf, *isSmock)
	log.Infof("config loaded: %+v", *config.Instance())
	workflow := buslogic.NewWorkflow()

	if *method == "reminder" {
		workflow.SendTomorrowReservationReminderSMS()
	} else if *method == "timetable" {
		workflow.SendTodayTimetableMail(*mailTo)
	} else if *method == "transfer-data-2016-11" {
		TransferDataForNov2016(workflow)
	} else if *method == "import-archive-file" {
		err := workflow.ImportArchiveFromCSVFile()
		if err != nil {
			log.Errorf("import archive file failed, err: %+v", err)
		}
	} else if *method == "reset-user-password" {
		if err := workflow.ResetUserPassword(*username, *userType, *password); err != nil {
			log.Errorf("reset user password failed, err: %+v", err)
			return
		}
		log.Info("Success")
	}
}
