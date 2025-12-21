package middleware

import (
	"net/http"
	"strings"

	"github.com/IndalAwalaikal/warung-pos/backend/config"
	"github.com/IndalAwalaikal/warung-pos/backend/model"
	"github.com/IndalAwalaikal/warung-pos/backend/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const ContextUserKey = "currentUser"

// AuthRequired checks Authorization header, validates JWT and attaches user to context
func AuthRequired(userRepo repository.UserRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        // allow token via query param for EventSource (SSE) clients which cannot set headers
        if auth == "" {
            q := c.Query("token")
            if q != "" {
                auth = "Bearer " + q
            }
        }
        if auth == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "missing authorization header"})
            c.Abort()
            return
        }
        parts := strings.SplitN(auth, " ", 2)
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "invalid authorization header"})
            c.Abort()
            return
        }
        tokenStr := parts[1]
        token, err := jwt.ParseWithClaims(tokenStr, &config.Claims{}, func(t *jwt.Token) (interface{}, error) {
            return config.JwtSecret(), nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "invalid token"})
            c.Abort()
            return
        }
        claims, ok := token.Claims.(*config.Claims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "invalid token claims"})
            c.Abort()
            return
        }

        // load user
        u, err := userRepo.FindByID(uint(claims.UserID))
        if err != nil || u == nil {
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "user not found"})
            c.Abort()
            return
        }

        c.Set(ContextUserKey, u)
        c.Next()
    }
}

// AdminRequired middleware ensures current user has role 'admin'
func AdminRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        v, exists := c.Get(ContextUserKey)
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "forbidden"})
            c.Abort()
            return
        }
        u, ok := v.(*model.User)
        if !ok || u.Role != "admin" {
            c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "admin only"})
            c.Abort()
            return
        }
        c.Next()
    }
}


