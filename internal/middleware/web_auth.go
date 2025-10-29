package middleware

import (
	"go-blog/internal/config"
	"go-blog/internal/service"
	"log" // Added log import
	"net/http"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const UserContextKey = "user"

// WebAuth populates the context with user information if a valid JWT cookie is found.
func WebAuth(userService service.UserService, cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("token")
			if err != nil {
				// If there's any error getting the cookie (including http.ErrNoCookie),
				// it means we don't have a valid token cookie to process.
				log.Printf("WebAuth: Error getting 'token' cookie: %v\n", err) // Debug
				return next(c)
			}

			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil {
			//	log.Printf("WebAuth: JWT parse error: %v\n", err) // Debug
				// Clear invalid cookie to prevent repeated errors
				c.SetCookie(&http.Cookie{Name: "token", Value: "", Expires: time.Unix(0, 0), Path: "/", HttpOnly: true, Secure: false})
				return next(c) // Invalid token, proceed without user
			}

			if token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					userID := int(claims["id"].(float64))
					user, err := userService.GetByID(userID)
					if err == nil { // User found
					//	log.Printf("WebAuth: User %s (ID: %d) found and set in context.\n", user.Username, user.ID) // Debug
					log.Printf("WebAuth: User: %d", userID) // Debug
					c.Set(UserContextKey, user) // Add user to context
					} else { // User not found in DB for valid token
					//	log.Printf("WebAuth: Error getting user by ID %d: %v\n", userID, err) // Debug
						// Clear cookie if user not found (e.g., user deleted)
						c.SetCookie(&http.Cookie{Name: "token", Value: "", Expires: time.Unix(0, 0), Path: "/", HttpOnly: true, Secure: false})
					}
				}
			}
			return next(c)
		}
	}
}