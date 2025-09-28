package scan

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Discover placeholders by scanning text for tokens like _name_ / _FOO_ / __Bar__
func DiscoverPlaceholders(root string, exts, excludeDirs, excludeFiles []string) ([]string, error) {
	set := map[string]struct{}{}

	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			for _, ex := range excludeDirs {
				if ex == d.Name() {
					return filepath.SkipDir
				}
			}
			return nil
		}
		base := d.Name()
		if skipFile(base, excludeFiles) {
			return nil
		}
		if len(exts) > 0 && !hasIncluded(base, exts) {
			return nil
		}
		b, rerr := os.ReadFile(p)
		if rerr != nil {
			return rerr
		}
		s := string(b)
		for i := 0; i < len(s); i++ {
			if s[i] != '_' {
				continue
			}
			j := strings.IndexByte(s[i+1:], '_')
			if j < 0 {
				break
			}
			j = i + 1 + j
			token := s[i : j+1]
			if validToken(token) {
				set[token] = struct{}{}
			}
			i = j
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i]) < strings.ToLower(out[j]) })
	return out, nil
}

/*
The string must be at least 3 characters long.
The first and last characters must be underscores (_).
Between the first and last underscores, there must be at least one character that is:
An uppercase or lowercase letter (A-Z, a-z)
A digit (0-9)
A dash (-)
An underscore (_)
If all these conditions are met, the function returns true; otherwise, it returns false.
*/
func validToken(t string) bool {
	if len(t) < 3 || t[0] != '_' || t[len(t)-1] != '_' {
		return false
	}
	for i := 1; i < len(t)-1; i++ {
		c := t[i]
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
			return true
		}
	}
	return false
}

func hasIncluded(base string, exts []string) bool {
	ext := strings.ToLower(filepath.Ext(base))
	for _, e := range exts {
		if strings.ToLower(e) == ext {
			return true
		}
	}
	return false
}

func skipFile(base string, excludeFiles []string) bool {
	lbase := strings.ToLower(base)
	for _, raw := range excludeFiles {
		p := strings.TrimSpace(raw)
		if p == "" {
			continue
		}
		lp := strings.ToLower(p)
		if strings.HasPrefix(lp, ".") && !strings.ContainsAny(lp, "*?[") {
			if strings.HasSuffix(lbase, lp) {
				return true
			}
			continue
		}
		if !strings.ContainsAny(lp, "*?[]/\\.") {
			if strings.HasSuffix(lbase, "."+lp) {
				return true
			}
			continue
		}
		if ok, _ := filepath.Match(lp, lbase); ok {
			return true
		}
		if lbase == lp {
			return true
		}
	}
	return false
}
