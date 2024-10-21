package user_responses

import "time"

type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TransactionResponse struct {
	ID       uint      `json:"id"`
	BranchID uint      `json:"branch_id"`
	Amount   float64   `json:"amount"`
	Date     time.Time `json:"date"`
}

type RewardResponse struct {
	ID         uint    `json:"id"`
	MerchantID uint    `json:"merchant_id"`
	Type       string  `json:"type"`
	Amount     float64 `json:"amount"`
}

type UserWithTransactionsResponse struct {
	ID           uint                  `json:"id"`
	Name         string                `json:"name"`
	Transactions []TransactionResponse `json:"transactions"`
}

type UserWithRewardsResponse struct {
	ID      uint             `json:"id"`
	Name    string           `json:"name"`
	Rewards []RewardResponse `json:"rewards"`
}
