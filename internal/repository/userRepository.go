package repository

import (
	"avitoMerch/internal/models"
	"database/sql"
	"log"
)

const (
	createUserQuery        = `INSERT INTO users (username, password, coins) VALUES ($1, $2, $3) RETURNING id, username, coins`
	getUserByUsernameQuery = `SELECT id, username, password, coins FROM users WHERE username = $1`
	updateUserCoinsQuery   = `UPDATE users SET coins = $1 WHERE id = $2`
	getUserByIDQuery       = `SELECT id, username, coins FROM users WHERE id = $1`
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(username, password string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(createUserQuery, username, password, 1000).Scan(&user.ID, &user.Username, &user.Coins)
	if err != nil {
		log.Printf("Ошибка при создании пользователя: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(getUserByUsernameQuery, username).Scan(&user.ID, &user.Username, &user.Password, &user.Coins)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUserCoinsTx(tx *sql.Tx, userID, newCoins int) error {
	_, err := tx.Exec(updateUserCoinsQuery, newCoins, userID)
	return err
}

func (r *UserRepository) GetUserByID(userID int) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(getUserByIDQuery, userID).Scan(&user.ID, &user.Username, &user.Coins)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}
