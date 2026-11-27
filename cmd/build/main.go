package main

import (
	"context"
	"encoding/xml"
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
	base, err := site.ValidateBase(os.Getenv("SITE_URL"))
	if err != nil {
		return err
	}
	if err := os.RemoveAll("dist"); err != nil {
		return fmt.Errorf("clear generated output: %w", err)
	}
	var allRoutes []string
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
			if err := write(path, c, kind, nil, prefix, base); err != nil {
				return err
			}
			allRoutes = append(allRoutes, path)
		}
		for i := range c.Projects {
			p := &c.Projects[i]
			path := prefix + "/projects/" + p.Slug + "/"
			if err := write(path, c, "project", p, prefix, base); err != nil {
				return err
			}
			allRoutes = append(allRoutes, path)
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
	b, err = os.ReadFile("public/favicon.svg")
	if err != nil {
		return err
	}
	if err = os.WriteFile("dist/favicon.svg", b, 0644); err != nil {
		return err
	}
	notFound := `<!doctype html><html lang='pt-BR'><meta charset='utf-8'><meta name='viewport' content='width=device-width,initial-scale=1'><meta name='robots' content='noindex'><title>404 — Neuri Eduardo</title><link rel='stylesheet' href='/assets/site.css'><main class='subpage'><h1>Página não encontrada / Page not found</h1><p>Escolha um idioma / Choose a language</p><a class='btn' href='/'>Português</a> <a class='btn' href='/en/'>English</a></main></html>`
	if err = os.WriteFile("dist/404.html", []byte(notFound), 0644); err != nil {
		return err
	}
	robots := "User-agent: *\nDisallow: /\n"
	if base != "" {
		robots = "User-agent: *\nAllow: /\nSitemap: " + base + "/sitemap.xml\n"
		type loc struct {
			XMLName xml.Name `xml:"url"`
			Loc     string   `xml:"loc"`
		}
		type sitemap struct {
			XMLName xml.Name `xml:"urlset"`
			Xmlns   string   `xml:"xmlns,attr"`
			URLs    []loc    `xml:"url"`
		}
		s := sitemap{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}
		for _, route := range allRoutes {
			s.URLs = append(s.URLs, loc{Loc: base + route})
		}
		b, err := xml.MarshalIndent(s, "", "  ")
		if err != nil {
			return err
		}
		if err = os.WriteFile("dist/sitemap.xml", append([]byte(xml.Header), b...), 0644); err != nil {
			return err
		}
	}
	return os.WriteFile("dist/robots.txt", []byte(robots), 0644)
}
func write(route string, c site.Content, kind string, p *site.Project, prefix string, base string) error {
	path := filepath.Join("dist", filepath.FromSlash(strings.TrimPrefix(route, "/")), "index.html")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return pages.Page(c, kind, p, prefix, base, route).Render(context.Background(), f)
}
