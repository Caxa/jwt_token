package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	signingKey = "claims"
	tokenTTL   = 12 * time.Hour
)

type TokenManager interface {
	GenerateToken(userId int) (string, error)
	ParseToken(token string) (int, error)
}

type JWTManager struct {
	secret string
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{secret: secret}
}

func (m *JWTManager) GenerateToken(userId int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(tokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secret))
}

func (m *JWTManager) ParseToken(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(m.secret), nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	idFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user id in token")
	}

	return int(idFloat), nil
}
