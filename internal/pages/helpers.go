package pages

import "portfolio/internal/site"

func nav(locale, en, pt string) string {
	if locale == "en" {
		return en
	}
	return pt
}
func first(s string) string {
	for _, r := range s {
		return string(r)
	}
	return "?"
}
func otherLocale(locale string) string {
	if locale == "en" {
		return "pt"
	}
	return "en"
}
func htmlLocale(locale string) string {
	if locale == "en" {
		return "en-US"
	}
	return "pt-BR"
}
func localeRoute(locale, kind string, p *site.Project) string {
	prefix := ""
	if locale == "en" {
		prefix = "/en"
	}
	if p != nil {
		return prefix + "/projects/" + p.Slug + "/"
	}
	if kind == "home" {
		return prefix + "/"
	}
	return prefix + "/" + kind + "/"
}
func localeLink(locale, kind string, p *site.Project) string {
	prefix := "/en"
	if locale == "en" {
		prefix = ""
	}
	if p != nil {
		return prefix + "/projects/" + p.Slug + "/"
	}
	if kind == "home" {
		return prefix + "/"
	}
	return prefix + "/" + kind + "/"
}
