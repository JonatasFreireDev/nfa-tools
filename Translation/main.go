package main

import (
	"fmt"
	"shared/log"
	"shared/spreadSheetReader"
	"time"
	"translation-tool/service"
)

func main() {
	startTime := time.Now()
	//entra em cada pasta referenciado em locale.
	//ler o arquivo que deve ser alterado
	// se o arquivo nao existir, criar com base em outro locale
	filesPath, err := service.FindFilesPath()

	if err != nil {
		log.WriteFile(err.Error())
	}

	//ler arquivo de traducao
	spreadSheet := spreadSheetReader.ReadXLS()

	//comparar o que deve ou nao ser alterado
	//alterar e salvar
	for folderName, filePath := range filesPath {
		service.UpdateLocales(folderName, filePath, spreadSheet)
	}

	endTime := time.Now()

	log.WriteFile("Tempo de execução: ")
	log.WriteFile(fmt.Sprintln(endTime.Sub(startTime)))
	fmt.Println("Done")
}
