package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// --- In-memory rate limiter (fallback) ---

type tokenBucket struct {
	tokens     int
	lastRefill time.Time
}

type inMemoryRateLimiter struct {
	mu       sync.Mutex
	requests map[string]*tokenBucket
	limit    int
	window   time.Duration
}

func newInMemoryRateLimiter(limit int, window time.Duration) *inMemoryRateLimiter {
	return &inMemoryRateLimiter{
		requests: make(map[string]*tokenBucket),
		limit:    limit,
		window:   window,
	}
}

func (rl *inMemoryRateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	bucket, exists := rl.requests[ip]
	if !exists || now.Sub(bucket.lastRefill) > rl.window {
		rl.requests[ip] = &tokenBucket{tokens: rl.limit - 1, lastRefill: now}
		return true
	}
	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}
	return false
}

// --- Rate limiting middleware ---

// RateLimit creates a rate limiter middleware using in-memory storage.
// Signature kept for backward compatibility with existing tests.
func RateLimit(limit int, duration time.Duration) gin.HandlerFunc {
	return RateLimitWithRedis(nil, limit, duration)
}

// RateLimitWithRedis creates a rate limiter middleware backed by Redis
// with automatic fallback to in-memory when rdb is nil.
func RateLimitWithRedis(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	var limiter *inMemoryRateLimiter
	if rdb == nil {
		limiter = newInMemoryRateLimiter(limit, window)
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:%s", ip)

		if rdb != nil {
			// Redis-backed rate limiting (sliding window via counter)
			ctx := c.Request.Context()
			val, err := rdb.Get(ctx, key).Int()
			if err == nil && val >= limit {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error":     "rate limit exceeded",
					"retry_after": int(window.Seconds()),
				})
				c.Abort()
				return
			}
			pipe := rdb.Pipeline()
			pipe.Incr(ctx, key)
			pipe.Expire(ctx, key, window)
			_, _ = pipe.Exec(ctx)
		} else {
			// In-memory fallback
			if !limiter.allow(ip) {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error":      "rate limit exceeded",
					"retry_after": int(window.Seconds()),
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// --- JWT authentication middleware ---

// Authenticate validates JWT tokens from the Authorization header.
// The `secret` parameter should be a string containing the JWT secret (HMAC).
// If secret is nil or empty, all tokens are rejected with 401 (secure default).
func Authenticate(secret interface{}) gin.HandlerFunc {
	var jwtSecret string
	if s, ok := secret.(string); ok {
		jwtSecret = s
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if jwtSecret == "" {
			log.Println("WARNING: JWT secret not configured, rejecting all tokens")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication not configured"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract claims and set them on the context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		if sid, ok := claims["sid"].(string); ok {
			c.Set("sid", sid)
		}
		if email, ok := claims["email"].(string); ok {
			c.Set("email", email)
		}
		if groups, ok := claims["groups"]; ok {
			c.Set("groups", groups)
		}
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireEntitlement checks that the authenticated user has the required
// group (e.g. "app_games") in their JWT claims.
// If embedHandoffSecret is empty, all authenticated users are allowed (dev mode).
// This middleware should be placed after the Authenticate middleware.
func RequireEntitlement(embedHandoffSecret string, requiredGroup string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Dev mode: if no embed secret configured, allow all
		if embedHandoffSecret == "" {
			c.Next()
			return
		}

		// Check groups from JWT claims (set by Authenticate middleware)
		groupsVal, exists := c.Get("groups")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "no groups found in token"})
			c.Abort()
			return
		}

		// groups may be []interface{} from JSON deserialization
		groups, ok := groupsVal.([]interface{})
		if !ok {
			// If groups is not a slice, try string type (single group)
			if gs, ok := groupsVal.(string); ok {
				if gs == requiredGroup {
					c.Next()
					return
				}
				c.JSON(http.StatusForbidden, gin.H{"error": "insufficient entitlement"})
				c.Abort()
				return
			}
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid groups format in token"})
			c.Abort()
			return
		}

		for _, g := range groups {
			if gs, ok := g.(string); ok && gs == requiredGroup {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient entitlement"})
		c.Abort()
	}
}
