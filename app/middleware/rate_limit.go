package middleware
import (
    "net/http"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

type visitor struct {
    limiter  *rate.Limiter
    lastSeen time.Time
}

var (
    mu       sync.Mutex
    visitors = make(map[string]*visitor)
)

func init() {
    go cleanupVisitors()
}

func getVisitor(ip string) *rate.Limiter {
    mu.Lock()
    defer mu.Unlock()

    v, exists := visitors[ip]
    if !exists {
        l := rate.NewLimiter(2, 3) // 2 req/detik, burst max 3
        visitors[ip] = &visitor{l, time.Now()}
        return l
    }
    v.lastSeen = time.Now()
    return v.limiter
}

func cleanupVisitors() {
    for {
        time.Sleep(time.Minute)
        mu.Lock()
        for ip, v := range visitors {
            if time.Since(v.lastSeen) > 3*time.Minute {
                delete(visitors, ip)
            }
        }
        mu.Unlock()
    }
}

func RateLimit() gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        if !getVisitor(ip).Allow() {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "error": "Too many requests",
            })
            return
        }
        c.Next()
    }
}