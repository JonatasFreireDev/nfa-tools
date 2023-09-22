package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"shared/config"
	"shared/log"
	"shared/utils"
	"sort"
	"strings"
	"time"
)

func main() {
	startTime := time.Now()
	if len(config.File.Locales) == 0 || config.File.Locales == nil {
		// return nil, errors.New("Theres is no Locale strings in config.json")
		//errr
	}

	defer func() {
		if msg := recover(); msg != nil {
			log.WriteFile(fmt.Errorf("%v", msg).Error())
		}
	}()

	// for _, localePath := range config.File.Locales {
	// 	localeFolderPaths := make(map[string]string)
	// 	folderName, err := utils.GetNfaFileName(localePath)

	// 	if err != nil {
	// 		// return nil, err
	// 	}

	// 	localesFolders, _ := os.ReadDir(localePath)

	// 	for _, localeFolder := range localesFolders {
	// 		filesInLocaleFolder, _ := os.ReadDir(filepath.Join(localePath, localeFolder.Name()))
	// 		foundLocales[folderName] = map[string]string{}
	// 		localeFolderPaths[localeFolder.Name()] = filepath.Join(localePath, localeFolder.Name(), filesInLocaleFolder[0].Name())
	// 	}

	// 	foundLocales[folderName] = localeFolderPaths
	// }

	for _, localePath := range config.File.Locales {
		localeFolderPaths := make(map[string]map[string]string)

		nfaFileName, err := utils.GetNfaFileName(localePath)

		if err != nil {
			log.WriteFile(err.Error())
		}

		localesFolders, err := os.ReadDir(localePath)

		if err != nil {
			log.WriteFile(err.Error())
		}

		for _, localeFolder := range localesFolders {
			teste := make(map[string]string)
			filesInLocaleFolder, err := os.ReadDir(filepath.Join(localePath, localeFolder.Name()))
			jsonFile, err := os.Open(filepath.Join(localePath, localeFolder.Name(), filesInLocaleFolder[0].Name()))
			defer jsonFile.Close()

			if err != nil {
				log.WriteFile(err.Error())
			}

			var jsonData any
			json.NewDecoder(jsonFile).Decode(&jsonData)

			updateJson(jsonData, []string{}, teste)

			localeFolderPaths[localeFolder.Name()] = teste
		}

		jsonDataEncoded, _ := json.MarshalIndent(localeFolderPaths, "", "  ")

		fmt.Println(string(jsonDataEncoded))

		log.WriteFile(string(jsonDataEncoded))

		time := time.Now().Format("02-01-2006_15:04")

		// Cria CSV para cada repo
		file, err := os.Create(nfaFileName + "_" + time + ".csv")

		if err != nil {
			log.WriteFile(err.Error())
		}

		defer file.Close()

		// Crie um buffer de bytes para adicionar o BOM UTF-8
		var buffer bytes.Buffer
		buffer.Write([]byte{0xEF, 0xBB, 0xBF}) // Adicione o BOM UTF-8

		writerCSV := csv.NewWriter(&buffer)
		writerCSV.Comma = ';'

		// Obtenha uma lista de todas as chaves externas (chave1, chave2, etc.)
		var locales []string
		for locale := range localeFolderPaths {
			locales = append(locales, locale)
		}

		// Ordene os cabeçalhos das colunas em ordem alfabética
		sort.Strings(locales)

		// Escreva os cabeçalhos de coluna no arquivo CSV
		if err := writerCSV.Write(append([]string{""}, locales...)); err != nil {
			log.WriteFile(err.Error())
		}

		// Percorra as subchaves (rótulos de linha) e escreva os valores no arquivo CSV
		for subchave := range localeFolderPaths[locales[0]] {
			linha := []string{subchave}
			for _, chaveExterna := range locales {
				valor := localeFolderPaths[chaveExterna][subchave]
				linha = append(linha, `"`+valor+`"`)
			}
			if err := writerCSV.Write(linha); err != nil {
				log.WriteFile(err.Error())
			}
			log.WriteFile(subchave)
		}

		// Certifique-se de liberar os recursos e escrever os dados no arquivo
		writerCSV.Flush()

		// Verificar erros após o flush
		if err := writerCSV.Error(); err != nil {
			log.WriteFile(err.Error())
		}

		// Escreva o buffer no arquivo
		file.Write(buffer.Bytes())

		if err != nil {
			log.WriteFile(err.Error())
		}

	}

	endTime := time.Now()

	log.WriteFile("Tempo de execução: ")
	log.WriteFile(fmt.Sprintln(endTime.Sub(startTime)))
	fmt.Println("Done")
}

func updateJson(data any, path []string, teste map[string]string) {
	switch v := data.(type) {
	case map[string]any:
		for key, val := range v {
			path = append(path, key)
			if value, ok := val.(string); ok {
				log.WriteFile(strings.Join(path, "."))

				teste[strings.Join(path, ".")] = value
			} else {
				updateJson(val, path, teste)
			}
			path = path[:len(path)-1]
		}
	}
}
