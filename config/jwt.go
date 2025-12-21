package config

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWT settings getters
func JwtSecret() []byte {
    return []byte(GetEnv("JWT_SECRET", "secret"))
}

func JwtExpiry() time.Duration {
    // minutes
    return time.Duration(60) * time.Minute
}

// Claims helper (can be extended)
type Claims struct {
    UserID uint `json:"user_id"`
    jwt.RegisteredClaims
}
