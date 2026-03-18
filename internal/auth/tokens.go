package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("невалидный токен")
	ErrTokenExpired = errors.New("токен истёк")
)

type AccessClaims struct {
	jwt.RegisteredClaims
	UserID uint `json:"user_id"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	UserID uint `json:"user_id"`
}

func IssueAccessToken(userID uint) (string, time.Duration, error) {
	ttl := accessTTL
	exp := time.Now().Add(ttl)

	claims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			Subject:   strconv.FormatUint(uint64(userID), 10),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(accessSecret)
	if err != nil {
		return "", 0, err
	}

	return signed, ttl, nil
}

func ValidateAccessToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return accessSecret, nil
	})

	if err != nil {
		return 0, ErrInvalidToken
	}

	claims, ok := token.Claims.(*AccessClaims)

	if !ok || !token.Valid {
		return 0, ErrInvalidToken
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return 0, ErrTokenExpired
	}

	return claims.UserID, nil
}

func IssueRefreshToken(userID uint) (string, time.Duration, error) {
	ttl := refreshTTL
	exp := time.Now().Add(ttl)

	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			Subject:   strconv.FormatUint(uint64(userID), 10),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(refreshSecret)
	if err != nil {
		return "", 0, err
	}

	return signed, ttl, nil
}

func ValidateRefreshToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return refreshSecret, nil
	})

	if err != nil {
		return 0, ErrInvalidToken
	}

	claims, ok := token.Claims.(*RefreshClaims)

	if !ok || !token.Valid {
		return 0, ErrInvalidToken
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return 0, ErrTokenExpired
	}

	return claims.UserID, nil
}
