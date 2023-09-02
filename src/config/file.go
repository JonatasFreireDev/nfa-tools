package config

import (
	"encoding/json"
	"os"
	mylog "translation-tool/src/log"
)

var File *Config

const configFileName = "config.json"

type Config struct {
	Locales           []string `json:"localesPath"`
	CopyFromLocale    string   `json:"CopyFromLocale"`
	ToLocale          string   `json:"toLocale"`
	TranslateFilePath string   `json:"translateFilePath"`
}

func init() {
	file, err := ReadFile()

	if err != nil {
		CreateFile()

		panic("Arquivo de configuracao criado, por favor configure-o.")
	}

	File = file
}

func CreateFile() {
	createdFile, err := os.Create(configFileName)

	if err != nil {
		panic("Nao foi possivel criar arquivo de configuracao.")
	}

	defer createdFile.Close()

	jsonString, err := json.MarshalIndent(Config{
		CopyFromLocale:    "xx-XX",
		ToLocale:          "xx-XX",
		TranslateFilePath: "pathToXLSX",
		Locales:           []string{"pathToLocaleString"},
	}, "", "")

	if err != nil {
		mylog.Write(err.Error())
	}

	_, err = createdFile.Write(jsonString)

	if err != nil {
		mylog.Write(err.Error())
	}

	mylog.Write("config file created")
}

func ReadFile() (*Config, error) {
	file, err := os.ReadFile(configFileName)

	if err != nil {
		mylog.Write(err.Error())
		return nil, err
	}

	var config Config

	json.Unmarshal(file, &config)

	return &config, nil
}
