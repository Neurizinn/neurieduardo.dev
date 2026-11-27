package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"portfolio/internal/pages"
	"portfolio/internal/site"
	"strings"
)

func main() {
	if err := build(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func build() error {
	for _, locale := range []string{"pt", "en"} {
		c, err := site.Load(locale)
		if err != nil {
			return err
		}
		prefix := ""
		if locale == "en" {
			prefix = "/en"
		}
		routes := []string{"home", "about", "experience", "projects", "ai-engineering", "stack", "resume", "contact", "blog"}
		for _, kind := range routes {
			path := prefix + "/"
			if kind != "home" {
				path += kind + "/"
			}
			if err := write(path, c, kind, nil, prefix); err != nil {
				return err
			}
		}
		for i := range c.Projects {
			p := &c.Projects[i]
			if err := write(prefix+"/projects/"+p.Slug+"/", c, "project", p, prefix); err != nil {
				return err
			}
		}
	}
	b, err := os.ReadFile("assets/css/site.css")
	if err != nil {
		return err
	}
	if err = os.MkdirAll("dist/assets", 0755); err != nil {
		return err
	}
	if err = os.WriteFile("dist/assets/site.css", b, 0644); err != nil {
		return err
	}
	b, err = os.ReadFile("assets/js/site.js")
	if err != nil {
		return err
	}
	if err = os.WriteFile("dist/assets/site.js", b, 0644); err != nil {
		return err
	}
	notFound := `<!doctype html><html lang='pt-BR'><meta charset='utf-8'><meta name='viewport' content='width=device-width,initial-scale=1'><title>404 — Neuri Eduardo</title><link rel='stylesheet' href='/assets/site.css'><main class='subpage'><h1>Página não encontrada / Page not found</h1><p>Escolha um idioma / Choose a language</p><a class='btn' href='/'>Português</a> <a class='btn' href='/en/'>English</a></main></html>`
	if err = os.WriteFile("dist/404.html", []byte(notFound), 0644); err != nil {
		return err
	}
	return os.WriteFile("dist/robots.txt", []byte("User-agent: *\nDisallow: /\n"), 0644)
}
func write(route string, c site.Content, kind string, p *site.Project, prefix string) error {
	path := filepath.Join("dist", filepath.FromSlash(strings.TrimPrefix(route, "/")), "index.html")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return pages.Page(c, kind, p, prefix).Render(context.Background(), f)
}
