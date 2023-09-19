package config

type Config struct {
	CopyFromLocale    string   `json:"copyFromLocale"`
	ToLocale          string   `json:"toLocale"`
	TranslateFilePath string   `json:"translateFilePath"`
	Locales           []string `json:"localesPath"`
	CustomFilesPaths  []string `json:"customFilesPaths"`
}
