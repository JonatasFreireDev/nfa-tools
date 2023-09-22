package service

import (
	"encoding/csv"
	"os"
	"shared/log"
	"sort"
)

func OrderByNameCSV(filePath string) {
	// Abrir o arquivo CSV para leitura
	file, err := os.Open(filePath)

	if err != nil {
		log.WriteFile(err.Error())
	}
	defer file.Close()

	// Ler os dados do arquivo CSV
	readerCSV := csv.NewReader(file)
	readerCSV.Comma = ';'

	data, err := readerCSV.ReadAll()
	if err != nil {
		log.WriteFile(err.Error())
	}

	// Separar a primeira linha dos dados
	dataHeader := data[0]
	data = data[1:]

	// Ordenar os dados com base na primeira coluna (coluna 0)
	sort.SliceStable(data, func(i, j int) bool {
		return data[i][0] < data[j][0]
	})

	file.Close()

	// Criar um arquivo CSV de saída
	outFile, err := os.Create(filePath)
	if err != nil {
		log.WriteFile(err.Error())
	}
	defer outFile.Close()

	// Criar um escritor CSV para o arquivo de saída
	writerCSV := csv.NewWriter(outFile)
	writerCSV.Comma = ';'

	// Adicione o BOM UTF-8 no início do arquivo
	bom := []byte{0xEF, 0xBB, 0xBF}
	if _, err := outFile.Write(bom); err != nil {
		log.WriteFile(err.Error())
	}

	// Escrever a primeira linha no arquivo de saída
	if err := writerCSV.Write(dataHeader); err != nil {
		log.WriteFile(err.Error())
	}

	// Escrever os dados ordenados no arquivo CSV de saída
	for _, row := range data {
		if err := writerCSV.Write(row); err != nil {
			log.WriteFile(err.Error())
		}
	}

	// Certificar-se de liberar os recursos e escrever os dados no arquivo de saída
	writerCSV.Flush()

	// Verificar erros após o flush
	if err := writerCSV.Error(); err != nil {
		log.WriteFile(err.Error())
	}
}
