package models

import "time"

type Transaction struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Type        string    `json:"type"`         // "buy", "send", "receive"
	RelatedUser int       `json:"related_user"` // (sender/receiver)
	FromUser    string    `json:"fromUser"`
	ToUser      string    `json:"toUser"`
	ItemID      string    `json:"item_id"`
	Amount      int       `json:"amount"`
	Timestamp   time.Time `json:"timestamp"`
}
