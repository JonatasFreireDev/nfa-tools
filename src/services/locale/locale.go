package locale

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"translation-tool/src/config"
	"translation-tool/src/services/log"
)

func CreateFile(nameDir string, nameCreatedFile string, nameCopyFile string) {
	//Cria diretorio
	err := os.Mkdir(nameDir, os.ModePerm)

	if err != nil {
		log.WriteFile(err.Error())
	}

	//Cria arquivo json
	createdLocaleFile, err := os.Create(nameCreatedFile)

	if err != nil {
		log.WriteFile(err.Error())
	}

	defer createdLocaleFile.Close()

	//Abre arquivo json de referencia.
	existingFile, err := os.Open(nameCopyFile)

	if err != nil {
		log.WriteFile(err.Error())
	}

	defer existingFile.Close()

	//Copia conteudo de um arquivo para o outro
	_, err = io.Copy(createdLocaleFile, existingFile)

	if err != nil {
		log.WriteFile(err.Error())
	}
}

// Pega o caminho de cada locale, se nao existir, o locale e criado
func FindFilesPath() (map[string]string, error) {
	foundLocales := map[string]string{}

	if len(config.File.Locales) == 0 || config.File.Locales == nil {
		return nil, errors.New("Theres is no Locale strings in config.json")
	}

	for _, localePath := range config.File.Locales {
		localeFolder, _ := os.ReadDir(localePath)
		fileWasFound := false
		folderName, err := GetNfaFileName(localePath)

		if err != nil {
			return nil, err
		}

		for _, folderInLocale := range localeFolder {
			//entra na pasta da lingua desejada
			if folderInLocale.Name() == config.File.ToLocale {
				localeTranslateFolder, _ := os.ReadDir(filepath.Join(localePath, folderInLocale.Name()))

				fileWasFound = true
				foundLocales[config.File.ToLocale+"/"+folderName] = filepath.Join(localePath, config.File.ToLocale, localeTranslateFolder[0].Name())
				// foundLocales = append(foundLocales, filepath.Join(localePath, config.File.ToLocale, localeTranslateFolder[0].Name()))
			}
		}

		//se o arquivo .json nao existe, sera criado aparti
		if fileWasFound == false {
			nameDir := filepath.Join(localePath, config.File.ToLocale)
			nameCreatedFile := filepath.Join(localePath, config.File.ToLocale, folderName+".json")
			nameCopyFile := filepath.Join(localePath, config.File.CopyFromLocale, folderName+".json")

			CreateFile(nameDir, nameCreatedFile, nameCopyFile)

			//Adiciona nos locales
			// foundLocales = append(foundLocales, nameCreatedFile)
			foundLocales[config.File.ToLocale+"/"+folderName] = nameCreatedFile
		}
	}

	if len(config.File.CustomFilesPaths) > 0 {
		for _, filePath := range config.File.CustomFilesPaths {
			isJson := strings.Contains(filePath, ".json")

			if !isJson {
				continue
			}

			folderName, _ := GetNfaFileName(filePath)
			fileName := GetJsonFileName(filePath)

			foundLocales[config.File.ToLocale+"/"+folderName+"/"+fileName] = filePath
		}
	}

	if len(foundLocales) > 0 {
		return foundLocales, nil
	}

	return nil, errors.New("File Locale not found")
}

func UpdateLocales(folderName string, filePath string, spreadSheet map[string]string) {
	log.WriteFile(fmt.Sprintf("File: ---------------------------------------------------- %s ", folderName))

	file, err := os.Open(filePath)

	if err != nil {
		log.WriteFile(err.Error())
	}

	defer file.Close()

	var jsonData any
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&jsonData); err != nil {
		log.WriteFile(err.Error())
	}

	updateJson(jsonData, spreadSheet, []string{})

	// Você pode codificar a estrutura de dados atualizada de volta em JSON
	jsonDataEncoded, err := json.MarshalIndent(jsonData, "", "")

	if err != nil {
		log.WriteFile(err.Error())
	}

	// Se desejar, você pode gravar os dados atualizados em um novo arquivo JSON
	outFile, err := os.Create(filePath)

	if err != nil {
		log.WriteFile(err.Error())
	}

	defer outFile.Close()

	_, err = outFile.Write(jsonDataEncoded)

	if err != nil {
		log.WriteFile(err.Error())
	}
}
