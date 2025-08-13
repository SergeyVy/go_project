package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"url-shorter/internal/config"
)

type Claims struct {
	UserID int64 `json:"uid"`
	jwt.RegisteredClaims
}

func NewAccess(userID int64) (string, error) { return sign(userID, "access", config.AccessTTL) }

func NewRefresh(userID int64) (string, error) { return sign(userID, "refresh", config.RefreshTTL) }

func sign(uid int64, sub string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(config.JWTSecret)
}
