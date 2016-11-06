package service

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/buslogic"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/config"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"github.com/mijia/sweb/log"
	"gopkg.in/redis.v5"
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

func Model() *model.Model {
	return wf.Model()
}

func RedisClient() *redis.Client {
	return wf.RedisClient()
}
