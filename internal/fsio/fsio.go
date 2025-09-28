package fsio

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type FileOp func(p string) error
type Done func() error

/*
 * Description: Filesystem input/output utilities
 * Usecase: Take external file types like: "tf, .yaml, .yml, txt" process them
 */

func Normalize(exts []string) []string {
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

// MatchAnyFilePattern returns true if base matches any of the exclude patterns.
// Supports: ".ext", "ext", globs (e.g., "*.bak"), and exact names.
func MatchAnyFilePattern(base string, patterns []string) bool {
	lbase := strings.ToLower(base)
	for _, raw := range patterns {
		p := strings.TrimSpace(raw)
		if p == "" {
			continue
		}
		lp := strings.ToLower(p)
		// ".tf" or ".yaml"
		if strings.HasPrefix(lp, ".") && !strings.ContainsAny(lp, "*?[") {
			if strings.HasSuffix(lbase, lp) {
				return true
			}
			continue
		}
		// bare token "tf" -> treat as extension
		if !strings.ContainsAny(lp, "*?[]/\\.") {
			if strings.HasSuffix(lbase, "."+lp) {
				return true
			}
			continue
		}
		// glob on basename
		if ok, _ := path.Match(lp, lbase); ok {
			return true
		}
		// exact match
		if lbase == lp {
			return true
		}
	}
	return false
}

// MatchesIncludedExt returns true if base has an extension listed in exts.
// If exts is empty, it returns true (i.e., include all).
func MatchesIncluded(base string, exts []string) bool {
	if len(exts) == 0 {
		return true
	}
	ext := strings.ToLower(filepath.Ext(base))
	if ext == "" {
		return false
	}
	for _, e := range exts {
		if strings.ToLower(e) == ext {
			return true
		}
	}
	return false
}

func FilesFiltered(root string, exts, excludeDirs, excludeFiles []string, fn FileOp, done Done) error {
	exts = Normalize(exts)
	exDir := map[string]struct{}{}
	for _, d := range excludeDirs {
		exDir[d] = struct{}{}
	}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		base := d.Name()
		if d.IsDir() {
			if _, skip := exDir[base]; skip {
				return filepath.SkipDir
			}
			return nil
		}
		if MatchAnyFilePattern(base, excludeFiles) {
			return nil
		}
		if !MatchesIncluded(base, exts) {
			return nil
		}
		return fn(p)
	})
	if err != nil {
		return err
	}
	if done != nil {
		return done()
	}
	return nil
}

func EnsureFreshDest(dst string, force bool, confirm func(string) (bool, error)) error {
	info, err := os.Stat(dst)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(dst, 0o755)
		}
		return err
	}
	if !info.IsDir() {
		return os.ErrInvalid
	}
	if !force {
		ok, err := confirm(dst)
		if err != nil {
			return err
		}
		if !ok {
			return os.ErrPermission
		}
	}
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	return os.MkdirAll(dst, 0o755)
}
