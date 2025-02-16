package tests

import (
	"avitoMerch/internal/handler"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuyItem(t *testing.T) {
	router, rr, db := setupApp()
	defer func() {
		if err := clearDatabaseBuyItemTest(db); err != nil {
			t.Fatalf("Failed to clear database: %v", err)
		}
		db.Close()
	}()

	// Тело теста
	authRequest := handler.AuthRequest{
		Username: "testuser",
		Password: "testpassword",
	}
	requestBody, err := json.Marshal(authRequest)
	if err != nil {
		t.Fatalf("Не удалось сериализовать запрос: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Не удалось создать запрос: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус код %v, но получили %v", http.StatusOK, status)
	}

	var authResponse handler.AuthResponse
	if err := json.NewDecoder(rr.Body).Decode(&authResponse); err != nil {
		t.Fatalf("Не удалось декодировать ответ: %v", err)
	}

	if authResponse.Token == "" {
		t.Error("Ожидался токен, но получили пустую строку")
	}

	testCases := []struct {
		item           string
		expectedStatus int
	}{
		{"t-shirt", http.StatusOK},
		{"cup", http.StatusOK},
		{"book", http.StatusOK},
		{"pen", http.StatusOK},
		{"hoody", http.StatusOK},
		{"umbrella", http.StatusOK},
		{"socks", http.StatusOK},
		{"wallet", http.StatusOK},
		{"nonexistent-item", http.StatusBadRequest}, // Тест на несуществующий товар
	}

	for _, tc := range testCases {
		reqBuy, err := http.NewRequest("GET", "/api/buy/"+tc.item, nil)
		if err != nil {
			t.Fatalf("Не удалось создать запрос: %v", err)
		}
		reqBuy.Header.Set("Authorization", "Bearer "+authResponse.Token)

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, reqBuy)

		if status := rr.Code; status != tc.expectedStatus {
			t.Errorf("Для товара '%s' ожидался статус код %v, но получили %v", tc.item, tc.expectedStatus, status)
		}

		if tc.expectedStatus == http.StatusBadRequest && !strings.Contains(rr.Body.String(), "item not found") {
			t.Errorf("Ожидалась ошибка о том, что товар не найден для '%s', но получили: %s", tc.item, rr.Body.String())
		}
	}
}
