package campaign_responses

import "time"

type CampaignResponse struct {
	ID         uint       `json:"id"`
	MerchantID uint       `json:"merchantId"`
	BranchID   *uint      `json:"branchId"`
	StartDate  time.Time  `json:"startDate"`
	EndDate    *time.Time `json:"endDate"`
	Type       string     `json:"type"`
	Value      float64    `json:"value"`
	MinAmount  *float64   `json:"minAmount"`
}
