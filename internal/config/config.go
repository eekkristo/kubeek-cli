package config

import (
	"encoding/json"
	"os"
	"strings"
)

type Config map[string]string

type AppConfig struct {
	Placeholders Config   `json:"placeholders"`
	Exts         []string `json:"exts"`
	ExcludeDirs  []string `json:"exclude_dirs"`
	ExcludeFiles []string `json:"exclude_files"`
	State        string   `json:"state"`

	// Behavior defaults (CLI may override)
	Unused       string `json:"unused"`        // keep|drop
	TemplateMode string `json:"template_mode"` // placeholder|values
}

func LoadAppConfig(path string) (AppConfig, bool, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return AppConfig{}, false, err
	}
	var ac AppConfig
	if err := json.Unmarshal(b, &ac); err == nil &&
		(ac.Placeholders != nil || ac.Exts != nil || ac.ExcludeDirs != nil || ac.ExcludeFiles != nil || ac.State != "" ||
			ac.Unused != "" || ac.TemplateMode != "") {
		return ac, true, nil
	}
	var legacy Config
	if err := json.Unmarshal(b, &legacy); err != nil {
		return AppConfig{}, false, err
	}
	return AppConfig{Placeholders: legacy}, false, nil
}

func SaveAppConfig(path string, ac AppConfig) error {
	b, err := json.MarshalIndent(ac, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func SaveConfig(path string, cfg Config) error {
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

type ResolvedOpts struct {
	StatePath    string
	Exts         []string
	ExcludeDirs  []string
	ExcludeFiles []string
}

func parseCSV(csv string) []string {
	if csv == "" {
		return nil
	}
	parts := strings.Split(csv, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func normalizeExts(exts []string) []string {
	out := make([]string, 0, len(exts))
	for _, e := range exts {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}
		if !strings.HasPrefix(e, ".") {
			e = "." + e
		}
		out = append(out, strings.ToLower(e))
	}
	return out
}

func ResolveOpts(flags map[string]string, ac AppConfig) ResolvedOpts {
	// defaults
	statePath := ".templater-state.json"
	exts := []string{".yaml", ".yml", ".tf"}
	excludeDirs := []string{".terraform"}
	excludeFiles := []string{}

	// config
	if ac.State != "" {
		statePath = ac.State
	}
	if len(ac.Exts) > 0 {
		exts = normalizeExts(ac.Exts)
	}
	if len(ac.ExcludeDirs) > 0 {
		excludeDirs = append([]string{}, ac.ExcludeDirs...)
	}
	if len(ac.ExcludeFiles) > 0 {
		excludeFiles = append([]string{}, ac.ExcludeFiles...)
	}

	// cli overrides
	if v := strings.TrimSpace(flags["state"]); v != "" {
		statePath = v
	}
	if v := strings.TrimSpace(flags["exts"]); v != "" {
		exts = normalizeExts(parseCSV(v))
	}
	if v := strings.TrimSpace(flags["exclude-dirs"]); v != "" {
		excludeDirs = parseCSV(v)
	}
	if v := strings.TrimSpace(flags["exclude-files"]); v != "" {
		excludeFiles = parseCSV(v)
	}

	return ResolvedOpts{
		StatePath:    statePath,
		Exts:         exts,
		ExcludeDirs:  excludeDirs,
		ExcludeFiles: excludeFiles,
	}
}

// Behavior normalization
func NormalizeUnused(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "drop", "remove", "delete":
		return "drop"
	case "", "keep":
		return "keep"
	default:
		return "keep"
	}
}
func NormalizeTemplateMode(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "values":
		return "values"
	case "", "placeholder":
		return "placeholder"
	default:
		return "placeholder"
	}
}
