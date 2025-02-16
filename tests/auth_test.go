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

func TestAuthenticate(t *testing.T) {
	router, rr, db := setupApp()
	defer func() {
		if err := clearDatabaseAuthTest(db); err != nil {
			t.Fatalf("Failed to clear database: %v", err)
		}
		db.Close()
	}()

	// Тело теста
	authRequestCreate := handler.AuthRequest{
		Username: "newuser",
		Password: "password",
	}
	requestBodyCreate, err := json.Marshal(authRequestCreate)
	if err != nil {
		t.Fatalf("Не удалось сериализовать запрос: %v", err)
	}

	reqCreate, err := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(requestBodyCreate))
	if err != nil {
		t.Fatalf("Не удалось создать запрос: %v", err)
	}
	reqCreate.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(rr, reqCreate)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус код %v при создании пользователя, но получили %v", http.StatusOK, status)
	}

	var authResponseCreate handler.AuthResponse
	if err := json.NewDecoder(rr.Body).Decode(&authResponseCreate); err != nil {
		t.Fatalf("Не удалось декодировать ответ: %v", err)
	}

	if authResponseCreate.Token == "" {
		t.Error("Ожидался токен при создании пользователя, но получили пустую строку")
	}

	authRequestLogin := handler.AuthRequest{
		Username: "newuser",
		Password: "password",
	}
	requestBodyLogin, err := json.Marshal(authRequestLogin)
	if err != nil {
		t.Fatalf("Не удалось сериализовать запрос: %v", err)
	}

	reqLogin, err := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(requestBodyLogin))
	if err != nil {
		t.Fatalf("Не удалось создать запрос: %v", err)
	}
	reqLogin.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, reqLogin)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус код %v при входе в существующего пользователя, но получили %v", http.StatusOK, status)
	}

	var authResponseLogin handler.AuthResponse
	if err := json.NewDecoder(rr.Body).Decode(&authResponseLogin); err != nil {
		t.Fatalf("Не удалось декодировать ответ: %v", err)
	}

	if authResponseLogin.Token == "" {
		t.Error("Ожидался токен при входе в существующего пользователя, но получили пустую строку")
	}

	authRequestWrongPassword := handler.AuthRequest{
		Username: "newuser",
		Password: "wrongpassword",
	}
	requestBodyWrongPassword, err := json.Marshal(authRequestWrongPassword)
	if err != nil {
		t.Fatalf("Не удалось сериализовать запрос: %v", err)
	}

	reqWrongPassword, err := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(requestBodyWrongPassword))
	if err != nil {
		t.Fatalf("Не удалось создать запрос: %v", err)
	}
	reqWrongPassword.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, reqWrongPassword)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Ожидался статус код %v при входе с неправильным паролем, но получили %v", http.StatusUnauthorized, status)
	}

	expectedErrorMessage := "неверные учетные данные"
	if !strings.Contains(rr.Body.String(), expectedErrorMessage) {
		t.Errorf("Ожидалось сообщение об ошибке '%s', получено '%s'", expectedErrorMessage, rr.Body.String())
	}
}
