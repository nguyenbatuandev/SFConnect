package service

import (
	"Order_Service/internal/entity"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

type JWTauthService struct {
	secretKey []byte
}

func NewJWTauthService(secretKey string) *JWTauthService {
	return &JWTauthService{
		secretKey: []byte(secretKey),
	}
}

func (s *JWTauthService) ValidateToken(tokenString string) (*entity.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("invalid exp in token")
	}
	if int64(exp) < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid UUID format")
	}

	roleStr, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid role in token")
	}

	return &entity.Claims{
		UserID: userID,
		Role:   entity.UserRole(roleStr),
	}, nil
}
