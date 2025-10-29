package middleware

import (
	"embed"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

const I18nContextKey = "i18n"

//go:embed locales/*.toml
var localeFS embed.FS

func I18n(defaultLang language.Tag) echo.MiddlewareFunc {
	bundle := i18n.NewBundle(defaultLang)
	// Register the unmarshaler for TOML files.
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFileFS(localeFS, "locales/en.toml") // This path is relative to the embed root
	bundle.LoadMessageFileFS(localeFS, "locales/vi.toml")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Determine language preference in order of priority:
			// 1. Query parameter `?lang=...`
			// 2. Cookie `lang=...`
			// 3. `Accept-Language` header
			langs := []string{}

			// Check query parameter first. If present, set the cookie.
			if qLang := c.QueryParam("lang"); qLang != "" {
				langs = append(langs, qLang)
				cookie := &http.Cookie{
					Name:    "lang",
					Value:   qLang,
					Path:    "/",
					Expires: time.Now().Add(365 * 24 * time.Hour), // 1 year
				}
				c.SetCookie(cookie)
			}

			// Check for the cookie if query param was not used.
			if cLang, err := c.Cookie("lang"); err == nil {
				langs = append(langs, cLang.Value)
			}

			localizer := i18n.NewLocalizer(bundle, langs...)

			// Store the localizer in the context
			c.Set(I18nContextKey, localizer)

			return next(c)
		}
	}
}