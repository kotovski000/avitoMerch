package handler

import (
	"avitoMerch/internal/middleware"
	"avitoMerch/internal/repository"
	"avitoMerch/internal/service"
	"github.com/gorilla/mux"
	"net/http"
)

type ItemHandler struct {
	itemService        *service.ItemService
	transactionService *service.TransactionService
	userRepo           *repository.UserRepository
}

func NewItemHandler(itemService *service.ItemService, transactionService *service.TransactionService, userRepo *repository.UserRepository) *ItemHandler {
	return &ItemHandler{itemService: itemService, transactionService: transactionService, userRepo: userRepo}
}

func (h *ItemHandler) BuyItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	item := vars["item"]

	middleware := middleware.AuthMiddleware{}
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.itemService.BuyItem(userID, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item purchased successfully"))
}
