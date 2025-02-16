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

func TestGetInfo(t *testing.T) {
	router, rr, db := setupApp()
	defer func() {
		if err := clearDatabaseInfoTest(db); err != nil {
			t.Fatalf("Failed to clear database: %v", err)
		}
		db.Close()
	}()

	// Тело теста
	authRequest := handler.AuthRequest{
		Username: "testuser1",
		Password: "testpassword1",
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

	token := authResponse.Token
	if token == "" {
		t.Error("Ожидался токен, но получили пустую строку")
	}

	reqInfo, err := http.NewRequest("GET", "/api/info", nil)
	if err != nil {
		t.Fatalf("Не удалось создать запрос: %v", err)
	}
	reqInfo.Header.Set("Authorization", "Bearer "+token)

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, reqInfo)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус код %v, но получили %v", http.StatusOK, status)
	}

	var infoResponse handler.InfoResponse
	if err := json.NewDecoder(rr.Body).Decode(&infoResponse); err != nil {
		t.Fatalf("Не удалось декодировать ответ: %v", err)
	}

	if infoResponse.Coins != 1000 {
		t.Error("Ожидалось 1000 монет")
	}

	if len(infoResponse.Inventory) > 0 {
		t.Logf("Инвентарь: %v", infoResponse.Inventory)
	}

	if len(infoResponse.CoinHistory.Received) > 0 {
		t.Logf("История полученных монет: %v", infoResponse.CoinHistory.Received)
	}
	if len(infoResponse.CoinHistory.Sent) > 0 {
		t.Logf("История отправленных монет: %v", infoResponse.CoinHistory.Sent)
	}

	testCases := []struct {
		name            string
		token           string
		expectedStatus  int
		expectedMessage string
	}{
		{"Неправильный токен", "invalid_token", http.StatusUnauthorized, "Некорректный токен"},
		{"Отсутствие заголовка Authorization", "", http.StatusUnauthorized, "Отсутствует заголовок Authorization"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqInfoWrongToken, err := http.NewRequest("GET", "/api/info", nil)
			if err != nil {
				t.Fatalf("Не удалось создать запрос: %v", err)
			}

			if tc.token != "" {
				reqInfoWrongToken.Header.Set("Authorization", "Bearer "+tc.token)
			}

			rrWrongToken := httptest.NewRecorder()
			router.ServeHTTP(rrWrongToken, reqInfoWrongToken)

			if status := rrWrongToken.Code; status != tc.expectedStatus {
				t.Errorf("Ожидался статус код %v, но получили %v", tc.expectedStatus, status)
			}

			if !strings.Contains(rrWrongToken.Body.String(), tc.expectedMessage) {
				t.Errorf("Ожидалось сообщение об ошибке '%s', получено '%s'", tc.expectedMessage, rrWrongToken.Body.String())
			}
		})
	}
}
