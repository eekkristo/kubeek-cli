package fsio

import (
	"os"
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
