package middleware

import (
	"avitoMerch/internal/repository"
	"avitoMerch/internal/utils"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	jwtService *utils.JWTService
	userRepo   *repository.UserRepository
}

func NewAuthMiddleware(jwtService *utils.JWTService, userRepo *repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService, userRepo: userRepo}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Отсутствует заголовок Authorization", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		if tokenString == "" {
			http.Error(w, "Некорректный формат заголовка Authorization", http.StatusUnauthorized)
			return
		}

		token, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Некорректный токен: "+err.Error(), http.StatusUnauthorized)
			return
		}

		userID, err := m.jwtService.GetUserIDFromToken(token)
		if err != nil {
			http.Error(w, "Некорректные claims токена", http.StatusUnauthorized)
			return
		}

		user, err := m.userRepo.GetUserByID(userID)
		if err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) GetUserIDFromContext(ctx context.Context) (int, error) {
	userIDValue := ctx.Value("userID")
	if userIDValue == nil {
		return 0, fmt.Errorf("ID пользователя не найден в контексте")
	}

	userIDFloat, ok := userIDValue.(int)
	if !ok {
		return 0, fmt.Errorf("неверный тип ID пользователя в контексте")
	}

	return userIDFloat, nil
}
