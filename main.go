package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"translation-tool/src/config"
	"translation-tool/src/log"

	"github.com/thedatashed/xlsxreader"
)

type Translation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func IdentifyOS() string {
	if runtime.GOOS == "windows" {
		return "\\"
	} else {
		return "/"
	}
}

func GetJsonNameFile(localePath string) string {
	localeSplite := strings.Split(localePath, IdentifyOS())
	folderName := localeSplite[len(localeSplite)-1]

	jsonNameSplite := strings.Split(folderName, "-")

	if len(jsonNameSplite) > 0 {
		restSlice := jsonNameSplite[1:]

		return strings.Join(restSlice, "-")
	}

	return jsonNameSplite[len(jsonNameSplite)-1]
}

// Pega o caminho de cara locale, se nao existir, o locale e criado
func FindFilesLocale() (map[string]string, error) {
	foundLocales := map[string]string{}

	for _, localePath := range config.File.Locales {
		localeFolder, _ := os.ReadDir(localePath)
		fileWasFound := false
		folderName := GetJsonNameFile(localePath)

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

			//Cria diretorio
			err := os.Mkdir(nameDir, os.ModePerm)

			if err != nil {
				log.Write(err.Error())
			}

			//Cria arquivo json
			createdLocaleFile, err := os.Create(nameCreatedFile)

			if err != nil {
				log.Write(err.Error())
			}

			defer createdLocaleFile.Close()

			//Abre arquivo json de referencia.
			existingFile, err := os.Open(nameCopyFile)

			if err != nil {
				log.Write(err.Error())
			}

			defer existingFile.Close()

			//Copia conteudo de um arquivo para o outro
			_, err = io.Copy(createdLocaleFile, existingFile)

			if err != nil {
				log.Write(err.Error())
			}

			//Adiciona nos locales
			// foundLocales = append(foundLocales, nameCreatedFile)
			foundLocales[config.File.ToLocale+"/"+folderName] = nameCreatedFile
		}
	}

	if len(foundLocales) > 0 {
		return foundLocales, nil
	}

	return nil, errors.New("File Locale not found")
}

func ReadTranslationXLS() map[string]string {
	xl, _ := xlsxreader.OpenFile(config.File.TranslateFilePath)

	defer xl.Close()

	planilha := make(map[string]string)

	for row := range xl.ReadRows(xl.Sheets[0]) {
		planilha[row.Cells[0].Value] = row.Cells[1].Value
	}

	return planilha
}

func updateJson(data any, keyValueMap map[string]string) {
	switch v := data.(type) {
	case map[string]any:
		for key, val := range v {
			fmt.Println(key, val)
			if strVal, ok := val.(string); ok {
				if newValue, exists := keyValueMap[strings.TrimSpace(strVal)]; exists {
					v[key] = newValue
					log.Write(fmt.Sprintf("Substituído: %s -> %s", val, newValue))
				}
			} else {
				updateJson(val, keyValueMap)
			}
		}
	case []interface{}:
		// O valor é uma lista; faça uma chamada recursiva para cada elemento da lista
		for i, val := range v {
			updateJson(val, keyValueMap)
			v[i] = val // Atualize o valor na lista
		}
	}
}

func main() {

	startTime := time.Now()
	//entra em cada pasta referenciado em locale.
	//ler o arquivo que deve ser alterado
	// se o arquivo nao existir, criar com base em outro locale
	filesPath, err := FindFilesLocale()

	if err != nil {
		log.Write(err.Error())
	}

	//ler arquivo de traducao
	planilha := ReadTranslationXLS()

	//comparar o que deve ou nao ser alterado
	//alterar e salvar
	for folderName, filePath := range filesPath {
		log.Write(fmt.Sprintf("File: ---------------------------------------------------- %s ", folderName))
		file, err := os.Open(filePath)

		if err != nil {
			log.Write(err.Error())
		}

		defer file.Close()

		var jsonData any
		decoder := json.NewDecoder(file)

		if err := decoder.Decode(&jsonData); err != nil {
			log.Write(err.Error())
		}

		updateJson(jsonData, planilha)

		// Você pode codificar a estrutura de dados atualizada de volta em JSON
		jsonDataEncoded, err := json.MarshalIndent(jsonData, "", "")

		if err != nil {
			log.Write(err.Error())
		}

		fmt.Println(string(jsonDataEncoded))

		// Se desejar, você pode gravar os dados atualizados em um novo arquivo JSON
		outFile, err := os.Create(filePath)

		if err != nil {
			log.Write(err.Error())
		}

		defer outFile.Close()

		_, err = outFile.Write(jsonDataEncoded)

		if err != nil {
			log.Write(err.Error())
		}

	}

	endTime := time.Now()

	log.Write("Tempo de execução: ")
	log.Write(fmt.Sprintln(endTime.Sub(startTime)))
}

//para cada pasta em locale, eu posso ler o arquivo todo
