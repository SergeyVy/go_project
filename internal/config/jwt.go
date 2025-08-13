package config

import (
	"os"
	"strings"
	"time"
)

var (
	JWTSecret  = []byte(getEnv("JWT_SECRET", "CHANGE_ME"))
	AccessTTL  = 15 * time.Minute
	RefreshTTL = 7 * 24 * time.Hour
)

func getEnv(k, def string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return def
}
