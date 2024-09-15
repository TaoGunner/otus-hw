package config

import (
	"encoding/json"
	"os"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger   CfgLogger   `json:"logger"`
	Database CfgDatabase `json:"database"`
	Rest     CfgRest     `json:"rest"`
}

type CfgLogger struct {
	Level       string `json:"level"`
	Format      string `json:"format"`
	IsAddSource bool   `json:"isAddSource"`
	LogFilename string `json:"logFilename"`
}

type CfgDatabase struct {
	Storage string `json:"storage"`
	DBPath  string `json:"dbPath,omitempty"`
}

type CfgRest struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

// NewConfig возвращает настройки, прочитанные из файла по указанному пути.
//
//	В случае ошибки будет возвращены настройки по-умолчанию
func NewConfig(cfgPath string) (Config, error) {
	// Настройки по-умолчанию
	cfg := Config{
		Logger: CfgLogger{
			Level:       "info",
			Format:      "logfmt",
			IsAddSource: false,
			LogFilename: "application.log",
		},
		Database: CfgDatabase{
			Storage: "sql",
			DBPath:  "database.sqlite3",
		},
		Rest: CfgRest{
			Host: "0.0.0.0",
			Port: 8000,
		},
	}

	cfgData, err := os.ReadFile(cfgPath)
	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(cfgData, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
