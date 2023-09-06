package locale

import (
	"fmt"
	"runtime"
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

func GetJsonFileName(localePath string) string {
	localeSplite := strings.Split(localePath, IdentifyOS())
	folderName := localeSplite[len(localeSplite)-3]

	jsonNameSplite := strings.Split(folderName, "-")

	if len(jsonNameSplite) > 0 {
		restSlice := jsonNameSplite[1:]

		return strings.Join(restSlice, "-")
	}

	return jsonNameSplite[len(jsonNameSplite)-3]
}

func updateJson(data any, keyValueMap map[string]string) {
	switch v := data.(type) {
	case map[string]any:
		for key, val := range v {
			if strVal, ok := val.(string); ok {
				if newValue, exists := keyValueMap[strings.TrimSpace(strVal)]; exists {
					v[key] = newValue
					log.WriteFile(fmt.Sprintf("Substituído: %s -> %s", val, newValue))
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
