package service

import (
	"avitoMerch/internal/models"
	"avitoMerch/internal/repository"
	"errors"
	"fmt"
)

type ItemService struct {
	itemRepo        *repository.ItemRepository
	userRepo        *repository.UserRepository
	transactionRepo *repository.TransactionRepository
}

func NewItemService(itemRepo *repository.ItemRepository, userRepo *repository.UserRepository, transactionRepo *repository.TransactionRepository) *ItemService {
	return &ItemService{itemRepo: itemRepo, userRepo: userRepo, transactionRepo: transactionRepo}
}

func (s *ItemService) BuyItem(userID int, itemName string) error {
	itemPrice, ok := s.itemRepo.GetItemPrice(itemName)
	if !ok {
		return errors.New("item not found")
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	if user.Coins < itemPrice {
		return errors.New("insufficient coins")
	}

	newBalance := user.Coins - itemPrice

	tx, err := s.userRepo.BeginTx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = s.userRepo.UpdateUserCoinsTx(tx, userID, newBalance)
	if err != nil {
		return fmt.Errorf("failed to update user coins: %w", err)
	}

	transaction := models.Transaction{
		UserID: userID,
		Type:   "buy",
		ItemID: itemName,
		Amount: itemPrice,
	}

	err = s.transactionRepo.CreateTransactionTx(tx, transaction)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	err = s.itemRepo.AddItemToInventory(tx, userID, itemName, 1)
	if err != nil {
		return fmt.Errorf("failed to add item to inventory: %w", err)
	}

	return nil
}

func (s *ItemService) GetUserInventory(userID int) (map[string]int, error) {
	return s.itemRepo.GetUserInventory(userID)
}
