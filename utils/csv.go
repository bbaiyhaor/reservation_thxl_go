package utils

import (
	"encoding/csv"
	"errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"os"
	"path/filepath"
)

const (
	DefaultStudentExportExcelFilename   = "student_export_template.xlsx"
	DefaultTimetableExportExcelFilename = "timetable_export_template.xlsx"
	ExportFolder                        = "../assets/export/"
	ExcelSuffix                         = ".xlsx"
	CsvSuffix                           = ".csv"
)

func WriteToCSV(data [][]string, filename string) error {
	// 写入文件
	fout, err := os.Create(filepath.FromSlash(ExportFolder + filename))
	if err != nil {
		return errors.New("建立文件失败")
	}
	defer fout.Close()
	w := csv.NewWriter(transform.NewWriter(fout, simplifiedchinese.GB18030.NewEncoder()))
	w.UseCRLF = true
	if err = w.WriteAll(data); err != nil {
		return errors.New("写入表数据失败")
	}
	w.Flush()
	return nil
}

func ReadFromCSV(filename string) ([][]string, error) {
	fin, err := os.Open(filepath.FromSlash(filename))
	if err != nil || fin == nil {
		return nil, errors.New("打开文件失败")
	}
	defer fin.Close()
	w := csv.NewReader(transform.NewReader(fin, simplifiedchinese.GB18030.NewDecoder()))
	data, err := w.ReadAll()
	if err != nil {
		return nil, errors.New("读取文件失败")
	}
	return data, nil
}
