package auth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const DefaultUsername = "armin"
const DefaultPassword = "dopadopa123"

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Middleware accepts HTTP Basic (username/password) or Bearer token equal to the password
// so OpenAI SDKs can send Authorization: Bearer <password>.
func Middleware(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" {
			c.Header("WWW-Authenticate", `Basic realm="cursor-cli-to-api"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": map[string]any{
				"message": "Missing Authorization header",
				"type":    "invalid_request_error",
				"code":    "unauthorized",
			}})
			return
		}

		ok := false
		if strings.HasPrefix(strings.ToLower(h), "basic ") {
			raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(h[6:]))
			if err == nil {
				parts := strings.SplitN(string(raw), ":", 2)
				if len(parts) == 2 && parts[0] == username && parts[1] == password {
					ok = true
				}
			}
		} else if strings.HasPrefix(strings.ToLower(h), "bearer ") {
			token := strings.TrimSpace(h[7:])
			if token == password || (token == username+":"+password) {
				ok = true
			}
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": map[string]any{
				"message": "Invalid credentials",
				"type":    "invalid_request_error",
				"code":    "unauthorized",
			}})
			return
		}
		c.Set("username", username)
		c.Next()
	}
}
