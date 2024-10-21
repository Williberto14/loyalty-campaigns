package branch_responses

import "time"

type BranchResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	MerchantID uint      `json:"merchant_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CampaignResponse struct {
	ID        uint      `json:"id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	MinAmount float64   `json:"min_amount"`
}

type BranchWithCampaignsResponse struct {
	ID         uint               `json:"id"`
	Name       string             `json:"name"`
	MerchantID uint               `json:"merchant_id"`
	Campaigns  []CampaignResponse `json:"campaigns"`
}
