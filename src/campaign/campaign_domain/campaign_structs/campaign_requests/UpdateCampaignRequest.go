package campaign_requests

import "time"

type UpdateCampaignRequest struct {
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	Type      string    `json:"type" binding:"required"`
	Value     float64   `json:"value" binding:"required"`
	MinAmount float64   `json:"min_amount" binding:"required"`
}
