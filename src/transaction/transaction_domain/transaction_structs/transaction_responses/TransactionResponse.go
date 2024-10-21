package transaction_responses

import "time"

type TransactionResponse struct {
	ID       uint      `json:"id"`
	UserID   uint      `json:"user_id"`
	BranchID uint      `json:"branch_id"`
	Amount   float64   `json:"amount"`
	Date     time.Time `json:"date"`
}
