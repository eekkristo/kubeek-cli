package common

import (
	"fmt"
	"path/filepath"
	"strings"

	cfg "kubeekcli/internal/config"
	"kubeekcli/internal/fsio"
	"kubeekcli/internal/prompt"
)

// EnsureDestFresh wipes or creates dst; asks before overwrite unless force.
func EnsureDestFresh(dst string, force bool) error {
	return fsio.EnsureFreshDest(dst, force, prompt.ConfirmOverwrite)
}

// AbsOrRelKey returns both absolute and relative keys for state lookup/cleanup.
func AbsOrRelKey(p string) (abs, rel string) {
	rel = p
	if a, err := filepath.Abs(p); err == nil {
		abs = a
	}
	return abs, rel
}

// BuildFinalConfig merges placeholder maps based on policy keep|drop.
func BuildFinalConfig(existing, found cfg.Config, policy string) cfg.Config {
	if cfg.NormalizeUnused(policy) == "drop" {
		return found
	}
	out := cfg.Config{}
	for k, v := range existing {
		out[k] = v
	}
	for k, v := range found {
		out[k] = v
	}
	return out
}

// CSV helpers (used occasionally by runners when flags are strings)
func ParseCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// Small UX helpers
func AbortUser() error { return fmt.Errorf("aborted by user") }
