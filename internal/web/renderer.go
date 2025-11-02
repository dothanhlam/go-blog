package web

import (
	"embed"
	"html/template"
	"io"

	"go-blog/internal/middleware"

	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

//go:embed all:template
var templatesFS embed.FS

// TemplateRenderer is a custom html/template renderer for the Echo framework.
type TemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer creates a new TemplateRenderer.
func NewTemplateRenderer() *TemplateRenderer {
	// Define a function map to be used in the templates.
	funcMap := template.FuncMap{
		"t": func(c echo.Context, messageID string) string {
			// The default value is the message ID itself
			defaultValue := messageID

			if c == nil {
				return defaultValue
			}

			// Get the localizer from the context
			localizer, ok := c.Get(middleware.I18nContextKey).(*i18n.Localizer)
			if !ok {
				return defaultValue
			}

			// Localize the message
			translated, err := localizer.Localize(&i18n.LocalizeConfig{
				MessageID: messageID,
			})
			if err != nil {
				return defaultValue // Fallback to message ID on error
			}
			return translated
		},
	}

	return &TemplateRenderer{
		// Parse all templates (pages and layouts) into a single template set.
		templates: template.Must(template.New("").Funcs(funcMap).ParseFS(templatesFS,
			"template/*.html", "template/layouts/*.html")),
	}
}

// Render renders a template document.
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}