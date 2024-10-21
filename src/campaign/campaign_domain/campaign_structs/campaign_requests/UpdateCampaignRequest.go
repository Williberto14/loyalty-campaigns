package campaign_requests

import "time"

type UpdateCampaignRequest struct {
	StartDate time.Time  `json:"startDate" binding:"required"`
	EndDate   *time.Time `json:"endDate"`
	Type      string     `json:"type" binding:"required"`
	Value     float64    `json:"value" binding:"required"`
	MinAmount *float64   `json:"minAmount"`
}
