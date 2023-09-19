package main

import (
	"fmt"
	"time"
	"translation-tool/src/services/locale"
	"translation-tool/src/services/log"
	"translation-tool/src/services/spreadSheetReader"
)

func main() {
	startTime := time.Now()
	log.WriteFile("Tempo de execução: ")
	log.WriteFile(fmt.Sprintln(endTime.Sub(startTime)))
	fmt.Println("Done")
}
