package handler

import (
	"avitoMerch/internal/middleware"
	"avitoMerch/internal/repository"
	"avitoMerch/internal/service"
	"encoding/json"
	"net/http"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
	userRepo           *repository.UserRepository
}

func NewTransactionHandler(transactionService *service.TransactionService, userRepo *repository.UserRepository) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService, userRepo: userRepo}
}

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

func (h *TransactionHandler) SendCoin(w http.ResponseWriter, r *http.Request) {
	var req SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	middleware := middleware.AuthMiddleware{}
	senderID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	receiver, err := h.userRepo.GetUserByUsername(req.ToUser)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if receiver == nil {
		http.Error(w, "Receiver not found", http.StatusBadRequest)
		return
	}

	err = h.transactionService.SendCoins(senderID, receiver.ID, req.Amount)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Coins sent successfully"))
}
