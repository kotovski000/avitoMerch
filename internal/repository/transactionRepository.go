package repository

import (
	"avitoMerch/internal/models"
	"database/sql"
	"log"
)

const (
	createTransactionQuery       = `INSERT INTO transactions (user_id, type, related_user, item_id, amount) VALUES ($1, $2, $3, $4, $5)`
	getReceivedTransactionsQuery = `
        SELECT 
            t.id, 
            t.user_id, 
            t.type, 
            t.related_user, 
            u.username as related_username,
            t.item_id, 
            t.amount, 
            t.timestamp
        FROM transactions t
        JOIN users u ON t.related_user = u.id 
        WHERE t.user_id = $1 AND t.type = 'receive'
    `
	getSentTransactionsQuery = `
        SELECT 
            t.id, 
            t.user_id, 
            t.type, 
            t.related_user,
            u.username as related_username, 
            t.item_id, 
            t.amount, 
            t.timestamp
        FROM transactions t
        JOIN users u ON t.related_user = u.id 
        WHERE t.user_id = $1 AND t.type = 'send'
    `
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransactionTx(tx *sql.Tx, transaction models.Transaction) error {
	_, err := tx.Exec(createTransactionQuery, transaction.UserID, transaction.Type, transaction.RelatedUser, transaction.ItemID, transaction.Amount)
	if err != nil {
		log.Printf("Error creating transaction: %v", err)
		return err
	}
	return nil
}

func (r *TransactionRepository) GetReceivedTransactions(userID int) ([]models.Transaction, error) {
	rows, err := r.db.Query(getReceivedTransactionsQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var relatedUsername string
		if err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Type,
			&t.RelatedUser,
			&relatedUsername,
			&t.ItemID,
			&t.Amount,
			&t.Timestamp,
		); err != nil {
			return nil, err
		}
		t.FromUser = relatedUsername
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetSentTransactions(userID int) ([]models.Transaction, error) {
	rows, err := r.db.Query(getSentTransactionsQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var relatedUsername string
		if err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Type,
			&t.RelatedUser,
			&relatedUsername,
			&t.ItemID,
			&t.Amount,
			&t.Timestamp,
		); err != nil {
			return nil, err
		}
		t.ToUser = relatedUsername
		transactions = append(transactions, t)
	}

	return transactions, nil
}
