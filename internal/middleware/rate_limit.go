package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// 100 Requests every 1 minute
func RateLimiter() gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(600*time.Millisecond), 100)
	return func(c *gin.Context) {

		if limiter.Allow() {
			c.Next()
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "Limit exceeded",
			})
			c.Abort()
		}

	}
}
