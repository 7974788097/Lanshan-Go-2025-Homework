package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IpRateLimiter struct {
	ips      map[string]*rate.Limiter
	mutex    *sync.RWMutex
	r        rate.Limit
	capacity int
}

func NewIpRateLimiter(r rate.Limit, capacity int) *IpRateLimiter {
	return &IpRateLimiter{
		ips:      make(map[string]*rate.Limiter),
		mutex:    &sync.RWMutex{},
		r:        r,
		capacity: capacity,
	}
}
func (i *IpRateLimiter) getLimiter(ip string) *rate.Limiter {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	limiter, ok := i.ips[ip]
	if !ok {
		limiter = rate.NewLimiter(i.r, i.capacity)
		i.ips[ip] = limiter
	}
	return limiter
}
func Limiter(limiter *IpRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.getLimiter(ip)
		if !l.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
