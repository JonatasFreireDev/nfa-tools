package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"shared/config"
	"shared/log"
	"shared/utils"
	"strings"
)

// Pega o caminho de cada locale, se nao existir, o locale e criado
func FindFilesPath() (map[string]string, error) {
	foundLocales := map[string]string{}

	if len(config.File.Locales) == 0 || config.File.Locales == nil {
		return nil, errors.New("Theres is no Locale strings in config.json")
	}

	for _, localePath := range config.File.Locales {
		localesFolders, _ := os.ReadDir(localePath)
		nfaName, err := utils.GetNfaFileName(localePath)
		fileWasFound := false

		if err != nil {
			return nil, err
		}

		for _, localeFolder := range localesFolders {
			//entra na pasta da lingua desejada
			if localeFolder.Name() == config.File.ToLocale {
				filesInLocaleFolder, _ := os.ReadDir(filepath.Join(localePath, localeFolder.Name()))

				fileWasFound = true
				foundLocales[config.File.ToLocale+"/"+nfaName] = filepath.Join(localePath, config.File.ToLocale, filesInLocaleFolder[0].Name())
				// foundLocales = append(foundLocales, filepath.Join(localePath, config.File.ToLocale, filesInLocaleFolder[0].Name()))
			}
		}

		//se o arquivo .json nao existe, sera criado aparti
		if fileWasFound == false {
			nameDir := filepath.Join(localePath, config.File.ToLocale)
			nameCreatedFile := filepath.Join(localePath, config.File.ToLocale, nfaName+".json")
			nameCopyFile := filepath.Join(localePath, config.File.CopyFromLocale, nfaName+".json")

			CreateFile(nameDir, nameCreatedFile, nameCopyFile)

			//Adiciona nos locales
			// foundLocales = append(foundLocales, nameCreatedFile)
			foundLocales[config.File.ToLocale+"/"+nfaName] = nameCreatedFile
		}
	}

	if len(config.File.CustomFilesPaths) > 0 {
		for _, filePath := range config.File.CustomFilesPaths {
			isJson := strings.Contains(filePath, ".json")

			if !isJson {
				continue
			}

			nfaName, _ := utils.GetNfaFileName(filePath)
			jsonFileName := utils.GetJsonFileName(filePath)

			foundLocales[config.File.ToLocale+"/"+nfaName+"/"+jsonFileName] = filePath
		}
	}

	if len(foundLocales) > 0 {
		return foundLocales, nil
	}

	return nil, errors.New("File Locale not found")
}

func CreateFile(nameDir string, nameCreatedFile string, nameCopyFile string) {
	defer func() {
		if msg := recover(); msg != nil {
			log.WriteFile(fmt.Errorf("%v", msg).Error())
		}
	}()

	//Cria diretorio
	os.Mkdir(nameDir, os.ModePerm)

	//Cria arquivo json
	createdLocaleFile, _ := os.Create(nameCreatedFile)
	defer createdLocaleFile.Close()

	//Abre arquivo json de referencia.
	existingFile, _ := os.Open(nameCopyFile)
	defer existingFile.Close()

	//Copia conteudo de um arquivo para o outro
	io.Copy(createdLocaleFile, existingFile)
}

func UpdateLocales(folderName string, filePath string, spreadSheet map[string]string) {
	defer func() {
		if msg := recover(); msg != nil {
			log.WriteFile(fmt.Errorf("%v", msg).Error())
		}
	}()

	log.WriteFile(fmt.Sprintf("File: ---------------------------------------------------- %s ", folderName))

	file, _ := os.Open(filePath)
	defer file.Close()

	var jsonData any
	json.NewDecoder(file).Decode(&jsonData)

	updateJson(jsonData, spreadSheet, []string{})

	// Você pode codificar a estrutura de dados atualizada de volta em JSON
	jsonDataEncoded, _ := json.MarshalIndent(jsonData, "", "  ")

	// Se desejar, você pode gravar os dados atualizados em um novo arquivo JSON
	outFile, _ := os.Create(filePath)

	defer outFile.Close()

	outFile.Write(jsonDataEncoded)
}

func updateJson(data any, keyValueMap map[string]string, path []string) {
	switch v := data.(type) {
	case map[string]any:
		for key, val := range v {
			path = append(path, key)
			if strVal, ok := val.(string); ok {
				if newValue, exists := keyValueMap[strings.TrimSpace(strVal)]; exists {
					v[key] = newValue
					log.WriteFile(fmt.Sprintf("Substituído: [%s] %s -> %s", strings.Join(path, "."), val, newValue))
				}
			} else {
				updateJson(val, keyValueMap, path)
			}

			path = path[:len(path)-1]
		}
	case []interface{}:
		// O valor é uma lista; faça uma chamada recursiva para cada elemento da lista
		for i, val := range v {
			updateJson(val, keyValueMap, path)
			v[i] = val // Atualize o valor na lista
		}
	}
}
