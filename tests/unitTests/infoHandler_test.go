package unitTests

import (
	"avitoMerch/internal/handler"
	"avitoMerch/internal/repository"
	"avitoMerch/internal/service"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestInfoHandler_GetInfo_Unauthorized(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	itemRepo := repository.NewItemRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	itemService := service.NewItemService(itemRepo, userRepo, transactionRepo)
	transactionService := service.NewTransactionService(transactionRepo, userRepo)

	infoHandler := handler.NewInfoHandler(transactionService, itemService, userRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)

	w := httptest.NewRecorder()
	infoHandler.GetInfo(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestInfoHandler_GetInfo_InvalidToken(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	itemRepo := repository.NewItemRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	itemService := service.NewItemService(itemRepo, userRepo, transactionRepo)
	transactionService := service.NewTransactionService(transactionRepo, userRepo)

	infoHandler := handler.NewInfoHandler(transactionService, itemService, userRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")

	w := httptest.NewRecorder()
	infoHandler.GetInfo(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
