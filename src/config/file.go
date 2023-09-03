package config

import (
	"encoding/json"
	"os"
	"translation-tool/src/services/log"
)

var File *Config

const configFileName = "config.json"

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
		log.WriteFile(err.Error())
	}

	_, err = createdFile.Write(jsonString)

	if err != nil {
		log.WriteFile(err.Error())
	}

	log.WriteFile("config file created")
}

func ReadFile() (*Config, error) {
	file, err := os.ReadFile(configFileName)

	if err != nil {
		log.WriteFile(err.Error())
		return nil, err
	}

	var config Config

	json.Unmarshal(file, &config)

	return &config, nil
}
