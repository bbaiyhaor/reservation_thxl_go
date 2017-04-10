package buslogic

import (
	"fmt"
	"github.com/shudiwsh2009/reservation_thxl_go/model"
	re "github.com/shudiwsh2009/reservation_thxl_go/rerror"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"github.com/tealeg/xlsx"
	"sort"
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
type FirstClassReport struct {
	FirstClass          string
	SecondClassReports  []*SecondClassReport
	MaleIdMap           map[string]int
	NumMale             int // 合计人数（男）
	CountMale           int // 合计人次（男）
	FemaleIdMap         map[string]int
	NumFemale           int
	CountFemale         int // 合计（女）
	UnderGraduateIdMap  map[string]int
	NumUnderGraduates   int
	CountUnderGraduates int // 合计（本科生）
	GraduateIdMap       map[string]int
	NumGraduates        int
	CountGraduates      int // 合计（研究生）
}

type SecondClassReport struct {
	SecondClass         string
	Grades              map[string]int // 本科生、硕士生、博士生
	Instructor          int            // 辅导员
	Teacher             int            // 教师
	Family              int            // 家属
	Others              int            // 其他
	MaleIdMap           map[string]int
	NumMale             int // 合计人数（男）
	CountMale           int // 合计人次（男）
	FemaleIdMap         map[string]int
	NumFemale           int
	CountFemale         int // 合计（女）
	UnderGraduateIdMap  map[string]int
	NumUnderGraduates   int
	CountUnderGraduates int // 合计（本科生）
	GraduateIdMap       map[string]int
	NumGraduates        int
	CountGraduates      int // 合计（研究生）
	TotalIdMap          map[string]int
	NumTotal            int
	CountTotal          int     // 会谈总计
	Ratio               float64 // 比例（需转成百分比）
}

func (w *Workflow) ExportReportToFile(reservations []*model.Reservation, path string) error {
	// 建立存储结构
	// 咨询情况汇总
	categoryFCReport := make([]*FirstClassReport, 0)
	categoryFCReportMap := make(map[string]*FirstClassReport)
	categorySCReportMap := make(map[string]*SecondClassReport)
	for fi, firstCategory := range model.FeedbackFirstCategoryMap {
		if fi == "" {
			continue
		}
		fcReport := &FirstClassReport{
			FirstClass:         firstCategory,
			SecondClassReports: make([]*SecondClassReport, 0),
			MaleIdMap:          make(map[string]int),
			FemaleIdMap:        make(map[string]int),
			UnderGraduateIdMap: make(map[string]int),
			GraduateIdMap:      make(map[string]int),
		}
		for si, secondCategory := range model.FeedbackSecondCategoryMap[fi] {
			if si == "" {
				continue
			}
			scReport := &SecondClassReport{
				SecondClass:        secondCategory,
				Grades:             make(map[string]int),
				MaleIdMap:          make(map[string]int),
				FemaleIdMap:        make(map[string]int),
				UnderGraduateIdMap: make(map[string]int),
				GraduateIdMap:      make(map[string]int),
				TotalIdMap:         make(map[string]int),
			}
			fcReport.SecondClassReports = append(fcReport.SecondClassReports, scReport)
			categorySCReportMap[si] = scReport
		}
		sort.Slice(fcReport.SecondClassReports, func(i, j int) bool {
			return fcReport.SecondClassReports[i].SecondClass < fcReport.SecondClassReports[j].SecondClass
		})
		categoryFCReportMap[fi] = fcReport
		categoryFCReport = append(categoryFCReport, fcReport)
	}
	sort.Slice(categoryFCReport, func(i, j int) bool {
		return categoryFCReport[i].FirstClass < categoryFCReport[j].FirstClass
	})
	categoryTotalReport := &SecondClassReport{
		SecondClass: "总计",
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
		MaleIdMap:          make(map[string]int),
		FemaleIdMap:        make(map[string]int),
		UnderGraduateIdMap: make(map[string]int),
		GraduateIdMap:      make(map[string]int),
		TotalIdMap:         make(map[string]int),
	}
	// 重点情况汇总
	emphasisSCReportMap := make(map[string]*SecondClassReport)
	severityFCReport := &FirstClassReport{
		FirstClass:         "严重程度",
		SecondClassReports: make([]*SecondClassReport, 0),
		MaleIdMap:          make(map[string]int),
		FemaleIdMap:        make(map[string]int),
		UnderGraduateIdMap: make(map[string]int),
		GraduateIdMap:      make(map[string]int),
	}
	for _, sc := range model.FeedbackSeverity {
		scReport := &SecondClassReport{
			SecondClass:        sc,
			Grades:             make(map[string]int),
			MaleIdMap:          make(map[string]int),
			FemaleIdMap:        make(map[string]int),
			UnderGraduateIdMap: make(map[string]int),
			GraduateIdMap:      make(map[string]int),
			TotalIdMap:         make(map[string]int),
		}
		severityFCReport.SecondClassReports = append(severityFCReport.SecondClassReports, scReport)
		emphasisSCReportMap[sc] = scReport
	}
	medicalDiagnosisFCReport := &FirstClassReport{
		FirstClass:         "疑似或明确的医疗诊断",
		SecondClassReports: make([]*SecondClassReport, 0),
		MaleIdMap:          make(map[string]int),
		FemaleIdMap:        make(map[string]int),
		UnderGraduateIdMap: make(map[string]int),
		GraduateIdMap:      make(map[string]int),
	}
	for _, sc := range model.FeedbackMedicalDiagnosis {
		scReport := &SecondClassReport{
			SecondClass:        sc,
			Grades:             make(map[string]int),
			MaleIdMap:          make(map[string]int),
			FemaleIdMap:        make(map[string]int),
			UnderGraduateIdMap: make(map[string]int),
			GraduateIdMap:      make(map[string]int),
			TotalIdMap:         make(map[string]int),
		}
		medicalDiagnosisFCReport.SecondClassReports = append(medicalDiagnosisFCReport.SecondClassReports, scReport)
		emphasisSCReportMap[sc] = scReport
	}
	crisisFCReport := &FirstClassReport{
		FirstClass:         "危急情况",
		SecondClassReports: make([]*SecondClassReport, 0),
		MaleIdMap:          make(map[string]int),
		FemaleIdMap:        make(map[string]int),
		UnderGraduateIdMap: make(map[string]int),
		GraduateIdMap:      make(map[string]int),
	}
	for _, sc := range model.FeedbackCrisis {
		scReport := &SecondClassReport{
			SecondClass:        sc,
			Grades:             make(map[string]int),
			MaleIdMap:          make(map[string]int),
			FemaleIdMap:        make(map[string]int),
			UnderGraduateIdMap: make(map[string]int),
			GraduateIdMap:      make(map[string]int),
			TotalIdMap:         make(map[string]int),
		}
		crisisFCReport.SecondClassReports = append(crisisFCReport.SecondClassReports, scReport)
		emphasisSCReportMap[sc] = scReport
	}
	emphasisTotalReport := &SecondClassReport{
		SecondClass: "总计",
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
		MaleIdMap:          make(map[string]int),
		FemaleIdMap:        make(map[string]int),
		UnderGraduateIdMap: make(map[string]int),
		GraduateIdMap:      make(map[string]int),
		TotalIdMap:         make(map[string]int),
	}

	// 分析咨询数据
	for _, r := range reservations {
		if r.TeacherFeedback.IsEmpty() {
			continue
		}
		student, err := w.mongoClient.GetStudentById(r.StudentId)
		if err != nil || student.UserType != model.USER_TYPE_STUDENT {
			continue
		}
		studentId := student.Id.Hex()
		grade, err := utils.ParseStudentId(student.Username)
		if err != nil {
			continue
		}
		// 咨询情况汇总
		category := r.TeacherFeedback.Category
		// 来访情况，H分类中的来访人员特殊处理
		switch category {
		case "H1":
			categorySCReportMap[category].Instructor++
			categoryTotalReport.Instructor++
		case "H2", "H6":
			categorySCReportMap[category].Teacher++
			categoryTotalReport.Teacher++
		case "H3":
			categorySCReportMap[category].Family++
			categoryTotalReport.Family++
		case "H4", "H5":
			categorySCReportMap[category].Others++
			categoryTotalReport.Others++
		default:
			categorySCReportMap[category].Grades[grade]++
			categoryTotalReport.Grades[grade]++
		}
		// 性别统计
		if student.Gender == "男" {
			if _, ok := categorySCReportMap[category].MaleIdMap[studentId]; !ok {
				categorySCReportMap[category].NumMale++
			}
			categorySCReportMap[category].CountMale++
			categorySCReportMap[category].MaleIdMap[studentId]++

			if _, ok := categoryFCReportMap[category[0:1]].MaleIdMap[studentId]; !ok {
				categoryFCReportMap[category[0:1]].NumMale++
			}
			categoryFCReportMap[category[0:1]].CountMale++
			categoryFCReportMap[category[0:1]].MaleIdMap[studentId]++

			if _, ok := categoryTotalReport.MaleIdMap[studentId]; !ok {
				categoryTotalReport.NumMale++
			}
			categoryTotalReport.CountMale++
			categoryTotalReport.MaleIdMap[studentId]++
		} else if student.Gender == "女" {
			if _, ok := categorySCReportMap[category].FemaleIdMap[studentId]; !ok {
				categorySCReportMap[category].NumFemale++
			}
			categorySCReportMap[category].CountFemale++
			categorySCReportMap[category].FemaleIdMap[studentId]++

			if _, ok := categoryFCReportMap[category[0:1]].FemaleIdMap[studentId]; !ok {
				categoryFCReportMap[category[0:1]].NumFemale++
			}
			categoryFCReportMap[category[0:1]].CountFemale++
			categoryFCReportMap[category[0:1]].FemaleIdMap[studentId]++

			if _, ok := categoryTotalReport.FemaleIdMap[studentId]; !ok {
				categoryTotalReport.NumFemale++
			}
			categoryTotalReport.CountFemale++
			categoryTotalReport.FemaleIdMap[studentId]++
		}
		// 本科生/研究生统计
		if strings.HasSuffix(grade, "级") {
			if _, ok := categorySCReportMap[category].UnderGraduateIdMap[studentId]; !ok {
				categorySCReportMap[category].NumUnderGraduates++
			}
			categorySCReportMap[category].CountUnderGraduates++
			categorySCReportMap[category].UnderGraduateIdMap[studentId]++

			if _, ok := categoryFCReportMap[category[0:1]].UnderGraduateIdMap[studentId]; !ok {
				categoryFCReportMap[category[0:1]].NumUnderGraduates++
			}
			categoryFCReportMap[category[0:1]].CountUnderGraduates++
			categoryFCReportMap[category[0:1]].UnderGraduateIdMap[studentId]++

			if _, ok := categoryTotalReport.UnderGraduateIdMap[studentId]; !ok {
				categoryTotalReport.NumUnderGraduates++
			}
			categoryTotalReport.CountUnderGraduates++
			categoryTotalReport.UnderGraduateIdMap[studentId]++
		} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
			if _, ok := categorySCReportMap[category].GraduateIdMap[studentId]; !ok {
				categorySCReportMap[category].NumGraduates++
			}
			categorySCReportMap[category].CountGraduates++
			categorySCReportMap[category].GraduateIdMap[studentId]++

			if _, ok := categoryFCReportMap[category[0:1]].GraduateIdMap[studentId]; !ok {
				categoryFCReportMap[category[0:1]].NumGraduates++
			}
			categoryFCReportMap[category[0:1]].CountGraduates++
			categoryFCReportMap[category[0:1]].GraduateIdMap[studentId]++

			if _, ok := categoryTotalReport.GraduateIdMap[studentId]; !ok {
				categoryTotalReport.NumGraduates++
			}
			categoryTotalReport.CountGraduates++
			categoryTotalReport.GraduateIdMap[studentId]++
		}
		// 总计
		if _, ok := categorySCReportMap[category].TotalIdMap[studentId]; !ok {
			categorySCReportMap[category].NumTotal++
		}
		categorySCReportMap[category].CountTotal++
		categorySCReportMap[category].TotalIdMap[studentId]++

		if _, ok := categoryTotalReport.TotalIdMap[studentId]; !ok {
			categoryTotalReport.NumTotal++
		}
		categoryTotalReport.CountTotal++
		categoryTotalReport.TotalIdMap[studentId]++

		// 重点情况汇总
		severity := r.TeacherFeedback.Severity
		medicalDiagnosis := r.TeacherFeedback.MedicalDiagnosis
		crisis := r.TeacherFeedback.Crisis
		for i, s := range severity {
			if s == 1 {
				emphasisSCReportMap[model.FeedbackSeverity[i]].Grades[grade]++
				emphasisTotalReport.Grades[grade]++
				if student.Gender == "男" {
					if _, ok := emphasisSCReportMap[model.FeedbackSeverity[i]].MaleIdMap[studentId]; !ok {
						emphasisSCReportMap[model.FeedbackSeverity[i]].NumMale++
					}
					emphasisSCReportMap[model.FeedbackSeverity[i]].CountMale++
					emphasisSCReportMap[model.FeedbackSeverity[i]].MaleIdMap[studentId]++

					if _, ok := emphasisTotalReport.MaleIdMap[studentId]; !ok {
						emphasisTotalReport.NumMale++
					}
					emphasisTotalReport.CountMale++
					emphasisTotalReport.MaleIdMap[studentId]++

					if _, ok := severityFCReport.MaleIdMap[studentId]; !ok {
						severityFCReport.NumMale++
					}
					severityFCReport.CountMale++
					severityFCReport.MaleIdMap[studentId]++
				} else if student.Gender == "女" {
					if _, ok := emphasisSCReportMap[model.FeedbackSeverity[i]].FemaleIdMap[studentId]; !ok {
						emphasisSCReportMap[model.FeedbackSeverity[i]].NumFemale++
					}
					emphasisSCReportMap[model.FeedbackSeverity[i]].CountFemale++
					emphasisSCReportMap[model.FeedbackSeverity[i]].FemaleIdMap[studentId]++

					if _, ok := emphasisTotalReport.FemaleIdMap[studentId]; !ok {
						emphasisTotalReport.NumFemale++
					}
					emphasisTotalReport.CountFemale++
					emphasisTotalReport.FemaleIdMap[studentId]++

					if _, ok := severityFCReport.FemaleIdMap[studentId]; !ok {
						severityFCReport.NumFemale++
					}
					severityFCReport.CountFemale++
					severityFCReport.FemaleIdMap[studentId]++
				}
				if strings.HasSuffix(grade, "级") {
					if _, ok := emphasisSCReportMap[model.FeedbackSeverity[i]].UnderGraduateIdMap[studentId]; !ok {
						emphasisSCReportMap[model.FeedbackSeverity[i]].NumUnderGraduates++
					}
					emphasisSCReportMap[model.FeedbackSeverity[i]].CountUnderGraduates++
					emphasisSCReportMap[model.FeedbackSeverity[i]].UnderGraduateIdMap[studentId]++

					if _, ok := emphasisTotalReport.UnderGraduateIdMap[studentId]; !ok {
						emphasisTotalReport.NumUnderGraduates++
					}
					emphasisTotalReport.CountUnderGraduates++
					emphasisTotalReport.UnderGraduateIdMap[studentId]++

					if _, ok := severityFCReport.UnderGraduateIdMap[studentId]; !ok {
						severityFCReport.NumUnderGraduates++
					}
					severityFCReport.CountUnderGraduates++
					severityFCReport.UnderGraduateIdMap[studentId]++
				} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
					if _, ok := emphasisSCReportMap[model.FeedbackSeverity[i]].GraduateIdMap[studentId]; !ok {
						emphasisSCReportMap[model.FeedbackSeverity[i]].NumGraduates++
					}
					emphasisSCReportMap[model.FeedbackSeverity[i]].CountGraduates++
					emphasisSCReportMap[model.FeedbackSeverity[i]].GraduateIdMap[studentId]++

					if _, ok := emphasisTotalReport.GraduateIdMap[studentId]; !ok {
						emphasisTotalReport.NumGraduates++
					}
					emphasisTotalReport.CountGraduates++
					emphasisTotalReport.GraduateIdMap[studentId]++

					if _, ok := severityFCReport.GraduateIdMap[studentId]; !ok {
						severityFCReport.NumGraduates++
					}
					severityFCReport.CountGraduates++
					severityFCReport.GraduateIdMap[studentId]++
				}
				if _, ok := emphasisSCReportMap[model.FeedbackSeverity[i]].TotalIdMap[studentId]; !ok {
					emphasisSCReportMap[model.FeedbackSeverity[i]].NumTotal++
				}
				emphasisSCReportMap[model.FeedbackSeverity[i]].CountTotal++
				emphasisSCReportMap[model.FeedbackSeverity[i]].TotalIdMap[studentId]++

				if _, ok := emphasisTotalReport.TotalIdMap[studentId]; !ok {
					emphasisTotalReport.NumTotal++
				}
				emphasisTotalReport.CountTotal++
				emphasisTotalReport.TotalIdMap[studentId]++
			}
		}
		for i, m := range medicalDiagnosis {
			if m == 1 {
				emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].Grades[grade]++
				emphasisTotalReport.Grades[grade]++
				if student.Gender == "男" {
					if _, ok := emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].MaleIdMap[studentId]; !ok {
						emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].NumMale++
					}
					emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].CountMale++
					emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].MaleIdMap[studentId]++

					if _, ok := emphasisTotalReport.MaleIdMap[studentId]; !ok {
						emphasisTotalReport.NumMale++
					}
					emphasisTotalReport.CountMale++
					emphasisTotalReport.MaleIdMap[studentId]++

					if _, ok := medicalDiagnosisFCReport.MaleIdMap[studentId]; !ok {
						medicalDiagnosisFCReport.NumMale++
					}
					medicalDiagnosisFCReport.CountMale++
					medicalDiagnosisFCReport.MaleIdMap[studentId]++
				} else if student.Gender == "女" {
					if _, ok := emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].FemaleIdMap[studentId]; !ok {
						emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].NumFemale++
					}
					emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].CountFemale++
					emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].FemaleIdMap[studentId]++

					if _, ok := emphasisTotalReport.FemaleIdMap[studentId]; !ok {
						emphasisTotalReport.NumFemale++
					}
					emphasisTotalReport.CountFemale++
					emphasisTotalReport.FemaleIdMap[studentId]++

					if _, ok := medicalDiagnosisFCReport.FemaleIdMap[studentId]; !ok {
						medicalDiagnosisFCReport.NumFemale++
					}
					medicalDiagnosisFCReport.CountFemale++
					medicalDiagnosisFCReport.FemaleIdMap[studentId]++
				}
				if strings.HasSuffix(grade, "级") {
					if _, ok := emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].UnderGraduateIdMap[studentId]; !ok {
						emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].NumUnderGraduates++
					}
					emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].CountUnderGraduates++
					emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].UnderGraduateIdMap[studentId]++

					if _, ok := emphasisTotalReport.UnderGraduateIdMap[studentId]; !ok {
						emphasisTotalReport.NumUnderGraduates++
					}
					emphasisTotalReport.CountUnderGraduates++
					emphasisTotalReport.UnderGraduateIdMap[studentId]++

					if _, ok := medicalDiagnosisFCReport.UnderGraduateIdMap[studentId]; !ok {
						medicalDiagnosisFCReport.NumUnderGraduates++
					}
					medicalDiagnosisFCReport.CountUnderGraduates++
					medicalDiagnosisFCReport.UnderGraduateIdMap[studentId]++
				} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
					if _, ok := emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].GraduateIdMap[studentId]; !ok {
						emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].NumGraduates++
					}
					emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].CountGraduates++
					emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].GraduateIdMap[studentId]++

					if _, ok := emphasisTotalReport.GraduateIdMap[studentId]; !ok {
						emphasisTotalReport.NumGraduates++
					}
					emphasisTotalReport.CountGraduates++
					emphasisTotalReport.GraduateIdMap[studentId]++

					if _, ok := medicalDiagnosisFCReport.GraduateIdMap[studentId]; !ok {
						medicalDiagnosisFCReport.NumGraduates++
					}
					medicalDiagnosisFCReport.CountGraduates++
					medicalDiagnosisFCReport.GraduateIdMap[studentId]++
				}
				if _, ok := emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].TotalIdMap[studentId]; !ok {
					emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].NumTotal++
				}
				emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].CountTotal++
				emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].TotalIdMap[studentId]++

				if _, ok := emphasisTotalReport.TotalIdMap[studentId]; !ok {
					emphasisTotalReport.NumTotal++
				}
				emphasisTotalReport.CountTotal++
				emphasisTotalReport.TotalIdMap[studentId]++
			}
		}
		for i, c := range crisis {
			if c == 1 {
				emphasisSCReportMap[model.FeedbackCrisis[i]].Grades[grade]++
				emphasisTotalReport.Grades[grade]++
				if student.Gender == "男" {
					if _, ok := emphasisSCReportMap[model.FeedbackCrisis[i]].MaleIdMap[studentId]; !ok {
						emphasisSCReportMap[model.FeedbackCrisis[i]].NumMale++
					}
					emphasisSCReportMap[model.FeedbackCrisis[i]].CountMale++
					emphasisSCReportMap[model.FeedbackCrisis[i]].MaleIdMap[studentId]++

					if _, ok := emphasisTotalReport.MaleIdMap[studentId]; !ok {
						emphasisTotalReport.NumMale++
					}
					emphasisTotalReport.CountMale++
					emphasisTotalReport.MaleIdMap[studentId]++

					if _, ok := crisisFCReport.MaleIdMap[studentId]; !ok {
						crisisFCReport.NumMale++
					}
					crisisFCReport.CountMale++
					crisisFCReport.MaleIdMap[studentId]++
				} else if student.Gender == "女" {
					if _, ok := emphasisSCReportMap[model.FeedbackCrisis[i]].FemaleIdMap[studentId]; !ok {
						emphasisSCReportMap[model.FeedbackCrisis[i]].NumFemale++
					}
					emphasisSCReportMap[model.FeedbackCrisis[i]].CountFemale++
					emphasisSCReportMap[model.FeedbackCrisis[i]].FemaleIdMap[studentId]++

					if _, ok := emphasisTotalReport.FemaleIdMap[studentId]; !ok {
						emphasisTotalReport.NumFemale++
					}
					emphasisTotalReport.CountFemale++
					emphasisTotalReport.FemaleIdMap[studentId]++

					if _, ok := crisisFCReport.FemaleIdMap[studentId]; !ok {
						crisisFCReport.NumFemale++
					}
					crisisFCReport.CountFemale++
					crisisFCReport.FemaleIdMap[studentId]++
				}
				if strings.HasSuffix(grade, "级") {
					if _, ok := emphasisSCReportMap[model.FeedbackCrisis[i]].UnderGraduateIdMap[studentId]; !ok {
						emphasisSCReportMap[model.FeedbackCrisis[i]].NumUnderGraduates++
					}
					emphasisSCReportMap[model.FeedbackCrisis[i]].CountUnderGraduates++
					emphasisSCReportMap[model.FeedbackCrisis[i]].UnderGraduateIdMap[studentId]++

					if _, ok := emphasisTotalReport.UnderGraduateIdMap[studentId]; !ok {
						emphasisTotalReport.NumUnderGraduates++
					}
					emphasisTotalReport.CountUnderGraduates++
					emphasisTotalReport.UnderGraduateIdMap[studentId]++

					if _, ok := crisisFCReport.UnderGraduateIdMap[studentId]; !ok {
						crisisFCReport.NumUnderGraduates++
					}
					crisisFCReport.CountUnderGraduates++
					crisisFCReport.UnderGraduateIdMap[studentId]++
				} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
					if _, ok := emphasisSCReportMap[model.FeedbackCrisis[i]].GraduateIdMap[studentId]; !ok {
						emphasisSCReportMap[model.FeedbackCrisis[i]].NumGraduates++
					}
					emphasisSCReportMap[model.FeedbackCrisis[i]].CountGraduates++
					emphasisSCReportMap[model.FeedbackCrisis[i]].GraduateIdMap[studentId]++

					if _, ok := emphasisTotalReport.GraduateIdMap[studentId]; !ok {
						emphasisTotalReport.NumGraduates++
					}
					emphasisTotalReport.CountGraduates++
					emphasisTotalReport.GraduateIdMap[studentId]++

					if _, ok := crisisFCReport.GraduateIdMap[studentId]; !ok {
						crisisFCReport.NumGraduates++
					}
					crisisFCReport.CountGraduates++
					crisisFCReport.GraduateIdMap[studentId]++
				}
				if _, ok := emphasisSCReportMap[model.FeedbackCrisis[i]].TotalIdMap[studentId]; !ok {
					emphasisSCReportMap[model.FeedbackCrisis[i]].NumTotal++
				}
				emphasisSCReportMap[model.FeedbackCrisis[i]].CountTotal++
				emphasisSCReportMap[model.FeedbackCrisis[i]].TotalIdMap[studentId]++

				if _, ok := emphasisTotalReport.TotalIdMap[studentId]; !ok {
					emphasisTotalReport.NumTotal++
				}
				emphasisTotalReport.CountTotal++
				emphasisTotalReport.TotalIdMap[studentId]++
			}
		}
	}
	for _, scReport := range categorySCReportMap {
		scReport.Ratio = float64(scReport.CountTotal) / float64(categoryTotalReport.CountTotal)
	}
	for _, scReport := range emphasisSCReportMap {
		scReport.Ratio = float64(scReport.CountTotal) / float64(emphasisTotalReport.CountTotal)
	}
	// 分析年级数据
	allGrades := make([]string, 0)
	for g := range categoryTotalReport.Grades {
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
	cell.SetValue("总计")
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
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
	cell.SetStyle(bgOrangeStyle)
	cell.SetValue("会谈总计")
	cell = row.AddCell()
	cell.SetStyle(bgOrangeStyle)
	cell.SetValue("百分比")
	// 咨询情况数据
	firstRowOfFcReport := 2 // 标记当前fcReport的行号，以便最后合并单元格
	for _, fcReport := range categoryFCReport {
		for _, scReport := range fcReport.SecondClassReports {
			row = sheet.AddRow()
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			cell.SetValue(strings.Split(fcReport.FirstClass, " ")[1])
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			cell.SetValue(scReport.SecondClass)
			for _, g := range allGrades {
				cell = row.AddCell()
				cell.SetStyle(borderStyle)
				if scReport.Grades[g] > 0 {
					cell.SetValue(scReport.Grades[g])
				}
			}
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			if scReport.Instructor > 0 {
				cell.SetValue(scReport.Instructor)
			}
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			if scReport.Teacher > 0 {
				cell.SetValue(scReport.Teacher)
			}
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			if scReport.Family > 0 {
				cell.SetValue(scReport.Family)
			}
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			if scReport.Others > 0 {
				cell.SetValue(scReport.Others)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGrayStyle)
			if scReport.CountMale > 0 {
				cell.SetValue(scReport.CountMale)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGrayStyle)
			if scReport.CountFemale > 0 {
				cell.SetValue(scReport.CountFemale)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGrayStyle)
			if scReport.CountUnderGraduates > 0 {
				cell.SetValue(scReport.CountUnderGraduates)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGrayStyle)
			if scReport.CountGraduates > 0 {
				cell.SetValue(scReport.CountGraduates)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			if fcReport.CountMale > 0 {
				cell.SetValue(fcReport.CountMale)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			if fcReport.CountFemale > 0 {
				cell.SetValue(fcReport.CountFemale)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			if fcReport.CountUnderGraduates > 0 {
				cell.SetValue(fcReport.CountUnderGraduates)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			if fcReport.CountGraduates > 0 {
				cell.SetValue(fcReport.CountGraduates)
			}
			cell = row.AddCell()
			cell.SetStyle(bgOrangeStyle)
			cell.SetValue(scReport.CountTotal)
			cell = row.AddCell()
			cell.SetStyle(bgOrangeStyle)
			cell.SetValue(fmt.Sprintf("%.2f%%", scReport.Ratio*100))
		}
		row = sheet.Rows[firstRowOfFcReport]
		cell = row.Cells[0]
		cell.Merge(0, len(fcReport.SecondClassReports)-1)
		cell.SetStyle(textCenterStyle)
		cell = row.Cells[len(allGrades)+10]
		cell.Merge(0, len(fcReport.SecondClassReports)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+11]
		cell.Merge(0, len(fcReport.SecondClassReports)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+12]
		cell.Merge(0, len(fcReport.SecondClassReports)-1)
		cell.SetStyle(textRightGreenStyle)
		cell = row.Cells[len(allGrades)+13]
		cell.Merge(0, len(fcReport.SecondClassReports)-1)
		cell.SetStyle(textRightGreenStyle)
		firstRowOfFcReport += len(fcReport.SecondClassReports)
	}
	// 总计
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.SecondClass)
	for _, g := range allGrades {
		cell = row.AddCell()
		cell.SetStyle(bgYellowStyle)
		cell.SetValue(categoryTotalReport.Grades[g])
	}
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.Instructor)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.Teacher)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.Family)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.Others)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.CountMale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.CountFemale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.CountUnderGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.CountGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.CountMale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.CountFemale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.CountUnderGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.CountGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgRedStyle)
	cell.SetValue(categoryTotalReport.CountTotal)
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
		cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.Grades[g])/float64(categoryTotalReport.CountTotal)*100))
	}
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.Instructor)/float64(categoryTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.Teacher)/float64(categoryTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.Family)/float64(categoryTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.Others)/float64(categoryTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.CountMale)/float64(categoryTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.CountFemale)/float64(categoryTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.CountUnderGraduates)/float64(categoryTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.CountGraduates)/float64(categoryTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.CountMale)/float64(categoryTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.CountFemale)/float64(categoryTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.CountUnderGraduates)/float64(categoryTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.CountGraduates)/float64(categoryTotalReport.CountTotal)*100))
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
	firstRowOfFcReport = 2
	// 严重程度
	for _, scReport := range severityFCReport.SecondClassReports {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(severityFCReport.FirstClass)
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(scReport.SecondClass)
		for _, g := range allGrades {
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			if scReport.Grades[g] > 0 {
				cell.SetValue(scReport.Grades[g])
			}
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.CountMale > 0 {
			cell.SetValue(scReport.CountMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.CountFemale > 0 {
			cell.SetValue(scReport.CountFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.CountUnderGraduates > 0 {
			cell.SetValue(scReport.CountUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.CountGraduates > 0 {
			cell.SetValue(scReport.CountGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if severityFCReport.CountMale > 0 {
			cell.SetValue(severityFCReport.CountMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if severityFCReport.CountFemale > 0 {
			cell.SetValue(severityFCReport.CountFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if severityFCReport.CountUnderGraduates > 0 {
			cell.SetValue(severityFCReport.CountUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if severityFCReport.CountGraduates > 0 {
			cell.SetValue(severityFCReport.CountGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgOrangeStyle)
		cell.SetValue(scReport.CountTotal)
		cell = row.AddCell()
		cell.SetStyle(bgOrangeStyle)
		cell.SetValue(fmt.Sprintf("%.2f%%", scReport.Ratio*100))
	}
	row = sheet.Rows[firstRowOfFcReport]
	cell = row.Cells[0]
	cell.Merge(0, len(severityFCReport.SecondClassReports)-1)
	cell.SetStyle(textCenterStyle)
	cell = row.Cells[len(allGrades)+6]
	cell.Merge(0, len(severityFCReport.SecondClassReports)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+7]
	cell.Merge(0, len(severityFCReport.SecondClassReports)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+8]
	cell.Merge(0, len(severityFCReport.SecondClassReports)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+9]
	cell.Merge(0, len(severityFCReport.SecondClassReports)-1)
	cell.SetStyle(textRightGreenStyle)
	firstRowOfFcReport += len(severityFCReport.SecondClassReports)
	// 疑似或明确的医疗诊断
	for _, scReport := range medicalDiagnosisFCReport.SecondClassReports {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(medicalDiagnosisFCReport.FirstClass)
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(scReport.SecondClass)
		for _, g := range allGrades {
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			if scReport.Grades[g] > 0 {
				cell.SetValue(scReport.Grades[g])
			}
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.CountMale > 0 {
			cell.SetValue(scReport.CountMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.CountFemale > 0 {
			cell.SetValue(scReport.CountFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.CountUnderGraduates > 0 {
			cell.SetValue(scReport.CountUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.CountGraduates > 0 {
			cell.SetValue(scReport.CountGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if medicalDiagnosisFCReport.CountMale > 0 {
			cell.SetValue(medicalDiagnosisFCReport.CountMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if medicalDiagnosisFCReport.CountFemale > 0 {
			cell.SetValue(medicalDiagnosisFCReport.CountFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if medicalDiagnosisFCReport.CountUnderGraduates > 0 {
			cell.SetValue(medicalDiagnosisFCReport.CountUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if medicalDiagnosisFCReport.CountGraduates > 0 {
			cell.SetValue(medicalDiagnosisFCReport.CountGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgOrangeStyle)
		cell.SetValue(scReport.CountTotal)
		cell = row.AddCell()
		cell.SetStyle(bgOrangeStyle)
		cell.SetValue(fmt.Sprintf("%.2f%%", scReport.Ratio*100))
	}
	row = sheet.Rows[firstRowOfFcReport]
	cell = row.Cells[0]
	cell.Merge(0, len(medicalDiagnosisFCReport.SecondClassReports)-1)
	cell.SetStyle(textCenterStyle)
	cell = row.Cells[len(allGrades)+6]
	cell.Merge(0, len(medicalDiagnosisFCReport.SecondClassReports)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+7]
	cell.Merge(0, len(medicalDiagnosisFCReport.SecondClassReports)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+8]
	cell.Merge(0, len(medicalDiagnosisFCReport.SecondClassReports)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+9]
	cell.Merge(0, len(medicalDiagnosisFCReport.SecondClassReports)-1)
	cell.SetStyle(textRightGreenStyle)
	firstRowOfFcReport += len(medicalDiagnosisFCReport.SecondClassReports)
	// 危急情况
	for _, scReport := range crisisFCReport.SecondClassReports {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(crisisFCReport.FirstClass)
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(scReport.SecondClass)
		for _, g := range allGrades {
			cell = row.AddCell()
			cell.SetStyle(borderStyle)
			if scReport.Grades[g] > 0 {
				cell.SetValue(scReport.Grades[g])
			}
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.CountMale > 0 {
			cell.SetValue(scReport.CountMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.CountFemale > 0 {
			cell.SetValue(scReport.CountFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.CountUnderGraduates > 0 {
			cell.SetValue(scReport.CountUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.CountGraduates > 0 {
			cell.SetValue(scReport.CountGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if crisisFCReport.CountMale > 0 {
			cell.SetValue(crisisFCReport.CountMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if crisisFCReport.CountFemale > 0 {
			cell.SetValue(crisisFCReport.CountFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if crisisFCReport.CountUnderGraduates > 0 {
			cell.SetValue(crisisFCReport.CountUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if crisisFCReport.CountGraduates > 0 {
			cell.SetValue(crisisFCReport.CountGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgOrangeStyle)
		cell.SetValue(scReport.CountTotal)
		cell = row.AddCell()
		cell.SetStyle(bgOrangeStyle)
		cell.SetValue(fmt.Sprintf("%.2f%%", scReport.Ratio*100))
	}
	row = sheet.Rows[firstRowOfFcReport]
	cell = row.Cells[0]
	cell.Merge(0, len(crisisFCReport.SecondClassReports)-1)
	cell.SetStyle(textCenterStyle)
	cell = row.Cells[len(allGrades)+6]
	cell.Merge(0, len(crisisFCReport.SecondClassReports)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+7]
	cell.Merge(0, len(crisisFCReport.SecondClassReports)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+8]
	cell.Merge(0, len(crisisFCReport.SecondClassReports)-1)
	cell.SetStyle(textRightGreenStyle)
	cell = row.Cells[len(allGrades)+9]
	cell.Merge(0, len(crisisFCReport.SecondClassReports)-1)
	cell.SetStyle(textRightGreenStyle)
	firstRowOfFcReport += len(crisisFCReport.SecondClassReports)
	// 总计
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalReport.SecondClass)
	for _, g := range allGrades {
		cell = row.AddCell()
		cell.SetStyle(bgYellowStyle)
		cell.SetValue(emphasisTotalReport.Grades[g])
	}
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalReport.CountMale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalReport.CountFemale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalReport.CountUnderGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalReport.CountGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalReport.CountMale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalReport.CountFemale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalReport.CountUnderGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalReport.CountGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgRedStyle)
	cell.SetValue(emphasisTotalReport.CountTotal)
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
		cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.Grades[g])/float64(emphasisTotalReport.CountTotal)*100))
	}
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.CountMale)/float64(emphasisTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.CountFemale)/float64(emphasisTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.CountUnderGraduates)/float64(emphasisTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.CountGraduates)/float64(emphasisTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.CountMale)/float64(emphasisTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.CountFemale)/float64(emphasisTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.CountUnderGraduates)/float64(emphasisTotalReport.CountTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.CountGraduates)/float64(emphasisTotalReport.CountTotal)*100))
	// 调整列宽
	sheet.SetColWidth(0, 0, 15.5)
	sheet.SetColWidth(1, 1, 24)
	sheet.SetColWidth(len(allGrades)+4, len(allGrades)+5, 11)
	sheet.SetColWidth(len(allGrades)+8, len(allGrades)+9, 11)

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
			if err != nil {
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
		if err != nil {
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
