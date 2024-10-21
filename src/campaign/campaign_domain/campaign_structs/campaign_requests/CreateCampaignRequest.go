package campaign_requests

import "time"

type CreateCampaignRequest struct {
	MerchantID uint      `json:"merchant_id" binding:"required"`
	BranchID   uint      `json:"branch_id" binding:"required"`
	StartDate  time.Time `json:"start_date" binding:"required"`
	EndDate    time.Time `json:"end_date" binding:"required"`
	Type       string    `json:"type" binding:"required"`
	Value      float64   `json:"value" binding:"required"`
	MinAmount  float64   `json:"min_amount" binding:"required"`
}
