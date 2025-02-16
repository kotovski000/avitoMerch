package service

import (
	"avitoMerch/internal/config"
	"avitoMerch/internal/repository"
	"avitoMerch/internal/utils"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtService *utils.JWTService
	bcryptCost int
}

func NewAuthService(userRepo *repository.UserRepository, jwtService *utils.JWTService, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, jwtService: jwtService, bcryptCost: cfg.BcryptCost}
}

func (s *AuthService) Authenticate(username, password string) (string, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	if user == nil {
		log.Printf("Пользователь не найден. Создаем нового пользователя: %s", username)
		hashedPassword, err := s.HashPassword(password)
		if err != nil {
			return "", err
		}
		user, err = s.userRepo.CreateUser(username, hashedPassword)
		if err != nil {
			return "", err
		}
	} else {
		err = s.verifyPassword(password, user.Password)
		if err != nil {
			return "", errors.New("неверные учетные данные")
		}
	}

	token, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *AuthService) verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
