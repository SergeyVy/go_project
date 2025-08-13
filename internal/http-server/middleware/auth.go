package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"url-shorter/internal/config"
)

type ctxKey string

const userIDKey ctxKey = "uid"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "missing bearer", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimSpace(parts[1])

		tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, jwt.ErrSignatureInvalid
			}
			return config.JWTSecret, nil
		})
		if err != nil || !tok.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid claims", http.StatusUnauthorized)
			return
		}
		fuid, ok := claims["uid"].(float64) // число из JSON приходит float64
		if !ok {
			http.Error(w, "no uid", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, int64(fuid))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserID(r *http.Request) (int64, bool) {
	v := r.Context().Value(userIDKey)
	if v == nil {
		return 0, false
	}
	id, ok := v.(int64)
	return id, ok
}
