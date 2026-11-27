package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run() error {
	dst := ".vercel/output/static"
	root, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	output, err := filepath.Abs(".vercel/output")
	if err != nil {
		return err
	}
	relOutput, err := filepath.Rel(root, output)
	if err != nil {
		return err
	}
	if relOutput == ".." || strings.HasPrefix(relOutput, ".."+string(filepath.Separator)) {
		return fmt.Errorf("output path escapes repository")
	}
	for _, candidate := range []string{".vercel", ".vercel/output"} {
		resolved, err := filepath.EvalSymlinks(candidate)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		resolved, err = filepath.Abs(resolved)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, resolved)
		if err != nil {
			return err
		}
		if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return fmt.Errorf("generated output symlink escapes repository")
		}
	}
	if err := os.RemoveAll(output); err != nil {
		return fmt.Errorf("clear generated Vercel output: %w", err)
	}
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	err = filepath.WalkDir("dist", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel("dist", path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		out, err := os.Create(target)
		if err != nil {
			in.Close()
			return err
		}
		_, copyErr := io.Copy(out, in)
		closeOutErr := out.Close()
		closeInErr := in.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeOutErr != nil {
			return closeOutErr
		}
		return closeInErr
	})
	if err != nil {
		return err
	}
	return os.WriteFile(".vercel/output/config.json", []byte(`{"version":3}`), 0644)
}
