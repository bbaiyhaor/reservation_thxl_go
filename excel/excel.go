package excel

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"github.com/tealeg/xlsx"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

const (
	DefaultStudentExportExcelFilename   = "student_export_template.xlsx"
	DefaultTimetableExportExcelFilename = "timetable_export_template.xlsx"
	ExportFolder                        = "assets/export/"
	ExcelSuffix                         = ".xlsx"
	CsvSuffix                           = ".csv"
)

func ExportStudent(student *models.Student, filename string) error {
	xl, err := xlsx.OpenFile(filepath.FromSlash(ExportFolder + DefaultStudentExportExcelFilename))
	if err != nil {
		return errors.New("导出失败：打开模板文件失败")
	}
	sheet := xl.Sheet["export"]
	if sheet == nil {
		return errors.New("导出失败：打开工作表失败")
	}
	var row *xlsx.Row
	var cell *xlsx.Cell

	// 学生基本信息
	row = sheet.Rows[0]
	cell = row.AddCell()
	cell.SetString(student.Username)

	row = sheet.Rows[1]
	cell = row.AddCell()
	cell.SetString(student.Fullname)

	row = sheet.Rows[2]
	cell = row.AddCell()
	cell.SetString(student.Gender)

	row = sheet.Rows[3]
	cell = row.AddCell()
	cell.SetString(student.Birthday)

	row = sheet.Rows[4]
	cell = row.AddCell()
	cell.SetString(student.School)

	row = sheet.Rows[5]
	cell = row.AddCell()
	cell.SetString(student.Grade)

	row = sheet.Rows[6]
	cell = row.AddCell()
	cell.SetString(student.CurrentAddress)

	row = sheet.Rows[7]
	cell = row.AddCell()
	cell.SetString(student.FamilyAddress)

	row = sheet.Rows[8]
	cell = row.AddCell()
	cell.SetString(student.Mobile)

	row = sheet.Rows[9]
	cell = row.AddCell()
	cell.SetString(student.Email)

	row = sheet.Rows[10]
	if !student.Experience.IsEmpty() {
		cell = row.AddCell()
		cell.SetString("时间")
		cell = row.AddCell()
		cell.SetString(student.Experience.Time)
		cell = row.AddCell()
		cell.SetString("地点")
		cell = row.AddCell()
		cell.SetString(student.Experience.Location)
		cell = row.AddCell()
		cell.SetString("咨询师姓名")
		cell = row.AddCell()
		cell.SetString(student.Experience.Teacher)
	} else {
		cell = row.AddCell()
		cell.SetString("无")
	}

	row = sheet.Rows[11]
	cell = row.AddCell()
	cell.SetString("年龄")
	cell = row.AddCell()
	cell.SetString(student.FatherAge)
	cell = row.AddCell()
	cell.SetString("职业")
	cell = row.AddCell()
	cell.SetString(student.FatherJob)
	cell = row.AddCell()
	cell.SetString("学历")
	cell = row.AddCell()
	cell.SetString(student.FatherEdu)

	row = sheet.Rows[12]
	cell = row.AddCell()
	cell.SetString("年龄")
	cell = row.AddCell()
	cell.SetString(student.MotherAge)
	cell = row.AddCell()
	cell.SetString("职业")
	cell = row.AddCell()
	cell.SetString(student.MotherJob)
	cell = row.AddCell()
	cell.SetString("学历")
	cell = row.AddCell()
	cell.SetString(student.MotherEdu)

	row = sheet.Rows[13]
	cell = row.AddCell()
	cell.SetString(student.ParentMarriage)

	row = sheet.Rows[14]
	cell = row.AddCell()
	cell.SetString(student.Significant)

	row = sheet.Rows[15]
	cell = row.AddCell()
	cell.SetString(student.Problem)

	row = sheet.Rows[16]
	bindedTeacher, err := models.GetTeacherById(student.BindedTeacherId)
	if err != nil {
		cell = row.AddCell()
		cell.SetString("无")
	} else {
		cell = row.AddCell()
		cell.SetString(bindedTeacher.Username)
		cell = row.AddCell()
		cell.SetString(bindedTeacher.Fullname)
	}
	row = sheet.AddRow()

	// 特别注意事项

	//咨询小结
	row = sheet.AddRow()
	if reservations, err := models.GetReservationsByStudentId(student.Id.Hex()); err == nil {
		for i, r := range reservations {
			teacher, err := models.GetTeacherById(r.TeacherId)
			if err != nil {
				continue
			}
			row = sheet.AddRow()
			cell = row.AddCell()
			cell.SetString("咨询小结" + strconv.Itoa(i+1))

			row = sheet.AddRow()
			cell = row.AddCell()
			cell.SetString("问题评估")
			cell = row.AddCell()
			cell.SetString(r.TeacherFeedback.Problem)

			row = sheet.AddRow()
			cell = row.AddCell()
			cell.SetString("咨询师")
			cell = row.AddCell()
			cell.SetString(teacher.Fullname)

			row = sheet.AddRow()
			cell = row.AddCell()
			cell.SetString("咨询日期")
			cell = row.AddCell()
			cell.SetString(r.StartTime.In(utils.Location).Format(utils.DATE_PATTERN))

			row = sheet.AddRow()
			cell = row.AddCell()
			cell.SetString("咨询记录")
			cell = row.AddCell()
			cell.SetString(r.TeacherFeedback.Record)

			row = sheet.AddRow()
			cell = row.AddCell()
			cell.SetString("来访者反馈")
			for _, s := range r.StudentFeedback.Scores {
				cell = row.AddCell()
				cell.SetInt(s)
			}
		}
	}

	err = xl.Save(filepath.FromSlash(ExportFolder + filename))
	if err != nil {
		return errors.New("导出失败：保存文件失败")
	}
	return nil
}

func ExportReservationTimetable(reservations []*models.Reservation, filename string) error {
	xl, err := xlsx.OpenFile(filepath.FromSlash(ExportFolder + DefaultTimetableExportExcelFilename))
	if err != nil {
		return errors.New("导出失败：打开模板文件失败")
	}
	sheet := xl.Sheet["export"]
	if sheet == nil {
		return errors.New("导出失败：打开工作表失败")
	}
	var row *xlsx.Row
	var cell *xlsx.Cell

	for _, r := range reservations {
		teacher, err := models.GetTeacherById(r.TeacherId)
		if err != nil {
			return nil
		}
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetString(r.StartTime.In(utils.Location).Format(utils.TIME_PATTERN))
		cell = row.AddCell()
		cell.SetString(r.EndTime.In(utils.Location).Format(utils.TIME_PATTERN))
		cell = row.AddCell()
		cell.SetString(teacher.Fullname)
	}

	err = xl.Save(filepath.FromSlash(ExportFolder + filename))
	if err != nil {
		return errors.New("导出失败：保存文件失败")
	}
	return nil
}

type CategoryReport struct {
	Category      string
	UnderGraduate map[string]int
	Master        int
	Doctor        int
	Male          int
	Female        int
	Parents       int
	Teacher       int
	Instructor    int
	Amount        int
}

func ExportMonthlyReport(reservations []*models.Reservation, filename string) error {
	report := make(map[string]*CategoryReport)
	for index, category := range models.FeedbackAllCategory {
		report[index] = &CategoryReport{
			Category:      models.FeedbackAllCategory[category],
			UnderGraduate: make(map[string]int),
		}
	}
	amount := &CategoryReport{
		UnderGraduate: make(map[string]int),
	}
	for _, r := range reservations {
		if r.TeacherFeedback.IsEmpty() || len(r.TeacherFeedback.Participants) != 4 {
			continue
		}
		category := r.TeacherFeedback.Category
		// 学生
		if r.TeacherFeedback.Participants[0] > 0 {
			student, err := models.GetStudentById(r.StudentId)
			if err != nil {
				continue
			}
			switch string(student.Username[4]) {
			case "0":
				grade := student.Username[2:4]
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
	}
	grades := make([]string, 0)
	for g, _ := range amount.UnderGraduate {
		grades = append(grades, g)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(grades)))
	categories := make([]string, 0)
	for c, _ := range models.FeedbackAllCategory {
		categories = append(categories, c)
	}
	sort.Sort(sort.StringSlice(categories))

	// csv的表头
	head := []string{""}
	for _, g := range grades {
		head = append(head, g)
	}
	head = append(head, "硕", "博", "合计（男）", "合计（女）", "家长", "教师", "辅导员", "总计", "百分比")
	// csv的数据
	data := make([][]string, 0)
	for _, category := range categories {
		line := []string{models.FeedbackAllCategory[category]}
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
		line = append(line, strconv.Itoa(report[category].Male))
		line = append(line, strconv.Itoa(report[category].Female))
		if report[category].Parents > 0 {
			line = append(line, strconv.Itoa(report[category].Parents))
		} else {
			line = append(line, "")
		}
		if report[category].Teacher > 0 {
			line = append(line, strconv.Itoa(report[category].Teacher))
		} else {
			line = append(line, "")
		}
		if report[category].Instructor > 0 {
			line = append(line, strconv.Itoa(report[category].Instructor))
		} else {
			line = append(line, "")
		}
		line = append(line, strconv.Itoa(report[category].Amount))
		line = append(line, fmt.Sprintf("%#.02f%%", float64(report[category].Amount)/(float64(amount.Amount)/float64(100))))
		data = append(data, line)
	}
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
	amountLine = append(amountLine, strconv.Itoa(amount.Parents))
	percentLine = append(percentLine, fmt.Sprintf("%#.02f%%", float64(amount.Parents)/(float64(amount.Amount)/float64(100))))
	amountLine = append(amountLine, strconv.Itoa(amount.Teacher))
	percentLine = append(percentLine, fmt.Sprintf("%#.02f%%", float64(amount.Teacher)/(float64(amount.Amount)/float64(100))))
	amountLine = append(amountLine, strconv.Itoa(amount.Instructor))
	percentLine = append(percentLine, fmt.Sprintf("%#.02f%%", float64(amount.Instructor)/(float64(amount.Amount)/float64(100))))
	amountLine = append(amountLine, strconv.Itoa(amount.Amount))
	data = append(data, amountLine)
	data = append(data, percentLine)
	// 写入文件
	fout, err := os.Create(filepath.FromSlash(ExportFolder + filename))
	if err != nil {
		return errors.New("建立月报文件失败")
	}
	defer fout.Close()
	w := csv.NewWriter(transform.NewWriter(fout, simplifiedchinese.GB18030.NewEncoder()))
	w.UseCRLF = true
	w.Write(head)
	w.WriteAll(data)
	w.Flush()
	return nil
}
