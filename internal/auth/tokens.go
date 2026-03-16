package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("невалидный токен")
)

type AccessClaims struct {
	jwt.RegisteredClaims
	UserID uint `json:"user_id"`
}

func IssueAccessToken(userID uint) (string, time.Duration, error) {
	ttl := accessTTL
	exp := time.Now.Add(ttl)

	claims := AccessClaims{
		RegisteredClaims: jws.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: Issuer,
			Subject: strconv.FormatUint(uint64(userID), 10),
		}
		UserID: userID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secretKey)
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
		return accessSecretKey, nil
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