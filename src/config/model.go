package config

type Config struct {
	Locales           []string `json:"localesPath"`
	CopyFromLocale    string   `json:"CopyFromLocale"`
	ToLocale          string   `json:"toLocale"`
	TranslateFilePath string   `json:"translateFilePath"`
}
