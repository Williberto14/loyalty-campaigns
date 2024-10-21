package loyalty_requests

type RedeemRewardsRequest struct {
	UserID     uint    `json:"userId" binding:"required"`
	MerchantID uint    `json:"merchantId" binding:"required"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
	RewardType string  `json:"rewardType" binding:"required,oneof=points cashback"`
}
