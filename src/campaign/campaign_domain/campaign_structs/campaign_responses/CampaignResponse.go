package campaign_responses

import "time"

type CampaignResponse struct {
	ID         uint      `json:"id"`
	MerchantID uint      `json:"merchant_id"`
	BranchID   uint      `json:"branch_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Type       string    `json:"type"`
	Value      float64   `json:"value"`
	MinAmount  float64   `json:"min_amount"`
}
