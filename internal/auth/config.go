package auth

import (
	"log"
	"os"
	"strconv"
	"time"
)

var (
	accessSecret []byte
	accessTTL    time.Duration
	issuer       string
)

const defaultAccessTTL = 30 * time.Minute
const defaultIssuer = "blog"

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
}