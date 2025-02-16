package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type JWTService struct {
	secretKey string
}

func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey: secretKey,
	}
}

func (j *JWTService) GenerateToken(userID int) (string, error) {
	claims := &jwt.MapClaims{
		"UserID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (j *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (j *JWTService) GetUserIDFromToken(token *jwt.Token) (int, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDFloat, ok := claims["UserID"].(float64)
		if !ok {
			return 0, fmt.Errorf("некорректный ID пользователя в токене")
		}
		userID, err := strconv.Atoi(fmt.Sprintf("%.0f", userIDFloat))
		if err != nil {
			return 0, fmt.Errorf("некорректный формат ID пользователя")
		}
		return userID, nil
	}
	return 0, fmt.Errorf("некорректные claims токена")
}
