package handler

import (
	"avitoMerch/internal/middleware"
	"avitoMerch/internal/models"
	"avitoMerch/internal/repository"
	"avitoMerch/internal/service"
	"encoding/json"
	"net/http"
)

type InfoHandler struct {
	transactionService *service.TransactionService
	itemService        *service.ItemService
	userRepo           *repository.UserRepository
}

func NewInfoHandler(transactionService *service.TransactionService, itemService *service.ItemService, userRepo *repository.UserRepository) *InfoHandler {
	return &InfoHandler{transactionService: transactionService, itemService: itemService, userRepo: userRepo}
}

type InfoResponse struct {
	Coins       int         `json:"coins"`
	Inventory   []Inventory `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

type Inventory struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type TransactionEntry struct {
	FromUser string `json:"fromUser,omitempty"`
	ToUser   string `json:"toUser,omitempty"`
	Amount   int    `json:"amount"`
}

type CoinHistory struct {
	Received []TransactionEntry `json:"received"`
	Sent     []TransactionEntry `json:"sent"`
}

func (h *InfoHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	middleware := middleware.AuthMiddleware{}
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	inventoryData, err := h.itemService.GetUserInventory(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	receivedTransactions, err := h.transactionService.GetReceivedTransactions(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	sentTransactions, err := h.transactionService.GetSentTransactions(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	inventory := []Inventory{}
	for itemType, quantity := range inventoryData {
		inventory = append(inventory, Inventory{
			Type:     itemType,
			Quantity: quantity,
		})
	}

	formatTransactionEntries := func(transactions []models.Transaction) []TransactionEntry {
		formatted := make([]TransactionEntry, len(transactions))
		for i, transaction := range transactions {
			entry := TransactionEntry{
				Amount: transaction.Amount,
			}
			if transaction.Type == "send" {
				entry.ToUser = transaction.ToUser
			} else if transaction.Type == "receive" {
				entry.FromUser = transaction.FromUser
			}
			formatted[i] = entry
		}
		return formatted
	}

	resp := InfoResponse{
		Coins:     user.Coins,
		Inventory: inventory,
		CoinHistory: CoinHistory{
			Received: formatTransactionEntries(receivedTransactions),
			Sent:     formatTransactionEntries(sentTransactions),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
