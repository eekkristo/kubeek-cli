package config

import (
	"encoding/json"
	"os"
)

type Config map[string]string

type AppConfig struct {
	Placeholders Config   `json:"placeholders"`
	Exts         []string `json:"exts"`
	ExcludeDirs  []string `json:"exclude_dirs"`
	ExcludeFiles []string `json:"exclude_files"`
	State        string   `json:"state"`
}

func LoadAppConfig(path string) (AppConfig, bool, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return AppConfig{}, false, err
	}
	var ac AppConfig
	if err := json.Unmarshal(b, &ac); err == nil &&
		(ac.Placeholders != nil || ac.Exts != nil || ac.ExcludeDirs != nil || ac.ExcludeFiles != nil || ac.State != "") {
		return ac, true, nil
	}
	var legacy Config
	if err := json.Unmarshal(b, &legacy); err != nil {
		return AppConfig{}, false, err
	}
	return AppConfig{Placeholders: legacy}, false, nil
}
