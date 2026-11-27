package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPackagingRemovesStaleFiles(t *testing.T) {
	root := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string]string{"dist/index.html": "ok", ".vercel/output/static/stale.pdf": "private"} {
		if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(name, []byte(data), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := run(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(".vercel/output/static/stale.pdf"); !os.IsNotExist(err) {
		t.Fatalf("stale file remains: %v", err)
	}
	if _, err := os.Stat(".vercel/output/static/index.html"); err != nil {
		t.Fatal(err)
	}
}
