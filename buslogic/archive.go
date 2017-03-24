package buslogic

import (
	"github.com/shudiwsh2009/reservation_thxl_go/model"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"path/filepath"
)

var IMPORT_ARCHIVE_FILE = filepath.Join("static", "import", "archive.csv")

func (w *Workflow) ImportArchiveFromCSVFile() error {
	data, err := utils.ReadFromCSV(IMPORT_ARCHIVE_FILE)
	if err != nil {
		return err
	}
	for i := 1; i < len(data); i++ {
		if count, err := w.mongoClient.CountByStudentUsername(data[i][2]); err == nil && count == 0 {
			archive := &model.Archive{
				StudentUsername: data[i][2],
				ArchiveCategory: data[i][0],
				ArchiveNumber:   data[i][1],
			}
			w.mongoClient.InsertArchive(archive)
		}
	}
	return nil
}
