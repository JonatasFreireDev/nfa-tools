package utils

import (
	"errors"
	"runtime"
	"slices"
	"strings"
)

func BarByOS() string {
	if runtime.GOOS == "windows" {
		return "\\"
	} else {
		return "/"
	}
}

func GetJsonFileName(filePath string) string {
	localeSplite := strings.Split(filePath, BarByOS())

	return localeSplite[len(localeSplite)-1]
}

func GetNfaFileName(filePath string) (string, error) {
	localeSplite := strings.Split(filePath, BarByOS())

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
