package middleware

import (
	"net/http"
	"sync"
	"time"

	"faizalmaulana/lsp/helper"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var clients sync.Map

func getLimiter(key string) *rate.Limiter {
	if v, ok := clients.Load(key); ok {
		cl := v.(*clientLimiter)
		cl.lastSeen = time.Now()
		return cl.limiter
	}

	limiter := rate.NewLimiter(rate.Every(time.Minute/5), 5)
	clients.Store(key, &clientLimiter{limiter: limiter, lastSeen: time.Now()})
	return limiter
}

func init() {
	go func() {
		for {
			time.Sleep(time.Minute)
			cutoff := time.Now().Add(-5 * time.Minute)
			clients.Range(func(key, value interface{}) bool {
				cl := value.(*clientLimiter)
				if cl.lastSeen.Before(cutoff) {
					clients.Delete(key)
				}
				return true
			})
		}
	}()
}

func LoginRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, helper.ErrorResponse("TOO_MANY_REQUESTS", "rate limit exceeded"))
			c.Abort()
			return
		}
		c.Next()
	}
}
