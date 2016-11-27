package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"github.com/mijia/sweb/log"
	"gopkg.in/redis.v5"
	"time"
)

type Workflow struct {
	mongoClient *model.MongoClient
	redisClient *redis.Client
}

func NewWorkflow() *Workflow {
	var err error
	if time.Local, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		log.Fatalf("初始化时区失败：%v", err)
	}
	ret := &Workflow{
		mongoClient: model.NewMongoClient(),
		redisClient: model.NewRedisClient(),
	}
	return ret
}

func (w *Workflow) MongoClient() *model.MongoClient {
	return w.mongoClient
}

func (w *Workflow) RedisClient() *redis.Client {
	return w.redisClient
}
