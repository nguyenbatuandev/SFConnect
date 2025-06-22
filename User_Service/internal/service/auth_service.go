package service

import (
	"errors"
	"time"
	"user_service/internal/entity"
	_interface "user_service/internal/interface"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type JWTauthService struct {
	secretKey       []byte
	expirationHours int
}

func NewJWTauthService(secretKey string, expirationHours int) *JWTauthService {
	return &JWTauthService{
		secretKey:       []byte(secretKey),
		expirationHours: expirationHours,
	}
}

func (s *JWTauthService) GenerateToken(user *entity.User) (string, error) {
	claim := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * time.Duration(s.expirationHours)).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(s.secretKey)
}

func (s *JWTauthService) ValidateToken(tokenString string) (*_interface.Claims, error) {
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

	// 🔒 Check expiration
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

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("invalid email in token")
	}

	roleStr, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid role in token")
	}

	return &_interface.Claims{
		UserID: userID,
		Email:  email,
		Role:   entity.UserRole(roleStr),
	}, nil
}



func (a *JWTauthService) HashPassword(password string) (string ,error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (a *JWTauthService) CheckPasswordHash(password, hashedPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
