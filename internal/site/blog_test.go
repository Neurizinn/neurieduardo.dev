package site

import (
	"strings"
	"testing"
)

func TestArticleSafety(t *testing.T) {
	data := []byte("---\nslug: safe-post\ntitle: Safe post\ndate: 2026-11-26\npublished: true\n---\n[bad](javascript:alert(1))\n<script>alert(1)</script>")
	article, published, err := ParseArticle("safe-post.md", data)
	if err != nil || !published {
		t.Fatalf("parse: %v", err)
	}
	if strings.Contains(article.HTML, "javascript:") || strings.Contains(article.HTML, "<script>") {
		t.Fatalf("unsafe HTML: %s", article.HTML)
	}
}
func TestArticleMetadata(t *testing.T) {
	data := []byte("---\nslug: ../bad\ntitle: Bad\ndate: 2026-11-26\npublished: true\n---\ntext")
	if _, _, err := ParseArticle("bad.md", data); err == nil {
		t.Fatal("unsafe slug accepted")
	}
}
func TestArticleTranslations(t *testing.T) {
	pt := []Article{{Slug: "post", Date: "2026-11-26"}}
	en := []Article{{Slug: "post", Date: "2026-11-25"}}
	if err := ValidateArticleTranslations(pt, en); err == nil {
		t.Fatal("mismatched dates accepted")
	}
	en[0].Date = pt[0].Date
	if err := ValidateArticleTranslations(pt, en); err != nil {
		t.Fatal(err)
	}
}
