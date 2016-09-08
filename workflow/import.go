package workflow

import (
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
)

const ArchiveFile = "../assets/import/archive.csv"

func ImportArchiveFromCSVFile() error {
	data, err := utils.ReadFromCSV(ArchiveFile)
	if err != nil {
		return err
	}
	for i := 1; i < len(data); i++ {
		if archive, err := models.GetArchiveByStudentUsername(data[i][2]); err != nil || archive == nil {
			models.AddArchive(data[i][2], data[i][0], data[i][1])
		}
	}
	return nil
}
