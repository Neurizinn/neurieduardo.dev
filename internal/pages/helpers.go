package pages

import "strings"

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
func equivalentRoute(route, locale string) string {
	path := strings.TrimPrefix(route, "/en")
	if locale == "en" {
		return "/en" + path
	}
	return path
}
