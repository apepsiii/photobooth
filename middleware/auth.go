package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"mastergin/core"
)

// RequireAuth is a Gin middleware that checks for a valid JWT in the
// "auth_token" HTTP-only cookie. If the token is missing or invalid, the
// request is redirected to /login. On success, the user ID and name are
// stored in the Gin context for downstream handlers.
func RequireAuth(cfg *core.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("auth_token")
		if err != nil || tokenString == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Store user info in context for handlers and templates.
		if uid, ok := claims["uid"].(float64); ok {
			c.Set("userID", uint(uid))
		}
		if name, ok := claims["name"].(string); ok {
			c.Set("userName", name)
		}
		if role, ok := claims["role"].(string); ok {
			c.Set("userRole", role)
		}

		c.Next()
	}
}

// CurrentUserName extracts the authenticated user's name from the Gin context.
// Returns an empty string if not set.
func CurrentUserName(c *gin.Context) string {
	if name, ok := c.Get("userName"); ok {
		if s, ok := name.(string); ok {
			return s
		}
	}
	return ""
}
