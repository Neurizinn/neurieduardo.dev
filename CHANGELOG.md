# Changelog

## Unreleased

- Added a static Go and Templ build with Portuguese and English routes, project pages, local styles, theme preference, and a bilingual 404 page.
- Added typed JSON content with slug validation and basic Go checks.
- Added a pinned local Tailwind build alongside the site's custom CSS.
- Added validated `SITE_URL` handling for canonical, hreflang, sitemap, and local noindex output.
- Added CI checks for code generation, Go tests, vet, and the static asset build.
- Added an internal link and anchor check for generated HTML.
- Added a keyboard accessible navigation search dialog.
- Added validated bilingual Markdown articles, sanitized rendering, and RSS generation when articles are published.
- Added opt-in public resume downloads gated by approved PDF files.
- Added local Geist and JetBrains Mono assets with their license files in the generated output.
- Added Chromium screenshots and navigation checks for four viewport widths.
- Improved keyboard labels, touch targets, and text contrast after local accessibility audit.
- Refined the stack and AI sections against the visual reference and withheld academic details pending validation.
- Added page-specific metadata and cleared stale Vercel output before packaging.
- Made the Chromium setup in CI retry transient installation failures.
- Kept publication and production SEO disabled while content and deployment remain under review.
