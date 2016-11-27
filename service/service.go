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
}

func Workflow() *buslogic.Workflow {
	return wf
}

func MongoClient() *model.MongoClient {
	return wf.MongoClient()
}

func RedisClient() *redis.Client {
	return wf.RedisClient()
}
