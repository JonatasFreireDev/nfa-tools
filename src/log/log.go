package log

import (
	"log"
	"os"
)

const logFileName = "logs.txt"

func createFile() *os.File {
	createdFile, err := os.Create(logFileName)

	if err != nil {
		panic(err.Error())
	}

	return createdFile
}

func Write(msg string) {
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0777)

	if err != nil {
		file = createFile()
	}

	defer file.Close()

	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime)

	log.Println(msg)
}
