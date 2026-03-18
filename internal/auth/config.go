package auth

import (
	"log"
	"os"
	"strconv"
	"time"
)

var (
	accessSecret  []byte
	accessTTL     time.Duration
	issuer        string
	refreshSecret []byte
	refreshTTL    time.Duration
)

const defaultIssuer = "blog"

func RefreshTTL() time.Duration {
	return refreshTTL
}

func InitJWT() {
	secret := os.Getenv("JWT_ACCESS_SECRET")
	if secret == "" {
		log.Println("JWT не задан, используем dev key")
		secret = "dev-key-xalupa-kvaska"
	}
	accessSecret = []byte(secret)

	ttlMin := os.Getenv("JWT_ACCESS_TTL")
	if ttlMin == "" {
		log.Println("JWT_ACCESS_TTL не задан, используем 30 минут")
		ttlMin = "30"
	}
	ttl, err := strconv.Atoi(ttlMin)
	if err != nil {
		log.Fatalf("Неверный формат JWT_ACCESS_TTL: %v", err)
	}
	accessTTL = time.Duration(ttl) * time.Minute

	issuer = os.Getenv("JWT_ISSUER")
	if issuer == "" {
		log.Println("JWT_ISSUER не задан, используем blog")
		issuer = defaultIssuer
	}

	refresh := os.Getenv("JWT_REFRESH_SECRET")
	if refresh == "" {
		log.Println("JWT_REFRESH_SECRET не задан, используем dev key (ТОЛЬКО для разработки)")
		refresh = "dev-refresh-secret-...-длинный-ключ"
	}
	refreshSecret = []byte(refresh)

	refreshDaysStr := os.Getenv("JWT_REFRESH_TTL_DAYS")
	if refreshDaysStr == "" {
		refreshDaysStr = "7"
	}
	days, err := strconv.Atoi(refreshDaysStr)
	if err != nil || days <= 0 {
		log.Fatalf("Неверный формат JWT_REFRESH_TTL_DAYS: %v", err)
	}
	refreshTTL = time.Duration(days) * 24 * time.Hour
}
