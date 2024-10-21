package loyalty_requests

import "time"

type ProcessTransactionRequest struct {
	UserID     uint      `json:"userId" binding:"required"`
	MerchantID uint      `json:"merchantId" binding:"required"`
	BranchID   uint      `json:"branchId" binding:"required"`
	Amount     float64   `json:"amount" binding:"required,gt=0"`
	Date       time.Time `json:"date" binding:"required"`
}
