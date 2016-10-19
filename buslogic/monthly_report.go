package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/utils"
	"fmt"
	"sort"
	"strconv"
)

type MonthlyReport struct {
	Category      string
	UnderGraduate map[string]int
	Master        int
	Doctor        int
	Male          int
	Female        int
	Parents       int
	Teacher       int
	Instructor    int
	Other         int
	Amount        int
}

func (w *Workflow) ExportReportFormToFile(reservations []*model.Reservation, filename string) error {
	report := make(map[string]*MonthlyReport)
	for index, category := range model.FeedbackAllCategory {
		report[index] = &MonthlyReport{
			Category:      model.FeedbackAllCategory[category],
			UnderGraduate: make(map[string]int),
		}
	}
	amount := &MonthlyReport{
		UnderGraduate: make(map[string]int),
	}
	for _, r := range reservations {
		if r.TeacherFeedback.IsEmpty() || len(r.TeacherFeedback.Participants) != len(model.PARTICIPANTS) {
			continue
		}
		category := r.TeacherFeedback.Category
		// 学生
		if r.TeacherFeedback.Participants[0] > 0 {
			student, err := w.model.GetStudentById(r.StudentId)
			if err != nil {
				continue
			}
			switch string(student.Username[4]) {
			case "0":
				grade := student.Username[2:4] + "级"
				if _, exist := report[category].UnderGraduate[grade]; !exist {
					report[category].UnderGraduate[grade] = 0
				}
				if _, exist := amount.UnderGraduate[grade]; !exist {
					amount.UnderGraduate[grade] = 0
				}
				report[category].UnderGraduate[grade]++
				report[category].Amount++
				amount.UnderGraduate[grade]++
				amount.Amount++
			case "2":
				report[category].Master++
				report[category].Amount++
				amount.Master++
				amount.Amount++
			case "3":
				report[category].Doctor++
				report[category].Amount++
				amount.Doctor++
				amount.Amount++
			}
			switch student.Gender {
			case "男":
				report[category].Male++
				amount.Male++
			case "女":
				report[category].Female++
				amount.Female++
			}
		}
		// 家长
		if r.TeacherFeedback.Participants[1] > 0 {
			report[category].Parents++
			report[category].Amount++
			amount.Parents++
			amount.Amount++
		}
		// 教师
		if r.TeacherFeedback.Participants[2] > 0 {
			report[category].Teacher++
			report[category].Amount++
			amount.Teacher++
			amount.Amount++
		}
		// 辅导员
		if r.TeacherFeedback.Participants[3] > 0 {
			report[category].Instructor++
			report[category].Amount++
			amount.Instructor++
			amount.Amount++
		}
		// 其他
		if r.TeacherFeedback.Participants[4] > 0 {
			report[category].Other++
			report[category].Amount++
			amount.Other++
			amount.Amount++
		}
	}
	grades := make([]string, 0)
	for g := range amount.UnderGraduate {
		grades = append(grades, g)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(grades)))
	categories := make([]string, 0)
	for c := range model.FeedbackAllCategory {
		categories = append(categories, c)
	}
	sort.Sort(sort.StringSlice(categories))

	data := make([][]string, 0)
	// csv的表头
	head := []string{""}
	for _, g := range grades {
		head = append(head, g)
	}
	head = append(head, "硕", "博", "合计（男）", "合计（女）", "男女总计", "家长", "教师/辅导员", "其他", "辅助总计", "百分比")
	data = append(data, head)
	// csv的数据
	for _, category := range categories {
		line := []string{model.FeedbackAllCategory[category]}
		for _, g := range grades {
			if value, exist := report[category].UnderGraduate[g]; exist && value > 0 {
				line = append(line, strconv.Itoa(value))
			} else {
				line = append(line, "")
			}
		}
		if report[category].Master > 0 {
			line = append(line, strconv.Itoa(report[category].Master))
		} else {
			line = append(line, "")
		}
		if report[category].Doctor > 0 {
			line = append(line, strconv.Itoa(report[category].Doctor))
		} else {
			line = append(line, "")
		}
		if report[category].Male > 0 {
			line = append(line, strconv.Itoa(report[category].Male))
		} else {
			line = append(line, "")
		}
		if report[category].Female > 0 {
			line = append(line, strconv.Itoa(report[category].Female))
		} else {
			line = append(line, "")
		}
		line = append(line, strconv.Itoa(report[category].Male+report[category].Female))
		if report[category].Parents > 0 {
			line = append(line, strconv.Itoa(report[category].Parents))
		} else {
			line = append(line, "")
		}
		if report[category].Teacher+report[category].Instructor > 0 {
			line = append(line, strconv.Itoa(report[category].Teacher+report[category].Instructor))
		} else {
			line = append(line, "")
		}
		if report[category].Other > 0 {
			line = append(line, strconv.Itoa(report[category].Other))
		} else {
			line = append(line, "")
		}
		line = append(line, strconv.Itoa(report[category].Amount))
		line = append(line, fmt.Sprintf("%#.02f%%", float64(report[category].Amount)/(float64(amount.Amount)/float64(100))))
		data = append(data, line)
	}
	// csv的总计行和百分比行
	amountLine := []string{"总计"}
	percentLine := []string{"百分比"}
	for _, g := range grades {
		amountLine = append(amountLine, strconv.Itoa(amount.UnderGraduate[g]))
		percentLine = append(percentLine, fmt.Sprintf("%#.02f%%", float64(amount.UnderGraduate[g])/(float64(amount.Amount)/float64(100))))
	}
	amountLine = append(amountLine, strconv.Itoa(amount.Master))
	percentLine = append(percentLine, fmt.Sprintf("%#.02f%%", float64(amount.Master)/(float64(amount.Amount)/float64(100))))
	amountLine = append(amountLine, strconv.Itoa(amount.Doctor))
	percentLine = append(percentLine, fmt.Sprintf("%#.02f%%", float64(amount.Doctor)/(float64(amount.Amount)/float64(100))))
	amountLine = append(amountLine, strconv.Itoa(amount.Male))
	percentLine = append(percentLine, fmt.Sprintf("%#.02f%%", float64(amount.Male)/(float64(amount.Male+amount.Female)/float64(100))))
	amountLine = append(amountLine, strconv.Itoa(amount.Female))
	percentLine = append(percentLine, fmt.Sprintf("%#.02f%%", float64(amount.Female)/(float64(amount.Male+amount.Female)/float64(100))))
	amountLine = append(amountLine, strconv.Itoa(amount.Male+amount.Female))
	percentLine = append(percentLine, "")
	amountLine = append(amountLine, strconv.Itoa(amount.Parents))
	percentLine = append(percentLine, fmt.Sprintf("%#.02f%%", float64(amount.Parents)/(float64(amount.Amount)/float64(100))))
	amountLine = append(amountLine, strconv.Itoa(amount.Teacher+amount.Instructor))
	percentLine = append(percentLine, fmt.Sprintf("%#.02f%%", float64(amount.Teacher+amount.Instructor)/(float64(amount.Amount)/float64(100))))
	amountLine = append(amountLine, strconv.Itoa(amount.Other))
	percentLine = append(percentLine, fmt.Sprintf("%#.02f%%", float64(amount.Other)/(float64(amount.Amount)/float64(100))))
	amountLine = append(amountLine, strconv.Itoa(amount.Amount))
	data = append(data, amountLine)
	data = append(data, percentLine)
	if err := utils.WriteToCSV(data, filename); err != nil {
		return err
	}
	return nil
}

//func ExportKeyCaseReport(reservations []*models.Reservation, filename string) error {
//	students := make(map[string]*models.Student)
//	for _, r := range reservations {
//		if r.TeacherFeedback.IsEmpty() {
//			continue
//		}
//		student, err := models.GetStudentById(r.StudentId)
//		if err != nil || student == nil {
//			continue
//		}
//		if student.CrisisLevel > 0 {
//			students[student.Id.Hex()] = student
//		}
//	}
//	keyCase := make(map[int]*MonthlyReport)
//	for index, category := range models.KEY_CASE {
//		keyCase[index] = &MonthlyReport{
//			Category:      category,
//			UnderGraduate: make(map[string]int),
//		}
//	}
//	medicalDiagnosis := make(map[int]*MonthlyReport)
//	for index, category := range models.MEDICAL_DIAGNOSIS {
//		medicalDiagnosis[index] = &MonthlyReport{
//			Category:      category,
//			UnderGraduate: make(map[string]int),
//		}
//	}
//	amount := &MonthlyReport{
//		UnderGraduate: make(map[string]int),
//	}
//	for _, student := range students {
//		if student.CrisisLevel == 0 {
//			continue
//		}
//		switch string(student.Username[4]) {
//		case "0":
//			grade := student.Username[2:4] + "级"
//			for index, value := range student.KeyCase {
//				if value > 0 {
//					if _, exist := keyCase[index].UnderGraduate[grade]; !exist {
//						keyCase[index].UnderGraduate[grade] = 0
//					}
//					if _, exist := amount.UnderGraduate[grade]; !exist {
//						amount.UnderGraduate[grade] = 0
//					}
//					keyCase[index].UnderGraduate[grade]++
//					keyCase[index].Amount++
//					amount.UnderGraduate[grade]++
//					amount.Amount++
//				}
//			}
//			for index, value := range student.MedicalDiagnosis {
//				if value > 0 {
//					if _, exist := medicalDiagnosis[index].UnderGraduate[grade]; !exist {
//						medicalDiagnosis[index].UnderGraduate[grade] = 0
//					}
//					if _, exist := amount.UnderGraduate[grade]; !exist {
//						amount.UnderGraduate[grade] = 0
//					}
//					medicalDiagnosis[index].UnderGraduate[grade]++
//					medicalDiagnosis[index].Amount++
//					amount.UnderGraduate[grade]++
//					amount.Amount++
//				}
//			}
//		case "2":
//			for index, value := range student.KeyCase {
//				if value > 0 {
//					keyCase[index].Master++
//					keyCase[index].Amount++
//					amount.Master++
//					amount.Amount++
//				}
//			}
//			for index, value := range student.MedicalDiagnosis {
//				if value > 0 {
//					medicalDiagnosis[index].Master++
//					medicalDiagnosis[index].Amount++
//					amount.Master++
//					amount.Amount++
//				}
//			}
//		case "3":
//			for index, value := range student.KeyCase {
//				if value > 0 {
//					keyCase[index].Doctor++
//					keyCase[index].Amount++
//					amount.Doctor++
//					amount.Amount++
//				}
//			}
//			for index, value := range student.MedicalDiagnosis {
//				if value > 0 {
//					medicalDiagnosis[index].Doctor++
//					medicalDiagnosis[index].Amount++
//					amount.Doctor++
//					amount.Amount++
//				}
//			}
//		}
//		switch student.Gender {
//		case "男":
//			for index, value := range student.KeyCase {
//				if value > 0 {
//					keyCase[index].Male++
//					amount.Male++
//				}
//			}
//			for index, value := range student.MedicalDiagnosis {
//				if value > 0 {
//					medicalDiagnosis[index].Male++
//					amount.Male++
//				}
//			}
//		case "女":
//			for index, value := range student.KeyCase {
//				if value > 0 {
//					keyCase[index].Female++
//					amount.Female++
//				}
//			}
//			for index, value := range student.MedicalDiagnosis {
//				if value > 0 {
//					medicalDiagnosis[index].Female++
//					amount.Female++
//				}
//			}
//		}
//	}
//	grades := make([]string, 0)
//	for g, _ := range amount.UnderGraduate {
//		grades = append(grades, g)
//	}
//	sort.Sort(sort.Reverse(sort.StringSlice(grades)))
//
//	data := make([][]string, 0)
//	head := []string{""}
//	for _, g := range grades {
//		head = append(head, g)
//	}
//	head = append(head, "硕", "博", "合计（男）", "合计（女）", "男女合计", "辅助总计", "百分比")
//	data = append(data, head)
//	for index, category := range models.KEY_CASE {
//		line := []string{category}
//		for _, g := range grades {
//			if value, exist := keyCase[index].UnderGraduate[g]; exist && value > 0 {
//				line = append(line, strconv.Itoa(value))
//			} else {
//				line = append(line, "")
//			}
//		}
//		if keyCase[index].Master > 0 {
//			line = append(line, strconv.Itoa(keyCase[index].Master))
//		} else {
//			line = append(line, "")
//		}
//		if keyCase[index].Doctor > 0 {
//			line = append(line, strconv.Itoa(keyCase[index].Doctor))
//		} else {
//			line = append(line, "")
//		}
//		if keyCase[index].Male > 0 {
//			line = append(line, strconv.Itoa(keyCase[index].Male))
//		} else {
//			line = append(line, "")
//		}
//		if keyCase[index].Female > 0 {
//			line = append(line, strconv.Itoa(keyCase[index].Female))
//		} else {
//			line = append(line, "")
//		}
//		line = append(line, strconv.Itoa(keyCase[index].Male+keyCase[index].Female))
//		line = append(line, strconv.Itoa(keyCase[index].Amount))
//		line = append(line, fmt.Sprintf("%#.02f%%", float64(keyCase[index].Amount)/(float64(amount.Amount)/float64(100))))
//		data = append(data, line)
//	}
//	data = append(data, []string{""})
//	for index, category := range models.MEDICAL_DIAGNOSIS {
//		line := []string{category}
//		for _, g := range grades {
//			if value, exist := medicalDiagnosis[index].UnderGraduate[g]; exist && value > 0 {
//				line = append(line, strconv.Itoa(value))
//			} else {
//				line = append(line, "")
//			}
//		}
//		if medicalDiagnosis[index].Master > 0 {
//			line = append(line, strconv.Itoa(medicalDiagnosis[index].Master))
//		} else {
//			line = append(line, "")
//		}
//		if medicalDiagnosis[index].Doctor > 0 {
//			line = append(line, strconv.Itoa(medicalDiagnosis[index].Doctor))
//		} else {
//			line = append(line, "")
//		}
//		if medicalDiagnosis[index].Male > 0 {
//			line = append(line, strconv.Itoa(medicalDiagnosis[index].Male))
//		} else {
//			line = append(line, "")
//		}
//		if medicalDiagnosis[index].Female > 0 {
//			line = append(line, strconv.Itoa(medicalDiagnosis[index].Female))
//		} else {
//			line = append(line, "")
//		}
//		line = append(line, strconv.Itoa(medicalDiagnosis[index].Male+medicalDiagnosis[index].Female))
//		line = append(line, strconv.Itoa(medicalDiagnosis[index].Amount))
//		line = append(line, fmt.Sprintf("%#.02f%%", float64(medicalDiagnosis[index].Amount)/(float64(amount.Amount)/float64(100))))
//		data = append(data, line)
//	}
//	amountLine := []string{"总计（人）"}
//	percentLine := []string{"百分比"}
//	for _, g := range grades {
//		amountLine = append(amountLine, strconv.Itoa(amount.UnderGraduate[g]))
//		percentLine = append(percentLine, fmt.Sprintf("%#.02f%%", float64(amount.UnderGraduate[g])/(float64(amount.Amount)/float64(100))))
//	}
//	amountLine = append(amountLine, strconv.Itoa(amount.Master))
//	percentLine = append(percentLine, fmt.Sprintf("%#.02f%%", float64(amount.Master)/(float64(amount.Amount)/float64(100))))
//	amountLine = append(amountLine, strconv.Itoa(amount.Doctor))
//	percentLine = append(percentLine, fmt.Sprintf("%#.02f%%", float64(amount.Doctor)/(float64(amount.Amount)/float64(100))))
//	amountLine = append(amountLine, strconv.Itoa(amount.Male))
//	percentLine = append(percentLine, fmt.Sprintf("%#.02f%%", float64(amount.Male)/(float64(amount.Male+amount.Female)/float64(100))))
//	amountLine = append(amountLine, strconv.Itoa(amount.Female))
//	percentLine = append(percentLine, fmt.Sprintf("%#.02f%%", float64(amount.Female)/(float64(amount.Male+amount.Female)/float64(100))))
//	amountLine = append(amountLine, strconv.Itoa(amount.Male+amount.Female))
//	percentLine = append(percentLine, "")
//	percentLine = append(percentLine, fmt.Sprintf("%#.02f%%", float64(amount.Amount)/(float64(amount.Amount)/float64(100))))
//	amountLine = append(amountLine, strconv.Itoa(amount.Amount))
//	data = append(data, amountLine)
//	data = append(data, percentLine)
//
//	// 学生列表
//	data = append(data, []string{""})
//	data = append(data, []string{""})
//	data = append(data, []string{"姓名", "学号", "个案类型"})
//	for _, student := range students {
//		line := []string{student.Fullname, student.Username}
//		for index, value := range student.KeyCase {
//			if value > 0 {
//				line = append(line, models.KEY_CASE[index])
//			}
//		}
//		for index, value := range student.MedicalDiagnosis {
//			if value > 0 {
//				line = append(line, models.MEDICAL_DIAGNOSIS[index])
//			}
//		}
//		data = append(data, line)
//	}
//	if err := utils.WriteToCSV(data, filename); err != nil {
//		return err
//	}
//	return nil
//}
