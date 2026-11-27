package site

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
)

type Article struct {
	Slug  string
	Title string
	Date  string
	HTML  string
}

func LoadArticles(locale string) ([]Article, error) {
	dir := filepath.Join("content", locale, "blog")
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read blog: %w", err)
	}
	var articles []Article
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		article, published, err := ParseArticle(entry.Name(), data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", entry.Name(), err)
		}
		if !published {
			continue
		}
		if seen[article.Slug] {
			return nil, fmt.Errorf("duplicate article slug %q", article.Slug)
		}
		seen[article.Slug] = true
		articles = append(articles, article)
	}
	return articles, nil
}

func ParseArticle(filename string, data []byte) (Article, bool, error) {
	source := strings.ReplaceAll(string(data), "\r\n", "\n")
	if !strings.HasPrefix(source, "---\n") {
		return Article{}, false, fmt.Errorf("missing front matter")
	}
	parts := strings.SplitN(source[4:], "\n---\n", 2)
	if len(parts) != 2 {
		return Article{}, false, fmt.Errorf("unterminated front matter")
	}
	meta := map[string]string{}
	for _, line := range strings.Split(parts[0], "\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			return Article{}, false, fmt.Errorf("invalid front matter line")
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		switch key {
		case "slug", "title", "date", "published":
		default:
			return Article{}, false, fmt.Errorf("unknown front matter key %q", key)
		}
		if meta[key] != "" || value == "" {
			return Article{}, false, fmt.Errorf("duplicate or empty %s", key)
		}
		meta[key] = value
	}
	if !ValidSlug(meta["slug"]) || strings.TrimSuffix(filename, ".md") != meta["slug"] {
		return Article{}, false, fmt.Errorf("invalid article slug")
	}
	if meta["title"] == "" {
		return Article{}, false, fmt.Errorf("missing title")
	}
	if _, err := time.Parse("2006-01-02", meta["date"]); err != nil {
		return Article{}, false, fmt.Errorf("invalid date: %w", err)
	}
	if meta["published"] != "true" && meta["published"] != "false" {
		return Article{}, false, fmt.Errorf("published must be true or false")
	}
	if meta["published"] == "false" {
		return Article{}, false, nil
	}
	var rendered bytes.Buffer
	if err := goldmark.Convert([]byte(parts[1]), &rendered); err != nil {
		return Article{}, false, err
	}
	safe := bluemonday.UGCPolicy().SanitizeBytes(rendered.Bytes())
	return Article{Slug: meta["slug"], Title: meta["title"], Date: meta["date"], HTML: string(safe)}, true, nil
}

func ValidateArticleTranslations(pt, en []Article) error {
	dates := map[string]string{}
	for _, a := range pt {
		dates[a.Slug] = a.Date
	}
	if len(pt) != len(en) {
		return fmt.Errorf("published article translations differ")
	}
	for _, a := range en {
		if dates[a.Slug] != a.Date {
			return fmt.Errorf("missing or inconsistent translation for %s", a.Slug)
		}
	}
	return nil
}
