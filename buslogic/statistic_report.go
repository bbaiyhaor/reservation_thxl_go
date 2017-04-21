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
	FirstClass         string
	SecondClassReports []*SecondClassReport
	NumMale            int // 合计（男）
	NumFemale          int // 合计（女）
	NumUnderGraduates  int // 合计（本科生）
	NumGraduates       int // 合计（研究生）
}

type SecondClassReport struct {
	SecondClass       string
	Grades            map[string]int // 本科生、硕士生、博士生
	Instructor        int            // 辅导员
	Teacher           int            // 教师
	Family            int            // 家属
	Others            int            // 其他
	NumMale           int            // 合计（男）
	NumFemale         int            // 合计（女）
	NumUnderGraduates int            // 合计（本科生）
	NumGraduates      int            // 合计（研究生）
	NumTotal          int            // 会谈总计
	Ratio             float64        // 比例（需转成百分比）
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
		}
		for si, secondCategory := range model.FeedbackSecondCategoryMap[fi] {
			if si == "" {
				continue
			}
			scReport := &SecondClassReport{
				SecondClass: secondCategory,
				Grades:      make(map[string]int),
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
	}
	// 重点情况汇总
	emphasisSCReportMap := make(map[string]*SecondClassReport)
	severityFCReport := &FirstClassReport{
		FirstClass:         "严重程度",
		SecondClassReports: make([]*SecondClassReport, 0),
	}
	for _, sc := range model.FeedbackSeverity {
		scReport := &SecondClassReport{
			SecondClass: sc,
			Grades:      make(map[string]int),
		}
		severityFCReport.SecondClassReports = append(severityFCReport.SecondClassReports, scReport)
		emphasisSCReportMap[sc] = scReport
	}
	medicalDiagnosisFCReport := &FirstClassReport{
		FirstClass:         "疑似或明确的医疗诊断",
		SecondClassReports: make([]*SecondClassReport, 0),
	}
	for _, sc := range model.FeedbackMedicalDiagnosis {
		scReport := &SecondClassReport{
			SecondClass: sc,
			Grades:      make(map[string]int),
		}
		medicalDiagnosisFCReport.SecondClassReports = append(medicalDiagnosisFCReport.SecondClassReports, scReport)
		emphasisSCReportMap[sc] = scReport
	}
	crisisFCReport := &FirstClassReport{
		FirstClass:         "危急情况",
		SecondClassReports: make([]*SecondClassReport, 0),
	}
	for _, sc := range model.FeedbackCrisis {
		scReport := &SecondClassReport{
			SecondClass: sc,
			Grades:      make(map[string]int),
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
	}

	// 分析咨询数据
	for _, r := range reservations {
		if r.TeacherFeedback.IsEmpty() {
			continue
		}
		student, err := w.mongoClient.GetStudentById(r.StudentId)
		if err != nil || student == nil || student.UserType != model.USER_TYPE_STUDENT {
			continue
		}
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
			categorySCReportMap[category].NumMale++
			categoryFCReportMap[category[0:1]].NumMale++
			categoryTotalReport.NumMale++
		} else if student.Gender == "女" {
			categorySCReportMap[category].NumFemale++
			categoryFCReportMap[category[0:1]].NumFemale++
			categoryTotalReport.NumFemale++
		}
		// 本科生/研究生统计
		if strings.HasSuffix(grade, "级") {
			categorySCReportMap[category].NumUnderGraduates++
			categoryFCReportMap[category[0:1]].NumUnderGraduates++
			categoryTotalReport.NumUnderGraduates++
		} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
			categorySCReportMap[category].NumGraduates++
			categoryFCReportMap[category[0:1]].NumGraduates++
			categoryTotalReport.NumGraduates++
		}
		// 总计
		categorySCReportMap[category].NumTotal++
		categoryTotalReport.NumTotal++

		// 重点情况汇总
		severity := r.TeacherFeedback.Severity
		medicalDiagnosis := r.TeacherFeedback.MedicalDiagnosis
		crisis := r.TeacherFeedback.Crisis
		for i, s := range severity {
			if s == 1 {
				emphasisSCReportMap[model.FeedbackSeverity[i]].Grades[grade]++
				emphasisTotalReport.Grades[grade]++
				if student.Gender == "男" {
					emphasisSCReportMap[model.FeedbackSeverity[i]].NumMale++
					emphasisTotalReport.NumMale++
					severityFCReport.NumMale++
				} else if student.Gender == "女" {
					emphasisSCReportMap[model.FeedbackSeverity[i]].NumFemale++
					emphasisTotalReport.NumFemale++
					severityFCReport.NumFemale++
				}
				if strings.HasSuffix(grade, "级") {
					emphasisSCReportMap[model.FeedbackSeverity[i]].NumUnderGraduates++
					emphasisTotalReport.NumUnderGraduates++
					severityFCReport.NumUnderGraduates++
				} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
					emphasisSCReportMap[model.FeedbackSeverity[i]].NumGraduates++
					emphasisTotalReport.NumGraduates++
					severityFCReport.NumGraduates++
				}
				emphasisSCReportMap[model.FeedbackSeverity[i]].NumTotal++
				emphasisTotalReport.NumTotal++
			}
		}
		for i, m := range medicalDiagnosis {
			if m == 1 {
				emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].Grades[grade]++
				emphasisTotalReport.Grades[grade]++
				if student.Gender == "男" {
					emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].NumMale++
					emphasisTotalReport.NumMale++
					medicalDiagnosisFCReport.NumMale++
				} else if student.Gender == "女" {
					emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].NumFemale++
					emphasisTotalReport.NumFemale++
					medicalDiagnosisFCReport.NumFemale++
				}
				if strings.HasSuffix(grade, "级") {
					emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].NumUnderGraduates++
					emphasisTotalReport.NumUnderGraduates++
					medicalDiagnosisFCReport.NumUnderGraduates++
				} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
					emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].NumGraduates++
					emphasisTotalReport.NumGraduates++
					medicalDiagnosisFCReport.NumGraduates++
				}
				emphasisSCReportMap[model.FeedbackMedicalDiagnosis[i]].NumTotal++
				emphasisTotalReport.NumTotal++
			}
		}
		for i, c := range crisis {
			if c == 1 {
				emphasisSCReportMap[model.FeedbackCrisis[i]].Grades[grade]++
				emphasisTotalReport.Grades[grade]++
				if student.Gender == "男" {
					emphasisSCReportMap[model.FeedbackCrisis[i]].NumMale++
					emphasisTotalReport.NumMale++
					crisisFCReport.NumMale++
				} else if student.Gender == "女" {
					emphasisSCReportMap[model.FeedbackCrisis[i]].NumFemale++
					emphasisTotalReport.NumFemale++
					crisisFCReport.NumFemale++
				}
				if strings.HasSuffix(grade, "级") {
					emphasisSCReportMap[model.FeedbackCrisis[i]].NumUnderGraduates++
					emphasisTotalReport.NumUnderGraduates++
					crisisFCReport.NumUnderGraduates++
				} else if strings.HasSuffix(grade, "硕") || strings.HasSuffix(grade, "博") {
					emphasisSCReportMap[model.FeedbackCrisis[i]].NumGraduates++
					emphasisTotalReport.NumGraduates++
					crisisFCReport.NumGraduates++
				}
				emphasisSCReportMap[model.FeedbackCrisis[i]].NumTotal++
				emphasisTotalReport.NumTotal++
			}
		}
	}
	for _, scReport := range categorySCReportMap {
		scReport.Ratio = float64(scReport.NumTotal) / float64(categoryTotalReport.NumTotal)
	}
	for _, scReport := range emphasisSCReportMap {
		scReport.Ratio = float64(scReport.NumTotal) / float64(emphasisTotalReport.NumTotal)
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
			if scReport.NumMale > 0 {
				cell.SetValue(scReport.NumMale)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGrayStyle)
			if scReport.NumFemale > 0 {
				cell.SetValue(scReport.NumFemale)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGrayStyle)
			if scReport.NumUnderGraduates > 0 {
				cell.SetValue(scReport.NumUnderGraduates)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGrayStyle)
			if scReport.NumGraduates > 0 {
				cell.SetValue(scReport.NumGraduates)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			if fcReport.NumMale > 0 {
				cell.SetValue(fcReport.NumMale)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			if fcReport.NumFemale > 0 {
				cell.SetValue(fcReport.NumFemale)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			if fcReport.NumUnderGraduates > 0 {
				cell.SetValue(fcReport.NumUnderGraduates)
			}
			cell = row.AddCell()
			cell.SetStyle(bgGreenStyle)
			if fcReport.NumGraduates > 0 {
				cell.SetValue(fcReport.NumGraduates)
			}
			cell = row.AddCell()
			cell.SetStyle(bgOrangeStyle)
			cell.SetValue(scReport.NumTotal)
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
	cell.SetValue(categoryTotalReport.NumMale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.NumFemale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.NumUnderGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.NumGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.NumMale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.NumFemale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.NumUnderGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(categoryTotalReport.NumGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgRedStyle)
	cell.SetValue(categoryTotalReport.NumTotal)
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
		cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.Grades[g])/float64(categoryTotalReport.NumTotal)*100))
	}
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.Instructor)/float64(categoryTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.Teacher)/float64(categoryTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.Family)/float64(categoryTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.Others)/float64(categoryTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.NumMale)/float64(categoryTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.NumFemale)/float64(categoryTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.NumUnderGraduates)/float64(categoryTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.NumGraduates)/float64(categoryTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.NumMale)/float64(categoryTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.NumFemale)/float64(categoryTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.NumUnderGraduates)/float64(categoryTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(categoryTotalReport.NumGraduates)/float64(categoryTotalReport.NumTotal)*100))
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
		if scReport.NumMale > 0 {
			cell.SetValue(scReport.NumMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.NumFemale > 0 {
			cell.SetValue(scReport.NumFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.NumUnderGraduates > 0 {
			cell.SetValue(scReport.NumUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.NumGraduates > 0 {
			cell.SetValue(scReport.NumGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if severityFCReport.NumMale > 0 {
			cell.SetValue(severityFCReport.NumMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if severityFCReport.NumFemale > 0 {
			cell.SetValue(severityFCReport.NumFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if severityFCReport.NumUnderGraduates > 0 {
			cell.SetValue(severityFCReport.NumUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if severityFCReport.NumGraduates > 0 {
			cell.SetValue(severityFCReport.NumGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgOrangeStyle)
		cell.SetValue(scReport.NumTotal)
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
		if scReport.NumMale > 0 {
			cell.SetValue(scReport.NumMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.NumFemale > 0 {
			cell.SetValue(scReport.NumFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.NumUnderGraduates > 0 {
			cell.SetValue(scReport.NumUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.NumGraduates > 0 {
			cell.SetValue(scReport.NumGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if medicalDiagnosisFCReport.NumMale > 0 {
			cell.SetValue(medicalDiagnosisFCReport.NumMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if medicalDiagnosisFCReport.NumFemale > 0 {
			cell.SetValue(medicalDiagnosisFCReport.NumFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if medicalDiagnosisFCReport.NumUnderGraduates > 0 {
			cell.SetValue(medicalDiagnosisFCReport.NumUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if medicalDiagnosisFCReport.NumGraduates > 0 {
			cell.SetValue(medicalDiagnosisFCReport.NumGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgOrangeStyle)
		cell.SetValue(scReport.NumTotal)
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
		if scReport.NumMale > 0 {
			cell.SetValue(scReport.NumMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.NumFemale > 0 {
			cell.SetValue(scReport.NumFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.NumUnderGraduates > 0 {
			cell.SetValue(scReport.NumUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGrayStyle)
		if scReport.NumGraduates > 0 {
			cell.SetValue(scReport.NumGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if crisisFCReport.NumMale > 0 {
			cell.SetValue(crisisFCReport.NumMale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if crisisFCReport.NumFemale > 0 {
			cell.SetValue(crisisFCReport.NumFemale)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if crisisFCReport.NumUnderGraduates > 0 {
			cell.SetValue(crisisFCReport.NumUnderGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgGreenStyle)
		if crisisFCReport.NumGraduates > 0 {
			cell.SetValue(crisisFCReport.NumGraduates)
		}
		cell = row.AddCell()
		cell.SetStyle(bgOrangeStyle)
		cell.SetValue(scReport.NumTotal)
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
	cell.SetValue(emphasisTotalReport.NumMale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalReport.NumFemale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalReport.NumUnderGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalReport.NumGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalReport.NumMale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalReport.NumFemale)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalReport.NumUnderGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(emphasisTotalReport.NumGraduates)
	cell = row.AddCell()
	cell.SetStyle(bgRedStyle)
	cell.SetValue(emphasisTotalReport.NumTotal)
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
		cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.Grades[g])/float64(emphasisTotalReport.NumTotal)*100))
	}
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.NumMale)/float64(emphasisTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.NumFemale)/float64(emphasisTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.NumUnderGraduates)/float64(emphasisTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.NumGraduates)/float64(emphasisTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.NumMale)/float64(emphasisTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.NumFemale)/float64(emphasisTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.NumUnderGraduates)/float64(emphasisTotalReport.NumTotal)*100))
	cell = row.AddCell()
	cell.SetStyle(bgYellowStyle)
	cell.SetValue(fmt.Sprintf("%.2f%%", float64(emphasisTotalReport.NumGraduates)/float64(emphasisTotalReport.NumTotal)*100))
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
