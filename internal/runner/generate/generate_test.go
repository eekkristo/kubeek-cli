package generate_test

import (
	"io"
	cfg "kubeekcli/internal/config"
	"kubeekcli/internal/runner/generate"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, p, s string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(p, []byte(s), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func readFile(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	return string(b)
}

func TestGenerate(t *testing.T) {
	td := t.TempDir()
	src := filepath.Join(td, "kubeekcli_template")
	dst := filepath.Join(td, "kubeekcli_generated")

	const testvalue = `clustername: "my-test-cluster"`

	// Template with one placeholder
	fp := filepath.Join(src, "app.yaml")
	writeFile(t, fp, `clustername: "_clustername_"`)

	ac := cfg.AppConfig{
		Placeholders: cfg.Config{"_clustername_": "my-test-cluster"}, // default used if user presses Enter
	}

	// Simulate user pressing Enter (accept default), so "my-test-cluster" is used.
	// RunInteractive uses prompt.AskLine (reading os.Stdin),
	// we set os.Stdin to a pipe that writes "\n".
	orig := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	defer func() { os.Stdin = orig }()

	// Per prompt placeholder, we need one newline answer.
	go func() {
		io.WriteString(w, "\n") // accept "overwrite prompt"
		io.WriteString(w, "\n") // accept default
		w.Close()
	}()

	if err := generate.RunInteractive(src, dst, ac, true); err != nil {
		t.Fatalf("generate.RunInteractive: %v", err)
	}

	// The generated file should contain the (default) value
	out := readFile(t, filepath.Join(dst, "app.yaml"))
	if !strings.Contains(out, testvalue) {
		t.Fatalf("TEST FAILED \n value(s) don't match with rendered key value: \n%s, \n GOT INSTEAD:\n%s", testvalue, out)
	}
}
