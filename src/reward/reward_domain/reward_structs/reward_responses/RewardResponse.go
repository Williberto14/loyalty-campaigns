package reward_responses

type RewardResponse struct {
	ID         uint    `json:"id"`
	UserID     uint    `json:"user_id"`
	MerchantID uint    `json:"merchant_id"`
	Type       string  `json:"type"`
	Amount     float64 `json:"amount"`
}

type TotalRewardsResponse struct {
	UserID        uint    `json:"user_id"`
	TotalPoints   float64 `json:"total_points"`
	TotalCashback float64 `json:"total_cashback"`
}
