package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ExtractToken(c)
		
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
		})

		// Handle token validation errors
		if err != nil || !token.Valid {
			fmt.Println("JWT error:", err.Error())

			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid or expired JWT",
			})

			c.Abort()
			return
		}

		// If token is valid, allow the request to proceed
		c.Set("user", token)
		c.Next()
	}
}

func ExtractToken(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
