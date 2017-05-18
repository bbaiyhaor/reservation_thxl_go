package buslogic

import (
	"fmt"
	"github.com/shudiwsh2009/reservation_thxl_go/model"
	re "github.com/shudiwsh2009/reservation_thxl_go/rerror"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"github.com/tealeg/xlsx"
	"sort"
	"strconv"
	"strings"
)

// 定义单元格样式
var (
	textCenterAlignment = xlsx.Alignment{
		Horizontal: "center",
		Vertical:   "center",
	}
	textRightAlignment = xlsx.Alignment{
		Horizontal: "right",
		Vertical:   "center",
	}
	border = xlsx.Border{
		Left:        "thin",
		LeftColor:   "000000",
		Right:       "thin",
		RightColor:  "000000",
		Top:         "thin",
		TopColor:    "000000",
		Bottom:      "thin",
		BottomColor: "000000",
	}

	grayFill   = *xlsx.NewFill("solid", "D9E2F3", "D9E2F3")
	greenFill  = *xlsx.NewFill("solid", "C5E0B2", "C5E0B2")
	orangeFill = *xlsx.NewFill("solid", "F4B082", "F4B082")
	yellowFill = *xlsx.NewFill("solid", "FFC000", "FFC000")
	redFill    = *xlsx.NewFill("solid", "FF0000", "FF0000")

	borderStyle *xlsx.Style = &xlsx.Style{
		ApplyBorder: true,
		Border:      border,
	}
	textCenterStyle *xlsx.Style = &xlsx.Style{
		ApplyAlignment: true,
		ApplyBorder:    true,
		Alignment:      textCenterAlignment,
		Border:         border,
	}
	textCenterGrayStyle *xlsx.Style = &xlsx.Style{
		ApplyAlignment: true,
		ApplyFill:      true,
		ApplyBorder:    true,
		Alignment:      textCenterAlignment,
		Fill:           grayFill,
		Border:         border,
	}
	textCenterGreenStyle *xlsx.Style = &xlsx.Style{
		ApplyAlignment: true,
		ApplyFill:      true,
		ApplyBorder:    true,
		Alignment:      textCenterAlignment,
		Fill:           greenFill,
		Border:         border,
	}
	textRightGreenStyle *xlsx.Style = &xlsx.Style{
		ApplyAlignment: true,
		ApplyFill:      true,
		ApplyBorder:    true,
		Alignment:      textRightAlignment,
		Fill:           greenFill,
		Border:         border,
	}
	textRightOrangeStyle *xlsx.Style = &xlsx.Style{
		ApplyAlignment: true,
		ApplyFill:      true,
		ApplyBorder:    true,
		Alignment:      textRightAlignment,
		Fill:           orangeFill,
		Border:         border,
	}
	textCenterOrangeStyle *xlsx.Style = &xlsx.Style{
		ApplyAlignment: true,
		ApplyFill:      true,
		ApplyBorder:    true,
		Alignment:      textCenterAlignment,
		Fill:           orangeFill,
		Border:         border,
	}
	bgGrayStyle *xlsx.Style = &xlsx.Style{
		ApplyFill:   true,
		ApplyBorder: true,
		Fill:        grayFill,
		Border:      border,
	}
	bgGreenStyle *xlsx.Style = &xlsx.Style{
		ApplyFill:   true,
		ApplyBorder: true,
		Fill:        greenFill,
		Border:      border,
	}
	bgOrangeStyle *xlsx.Style = &xlsx.Style{
		ApplyFill:   true,
		ApplyBorder: true,
		Fill:        orangeFill,
		Border:      border,
	}
	bgYellowStyle *xlsx.Style = &xlsx.Style{
		ApplyFill:   true,
		ApplyBorder: true,
		Fill:        yellowFill,
		Border:      border,
	}
	bgRedStyle *xlsx.Style = &xlsx.Style{
		ApplyFill:   true,
		ApplyBorder: true,
		Fill:        redFill,
		Border:      border,
	}
)

//================================================================
//===========================咨询月报==============================
//================================================================
type FeedbackGroup struct {
	GroupName           string
	SecondaryGroups     []*FeedbackGroup
	Grades              map[string]int // 本科生、硕士生、博士生
	Instructor          int            // 辅导员
	Teacher             int            // 教师
	Family              int            // 家属
	Others              int            // 其他
	MaleIdMap           map[string]int // 学生（男）咨询次数表
	NumMale             int            // 合计人数（男）
	CountMale           int            // 合计人次（男）
	FemaleIdMap         map[string]int
	NumFemale           int
	CountFemale         int
	UnderGraduateIdMap  map[string]int
	NumUnderGraduates   int
	CountUnderGraduates int
	GraduateIdMap       map[string]int
	NumGraduates        int
	CountGraduates      int
	TotalIdMap          map[string]int // 学生咨询次数表
	NumTotal            int            // 会谈总人数
	CountTotal          int            // 会谈总人次
	Ratio               float64        // 比例（需转成百分比）
	// 交叉统计
	HasEmphasisUnderGraduateIdMap        map[string]int // 含有重点标记的本科生次数表
	NumHasEmphasisUnderGraduate          int            // 含有重点标记的本科生人数
	CountHasEmphasisUnderGraduate        int            // 含有重点标记的本科生人次
	HasEmphasisGraduateIdMap             map[string]int
	NumHasEmphasisGraduate               int
	CountHasEmphasisGraduate             int
	HasEmphasisIdMap                     map[string]int // 含有重点标记的学生次数表
	NumHasEmphasis                       int            // 含有重点标记的咨询人数
	CountHasEmphasis                     int            // 含有重点标记的咨询人次
	CountUnderGraduateEmphasisInCategory int            // 评估分类中本科生重点情况的频次
	CountGraduateEmphasisInCategory      int            // 评估分类中研究生重点情况的频次
	CountEmphasisInCategory              int
}

// 来访&危机情况
type ReservationCrisis struct {
	CountOnlyReservated int // 有预约，无危机
	CountOnlyCrisis     int // 无预约，有危机
	CountNeither        int // 无预约，无危机
	CountBoth           int // 有预约，有危机
	CountSendNotify     int // 发危机通报人次
}

type AdvisoryConsultation struct {
	Advisory     ReservationCrisis // 咨询
	Consultation ReservationCrisis // 会商
	Total        ReservationCrisis // 总数
}

// 会商与危机干预
type ConsultationCrisis struct {
	Date                 string   // 日期
	Fullname             string   // 姓名
	Username             string   // 学号
	Gender               string   // 性别
	Academic             string   // 学历
	School               string   // 院系
	TeacherFullname      string   // 接待咨询师
	SchoolContact        string   // 院系联系人
	ConsultationOrCrisis []string // 会商or危机处理
	ReservationStatus    string   // 来访情况（是否预约）
	Category             string   // 评估分类
	EmphasisStr          string   // 重点标记
	CrisisLevel          string   // 星级
	IsSendNotify         string   // 是否发危机通报
}

func (w *Workflow) ExportReservationFeedbackReportToFile(reservations []*model.Reservation, path string) error {
	// 建立存储结构
	// 咨询情况汇总
	categoryFCGroup := make([]*FeedbackGroup, 0)
	categoryFCGroupMap := make(map[string]*FeedbackGroup)
	categorySCGroupMap := make(map[string]*FeedbackGroup)
	for fi, firstCategory := range model.FeedbackFirstCategoryMap {
		if fi == "" {
			continue
		}
		fcGroup := &FeedbackGroup{
			GroupName:                     firstCategory,
			SecondaryGroups:               make([]*FeedbackGroup, 0),
			MaleIdMap:                     make(map[string]int),
			FemaleIdMap:                   make(map[string]int),
			UnderGraduateIdMap:            make(map[string]int),
			GraduateIdMap:                 make(map[string]int),
			TotalIdMap:                    make(map[string]int),
			HasEmphasisUnderGraduateIdMap: make(map[string]int),
			HasEmphasisGraduateIdMap:      make(map[string]int),
			HasEmphasisIdMap:              make(map[string]int),
		}
		for si, secondCategory := range model.FeedbackSecondCategoryMap[fi] {
			if si == "" {
				continue
			}
			scGroup := &FeedbackGroup{
				GroupName:                     secondCategory,
				Grades:                        make(map[string]int),
				MaleIdMap:                     make(map[string]int),
				FemaleIdMap:                   make(map[string]int),
				UnderGraduateIdMap:            make(map[string]int),
				GraduateIdMap:                 make(map[string]int),
				TotalIdMap:                    make(map[string]int),
				HasEmphasisUnderGraduateIdMap: make(map[string]int),
				HasEmphasisGraduateIdMap:      make(map[string]int),
				HasEmphasisIdMap:              make(map[string]int),
			}
			fcGroup.SecondaryGroups = append(fcGroup.SecondaryGroups, scGroup)
			categorySCGroupMap[si] = scGroup
		}
		sort.Slice(fcGroup.SecondaryGroups, func(i, j int) bool {
			return fcGroup.SecondaryGroups[i].GroupName < fcGroup.SecondaryGroups[j].GroupName
		})
		categoryFCGroupMap[fi] = fcGroup
		categoryFCGroup = append(categoryFCGroup, fcGroup)
	}
	sort.Slice(categoryFCGroup, func(i, j int) bool {
		return categoryFCGroup[i].GroupName < categoryFCGroup[j].GroupName
	})
	categoryTotalGroup := &FeedbackGroup{
		GroupName: "总计",
		Grades: map[string]int{
			"16级": 0,
			"15级": 0,
			"14级": 0,
			"13级": 0,
			"12级": 0,
			"16硕": 0,
			"15硕": 0,
			"14硕": 0,
			"16博": 0,
			"15博": 0,
			"14博": 0,
		},
		MaleIdMap:                     make(map[string]int),
		FemaleIdMap:                   make(map[string]int),
		UnderGraduateIdMap:            make(map[string]int),
		GraduateIdMap:                 make(map[string]int),
		TotalIdMap:                    make(map[string]int),
		HasEmphasisUnderGraduateIdMap: make(map[string]int),
		HasEmphasisGraduateIdMap:      make(map[string]int),
		HasEmphasisIdMap:              make(map[string]int),
	}
	// 重点情况汇总
	emphasisSCGroupMap := make(map[string]*FeedbackGroup)
	severityFCGroup := &FeedbackGroup{
		GroupName:                     "严重程度",
		SecondaryGroups:               make([]*FeedbackGroup, 0),
		MaleIdMap:                     make(map[string]int),
		FemaleIdMap:                   make(map[string]int),
		UnderGraduateIdMap:            make(map[string]int),
		GraduateIdMap:                 make(map[string]int),
		TotalIdMap:                    make(map[string]int),
		HasEmphasisUnderGraduateIdMap: make(map[string]int),
		HasEmphasisGraduateIdMap:      make(map[string]int),
		HasEmphasisIdMap:              make(map[string]int),
	}
	for _, sc := range model.FeedbackSeverity {
		scGroup := &FeedbackGroup{
			GroupName:                     sc,
			Grades:                        make(map[string]int),
			MaleIdMap:                     make(map[string]int),
			FemaleIdMap:                   make(map[string]int),
			UnderGraduateIdMap:            make(map[string]int),
			GraduateIdMap:                 make(map[string]int),
			TotalIdMap:                    make(map[string]int),
			HasEmphasisUnderGraduateIdMap: make(map[string]int),
			HasEmphasisGraduateIdMap:      make(map[string]int),
			HasEmphasisIdMap:              make(map[string]int),
		}
		severityFCGroup.SecondaryGroups = append(severityFCGroup.SecondaryGroups, scGroup)
		emphasisSCGroupMap[sc] = scGroup
	}
	medicalDiagnosisFCGroup := &FeedbackGroup{
		GroupName:                     "疑似或明确的医疗诊断",
		SecondaryGroups:               make([]*FeedbackGroup, 0),
		MaleIdMap:                     make(map[string]int),
		FemaleIdMap:                   make(map[string]int),
		UnderGraduateIdMap:            make(map[string]int),
		GraduateIdMap:                 make(map[string]int),
		TotalIdMap:                    make(map[string]int),
		HasEmphasisUnderGraduateIdMap: make(map[string]int),
		HasEmphasisGraduateIdMap:      make(map[string]int),
		HasEmphasisIdMap:              make(map[string]int),
	}
	for _, sc := range model.FeedbackMedicalDiagnosis {
		scGroup := &FeedbackGroup{
			GroupName:                     sc,
			Grades:                        make(map[string]int),
			MaleIdMap:                     make(map[string]int),
			FemaleIdMap:                   make(map[string]int),
			UnderGraduateIdMap:            make(map[string]int),
			GraduateIdMap:                 make(map[string]int),
			TotalIdMap:                    make(map[string]int),
			HasEmphasisUnderGraduateIdMap: make(map[string]int),
			HasEmphasisGraduateIdMap:      make(map[string]int),
			HasEmphasisIdMap:              make(map[string]int),
		}
		medicalDiagnosisFCGroup.SecondaryGroups = append(medicalDiagnosisFCGroup.SecondaryGroups, scGroup)
		emphasisSCGroupMap[sc] = scGroup
	}
	crisisFCGroup := &FeedbackGroup{
		GroupName:                     "危急情况",
		SecondaryGroups:               make([]*FeedbackGroup, 0),
		MaleIdMap:                     make(map[string]int),
		FemaleIdMap:                   make(map[string]int),
		UnderGraduateIdMap:            make(map[string]int),
		GraduateIdMap:                 make(map[string]int),
		TotalIdMap:                    make(map[string]int),
		HasEmphasisUnderGraduateIdMap: make(map[string]int),
		HasEmphasisGraduateIdMap:      make(map[string]int),
		HasEmphasisIdMap:              make(map[string]int),
	}
	for _, sc := range model.FeedbackCrisis {
		scGroup := &FeedbackGroup{
			GroupName:                     sc,
			Grades:                        make(map[string]int),
			MaleIdMap:                     make(map[string]int),
			FemaleIdMap:                   make(map[string]int),
			UnderGraduateIdMap:            make(map[string]int),
			GraduateIdMap:                 make(map[string]int),
			TotalIdMap:                    make(map[string]int),
			HasEmphasisUnderGraduateIdMap: make(map[string]int),
			HasEmphasisGraduateIdMap:      make(map[string]int),
			HasEmphasisIdMap:              make(map[string]int),
		}
		crisisFCGroup.SecondaryGroups = append(crisisFCGroup.SecondaryGroups, scGroup)
		emphasisSCGroupMap[sc] = scGroup
	}
	emphasisTotalGroup := &FeedbackGroup{
		GroupName: "总计",
		Grades: map[string]int{
			"16级": 0,
			"15级": 0,
			"14级": 0,
			"13级": 0,
			"12级": 0,
			"16硕": 0,
			"15硕": 0,
			"14硕": 0,
			"16博": 0,
			"15博": 0,
			"14博": 0,
		},
		MaleIdMap:                     make(map[string]int),
		FemaleIdMap:                   make(map[string]int),
		UnderGraduateIdMap:            make(map[string]int),
		GraduateIdMap:                 make(map[string]int),
		TotalIdMap:                    make(map[string]int),
		HasEmphasisUnderGraduateIdMap: make(map[string]int),
		HasEmphasisGraduateIdMap:      make(map[string]int),
		HasEmphasisIdMap:              make(map[string]int),
	}
	// 来访&危机
	var advisoryConsulation AdvisoryConsultation
	// 会商与危机干预
	consultationCrisisList := make([]*ConsultationCrisis, 0)
	// 学生咨询情况
	studentMap := make(map[string]*model.Student)
	teacherMap := make(map[string]*model.Teacher)
	studentReservationsMap := make(map[string][]*model.Reservation)

	// 分析咨询数据
	for _, r := range reservations {
		if r.TeacherFeedback.IsEmpty() {
			continue
		}
		student, err := w.mongoClient.GetStudentById(r.StudentId)
		if err != nil || student == nil || student.UserType != model.USER_TYPE_STUDENT {
			continue
		}
		studentId := student.Id.Hex()
		grade, err := utils.ParseStudentId(student.Username)
		if err != nil {
			continue
		}
		teacher, err := w.mongoClient.GetTeacherById(r.TeacherId)
		if err != nil || teacher == nil || teacher.UserType != model.USER_TYPE_TEACHER {
			continue
		}
		// 学生咨询情况
		studentMap[studentId] = student
		teacherMap[r.TeacherId] = teacher
		if _, ok := studentReservationsMap[studentId]; !ok {
			studentReservationsMap[studentId] = make([]*model.Reservation, 0)
		}
		studentReservationsMap[studentId] = append(studentReservationsMap[studentId], r)
		// 咨询情况汇总
		category := r.TeacherFeedback.Category
		// 来访情况，H分类中的来访人员特殊处理
		switch category {
		case "H1":
			categorySCGroupMap[category].Instructor++
			categoryTotalGroup.Instructor++
		case "H2", "H6":
			categorySCGroupMap[category].Teacher++
			categoryTotalGroup.Teacher++
		case "H3":
			categorySCGroupMap[category].Family++
			categoryTotalGroup.Family++
		case "H4", "H5":
			categorySCGroupMap[category].Others++
			categoryTotalGroup.Others++
		default:
			categorySCGroupMap[category].Grades[grade]++
			categoryTotalGroup.Grades[grade]++
		}
		// 性别统计
		if student.Gender == "男" {
			if _, ok := categorySCGroupMap[category].MaleIdMap[studentId]; !ok {
				categorySCGroupMap[category].NumMale++
			}
			categorySCGroupMap[category].CountMale++
			categorySCGroupMap[category].MaleIdMap[studentId]++

			if _, ok := categoryFCGroupMap[category[0:1]].MaleIdMap[studentId]; !ok {
				categoryFCGroupMap[category[0:1]].NumMale++
			}
			categoryFCGroupMap[category[0:1]].CountMale++
			categoryFCGroupMap[category[0:1]].MaleIdMap[studentId]++

			if _, ok := categoryTotalGroup.MaleIdMap[studentId]; !ok {
				categoryTotalGroup.NumMale++
			}
			categoryTotalGroup.CountMale++
			categoryTotalGroup.MaleIdMap[studentId]++
		} else if student.Gender == "女" {
			if _, ok := categorySCGroupMap[category].FemaleIdMap[studentId]; !ok {
				categorySCGroupMap[category].NumFemale++
			}
			categorySCGroupMap[category].CountFemale++
			categorySCGroupMap[category].FemaleIdMap[studentId]++

			if _, ok := categoryFCGroupMap[category[0:1]].FemaleIdMap[studentId]; !ok {
				categoryFCGroupMap[category[0:1]].NumFemale++
			}
			categoryFCGroupMap[category[0:1]].CountFemale++
			categoryFCGroupMap[category[0:1]].FemaleIdMap[studentId]++

			if _, ok := categoryTotalGroup.FemaleIdMap[studentId]; !ok {
				categoryTotalGroup.NumFemale++
			}
			categoryTotalGroup.CountFemale++
			categoryTotalGroup.FemaleIdMap[studentId]++
		}
		// 本科生/研究生统计
		if strings.HasSuffix(grade, "级") {
			if _, ok := categorySCGroupMap[category].UnderGraduateIdMap[studentId]; !ok {
				categorySCGroupMap[category].NumUnderGraduates++
			}
			categorySCGroupMap[category].CountUnderGraduates++
			categorySCGroupMap[category].UnderGraduateIdMap[studentId]++

			if _, ok := categoryFCGroupMap[category[0:1]].UnderGraduateIdMap[studentId]; !ok {
				categoryFCGroupMap[category[0:1]].NumUnderGraduates++
			}
			categoryFCGroupMap[category[0:1]].CountUnderGraduates++
			categoryFCGroupMap[category[0:1]].UnderGraduateIdMap[studentId]++

			if _, ok := categoryTotalGroup.UnderGraduateIdMap[studentId]; !ok {
				categoryTotalGroup.NumUnderGraduates++
			}
			categoryTotalGroup.CountUnderGraduates++
			categoryTotalGroup.UnderGraduateIdMap[studentId]++
		} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
			if _, ok := categorySCGroupMap[category].GraduateIdMap[studentId]; !ok {
				categorySCGroupMap[category].NumGraduates++
			}
			categorySCGroupMap[category].CountGraduates++
			categorySCGroupMap[category].GraduateIdMap[studentId]++

			if _, ok := categoryFCGroupMap[category[0:1]].GraduateIdMap[studentId]; !ok {
				categoryFCGroupMap[category[0:1]].NumGraduates++
			}
			categoryFCGroupMap[category[0:1]].CountGraduates++
			categoryFCGroupMap[category[0:1]].GraduateIdMap[studentId]++

			if _, ok := categoryTotalGroup.GraduateIdMap[studentId]; !ok {
				categoryTotalGroup.NumGraduates++
			}
			categoryTotalGroup.CountGraduates++
			categoryTotalGroup.GraduateIdMap[studentId]++
		}
		// 总计
		if _, ok := categorySCGroupMap[category].TotalIdMap[studentId]; !ok {
			categorySCGroupMap[category].NumTotal++
		}
		categorySCGroupMap[category].CountTotal++
		categorySCGroupMap[category].TotalIdMap[studentId]++

		if _, ok := categoryFCGroupMap[category[0:1]].TotalIdMap[studentId]; !ok {
			categoryFCGroupMap[category[0:1]].NumTotal++
		}
		categoryFCGroupMap[category[0:1]].CountTotal++
		categoryFCGroupMap[category[0:1]].TotalIdMap[studentId]++

		if _, ok := categoryTotalGroup.TotalIdMap[studentId]; !ok {
			categoryTotalGroup.NumTotal++
		}
		categoryTotalGroup.CountTotal++
		categoryTotalGroup.TotalIdMap[studentId]++

		// 重点情况汇总
		hasEmphasis := false
		severity := r.TeacherFeedback.Severity
		medicalDiagnosis := r.TeacherFeedback.MedicalDiagnosis
		crisis := r.TeacherFeedback.Crisis
		for i, s := range severity {
			if s == 1 {
				emphasisSCGroupMap[model.FeedbackSeverity[i]].Grades[grade]++
				emphasisTotalGroup.Grades[grade]++
				if student.Gender == "男" {
					if _, ok := emphasisSCGroupMap[model.FeedbackSeverity[i]].MaleIdMap[studentId]; !ok {
						emphasisSCGroupMap[model.FeedbackSeverity[i]].NumMale++
					}
					emphasisSCGroupMap[model.FeedbackSeverity[i]].CountMale++
					emphasisSCGroupMap[model.FeedbackSeverity[i]].MaleIdMap[studentId]++

					if _, ok := emphasisTotalGroup.MaleIdMap[studentId]; !ok {
						emphasisTotalGroup.NumMale++
					}
					emphasisTotalGroup.CountMale++
					emphasisTotalGroup.MaleIdMap[studentId]++

					if _, ok := severityFCGroup.MaleIdMap[studentId]; !ok {
						severityFCGroup.NumMale++
					}
					severityFCGroup.CountMale++
					severityFCGroup.MaleIdMap[studentId]++
				} else if student.Gender == "女" {
					if _, ok := emphasisSCGroupMap[model.FeedbackSeverity[i]].FemaleIdMap[studentId]; !ok {
						emphasisSCGroupMap[model.FeedbackSeverity[i]].NumFemale++
					}
					emphasisSCGroupMap[model.FeedbackSeverity[i]].CountFemale++
					emphasisSCGroupMap[model.FeedbackSeverity[i]].FemaleIdMap[studentId]++

					if _, ok := emphasisTotalGroup.FemaleIdMap[studentId]; !ok {
						emphasisTotalGroup.NumFemale++
					}
					emphasisTotalGroup.CountFemale++
					emphasisTotalGroup.FemaleIdMap[studentId]++

					if _, ok := severityFCGroup.FemaleIdMap[studentId]; !ok {
						severityFCGroup.NumFemale++
					}
					severityFCGroup.CountFemale++
					severityFCGroup.FemaleIdMap[studentId]++
				}
				if strings.HasSuffix(grade, "级") {
					if _, ok := emphasisSCGroupMap[model.FeedbackSeverity[i]].UnderGraduateIdMap[studentId]; !ok {
						emphasisSCGroupMap[model.FeedbackSeverity[i]].NumUnderGraduates++
					}
					emphasisSCGroupMap[model.FeedbackSeverity[i]].CountUnderGraduates++
					emphasisSCGroupMap[model.FeedbackSeverity[i]].UnderGraduateIdMap[studentId]++

					if _, ok := emphasisTotalGroup.UnderGraduateIdMap[studentId]; !ok {
						emphasisTotalGroup.NumUnderGraduates++
					}
					emphasisTotalGroup.CountUnderGraduates++
					emphasisTotalGroup.UnderGraduateIdMap[studentId]++

					if _, ok := severityFCGroup.UnderGraduateIdMap[studentId]; !ok {
						severityFCGroup.NumUnderGraduates++
					}
					severityFCGroup.CountUnderGraduates++
					severityFCGroup.UnderGraduateIdMap[studentId]++
				} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
					if _, ok := emphasisSCGroupMap[model.FeedbackSeverity[i]].GraduateIdMap[studentId]; !ok {
						emphasisSCGroupMap[model.FeedbackSeverity[i]].NumGraduates++
					}
					emphasisSCGroupMap[model.FeedbackSeverity[i]].CountGraduates++
					emphasisSCGroupMap[model.FeedbackSeverity[i]].GraduateIdMap[studentId]++

					if _, ok := emphasisTotalGroup.GraduateIdMap[studentId]; !ok {
						emphasisTotalGroup.NumGraduates++
					}
					emphasisTotalGroup.CountGraduates++
					emphasisTotalGroup.GraduateIdMap[studentId]++

					if _, ok := severityFCGroup.GraduateIdMap[studentId]; !ok {
						severityFCGroup.NumGraduates++
					}
					severityFCGroup.CountGraduates++
					severityFCGroup.GraduateIdMap[studentId]++
				}
				if _, ok := emphasisSCGroupMap[model.FeedbackSeverity[i]].TotalIdMap[studentId]; !ok {
					emphasisSCGroupMap[model.FeedbackSeverity[i]].NumTotal++
				}
				emphasisSCGroupMap[model.FeedbackSeverity[i]].CountTotal++
				emphasisSCGroupMap[model.FeedbackSeverity[i]].TotalIdMap[studentId]++

				if _, ok := severityFCGroup.TotalIdMap[studentId]; !ok {
					severityFCGroup.NumTotal++
				}
				severityFCGroup.CountTotal++
				severityFCGroup.TotalIdMap[studentId]++

				if _, ok := emphasisTotalGroup.TotalIdMap[studentId]; !ok {
					emphasisTotalGroup.NumTotal++
				}
				emphasisTotalGroup.CountTotal++
				emphasisTotalGroup.TotalIdMap[studentId]++

				hasEmphasis = true
				if strings.HasSuffix(grade, "级") {
					categorySCGroupMap[category].CountUnderGraduateEmphasisInCategory++
					categoryFCGroupMap[category[0:1]].CountUnderGraduateEmphasisInCategory++
				} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
					categorySCGroupMap[category].CountGraduateEmphasisInCategory++
					categoryFCGroupMap[category[0:1]].CountGraduateEmphasisInCategory++
				}
				categorySCGroupMap[category].CountEmphasisInCategory++
				categoryFCGroupMap[category[0:1]].CountEmphasisInCategory++
			}
		}
		for i, m := range medicalDiagnosis {
			if m == 1 {
				emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].Grades[grade]++
				emphasisTotalGroup.Grades[grade]++
				if student.Gender == "男" {
					if _, ok := emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].MaleIdMap[studentId]; !ok {
						emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].NumMale++
					}
					emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].CountMale++
					emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].MaleIdMap[studentId]++

					if _, ok := emphasisTotalGroup.MaleIdMap[studentId]; !ok {
						emphasisTotalGroup.NumMale++
					}
					emphasisTotalGroup.CountMale++
					emphasisTotalGroup.MaleIdMap[studentId]++

					if _, ok := medicalDiagnosisFCGroup.MaleIdMap[studentId]; !ok {
						medicalDiagnosisFCGroup.NumMale++
					}
					medicalDiagnosisFCGroup.CountMale++
					medicalDiagnosisFCGroup.MaleIdMap[studentId]++
				} else if student.Gender == "女" {
					if _, ok := emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].FemaleIdMap[studentId]; !ok {
						emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].NumFemale++
					}
					emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].CountFemale++
					emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].FemaleIdMap[studentId]++

					if _, ok := emphasisTotalGroup.FemaleIdMap[studentId]; !ok {
						emphasisTotalGroup.NumFemale++
					}
					emphasisTotalGroup.CountFemale++
					emphasisTotalGroup.FemaleIdMap[studentId]++

					if _, ok := medicalDiagnosisFCGroup.FemaleIdMap[studentId]; !ok {
						medicalDiagnosisFCGroup.NumFemale++
					}
					medicalDiagnosisFCGroup.CountFemale++
					medicalDiagnosisFCGroup.FemaleIdMap[studentId]++
				}
				if strings.HasSuffix(grade, "级") {
					if _, ok := emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].UnderGraduateIdMap[studentId]; !ok {
						emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].NumUnderGraduates++
					}
					emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].CountUnderGraduates++
					emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].UnderGraduateIdMap[studentId]++

					if _, ok := emphasisTotalGroup.UnderGraduateIdMap[studentId]; !ok {
						emphasisTotalGroup.NumUnderGraduates++
					}
					emphasisTotalGroup.CountUnderGraduates++
					emphasisTotalGroup.UnderGraduateIdMap[studentId]++

					if _, ok := medicalDiagnosisFCGroup.UnderGraduateIdMap[studentId]; !ok {
						medicalDiagnosisFCGroup.NumUnderGraduates++
					}
					medicalDiagnosisFCGroup.CountUnderGraduates++
					medicalDiagnosisFCGroup.UnderGraduateIdMap[studentId]++
				} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
					if _, ok := emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].GraduateIdMap[studentId]; !ok {
						emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].NumGraduates++
					}
					emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].CountGraduates++
					emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].GraduateIdMap[studentId]++

					if _, ok := emphasisTotalGroup.GraduateIdMap[studentId]; !ok {
						emphasisTotalGroup.NumGraduates++
					}
					emphasisTotalGroup.CountGraduates++
					emphasisTotalGroup.GraduateIdMap[studentId]++

					if _, ok := medicalDiagnosisFCGroup.GraduateIdMap[studentId]; !ok {
						medicalDiagnosisFCGroup.NumGraduates++
					}
					medicalDiagnosisFCGroup.CountGraduates++
					medicalDiagnosisFCGroup.GraduateIdMap[studentId]++
				}
				if _, ok := emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].TotalIdMap[studentId]; !ok {
					emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].NumTotal++
				}
				emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].CountTotal++
				emphasisSCGroupMap[model.FeedbackMedicalDiagnosis[i]].TotalIdMap[studentId]++

				if _, ok := medicalDiagnosisFCGroup.TotalIdMap[studentId]; !ok {
					medicalDiagnosisFCGroup.NumTotal++
				}
				medicalDiagnosisFCGroup.CountTotal++
				medicalDiagnosisFCGroup.TotalIdMap[studentId]++

				if _, ok := emphasisTotalGroup.TotalIdMap[studentId]; !ok {
					emphasisTotalGroup.NumTotal++
				}
				emphasisTotalGroup.CountTotal++
				emphasisTotalGroup.TotalIdMap[studentId]++

				hasEmphasis = true
				if strings.HasSuffix(grade, "级") {
					categorySCGroupMap[category].CountUnderGraduateEmphasisInCategory++
					categoryFCGroupMap[category[0:1]].CountUnderGraduateEmphasisInCategory++
				} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
					categorySCGroupMap[category].CountGraduateEmphasisInCategory++
					categoryFCGroupMap[category[0:1]].CountGraduateEmphasisInCategory++
				}
				categorySCGroupMap[category].CountEmphasisInCategory++
				categoryFCGroupMap[category[0:1]].CountEmphasisInCategory++
			}
		}
		for i, c := range crisis {
			if c == 1 {
				emphasisSCGroupMap[model.FeedbackCrisis[i]].Grades[grade]++
				emphasisTotalGroup.Grades[grade]++
				if student.Gender == "男" {
					if _, ok := emphasisSCGroupMap[model.FeedbackCrisis[i]].MaleIdMap[studentId]; !ok {
						emphasisSCGroupMap[model.FeedbackCrisis[i]].NumMale++
					}
					emphasisSCGroupMap[model.FeedbackCrisis[i]].CountMale++
					emphasisSCGroupMap[model.FeedbackCrisis[i]].MaleIdMap[studentId]++

					if _, ok := emphasisTotalGroup.MaleIdMap[studentId]; !ok {
						emphasisTotalGroup.NumMale++
					}
					emphasisTotalGroup.CountMale++
					emphasisTotalGroup.MaleIdMap[studentId]++

					if _, ok := crisisFCGroup.MaleIdMap[studentId]; !ok {
						crisisFCGroup.NumMale++
					}
					crisisFCGroup.CountMale++
					crisisFCGroup.MaleIdMap[studentId]++
				} else if student.Gender == "女" {
					if _, ok := emphasisSCGroupMap[model.FeedbackCrisis[i]].FemaleIdMap[studentId]; !ok {
						emphasisSCGroupMap[model.FeedbackCrisis[i]].NumFemale++
					}
					emphasisSCGroupMap[model.FeedbackCrisis[i]].CountFemale++
					emphasisSCGroupMap[model.FeedbackCrisis[i]].FemaleIdMap[studentId]++

					if _, ok := emphasisTotalGroup.FemaleIdMap[studentId]; !ok {
						emphasisTotalGroup.NumFemale++
					}
					emphasisTotalGroup.CountFemale++
					emphasisTotalGroup.FemaleIdMap[studentId]++

					if _, ok := crisisFCGroup.FemaleIdMap[studentId]; !ok {
						crisisFCGroup.NumFemale++
					}
					crisisFCGroup.CountFemale++
					crisisFCGroup.FemaleIdMap[studentId]++
				}
				if strings.HasSuffix(grade, "级") {
					if _, ok := emphasisSCGroupMap[model.FeedbackCrisis[i]].UnderGraduateIdMap[studentId]; !ok {
						emphasisSCGroupMap[model.FeedbackCrisis[i]].NumUnderGraduates++
					}
					emphasisSCGroupMap[model.FeedbackCrisis[i]].CountUnderGraduates++
					emphasisSCGroupMap[model.FeedbackCrisis[i]].UnderGraduateIdMap[studentId]++

					if _, ok := emphasisTotalGroup.UnderGraduateIdMap[studentId]; !ok {
						emphasisTotalGroup.NumUnderGraduates++
					}
					emphasisTotalGroup.CountUnderGraduates++
					emphasisTotalGroup.UnderGraduateIdMap[studentId]++

					if _, ok := crisisFCGroup.UnderGraduateIdMap[studentId]; !ok {
						crisisFCGroup.NumUnderGraduates++
					}
					crisisFCGroup.CountUnderGraduates++
					crisisFCGroup.UnderGraduateIdMap[studentId]++
				} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
					if _, ok := emphasisSCGroupMap[model.FeedbackCrisis[i]].GraduateIdMap[studentId]; !ok {
						emphasisSCGroupMap[model.FeedbackCrisis[i]].NumGraduates++
					}
					emphasisSCGroupMap[model.FeedbackCrisis[i]].CountGraduates++
					emphasisSCGroupMap[model.FeedbackCrisis[i]].GraduateIdMap[studentId]++

					if _, ok := emphasisTotalGroup.GraduateIdMap[studentId]; !ok {
						emphasisTotalGroup.NumGraduates++
					}
					emphasisTotalGroup.CountGraduates++
					emphasisTotalGroup.GraduateIdMap[studentId]++

					if _, ok := crisisFCGroup.GraduateIdMap[studentId]; !ok {
						crisisFCGroup.NumGraduates++
					}
					crisisFCGroup.CountGraduates++
					crisisFCGroup.GraduateIdMap[studentId]++
				}
				if _, ok := emphasisSCGroupMap[model.FeedbackCrisis[i]].TotalIdMap[studentId]; !ok {
					emphasisSCGroupMap[model.FeedbackCrisis[i]].NumTotal++
				}
				emphasisSCGroupMap[model.FeedbackCrisis[i]].CountTotal++
				emphasisSCGroupMap[model.FeedbackCrisis[i]].TotalIdMap[studentId]++

				if _, ok := crisisFCGroup.TotalIdMap[studentId]; !ok {
					crisisFCGroup.NumTotal++
				}
				crisisFCGroup.CountTotal++
				crisisFCGroup.TotalIdMap[studentId]++

				if _, ok := emphasisTotalGroup.TotalIdMap[studentId]; !ok {
					emphasisTotalGroup.NumTotal++
				}
				emphasisTotalGroup.CountTotal++
				emphasisTotalGroup.TotalIdMap[studentId]++

				hasEmphasis = true
				if strings.HasSuffix(grade, "级") {
					categorySCGroupMap[category].CountUnderGraduateEmphasisInCategory++
					categoryFCGroupMap[category[0:1]].CountUnderGraduateEmphasisInCategory++
				} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
					categorySCGroupMap[category].CountGraduateEmphasisInCategory++
					categoryFCGroupMap[category[0:1]].CountGraduateEmphasisInCategory++
				}
				categorySCGroupMap[category].CountEmphasisInCategory++
				categoryFCGroupMap[category[0:1]].CountEmphasisInCategory++
			}
		}
		if hasEmphasis {
			if strings.HasSuffix(grade, "级") {
				if _, ok := categorySCGroupMap[category].HasEmphasisUnderGraduateIdMap[studentId]; !ok {
					categorySCGroupMap[category].NumHasEmphasisUnderGraduate++
				}
				categorySCGroupMap[category].CountHasEmphasisUnderGraduate++
				categorySCGroupMap[category].HasEmphasisUnderGraduateIdMap[studentId]++

				if _, ok := categoryFCGroupMap[category[0:1]].HasEmphasisUnderGraduateIdMap[studentId]; !ok {
					categoryFCGroupMap[category[0:1]].NumHasEmphasisUnderGraduate++
				}
				categoryFCGroupMap[category[0:1]].CountHasEmphasisUnderGraduate++
				categoryFCGroupMap[category[0:1]].HasEmphasisUnderGraduateIdMap[studentId]++
			} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
				if _, ok := categorySCGroupMap[category].HasEmphasisGraduateIdMap[studentId]; !ok {
					categorySCGroupMap[category].NumHasEmphasisGraduate++
				}
				categorySCGroupMap[category].CountHasEmphasisGraduate++
				categorySCGroupMap[category].HasEmphasisGraduateIdMap[studentId]++

				if _, ok := categoryFCGroupMap[category[0:1]].HasEmphasisGraduateIdMap[studentId]; !ok {
					categoryFCGroupMap[category[0:1]].NumHasEmphasisGraduate++
				}
				categoryFCGroupMap[category[0:1]].CountHasEmphasisGraduate++
				categoryFCGroupMap[category[0:1]].HasEmphasisGraduateIdMap[studentId]++
			}

			if _, ok := categorySCGroupMap[category].HasEmphasisIdMap[studentId]; !ok {
				categorySCGroupMap[category].NumHasEmphasis++
			}
			categorySCGroupMap[category].CountHasEmphasis++
			categorySCGroupMap[category].HasEmphasisIdMap[studentId]++

			if _, ok := categoryFCGroupMap[category[0:1]].HasEmphasisIdMap[studentId]; !ok {
				categoryFCGroupMap[category[0:1]].NumHasEmphasis++
			}
			categoryFCGroupMap[category[0:1]].CountHasEmphasis++
			categoryFCGroupMap[category[0:1]].HasEmphasisIdMap[studentId]++
		}

		// 来访&危机
		if r.TeacherFeedback.HasReservated && r.TeacherFeedback.HasCrisis {
			if strings.HasPrefix(r.TeacherFeedback.Category, "H") {
				advisoryConsulation.Consultation.CountBoth++
			} else {
				advisoryConsulation.Advisory.CountBoth++
			}
			advisoryConsulation.Total.CountBoth++
		} else if r.TeacherFeedback.HasReservated && !r.TeacherFeedback.HasCrisis {
			if strings.HasPrefix(r.TeacherFeedback.Category, "H") {
				advisoryConsulation.Consultation.CountOnlyReservated++
			} else {
				advisoryConsulation.Advisory.CountOnlyReservated++
			}
			advisoryConsulation.Total.CountOnlyReservated++
		} else if !r.TeacherFeedback.HasReservated && r.TeacherFeedback.HasCrisis {
			if strings.HasPrefix(r.TeacherFeedback.Category, "H") {
				advisoryConsulation.Consultation.CountOnlyCrisis++
			} else {
				advisoryConsulation.Advisory.CountOnlyCrisis++
			}
			advisoryConsulation.Total.CountOnlyCrisis++
		} else if !r.TeacherFeedback.HasReservated && !r.TeacherFeedback.HasCrisis {
			if strings.HasPrefix(r.TeacherFeedback.Category, "H") {
				advisoryConsulation.Consultation.CountNeither++
			} else {
				advisoryConsulation.Advisory.CountNeither++
			}
			advisoryConsulation.Total.CountNeither++
		}
		if r.TeacherFeedback.IsSendNotify {
			if strings.HasPrefix(r.TeacherFeedback.Category, "H") {
				advisoryConsulation.Consultation.CountSendNotify++
			} else {
				advisoryConsulation.Advisory.CountSendNotify++
			}
			advisoryConsulation.Total.CountSendNotify++
		}

		// 会商与危机干预
		if strings.HasPrefix(r.TeacherFeedback.Category, "H") || r.TeacherFeedback.HasCrisis {
			cc := &ConsultationCrisis{
				Date:                 r.StartTime.Format("2006/01/02"),
				Fullname:             student.Fullname,
				Username:             student.Username,
				Gender:               student.Gender,
				School:               student.School,
				TeacherFullname:      teacher.Fullname,
				SchoolContact:        r.TeacherFeedback.SchoolContact,
				ConsultationOrCrisis: make([]string, 0),
				Category:             r.TeacherFeedback.Category,
				EmphasisStr:          r.TeacherFeedback.GetEmphasisStr(),
				CrisisLevel:          strconv.Itoa(student.CrisisLevel),
			}
			if strings.HasSuffix(grade, "级") {
				cc.Academic = "本科生"
			} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
				cc.Academic = "研究生"
			}
			if strings.HasPrefix(r.TeacherFeedback.Category, "H") {
				cc.ConsultationOrCrisis = append(cc.ConsultationOrCrisis, "会商")
			}
			if r.TeacherFeedback.HasCrisis {
				cc.ConsultationOrCrisis = append(cc.ConsultationOrCrisis, "危机处理")
			}
			if r.TeacherFeedback.HasReservated {
				cc.ReservationStatus = "有预约"
			} else {
				cc.ReservationStatus = "未预约"
			}
			if r.TeacherFeedback.IsSendNotify {
				cc.IsSendNotify = "是"
			} else {
				cc.IsSendNotify = "否"
			}
			consultationCrisisList = append(consultationCrisisList, cc)
		}
	}
	for _, scGroup := range categorySCGroupMap {
		scGroup.Ratio = float64(scGroup.CountTotal) / float64(categoryTotalGroup.CountTotal)
	}
	for _, fcGroup := range categoryFCGroupMap {
		fcGroup.Ratio = float64(fcGroup.CountTotal) / float64(categoryTotalGroup.CountTotal)
	}
	for _, scGroup := range emphasisSCGroupMap {
		scGroup.Ratio = float64(scGroup.CountTotal) / float64(emphasisTotalGroup.CountTotal)
	}
	// 分析年级数据
	allGrades := make([]string, 0)
	for g := range categoryTotalGroup.Grades {
		allGrades = append(allGrades, g)
	}
	sort.Slice(allGrades, func(i, j int) bool {
		if strings.HasSuffix(allGrades[i], "级") {
			if strings.HasSuffix(allGrades[j], "级") {
				return allGrades[i] > allGrades[j]
			} else {
				return true
			}
		} else if strings.HasSuffix(allGrades[i], "硕") {
			if strings.HasSuffix(allGrades[j], "级") {
				return false
			} else if strings.HasSuffix(allGrades[j], "硕") {
				return allGrades[i] > allGrades[j]
			} else if strings.HasSuffix(allGrades[j], "博") {
				return true
			}
		} else if strings.HasSuffix(allGrades[i], "博") {
			if strings.HasSuffix(allGrades[j], "博") {
				return allGrades[i] > allGrades[j]
			} else {
				return false
			}
		}
		return true
	})
	// 开始写入文件
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error
	xlsx.SetDefaultFont(11, "宋体")
	file = xlsx.NewFile()
	//==========咨询情况汇总表=========
	sheet, err = file.AddSheet("咨询情况汇总")
	if err != nil {
		return re.NewRError("fail to create category sheet", err)
	}
	// 第一表头
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("评估分类")
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("来访情况分项")
	for i := 3; i <= len(allGrades)+5; i++ {
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
	}
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("分项合计")
	for i := 0; i < 3; i++ {
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
	}
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("大类分项合计")
	for i := 0; i < 3; i++ {
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
	}
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("大类人数统计")
	for i := 0; i < 2; i++ {
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
	}
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("总计")
	for i := 0; i < 3; i++ {
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
	}
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("重点情况频次大类分项统计")
	for i := 0; i < 2; i++ {
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
	}
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("平均标记值")
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("重点情况人次统计")
	for i := 0; i < 2; i++ {
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
	}
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("重点情况人数统计")
	for i := 0; i < 2; i++ {
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
	}
	// 合并第一表头
	cell = row.Cells[0]
	cell.Merge(1, 0)
	cell.SetStyle(textCenterStyle)
	cell = row.Cells[2]
	cell.Merge(len(allGrades)+3, 0)
	cell.SetStyle(textCenterStyle)
	cell = row.Cells[len(allGrades)+6]
	cell.Merge(3, 0)
	cell.SetStyle(textCenterGrayStyle)
	cell = row.Cells[len(allGrades)+10]
	cell.Merge(3, 0)
	cell.SetStyle(textCenterGreenStyle)
	cell = row.Cells[len(allGrades)+14]
	cell.Merge(2, 0)
	cell.SetStyle(textCenterGreenStyle)
	cell = row.Cells[len(allGrades)+17]
	cell.Merge(3, 0)
	cell.SetStyle(textCenterOrangeStyle)
	cell = row.Cells[len(allGrades)+21]
	cell.Merge(2, 0)
	cell.SetStyle(textCenterGreenStyle)
	cell = row.Cells[len(allGrades)+24]
	cell.SetStyle(textCenterGreenStyle)
	cell = row.Cells[len(allGrades)+25]
	cell.Merge(2, 0)
	cell.SetStyle(textCenterGreenStyle)
	cell = row.Cells[len(allGrades)+28]
	cell.Merge(2, 0)
	cell.SetStyle(textCenterGreenStyle)
	// 第二表头
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(textCenterStyle)
	cell.SetValue("大类")
	cell = row.AddCell()
	cell.SetStyle(textCenterStyle)
	cell.SetValue("评估")
	for _, g := range allGrades {
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(g)
	}
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("辅导员")
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("教师")
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("家属")
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("其他")
	cell = row.AddCell()
	cell.SetStyle(bgGrayStyle)
	cell.SetValue("合计（男）")
	cell = row.AddCell()
	cell.SetStyle(bgGrayStyle)
	cell.SetValue("合计（女）")
	cell = row.AddCell()
	cell.SetStyle(bgGrayStyle)
	cell.SetValue("合计（本科生）")
	cell = row.AddCell()
	cell.SetStyle(bgGrayStyle)
	cell.SetValue("合计（研究生）")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("合计（男）")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("合计（女）")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("合计（本科生）")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("合计（研究生）")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("本科生")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("研究生")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("小计")
	cell = row.AddCell()
	cell.SetStyle(bgOrangeStyle)
	cell.SetValue("会谈总计")
	cell = row.AddCell()
	cell.SetStyle(bgOrangeStyle)
	cell.SetValue("百分比")
	cell = row.AddCell()
	cell.SetStyle(bgOrangeStyle)
	cell.SetValue("大类总计")
	cell = row.AddCell()
	cell.SetStyle(bgOrangeStyle)
	cell.SetValue("大类百分比")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("本科生")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("研究生")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("小计")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("频次/会谈")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("本科生")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("研究生")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("小计")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("本科生")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("研究生")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("小计")
	// 咨询情况数据
	firstRowOfFcGroup := 2 // 标记当前fcGroup的行号，以便最后合并单元格
	for _, fcGroup := range categoryFCGroup {
		for _, scGroup := range fcGroup.SecondaryGroups {
			row = sheet.AddRow()
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			cell.SetValue(strings.Split(fcGroup.GroupName, " ")[1])
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			cell.SetValue(scGroup.GroupName)
			for _, g := range allGrades {
				cell = row.AddCell()
				cell.SetStyle(borderStyle)
				if scGroup.Grades[g] > 0 {
					cell.SetValue(scGroup.Grades[g])
				}
			}
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			if scGroup.Instructor > 0 {
				cell.SetValue(scGroup.Instructor)
			}
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			if scGroup.Teacher > 0 {
				cell.SetValue(scGroup.Teacher)
			}
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			if scGroup.Family > 0 {
				cell.SetValue(scGroup.Family)
			}
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			if scGroup.Others > 0 {
				cell.SetValue(scGroup.Others)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGrayStyle)
			if scGroup.CountMale > 0 {
				cell.SetValue(scGroup.CountMale)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGrayStyle)
			if scGroup.CountFemale > 0 {
				cell.SetValue(scGroup.CountFemale)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGrayStyle)
			if scGroup.CountUnderGraduates > 0 {
				cell.SetValue(scGroup.CountUnderGraduates)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGrayStyle)
			if scGroup.CountGraduates > 0 {
				cell.SetValue(scGroup.CountGraduates)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			if fcGroup.CountMale > 0 {
				cell.SetValue(fcGroup.CountMale)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			if fcGroup.CountFemale > 0 {
				cell.SetValue(fcGroup.CountFemale)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			if fcGroup.CountUnderGraduates > 0 {
				cell.SetValue(fcGroup.CountUnderGraduates)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			if fcGroup.CountGraduates > 0 {
				cell.SetValue(fcGroup.CountGraduates)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			if fcGroup.NumUnderGraduates > 0 {
				cell.SetValue(fcGroup.NumUnderGraduates)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			if fcGroup.NumGraduates > 0 {
				cell.SetValue(fcGroup.NumGraduates)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			if fcGroup.NumTotal > 0 {
				cell.SetValue(fcGroup.NumTotal)
			}
			cell = row.AddCell()
			cell.SetStyle(bgOrangeStyle)
			cell.SetValue(scGroup.CountTotal)
			cell = row.AddCell()
			cell.SetStyle(bgOrangeStyle)
			cell.SetValue(fmt.Sprintf("%.2f%%", scGroup.Ratio*100))
			cell = row.AddCell()
			cell.SetStyle(bgOrangeStyle)
			cell.SetValue(fcGroup.CountTotal)
			cell = row.AddCell()
			cell.SetStyle(bgOrangeStyle)
			cell.SetValue(fmt.Sprintf("%.2f%%", fcGroup.Ratio*100))
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			cell.SetValue(fcGroup.CountUnderGraduateEmphasisInCategory)
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			cell.SetValue(fcGroup.CountGraduateEmphasisInCategory)
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			cell.SetValue(fcGroup.CountEmphasisInCategory)
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			cell.SetValue(fmt.Sprintf("%.2f", float64(fcGroup.CountEmphasisInCategory)/float64(fcGroup.CountHasEmphasis)))
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			cell.SetValue(fcGroup.CountHasEmphasisUnderGraduate)
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			cell.SetValue(fcGroup.CountHasEmphasisGraduate)
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			cell.SetValue(fcGroup.CountHasEmphasis)
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			cell.SetValue(fcGroup.NumHasEmphasisUnderGraduate)
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			cell.SetValue(fcGroup.NumHasEmphasisGraduate)
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			cell.SetValue(fcGroup.NumHasEmphasis)
		}
		row = sheet.Rows[firstRowOfFcGroup]
		cell = row.Cells[0]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textCenterStyle)
		cell = row.Cells[len(allGrades)+10]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+11]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+12]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+13]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+14]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+15]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+16]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+19]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightOrangeStyle)
		cell = row.Cells[len(allGrades)+20]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightOrangeStyle)
		cell = row.Cells[len(allGrades)+21]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+22]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+23]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+24]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+25]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+26]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+27]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+28]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+29]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+30]
		cell.Merge(0, len(fcGroup.SecondaryGroups)-1)
		cell.SetStyle(textRightGreenStyle)
		firstRowOfFcGroup += len(fcGroup.SecondaryGroups)
	}
	// 总计
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalGroup.GroupName)
	for _, g := range allGrades {
		cell = row.AddCell()
		cell.SetStyle(bgYellowStyle)
		cell.SetValue(categoryTotalGroup.Grades[g])
	}
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalGroup.Instructor)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalGroup.Teacher)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalGroup.Family)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalGroup.Others)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalGroup.CountMale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalGroup.CountFemale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalGroup.CountUnderGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalGroup.CountGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalGroup.CountMale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalGroup.CountFemale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalGroup.CountUnderGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalGroup.CountGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalGroup.NumUnderGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalGroup.NumGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalGroup.NumTotal)
	cell = row.AddCell()
	cell.SetStyle(bgRedStyle)
	cell.SetValue(categoryTotalGroup.CountTotal)
	// 百分比
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue("百分比")
	for _, g := range allGrades {
		cell = row.AddCell()
		cell.SetStyle(bgYellowStyle)
		cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalGroup.Grades[g])/float64(categoryTotalGroup.CountTotal)*100))
	}
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalGroup.Instructor)/float64(categoryTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalGroup.Teacher)/float64(categoryTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalGroup.Family)/float64(categoryTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalGroup.Others)/float64(categoryTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalGroup.CountMale)/float64(categoryTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalGroup.CountFemale)/float64(categoryTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalGroup.CountUnderGraduates)/float64(categoryTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalGroup.CountGraduates)/float64(categoryTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalGroup.CountMale)/float64(categoryTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalGroup.CountFemale)/float64(categoryTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalGroup.CountUnderGraduates)/float64(categoryTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalGroup.CountGraduates)/float64(categoryTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalGroup.NumUnderGraduates)/float64(categoryTotalGroup.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalGroup.NumGraduates)/float64(categoryTotalGroup.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalGroup.NumUnderGraduates+categoryTotalGroup.NumGraduates)/float64(categoryTotalGroup.NumTotal)*100))
	// 调整列宽
	sheet.SetColWidth(0, 0, 11)
	sheet.SetColWidth(1, 1, 22)
	sheet.SetColWidth(len(allGrades)+8, len(allGrades)+9, 11)
	sheet.SetColWidth(len(allGrades)+12, len(allGrades)+13, 11)

	//==========重点情况汇总==========
	sheet, err = file.AddSheet("重点情况汇总")
	if err != nil {
		return re.NewRError("fail to create emphasis sheet", err)
	}
	// 第一表头
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("评估分类")
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("来访情况分项")
	for i := 3; i <= len(allGrades)+1; i++ {
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
	}
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("分项合计")
	for i := 0; i < 3; i++ {
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
	}
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("大类分项合计")
	for i := 0; i < 3; i++ {
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
	}
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("总计")
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	// 合并第一表头
	cell = row.Cells[0]
	cell.Merge(1, 0)
	cell.SetStyle(textCenterStyle)
	cell = row.Cells[2]
	cell.Merge(len(allGrades)-1, 0)
	cell.SetStyle(textCenterStyle)
	cell = row.Cells[len(allGrades)+2]
	cell.Merge(3, 0)
	cell.SetStyle(textCenterGrayStyle)
	cell = row.Cells[len(allGrades)+6]
	cell.Merge(3, 0)
	cell.SetStyle(textCenterGreenStyle)
	cell = row.Cells[len(allGrades)+10]
	cell.Merge(1, 0)
	cell.SetStyle(textCenterOrangeStyle)
	// 第二表头
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(textCenterStyle)
	cell.SetValue("大类")
	cell = row.AddCell()
	cell.SetStyle(textCenterStyle)
	cell.SetValue("评估")
	for _, g := range allGrades {
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(g)
	}
	cell = row.AddCell()
	cell.SetStyle(bgGrayStyle)
	cell.SetValue("合计（男）")
	cell = row.AddCell()
	cell.SetStyle(bgGrayStyle)
	cell.SetValue("合计（女）")
	cell = row.AddCell()
	cell.SetStyle(bgGrayStyle)
	cell.SetValue("合计（本科生）")
	cell = row.AddCell()
	cell.SetStyle(bgGrayStyle)
	cell.SetValue("合计（研究生）")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("合计（男）")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("合计（女）")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("合计（本科生）")
	cell = row.AddCell()
	cell.SetStyle(bgGreenStyle)
	cell.SetValue("合计（研究生）")
	cell = row.AddCell()
	cell.SetStyle(bgOrangeStyle)
	cell.SetValue("频次总计")
	cell = row.AddCell()
	cell.SetStyle(bgOrangeStyle)
	cell.SetValue("百分比")
	// 重点情况数据
	firstRowOfFcGroup = 2
	// 严重程度
	for _, scGroup := range severityFCGroup.SecondaryGroups {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(severityFCGroup.GroupName)
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(scGroup.GroupName)
		for _, g := range allGrades {
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			if scGroup.Grades[g] > 0 {
				cell.SetValue(scGroup.Grades[g])
			}
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scGroup.CountMale > 0 {
			cell.SetValue(scGroup.CountMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scGroup.CountFemale > 0 {
			cell.SetValue(scGroup.CountFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scGroup.CountUnderGraduates > 0 {
			cell.SetValue(scGroup.CountUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scGroup.CountGraduates > 0 {
			cell.SetValue(scGroup.CountGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if severityFCGroup.CountMale > 0 {
			cell.SetValue(severityFCGroup.CountMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if severityFCGroup.CountFemale > 0 {
			cell.SetValue(severityFCGroup.CountFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if severityFCGroup.CountUnderGraduates > 0 {
			cell.SetValue(severityFCGroup.CountUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if severityFCGroup.CountGraduates > 0 {
			cell.SetValue(severityFCGroup.CountGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgOrangeStyle)
		cell.SetValue(scGroup.CountTotal)
		cell = row.AddCell()
		cell.SetStyle(bgOrangeStyle)
		cell.SetValue(fmt.Sprintf("%.2f%%", scGroup.Ratio*100))
	}
	row = sheet.Rows[firstRowOfFcGroup]
	cell = row.Cells[0]
	cell.Merge(0, len(severityFCGroup.SecondaryGroups)-1)
	cell.SetStyle(textCenterStyle)
	cell = row.Cells[len(allGrades)+6]
	cell.Merge(0, len(severityFCGroup.SecondaryGroups)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+7]
	cell.Merge(0, len(severityFCGroup.SecondaryGroups)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+8]
	cell.Merge(0, len(severityFCGroup.SecondaryGroups)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+9]
	cell.Merge(0, len(severityFCGroup.SecondaryGroups)-1)
	cell.SetStyle(textRightGreenStyle)
	firstRowOfFcGroup += len(severityFCGroup.SecondaryGroups)
	// 疑似或明确的医疗诊断
	for _, scGroup := range medicalDiagnosisFCGroup.SecondaryGroups {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(medicalDiagnosisFCGroup.GroupName)
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(scGroup.GroupName)
		for _, g := range allGrades {
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			if scGroup.Grades[g] > 0 {
				cell.SetValue(scGroup.Grades[g])
			}
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scGroup.CountMale > 0 {
			cell.SetValue(scGroup.CountMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scGroup.CountFemale > 0 {
			cell.SetValue(scGroup.CountFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scGroup.CountUnderGraduates > 0 {
			cell.SetValue(scGroup.CountUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scGroup.CountGraduates > 0 {
			cell.SetValue(scGroup.CountGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if medicalDiagnosisFCGroup.CountMale > 0 {
			cell.SetValue(medicalDiagnosisFCGroup.CountMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if medicalDiagnosisFCGroup.CountFemale > 0 {
			cell.SetValue(medicalDiagnosisFCGroup.CountFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if medicalDiagnosisFCGroup.CountUnderGraduates > 0 {
			cell.SetValue(medicalDiagnosisFCGroup.CountUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if medicalDiagnosisFCGroup.CountGraduates > 0 {
			cell.SetValue(medicalDiagnosisFCGroup.CountGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgOrangeStyle)
		cell.SetValue(scGroup.CountTotal)
		cell = row.AddCell()
		cell.SetStyle(bgOrangeStyle)
		cell.SetValue(fmt.Sprintf("%.2f%%", scGroup.Ratio*100))
	}
	row = sheet.Rows[firstRowOfFcGroup]
	cell = row.Cells[0]
	cell.Merge(0, len(medicalDiagnosisFCGroup.SecondaryGroups)-1)
	cell.SetStyle(textCenterStyle)
	cell = row.Cells[len(allGrades)+6]
	cell.Merge(0, len(medicalDiagnosisFCGroup.SecondaryGroups)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+7]
	cell.Merge(0, len(medicalDiagnosisFCGroup.SecondaryGroups)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+8]
	cell.Merge(0, len(medicalDiagnosisFCGroup.SecondaryGroups)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+9]
	cell.Merge(0, len(medicalDiagnosisFCGroup.SecondaryGroups)-1)
	cell.SetStyle(textRightGreenStyle)
	firstRowOfFcGroup += len(medicalDiagnosisFCGroup.SecondaryGroups)
	// 危急情况
	for _, scGroup := range crisisFCGroup.SecondaryGroups {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(crisisFCGroup.GroupName)
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(scGroup.GroupName)
		for _, g := range allGrades {
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			if scGroup.Grades[g] > 0 {
				cell.SetValue(scGroup.Grades[g])
			}
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scGroup.CountMale > 0 {
			cell.SetValue(scGroup.CountMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scGroup.CountFemale > 0 {
			cell.SetValue(scGroup.CountFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scGroup.CountUnderGraduates > 0 {
			cell.SetValue(scGroup.CountUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scGroup.CountGraduates > 0 {
			cell.SetValue(scGroup.CountGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if crisisFCGroup.CountMale > 0 {
			cell.SetValue(crisisFCGroup.CountMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if crisisFCGroup.CountFemale > 0 {
			cell.SetValue(crisisFCGroup.CountFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if crisisFCGroup.CountUnderGraduates > 0 {
			cell.SetValue(crisisFCGroup.CountUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if crisisFCGroup.CountGraduates > 0 {
			cell.SetValue(crisisFCGroup.CountGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgOrangeStyle)
		cell.SetValue(scGroup.CountTotal)
		cell = row.AddCell()
		cell.SetStyle(bgOrangeStyle)
		cell.SetValue(fmt.Sprintf("%.2f%%", scGroup.Ratio*100))
	}
	row = sheet.Rows[firstRowOfFcGroup]
	cell = row.Cells[0]
	cell.Merge(0, len(crisisFCGroup.SecondaryGroups)-1)
	cell.SetStyle(textCenterStyle)
	cell = row.Cells[len(allGrades)+6]
	cell.Merge(0, len(crisisFCGroup.SecondaryGroups)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+7]
	cell.Merge(0, len(crisisFCGroup.SecondaryGroups)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+8]
	cell.Merge(0, len(crisisFCGroup.SecondaryGroups)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+9]
	cell.Merge(0, len(crisisFCGroup.SecondaryGroups)-1)
	cell.SetStyle(textRightGreenStyle)
	firstRowOfFcGroup += len(crisisFCGroup.SecondaryGroups)
	// 总计
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalGroup.GroupName)
	for _, g := range allGrades {
		cell = row.AddCell()
		cell.SetStyle(bgYellowStyle)
		cell.SetValue(emphasisTotalGroup.Grades[g])
	}
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalGroup.CountMale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalGroup.CountFemale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalGroup.CountUnderGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalGroup.CountGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalGroup.CountMale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalGroup.CountFemale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalGroup.CountUnderGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalGroup.CountGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgRedStyle)
	cell.SetValue(emphasisTotalGroup.CountTotal)
	// 百分比
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue("百分比")
	for _, g := range allGrades {
		cell = row.AddCell()
		cell.SetStyle(bgYellowStyle)
		cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalGroup.Grades[g])/float64(emphasisTotalGroup.CountTotal)*100))
	}
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalGroup.CountMale)/float64(emphasisTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalGroup.CountFemale)/float64(emphasisTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalGroup.CountUnderGraduates)/float64(emphasisTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalGroup.CountGraduates)/float64(emphasisTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalGroup.CountMale)/float64(emphasisTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalGroup.CountFemale)/float64(emphasisTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalGroup.CountUnderGraduates)/float64(emphasisTotalGroup.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalGroup.CountGraduates)/float64(emphasisTotalGroup.CountTotal)*100))
	// 调整列宽
	sheet.SetColWidth(0, 0, 15.5)
	sheet.SetColWidth(1, 1, 24)
	sheet.SetColWidth(len(allGrades)+4, len(allGrades)+5, 11)
	sheet.SetColWidth(len(allGrades)+8, len(allGrades)+9, 11)

	//==========来访&危机情况表==========
	sheet, err = file.AddSheet("来访&危机情况表")
	if err != nil {
		return re.NewRError("fail to create reservation crisis sheet", err)
	}

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetValue("咨询")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell = row.AddCell()
	cell.SetValue("无危机")
	cell = row.AddCell()
	cell.SetValue("有危机")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetValue("有预约")
	cell = row.AddCell()
	cell.SetValue(advisoryConsulation.Advisory.CountOnlyReservated)
	cell = row.AddCell()
	cell.SetValue(advisoryConsulation.Advisory.CountBoth)
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetValue("无预约")
	cell = row.AddCell()
	cell.SetValue(advisoryConsulation.Advisory.CountNeither)
	cell = row.AddCell()
	cell.SetValue(advisoryConsulation.Advisory.CountOnlyCrisis)
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetValue("发危机通报人次")
	cell = row.AddCell()
	cell.SetValue(advisoryConsulation.Advisory.CountSendNotify)

	row = sheet.AddRow()
	row = sheet.AddRow()

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetValue("会商")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell = row.AddCell()
	cell.SetValue("无危机")
	cell = row.AddCell()
	cell.SetValue("有危机")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetValue("有预约")
	cell = row.AddCell()
	cell.SetValue(advisoryConsulation.Consultation.CountOnlyReservated)
	cell = row.AddCell()
	cell.SetValue(advisoryConsulation.Consultation.CountBoth)
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetValue("无预约")
	cell = row.AddCell()
	cell.SetValue(advisoryConsulation.Consultation.CountNeither)
	cell = row.AddCell()
	cell.SetValue(advisoryConsulation.Consultation.CountOnlyCrisis)
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetValue("发危机通报人次")
	cell = row.AddCell()
	cell.SetValue(advisoryConsulation.Consultation.CountSendNotify)

	row = sheet.AddRow()
	row = sheet.AddRow()

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetValue("总数")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell = row.AddCell()
	cell.SetValue("无危机")
	cell = row.AddCell()
	cell.SetValue("有危机")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetValue("有预约")
	cell = row.AddCell()
	cell.SetValue(advisoryConsulation.Total.CountOnlyReservated)
	cell = row.AddCell()
	cell.SetValue(advisoryConsulation.Total.CountBoth)
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetValue("无预约")
	cell = row.AddCell()
	cell.SetValue(advisoryConsulation.Total.CountNeither)
	cell = row.AddCell()
	cell.SetValue(advisoryConsulation.Total.CountOnlyCrisis)
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetValue("发危机通报人次")
	cell = row.AddCell()
	cell.SetValue(advisoryConsulation.Total.CountSendNotify)
	// 调整列宽
	sheet.SetColWidth(0, 0, 15)

	//==========会商与危机干预表==========
	sheet, err = file.AddSheet("会商与危机干预表")
	if err != nil {
		return re.NewRError("fail to create consultation and crisis sheet", err)
	}
	// 表头
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetValue("日期")
	cell = row.AddCell()
	cell.SetValue("姓名")
	cell = row.AddCell()
	cell.SetValue("学号")
	cell = row.AddCell()
	cell.SetValue("性别")
	cell = row.AddCell()
	cell.SetValue("学历")
	cell = row.AddCell()
	cell.SetValue("院系")
	cell = row.AddCell()
	cell.SetValue("接待咨询师")
	cell = row.AddCell()
	cell.SetValue("院系联系人")
	cell = row.AddCell()
	cell.SetValue("会商or危机处理")
	cell = row.AddCell()
	cell.SetValue("来访情况")
	cell = row.AddCell()
	cell.SetValue("评估分类")
	cell = row.AddCell()
	cell.SetValue("重点标记")
	cell = row.AddCell()
	cell.SetValue("星级")
	cell = row.AddCell()
	cell.SetValue("是否发危机通报")
	for _, cc := range consultationCrisisList {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetValue(cc.Date)
		cell = row.AddCell()
		cell.SetValue(cc.Fullname)
		cell = row.AddCell()
		cell.SetValue(cc.Username)
		cell = row.AddCell()
		cell.SetValue(cc.Gender)
		cell = row.AddCell()
		cell.SetValue(cc.Academic)
		cell = row.AddCell()
		cell.SetValue(cc.School)
		cell = row.AddCell()
		cell.SetValue(cc.TeacherFullname)
		cell = row.AddCell()
		cell.SetValue(cc.SchoolContact)
		cell = row.AddCell()
		cell.SetValue(strings.Join(cc.ConsultationOrCrisis, "、"))
		cell = row.AddCell()
		cell.SetValue(cc.ReservationStatus)
		cell = row.AddCell()
		cell.SetValue(model.FeedbackAllCategoryMap[cc.Category])
		cell = row.AddCell()
		cell.SetValue(cc.EmphasisStr)
		cell = row.AddCell()
		cell.SetValue(cc.CrisisLevel)
		cell = row.AddCell()
		cell.SetValue(cc.IsSendNotify)
	}
	// 调整列宽
	sheet.SetColWidth(0, 0, 10)
	sheet.SetColWidth(2, 2, 10)
	sheet.SetColWidth(5, 5, 15)
	sheet.SetColWidth(6, 6, 10)
	sheet.SetColWidth(7, 7, 10)
	sheet.SetColWidth(8, 8, 13.5)
	sheet.SetColWidth(10, 10, 20)
	sheet.SetColWidth(11, 11, 20)

	//==========学生咨询情况统计表==========
	sheet, err = file.AddSheet("学生咨询情况统计表")
	if err != nil {
		return re.NewRError("fail to create student reservations sheet", err)
	}
	// 表头
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetValue("姓名")
	cell = row.AddCell()
	cell.SetValue("学号")
	cell = row.AddCell()
	cell.SetValue("性别")
	cell = row.AddCell()
	cell.SetValue("院系")
	cell = row.AddCell()
	cell.SetValue("最近一次咨询分类")
	cell = row.AddCell()
	cell.SetValue("最后一次重点标记")
	cell = row.AddCell()
	cell.SetValue("最后一次来访&危机情况")
	cell = row.AddCell()
	cell.SetValue("最后一次的咨询师")
	cell = row.AddCell()
	cell.SetValue("星级")
	for sid, stu := range studentMap {
		reservations := studentReservationsMap[sid]
		if len(reservations) == 0 {
			continue
		}
		sort.Sort(sort.Reverse(ByStartTimeOfReservation(reservations)))
		r := reservations[0]
		if r.TeacherFeedback.IsEmpty() {
			continue
		}
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetValue(stu.Fullname)
		cell = row.AddCell()
		cell.SetValue(stu.Username)
		cell = row.AddCell()
		cell.SetValue(stu.Gender)
		cell = row.AddCell()
		cell.SetValue(stu.School)
		cell = row.AddCell()
		cell.SetValue(model.FeedbackAllCategoryMap[r.TeacherFeedback.Category])
		cell = row.AddCell()
		cell.SetValue(r.TeacherFeedback.GetEmphasisStr())
		cell = row.AddCell()
		s := ""
		if r.TeacherFeedback.HasReservated {
			s += "有预约 "
		} else {
			s += "无预约 "
		}
		if r.TeacherFeedback.HasCrisis {
			s += "有危机"
		} else {
			s += "无危机"
		}
		cell.SetValue(s)
		cell = row.AddCell()
		if t, ok := teacherMap[r.TeacherId]; ok {
			cell.SetValue(t.Fullname)
		}
		cell = row.AddCell()
		cell.SetValue(stu.CrisisLevel)
	}
	// 调整列宽
	sheet.SetColWidth(1, 1, 10)
	sheet.SetColWidth(3, 3, 15)
	sheet.SetColWidth(4, 4, 20)
	sheet.SetColWidth(5, 5, 20)
	sheet.SetColWidth(6, 6, 20)
	sheet.SetColWidth(7, 7, 16)

	err = file.Save(path)
	if err != nil {
		return re.NewRError("fail to save file to path", err)
	}
	return nil
}

//================================================================
//====================咨询师工作量统计==============================
//================================================================
type TeacherWorkload struct {
	TeacherId           string
	Fullname            string
	StudentIdMap        map[string]int
	TotalNum            int
	TotalCount          int
	UnderGraduateIdMap  map[string]int
	NumUnderGraduates   int // 本科生人数
	CountUnderGraduates int // 本科生人次
	GraduateIdMap       map[string]int
	NumGraduates        int // 研究生人数
	CountGraduates      int // 研究生人次
	FirstClassWorkloads map[string]*FirstClassWorkload
}

type FirstClassWorkload struct {
	FirstClass          string
	UnderGraduateIdMap  map[string]int
	NumUnderGraduates   int // 本科生人数
	CountUnderGraduates int // 本科生人次
	GraduateIdMap       map[string]int
	NumGraduates        int // 研究生人数
	CountGraduates      int // 研究生人次
}

func (w *Workflow) ExportWorkloadToFile(reservations []*model.Reservation, path string) error {
	// 建立存储结构
	teacherWorkloads := make([]*TeacherWorkload, 0)
	teacherWorkloadMap := make(map[string]*TeacherWorkload)
	for _, r := range reservations {
		if r.Status != model.RESERVATION_STATUS_RESERVATED || r.TeacherFeedback.IsEmpty() {
			continue
		}
		if tWork, ok := teacherWorkloadMap[r.TeacherId]; !ok {
			teacher, err := w.mongoClient.GetTeacherById(r.TeacherId)
			if err != nil || teacher == nil || teacher.UserType != model.USER_TYPE_TEACHER {
				continue
			}
			tWork = &TeacherWorkload{
				TeacherId:           r.TeacherId,
				Fullname:            teacher.Fullname,
				StudentIdMap:        make(map[string]int),
				UnderGraduateIdMap:  make(map[string]int),
				GraduateIdMap:       make(map[string]int),
				FirstClassWorkloads: make(map[string]*FirstClassWorkload),
			}
			for fi, firstCategory := range model.FeedbackFirstCategoryMap {
				if fi == "" {
					continue
				}
				tWork.FirstClassWorkloads[fi] = &FirstClassWorkload{
					FirstClass:         firstCategory,
					UnderGraduateIdMap: make(map[string]int),
					GraduateIdMap:      make(map[string]int),
				}
			}
			teacherWorkloads = append(teacherWorkloads, tWork)
			teacherWorkloadMap[r.TeacherId] = tWork
		}
		student, err := w.mongoClient.GetStudentById(r.StudentId)
		if err != nil || student == nil || student.UserType != model.USER_TYPE_STUDENT {
			continue
		}
		grade, err := utils.ParseStudentId(student.Username)
		if err != nil {
			continue
		}
		tWork := teacherWorkloadMap[r.TeacherId]
		// 总人数/人次
		if _, ok := tWork.StudentIdMap[r.StudentId]; !ok {
			tWork.TotalNum++
		}
		tWork.TotalCount++
		tWork.StudentIdMap[r.StudentId]++
		// 本科生/研究生人数/人次
		if strings.HasSuffix(grade, "级") {
			if _, ok := tWork.UnderGraduateIdMap[r.StudentId]; !ok {
				tWork.NumUnderGraduates++
			}
			tWork.CountUnderGraduates++
			tWork.UnderGraduateIdMap[r.StudentId]++
		} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
			if _, ok := tWork.GraduateIdMap[r.StudentId]; !ok {
				tWork.NumGraduates++
			}
			tWork.CountGraduates++
			tWork.GraduateIdMap[r.StudentId]++
		}
		// 分类人次
		if r.TeacherFeedback.Category != "" {
			fc := r.TeacherFeedback.Category[0:1]
			if fcWork, ok := tWork.FirstClassWorkloads[fc]; ok {
				if strings.HasSuffix(grade, "级") {
					if _, ok := fcWork.UnderGraduateIdMap[r.StudentId]; !ok {
						fcWork.NumUnderGraduates++
					}
					fcWork.CountUnderGraduates++
					fcWork.UnderGraduateIdMap[r.StudentId]++
				} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
					if _, ok := fcWork.GraduateIdMap[r.StudentId]; !ok {
						fcWork.NumGraduates++
					}
					fcWork.CountGraduates++
					fcWork.GraduateIdMap[r.StudentId]++
				}
			}
		}
	}
	sort.Slice(teacherWorkloads, func(i, j int) bool {
		return teacherWorkloads[i].TeacherId < teacherWorkloads[j].TeacherId
	})
	// 写入文件
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error
	xlsx.SetDefaultFont(11, "宋体")
	file = xlsx.NewFile()
	sheet, err = file.AddSheet("咨询师工作量汇总")
	if err != nil {
		return re.NewRError("fail to create workload sheet", err)
	}
	// 第一表头
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("咨询师")
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("人数")
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("人次")
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("本科生")
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("研究生")
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	for _, fi := range model.FeedbackFirstCategoryKeys {
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(model.FeedbackFirstCategoryMap[fi])
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
	}
	// 合并第一表头
	cell = row.Cells[3]
	cell.Merge(1, 0)
	cell.SetStyle(textCenterStyle)
	cell = row.Cells[5]
	cell.Merge(1, 0)
	cell.SetStyle(textCenterStyle)
	for i := 0; i < len(model.FeedbackFirstCategoryKeys); i++ {
		cell = row.Cells[i*2+7]
		cell.Merge(1, 0)
		cell.SetStyle(textCenterStyle)
	}
	// 第二表头
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("人数")
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("人次")
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("人数")
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("人次")
	for i := 0; i < len(model.FeedbackFirstCategoryMap)-1; i++ {
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue("本科生人次")
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue("研究生人次")
	}
	// 合并第二表头
	row = sheet.Rows[0]
	cell = row.Cells[0]
	cell.Merge(0, 1)
	cell.SetStyle(textCenterStyle)
	cell = row.Cells[1]
	cell.Merge(0, 1)
	cell.SetStyle(textCenterStyle)
	cell = row.Cells[2]
	cell.Merge(0, 1)
	cell.SetStyle(textCenterStyle)
	// 工作量汇总
	for _, tWork := range teacherWorkloads {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(tWork.Fullname)
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(tWork.TotalNum)
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(tWork.TotalCount)
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(tWork.NumUnderGraduates)
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(tWork.CountUnderGraduates)
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(tWork.NumGraduates)
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(tWork.CountGraduates)
		for _, fi := range model.FeedbackFirstCategoryKeys {
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			if tWork.FirstClassWorkloads[fi].CountUnderGraduates > 0 {
				cell.SetValue(tWork.FirstClassWorkloads[fi].CountUnderGraduates)
			}
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			if tWork.FirstClassWorkloads[fi].CountGraduates > 0 {
				cell.SetValue(tWork.FirstClassWorkloads[fi].CountGraduates)
			}
		}
	}
	// 调整列宽
	sheet.SetColWidth(7, 7+2*len(model.FeedbackFirstCategoryKeys), 10)

	err = file.Save(path)
	if err != nil {
		return re.NewRError("fail to save file to path", err)
	}
	return nil
}
