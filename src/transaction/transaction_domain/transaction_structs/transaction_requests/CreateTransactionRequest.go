package transaction_requests

import "time"

type CreateTransactionRequest struct {
	UserID   uint      `json:"user_id" binding:"required"`
	BranchID uint      `json:"branch_id" binding:"required"`
	Amount   float64   `json:"amount" binding:"required"`
	Date     time.Time `json:"date" binding:"required"`
}
