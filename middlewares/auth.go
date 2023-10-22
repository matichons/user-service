package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ErrNoAuthHeader         = "No Authorization header provided"
	ErrIncorrectTokenFormat = "Incorrect Format of Authorization Token"
)

func extractTokenFromHeader(header string) (string, error) {
	if header == "" {
		return "", errors.New(ErrNoAuthHeader)
	}
	parts := strings.Split(header, "Bearer ")
	if len(parts) != 2 {
		return "", errors.New(ErrIncorrectTokenFormat)
	}
	return strings.TrimSpace(parts[1]), nil
}

func Authz(secretKey, issuer string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractTokenFromHeader(c.Request.Header.Get("Authorization"))
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			c.Abort()
			return
		}

		jwtWrapper := JwtWrapper{
			SecretKey: secretKey,
			Issuer:    issuer,
		}

		claims, err := jwtWrapper.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		c.Set("id", claims.Id)
		c.Next()
	}
}
