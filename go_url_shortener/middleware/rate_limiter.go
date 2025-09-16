package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var ctx = context.Background()

// RateLimiterMiddleware aplica limite de requisições por IP
func RateLimiterMiddleware(rdb *redis.Client, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extrai IP do cliente
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = strings.Split(r.RemoteAddr, ":")[0]
			}

			key := "rate_limit:" + ip
			count, err := rdb.Incr(ctx, key).Result()
			if err != nil {
				http.Error(w, "internal rate limiter error", http.StatusInternalServerError)
				return
			}

			if count == 1 {
				rdb.Expire(ctx, key, window)
			}

			if int(count) > limit {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
