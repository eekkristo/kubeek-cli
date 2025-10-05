package config

import (
	"encoding/json"
	"os"
)

type Config map[string]string

type AppConfig struct {
	Placeholders Config   `json:"placeholders"`  // Placeholder values for rendering template
	Exts         []string `json:"exts"`          // Extensions to read from
	ExcludeDirs  []string `json:"exclude_dirs"`  // Excluded folders from rendering
	ExcludeFiles []string `json:"exclude_files"` // Excluded files from rendering
}

func DefaultAppConfig() AppConfig {
	return AppConfig{
		Placeholders: make(Config),
		Exts:         []string{},
	}
}

func LoadAppConfig(path string) (AppConfig, bool, error) {
	ac := DefaultAppConfig()
	b, err := os.ReadFile(path)
	if err != nil {
		return ac, false, err
	}

	if ac.Placeholders == nil {
		ac.Placeholders = make(Config)
	}

	if ac.Exts == nil {
		ac.Exts = []string{}
	}

	if err = json.Unmarshal(b, &ac); err != nil {
		return ac, false, err
	}

	return ac, true, nil
}
