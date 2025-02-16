package service

import (
	"avitoMerch/internal/models"
	"avitoMerch/internal/repository"
	"errors"
	"fmt"
	"log"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	userRepo        *repository.UserRepository
}

func NewTransactionService(transactionRepo *repository.TransactionRepository, userRepo *repository.UserRepository) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo, userRepo: userRepo}
}

func (s *TransactionService) SendCoins(senderID, receiverID int, amount int) error {
	if senderID == receiverID {
		return errors.New("cannot send coins to yourself")
	}

	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	sender, err := s.userRepo.GetUserByID(senderID)
	if err != nil {
		return err
	}

	if sender == nil {
		return errors.New("sender not found")
	}

	receiver, err := s.userRepo.GetUserByID(receiverID)
	if err != nil {
		return err
	}

	if receiver == nil {
		return errors.New("receiver not found")
	}

	if sender.Coins < amount {
		return errors.New("insufficient coins")
	}

	newSenderBalance := sender.Coins - amount

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

	err = s.userRepo.UpdateUserCoinsTx(tx, senderID, newSenderBalance)
	if err != nil {
		return fmt.Errorf("failed to update sender coins: %w", err)
	}

	newReceiverBalance := receiver.Coins + amount
	err = s.userRepo.UpdateUserCoinsTx(tx, receiverID, newReceiverBalance)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			log.Printf("Failed to rollback transaction: %v", rErr)
		}
		return fmt.Errorf("failed to update receiver balance: %w", err)
	}

	transactionSender := models.Transaction{
		UserID:      senderID,
		Type:        "send",
		RelatedUser: receiverID,
		Amount:      amount,
	}

	err = s.transactionRepo.CreateTransactionTx(tx, transactionSender)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			log.Printf("Failed to rollback transaction: %v", rErr)
		}
		return fmt.Errorf("failed to record sender transaction: %w", err)
	}

	transactionReceiver := models.Transaction{
		UserID:      receiverID,
		Type:        "receive",
		RelatedUser: senderID,
		Amount:      amount,
	}

	err = s.transactionRepo.CreateTransactionTx(tx, transactionReceiver)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			log.Printf("Failed to rollback transaction: %v", rErr)
		}
		return fmt.Errorf("failed to record receiver transaction: %w", err)
	}

	return nil
}

func (s *TransactionService) GetReceivedTransactions(userID int) ([]models.Transaction, error) {
	return s.transactionRepo.GetReceivedTransactions(userID)
}

func (s *TransactionService) GetSentTransactions(userID int) ([]models.Transaction, error) {
	return s.transactionRepo.GetSentTransactions(userID)
}
