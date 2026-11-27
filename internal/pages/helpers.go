package pages

import (
	"portfolio/internal/site"
	"strings"
)

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
func navHref(kind, prefix, target string) string {
	if kind == "home" {
		return prefix + "/#" + target + "-n"
	}
	return prefix + "/" + target + "/"
}
func navCurrent(kind, target string) string {
	if kind == target || (kind == "project" && target == "projects") {
		return "page"
	}
	return "false"
}
func pageTitle(c site.Content, kind string, project *site.Project) string {
	name := "Portfolio"
	if project != nil {
		name = project.Name
	} else if c.CurrentArticle != nil {
		name = c.CurrentArticle.Title
	} else {
		switch kind {
		case "about":
			name = nav(c.Locale, "About", "Sobre")
		case "experience":
			name = nav(c.Locale, "Experience", "Experiência")
		case "projects":
			name = nav(c.Locale, "Projects", "Projetos")
		case "ai-engineering":
			name = "AI Engineering"
		case "stack":
			name = "Stack"
		case "resume":
			name = nav(c.Locale, "Resume", "Currículo")
		case "contact":
			name = nav(c.Locale, "Contact", "Contato")
		case "blog":
			name = "Blog"
		}
	}
	return name + " | Neuri Eduardo de Lima"
}
func pageDescription(c site.Content, kind string, project *site.Project) string {
	if project != nil {
		return project.Description
	}
	switch kind {
	case "about":
		return c.About
	case "ai-engineering":
		return c.AI
	case "contact":
		return c.Contact
	default:
		return c.Intro
	}
}
