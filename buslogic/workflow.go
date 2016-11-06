package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"gopkg.in/redis.v5"
	"log"
	"time"
)

type Workflow struct {
	model       *model.Model
	redisClient *redis.Client
}

func NewWorkflow() *Workflow {
	var err error
	if time.Local, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		log.Fatalf("初始化时区失败：%v", err)
	}
	ret := &Workflow{
		model:       model.NewModel(),
		redisClient: model.NewRedisClient(),
	}
	return ret
}

func (w *Workflow) Model() *model.Model {
	return w.model
}

func (w *Workflow) RedisClient() *redis.Client {
	return w.redisClient
}
