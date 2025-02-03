package models

import "time"

type TransactionType string

const (
	Deposit  TransactionType = "deposit"
	Transfer TransactionType = "transfer"
	Receive  TransactionType = "receive"
)

type Transaction struct {
	ID        int             `json:"id"`
	UserId    int             `json:"id"`
	Type      TransactionType `json:"type"`
	Amount    int             `json:"amount"`
	FromId    *int            `json:"from_id,omitempty"`
	ToId      *int            `json:"to_id,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}
