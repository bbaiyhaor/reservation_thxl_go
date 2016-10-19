package service

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/buslogic"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/config"
	"github.com/mijia/sweb/log"
)

var wf *buslogic.Workflow

func InitService(confPath string, isSmock bool) {
	config.InitWithParams(confPath, isSmock)
	log.Infof("config loaded: %+v", *config.Instance())
	wf = buslogic.NewWorkflow()

	if err := wf.ImportArchiveFromCSVFile(); err != nil {
		log.Fatalf("初始化档案失败：%v", err)
	}
}

func Workflow() *buslogic.Workflow {
	return wf
}
