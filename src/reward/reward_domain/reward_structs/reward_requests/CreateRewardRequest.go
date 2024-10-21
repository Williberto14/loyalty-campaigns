package reward_requests

type CreateRewardRequest struct {
	UserID     uint    `json:"user_id" binding:"required"`
	MerchantID uint    `json:"merchant_id" binding:"required"`
	Type       string  `json:"type" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
}
