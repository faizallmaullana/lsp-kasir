package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"faizalmaulana/lsp/conf"
	"faizalmaulana/lsp/helper"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(cfg *conf.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
			c.Abort()
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
			c.Abort()
			return
		}

		tokenStr := parts[1]
		fmt.Println("DEBUG: Authorization token:", tokenStr)

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			fmt.Println("DEBUG: token header alg:", t.Header["alg"])
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			if cfg.JWTSecret == "" {
				return nil, fmt.Errorf("jwt secret not configured")
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil {
			fmt.Println("DEBUG: parse error:", err)
			c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
			c.Abort()
			return
		}

		fmt.Println("DEBUG: parsed token.Valid =", token.Valid)
		fmt.Println("DEBUG: parsed claims =", claims)

		if !token.Valid {
			fmt.Println("DEBUG: token is not valid")
			c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
			c.Abort()
			return
		}

		fmt.Println("JWT Claims:", token.Claims)

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("claims", claims)
		}

		fmt.Println("JWT Middleware passed")

		c.Next()
	}
}
