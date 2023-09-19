package spreadSheetReader

import (
	"translation-tool/src/config"

	"github.com/thedatashed/xlsxreader"
)

func ReadXLS() map[string]string {
	xl, _ := xlsxreader.OpenFile(config.File.TranslateFilePath)

	defer xl.Close()

	planilha := make(map[string]string)

	for row := range xl.ReadRows(xl.Sheets[0]) {
		planilha[row.Cells[0].Value] = row.Cells[1].Value
	}

	return planilha
}
