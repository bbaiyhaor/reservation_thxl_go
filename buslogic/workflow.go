package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"log"
	"time"
)

type Workflow struct {
	model *model.Model
}

func NewWorkflow() *Workflow {
	var err error
	if time.Local, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		log.Fatalf("初始化时区失败：%v", err)
	}
	ret := &Workflow{
		model: model.NewModel(),
	}
	return ret
}

func (w *Workflow) Model() *model.Model {
	return w.model
}
