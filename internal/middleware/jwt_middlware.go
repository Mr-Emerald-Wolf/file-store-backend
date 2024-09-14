package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/services"
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

			c.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Invalid or expired JWT",
			})

			c.Abort()
			return
		}


		email, err := token.Claims.GetSubject()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to parse email from claims",
			})
			c.Abort()
			return
		}

		// Find User
		user, err := services.FindUserByEmail(email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not find user",
			})
			c.Abort()
			return
		}
		
		c.Set("user", user)
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
