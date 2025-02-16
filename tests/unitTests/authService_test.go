package unitTests

import (
	"avitoMerch/internal/config"
	"avitoMerch/internal/repository"
	"avitoMerch/internal/service"
	"avitoMerch/internal/utils"
	"log"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_Authenticate_Success_NewUser(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	cfg := config.LoadConfig()
	userRepo := repository.NewUserRepository(db)
	jwtService := utils.NewJWTService(cfg.JWTSecretKey)

	authService := service.NewAuthService(userRepo, jwtService, cfg)

	testUser := TestUser{
		Username: "testuser",
		Password: "password",
	}

	token, err := authService.Authenticate(testUser.Username, testUser.Password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	user, _ := userRepo.GetUserByUsername(testUser.Username)
	assert.NotNil(t, user)
	assert.Equal(t, testUser.Username, user.Username)

}

func TestAuthService_Authenticate_Success_ExistingUser(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	cfg := config.LoadConfig()
	userRepo := repository.NewUserRepository(db)
	jwtService := utils.NewJWTService(cfg.JWTSecretKey)

	authService := service.NewAuthService(userRepo, jwtService, cfg)

	testUser := TestUser{
		Username: "testuser",
		Password: "password",
	}

	token, err := authService.Authenticate(testUser.Username, testUser.Password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	user, _ := userRepo.GetUserByUsername(testUser.Username)
	assert.NotNil(t, user)
	assert.Equal(t, testUser.Username, user.Username)

}

func TestAuthService_Authenticate_InvalidCredentials(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	cfg := config.LoadConfig()
	userRepo := repository.NewUserRepository(db)
	jwtService := utils.NewJWTService(cfg.JWTSecretKey)

	authService := service.NewAuthService(userRepo, jwtService, cfg)

	testUser := TestUser{
		Username: "testuser",
		Password: "password",
	}

	token, err := authService.Authenticate(testUser.Username, "wrongpassword")

	assert.Error(t, err)
	assert.Equal(t, "неверные учетные данные", err.Error())
	assert.Empty(t, token)

}
