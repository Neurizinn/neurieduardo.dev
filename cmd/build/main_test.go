package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildWithPublishedArticle(t *testing.T) {
	root := t.TempDir()
	source, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repo := filepath.Clean(filepath.Join(source, "../.."))
	for _, name := range []string{"content/pt/site.json", "content/en/site.json", "assets/css/site.css", "assets/css/fonts.css", "assets/js/site.js", "public/favicon.svg"} {
		data, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(name)))
		if err != nil {
			t.Fatal(err)
		}
		dst := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(dst, data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	for _, locale := range []string{"pt", "en"} {
		dst := filepath.Join(root, "content", locale, "blog", "sample.md")
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			t.Fatal(err)
		}
		markdown := "---\nslug: sample\ntitle: Sample\ndate: 2026-11-26\npublished: true\n---\n[bad](javascript:alert(1))"
		if err := os.WriteFile(dst, []byte(markdown), 0644); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() { _ = os.Chdir(source) })
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SITE_URL", "https://example.org")
	if err := build(); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"dist/blog/sample/index.html", "dist/en/blog/sample/index.html", "dist/feed.xml", "dist/en/feed.xml", "dist/sitemap.xml"} {
		if _, err := os.Stat(name); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
	html, err := os.ReadFile("dist/blog/sample/index.html")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(html), "javascript:") || !strings.Contains(string(html), "rel=\"canonical\"") {
		t.Fatal("article HTML is unsafe or SEO metadata is missing")
	}
	resume, err := os.ReadFile("dist/resume/index.html")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(resume), "resume-approved-pt.pdf") {
		t.Fatal("resume link shown without approval")
	}
	t.Setenv("PUBLISH_RESUMES", "true")
	if err := build(); err == nil {
		t.Fatal("missing approved PDFs accepted")
	}
}
