package main

import (
	locale "command-line-argumentsD:\\code\\nfa-translation-tool\\Translation\\services\\locale.go"
	"fmt"
	"time"
)

func main() {
	startTime := time.Now()
	log.WriteFile("Tempo de execução: ")

	locale.FindFilesPath()
	log.WriteFile(fmt.Sprintln(endTime.Sub(startTime)))
	fmt.Println("Done")
}
