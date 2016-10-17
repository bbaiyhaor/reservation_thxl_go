package service

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/buslogic"
	"log"
)

type Service struct {
	w *buslogic.Workflow
}

func NewService() *Service {
	s := &Service{
		w: buslogic.NewWorkflow(),
	}

	if err := s.w.ImportArchiveFromCSVFile(); err != nil {
		log.Fatalf("初始化档案失败：%v", err)
	}

	return s
}
