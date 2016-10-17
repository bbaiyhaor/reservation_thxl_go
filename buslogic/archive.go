package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/util"
)

const ArchiveFile = "../assets/import/archive.csv"

func (w *Workflow) ImportArchiveFromCSVFile() error {
	data, err := util.ReadFromCSV(ArchiveFile)
	if err != nil {
		return err
	}
	for i := 1; i < len(data); i++ {
		if archive, err := w.model.GetArchiveByStudentUsername(data[i][2]); err != nil || archive == nil {
			w.model.AddArchive(data[i][2], data[i][0], data[i][1])
		}
	}
	return nil
}
