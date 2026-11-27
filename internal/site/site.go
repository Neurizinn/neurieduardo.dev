package site

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Project struct {
	Slug        string   `json:"slug"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Detail      string   `json:"detail"`
	Tags        []string `json:"tags"`
	Diagram     []string `json:"diagram"`
}
type Entry struct {
	Date         string `json:"date"`
	Title        string `json:"title"`
	Organization string `json:"organization"`
	Description  string `json:"description"`
}
type Content struct {
	Locale         string    `json:"locale"`
	Intro          string    `json:"intro"`
	About          string    `json:"about"`
	Projects       []Project `json:"projects"`
	Experience     []Entry   `json:"experience"`
	Stack          []Entry   `json:"stack"`
	AI             string    `json:"ai"`
	AITags         []string  `json:"ai_tags"`
	AITools        []string  `json:"ai_tools"`
	Education      string    `json:"education"`
	Contact        string    `json:"contact"`
	Articles       []Article `json:"-"`
	CurrentArticle *Article  `json:"-"`
	ResumeURL      string    `json:"-"`
}

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func ValidSlug(s string) bool { return slugPattern.MatchString(s) }
func ValidateBase(raw string) (string, error) {
	if raw == "" {
		return "", nil
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
		return "", fmt.Errorf("SITE_URL must be an HTTPS origin")
	}
	return strings.TrimSuffix(u.String(), "/"), nil
}
func Load(locale string) (Content, error) {
	if locale != "pt" && locale != "en" {
		return Content{}, fmt.Errorf("invalid locale %q", locale)
	}
	b, err := os.ReadFile(filepath.Join("content", locale, "site.json"))
	if err != nil {
		return Content{}, fmt.Errorf("load %s: %w", locale, err)
	}
	var c Content
	if err = json.Unmarshal(b, &c); err != nil {
		return c, fmt.Errorf("parse %s: %w", locale, err)
	}
	if c.Locale != locale {
		return c, fmt.Errorf("locale mismatch: %s", c.Locale)
	}
	seen := map[string]bool{}
	for _, p := range c.Projects {
		if !ValidSlug(p.Slug) || seen[p.Slug] {
			return c, fmt.Errorf("invalid or duplicate project slug %q", p.Slug)
		}
		seen[p.Slug] = true
	}
	c.Articles, err = LoadArticles(locale)
	if err != nil {
		return c, err
	}
	return c, nil
}
