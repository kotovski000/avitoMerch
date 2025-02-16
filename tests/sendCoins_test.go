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

func TestSendCoin(t *testing.T) {
	router, rr, db := setupApp()
	defer func() {
		if err := clearDatabaseSendCoinTest(db); err != nil {
			t.Fatalf("Failed to clear database: %v", err)
		}
		db.Close()
	}()
	// Тело теста
	var (
		senderToken    string
		recipientToken string
	)

	testUsers := []struct {
		username string
		password string
	}{
		{"sender", "password1"},
		{"recipient", "password2"},
	}

	for _, user := range testUsers {
		authRequest := handler.AuthRequest{
			Username: user.username,
			Password: user.password,
		}
		requestBody, err := json.Marshal(authRequest)
		if err != nil {
			t.Fatalf("Не удалось сериализовать запрос: %v", err)
		}

		reqAuth, err := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("Не удалось создать запрос: %v", err)
		}
		reqAuth.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, reqAuth)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Ожидался статус код %v для пользователя %s, но получили %v",
				http.StatusOK,
				user.username,
				status,
			)
		}

		var authResponse handler.AuthResponse
		if err = json.Unmarshal(rr.Body.Bytes(), &authResponse); err != nil || authResponse.Token == "" {
			t.Errorf("Не удалось получить токен для пользователя %s.", user.username)
		}

		if user.username == "sender" {
			senderToken = authResponse.Token
		} else if user.username == "recipient" {
			recipientToken = authResponse.Token
		}

		rr = httptest.NewRecorder()
	}

	sendCoinRequests := []struct {
		fromUser       string
		toUser         string
		amount         int
		expectedStatus int
		expectedError  string
		expectSent     bool
		expectReceived bool
	}{
		{"sender", "recipient", 100, http.StatusOK, "", true, true},
		{"sender", "recipient", 5000, http.StatusBadRequest, "insufficient coins", false, false},
		{"sender", "nonexistent-user", 100, http.StatusBadRequest, "Receiver not found", false, false},
		{"recipient", "sender", 50, http.StatusOK, "", true, true},
	}

	for _, sendCoinReq := range sendCoinRequests {
		var token string

		if sendCoinReq.fromUser == "sender" {
			token = senderToken
		} else if sendCoinReq.fromUser == "recipient" {
			token = recipientToken
		}

		sendCoinRequestBody, _ := json.Marshal(handler.SendCoinRequest{
			ToUser: sendCoinReq.toUser,
			Amount: sendCoinReq.amount,
		})

		reqSendCoin, _ := http.NewRequest("POST", "/api/sendCoin",
			bytes.NewBuffer(sendCoinRequestBody))

		reqSendCoin.Header.Set("Content-Type", "application/json")
		reqSendCoin.Header.Set("Authorization", "Bearer "+token)

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, reqSendCoin)

		if status := rr.Code; status != sendCoinReq.expectedStatus {
			t.Errorf("Ожидался статус код %v для отправки от '%s' к '%s', но получили %v",
				sendCoinReq.expectedStatus,
				sendCoinReq.fromUser,
				sendCoinReq.toUser,
				status,
			)
		}

		if sendCoinReq.expectedStatus == http.StatusBadRequest && !strings.Contains(rr.Body.String(), sendCoinReq.expectedError) {
			t.Errorf("Ожидалась ошибка '%s' при отправке от '%s' к '%s', но получили: %s",
				sendCoinReq.expectedError,
				sendCoinReq.fromUser,
				sendCoinReq.toUser,
				rr.Body.String(),
			)
		}

		reqInfoSender, _ := http.NewRequest("GET", "/api/info", nil)
		reqInfoSender.Header.Set("Authorization", "Bearer "+token)

		rrInfoSender := httptest.NewRecorder()
		router.ServeHTTP(rrInfoSender, reqInfoSender)

		if status := rrInfoSender.Code; status != http.StatusOK {
			t.Errorf("Ошибка при получении информации об отправителе: %v", status)
		}

		var infoResponseSender handler.InfoResponse
		if err := json.NewDecoder(rrInfoSender.Body).Decode(&infoResponseSender); err != nil {
			t.Fatalf("Не удалось декодировать ответ: %v", err)
		}

		if sendCoinReq.expectSent {
			found := false
			for _, transaction := range infoResponseSender.CoinHistory.Sent {
				if transaction.ToUser == sendCoinReq.toUser && transaction.Amount == sendCoinReq.amount {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Не найдена запись в CoinHistory.Sent для отправителя '%s' к '%s'", sendCoinReq.fromUser, sendCoinReq.toUser)
			}
		}

		if sendCoinReq.toUser != "nonexistent-user" {
			var recipientTokenToUse string
			if sendCoinReq.toUser == "sender" {
				recipientTokenToUse = senderToken
			} else {
				recipientTokenToUse = recipientToken
			}

			reqInfoRecipient, _ := http.NewRequest("GET", "/api/info", nil)
			reqInfoRecipient.Header.Set("Authorization", "Bearer "+recipientTokenToUse)

			rrInfoRecipient := httptest.NewRecorder()
			router.ServeHTTP(rrInfoRecipient, reqInfoRecipient)

			if status := rrInfoRecipient.Code; status != http.StatusOK {
				t.Errorf("Ошибка при получении информации о получателе: %v", status)
			}

			var infoResponseRecipient handler.InfoResponse
			if err := json.NewDecoder(rrInfoRecipient.Body).Decode(&infoResponseRecipient); err != nil {
				t.Fatalf("Не удалось декодировать ответ: %v", err)
			}

			if sendCoinReq.expectReceived {
				found := false
				for _, transaction := range infoResponseRecipient.CoinHistory.Received {
					if transaction.FromUser == sendCoinReq.fromUser && transaction.Amount == sendCoinReq.amount {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Не найдена запись в CoinHistory.Received для получателя '%s' от '%s'", sendCoinReq.toUser, sendCoinReq.fromUser)
				}
			}
		}

		rr = httptest.NewRecorder()
	}
}
