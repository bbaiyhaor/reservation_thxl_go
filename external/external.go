package main

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/buslogic"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/config"
	"flag"
	"github.com/mijia/sweb/log"
)

func main() {
	conf := flag.String("conf", "deploy/thxl.conf", "conf file path")
	isSmock := flag.Bool("smock", true, "is smock server")
	method := flag.String("method", "", "method")
	mailTo := flag.String("mail-to", "shudiwsh2009@gmail.com", "mail to list")
	flag.Parse()

	config.InitWithParams(*conf, *isSmock)
	log.Infof("config loaded: %+v", *config.Instance())
	workflow := buslogic.NewWorkflow()

	if *method == "reminder" {
		workflow.SendTomorrowReservationReminderSMS()
	} else if *method == "timetable" {
		workflow.SendTodayTimetableMail(*mailTo)
	}
}
