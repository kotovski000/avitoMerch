package repository

import (
	"avitoMerch/internal/models"
	"database/sql"
)

const (
	getItemInventoryQuery = `SELECT item_id, quantity FROM inventory WHERE user_id = $1`
	updateInventoryQuery  = `UPDATE inventory SET quantity = quantity + $3 WHERE user_id = $1 AND item_id = $2`
	insertInventoryQuery  = `INSERT INTO inventory (user_id, item_id, quantity) VALUES ($1, $2, $3)`
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) GetItemPrice(itemName string) (int, bool) {
	price, ok := models.ItemPrices[itemName]
	return price, ok
}

func (r *ItemRepository) GetUserInventory(userID int) (map[string]int, error) {
	rows, err := r.db.Query(getItemInventoryQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventory := make(map[string]int)
	for rows.Next() {
		var itemID string
		var quantity int
		if err := rows.Scan(&itemID, &quantity); err != nil {
			return nil, err
		}
		inventory[itemID] = quantity
	}

	return inventory, nil
}

func (r *ItemRepository) AddItemToInventory(tx *sql.Tx, userID int, itemID string, quantity int) error {
	result, err := tx.Exec(updateInventoryQuery, userID, itemID, quantity)

	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		_, err = tx.Exec(insertInventoryQuery, userID, itemID, quantity)
		if err != nil {
			return err
		}
	}
	return nil
}
