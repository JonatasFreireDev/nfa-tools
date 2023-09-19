package locale

import (
	"errors"
	"fmt"
	"runtime"
	"slices"
	"strings"
	"translation-tool/src/services/log"
)

func IdentifyOS() string {
	if runtime.GOOS == "windows" {
		return "\\"
	} else {
		return "/"
	}
}

func GetJsonFileName(filePath string) string {
	localeSplite := strings.Split(filePath, IdentifyOS())

	return localeSplite[len(localeSplite)-1]
}

func GetNfaFileName(filePath string) (string, error) {
	localeSplite := strings.Split(filePath, IdentifyOS())

	indexFolderName := slices.IndexFunc(localeSplite, func(s string) bool {
		foundFolder := strings.Contains(s, "nfa-")

		if foundFolder {
			return true
		} else {
			return false
		}
	})

	if indexFolderName < 0 {
		return "", errors.New("nfa folder not found")
	}

	jsonNameSplite := strings.Split(localeSplite[indexFolderName], "-")

	if len(jsonNameSplite) > 0 {
		return strings.Join(jsonNameSplite[1:], "-"), nil
	} else {
		return jsonNameSplite[0], nil
	}
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
