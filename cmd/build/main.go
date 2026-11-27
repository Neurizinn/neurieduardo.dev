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
	"time"
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
	contents := map[string]site.Content{}
	resumes := map[string][]byte{}
	publishResumes := os.Getenv("PUBLISH_RESUMES")
	if publishResumes != "" && publishResumes != "true" {
		return fmt.Errorf("PUBLISH_RESUMES must be true or unset")
	}
	for _, locale := range []string{"pt", "en"} {
		contents[locale], err = site.Load(locale)
		if err != nil {
			return err
		}
		if publishResumes == "true" {
			name := "resume-approved-" + locale + ".pdf"
			data, err := os.ReadFile(filepath.Join("public", name))
			if err != nil {
				return fmt.Errorf("approved resume %s: %w", locale, err)
			}
			if len(data) < 5 || string(data[:5]) != "%PDF-" || len(data) > 10<<20 {
				return fmt.Errorf("approved resume %s is not a PDF under 10 MB", locale)
			}
			resumes[locale] = data
			c := contents[locale]
			c.ResumeURL = "/" + name
			contents[locale] = c
		}
	}
	if err := site.ValidateArticleTranslations(contents["pt"].Articles, contents["en"].Articles); err != nil {
		return err
	}
	if err := os.RemoveAll("dist"); err != nil {
		return fmt.Errorf("clear generated output: %w", err)
	}
	var allRoutes []string
	for _, locale := range []string{"pt", "en"} {
		c := contents[locale]
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
		for i := range c.Articles {
			c.CurrentArticle = &c.Articles[i]
			path := prefix + "/blog/" + c.CurrentArticle.Slug + "/"
			if err := write(path, c, "article", nil, prefix, base); err != nil {
				return err
			}
			allRoutes = append(allRoutes, path)
		}
		if base != "" && len(c.Articles) > 0 {
			if err := writeRSS(base, prefix, c); err != nil {
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
	b, err = os.ReadFile("public/favicon.svg")
	if err != nil {
		return err
	}
	if err = os.WriteFile("dist/favicon.svg", b, 0644); err != nil {
		return err
	}
	for locale, data := range resumes {
		if err := os.WriteFile(filepath.Join("dist", "resume-approved-"+locale+".pdf"), data, 0644); err != nil {
			return err
		}
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
func writeRSS(base, prefix string, c site.Content) error {
	type item struct {
		Title   string `xml:"title"`
		Link    string `xml:"link"`
		PubDate string `xml:"pubDate"`
	}
	type channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Items       []item `xml:"item"`
	}
	type rss struct {
		XMLName xml.Name `xml:"rss"`
		Version string   `xml:"version,attr"`
		Channel channel  `xml:"channel"`
	}
	feed := rss{Version: "2.0", Channel: channel{Title: "Neuri Eduardo de Lima — Blog", Link: base + prefix + "/blog/", Description: c.Intro}}
	for _, a := range c.Articles {
		date, _ := time.Parse("2006-01-02", a.Date)
		feed.Channel.Items = append(feed.Channel.Items, item{Title: a.Title, Link: base + prefix + "/blog/" + a.Slug + "/", PubDate: date.Format(time.RFC1123Z)})
	}
	b, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join("dist", strings.TrimPrefix(prefix, "/"), "feed.xml"), append([]byte(xml.Header), b...), 0644)
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
