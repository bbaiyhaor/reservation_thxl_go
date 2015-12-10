package excel

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/tealeg/xlsx"
	"path/filepath"
	"strconv"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
)

const (
	DefaultStudentExportExcelFilename   = "student_export_template.xlsx"
	DefaultTimetableExportExcelFilename = "timetable_export_template.xlsx"
	ExportFolder                        = "assets/export/"
	ExcelSuffix                         = ".xlsx"
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
				cell.SetString(s)
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
