package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cfg "kubeekcli/internal/config"
	"kubeekcli/internal/engine/replace"
	"kubeekcli/internal/engine/scan"
	st "kubeekcli/internal/engine/state"
	"kubeekcli/internal/fsio"
	"kubeekcli/internal/prompt"
)

// RunInteractive prompts, renders into dst, and optionally emits BOTH config+state.
// State keys are RELATIVE to dst.
func RunInteractive(src, dst string, ac cfg.AppConfig, opts cfg.ResolvedOpts, force bool, meta *bool) error {

	if info, err := os.Stat(dst); err == nil && info.IsDir() {
		if err := os.RemoveAll(dst); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}

	found, err := scan.DiscoverPlaceholders(src, opts.Exts, opts.ExcludeDirs, opts.ExcludeFiles)
	if err != nil {
		return err
	}
	seen := map[string]struct{}{}
	for _, ph := range found {
		seen[ph] = struct{}{}
	}
	for ph := range ac.Placeholders {
		if _, ok := seen[ph]; !ok {
			found = append(found, ph)
		}
	}

	answers := cfg.Config{}
	if len(found) > 0 {
		fmt.Println("Provide values for placeholders (press Enter to accept default).")
		for _, ph := range found {
			def := ac.Placeholders[ph]
			promptTxt := ph
			if def != "" {
				promptTxt += " [" + def + "]: "
			} else {
				promptTxt += ": "
			}
			line, _ := prompt.Prompt(promptTxt)
			line = strings.TrimSpace(line)
			if line == "" {
				answers[ph] = def
			} else {
				answers[ph] = line
			}
		}
	}

	genState := st.State{}
	err = filepath.WalkDir(src, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if d.IsDir() {
			exDir := map[string]struct{}{}
			for _, ed := range opts.ExcludeDirs {
				exDir[ed] = struct{}{}
			}
			if _, skip := exDir[d.Name()]; skip {
				return filepath.SkipDir
			}
			return os.MkdirAll(target, 0o755)
		}

		base := d.Name()
		if fsio.MatchAnyFilePattern(base, opts.ExcludeFiles) {
			return nil
		}

		data, rerr := os.ReadFile(p)
		if rerr != nil {
			return rerr
		}
		content := string(data)

		if len(opts.Exts) > 0 && !fsio.MatchesIncluded(base, opts.Exts) {
			return os.WriteFile(target, []byte(content), 0644)
		}

		lines := strings.Split(content, "\n")
		newLines, entries := replace.RenderLinesWithState(lines, answers)
		if err := os.WriteFile(target, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
			return err
		}
		if len(entries) > 0 {
			genState[rel] = entries // RELATIVE key
		}
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Println("Generated into:", dst)

	return nil
}
