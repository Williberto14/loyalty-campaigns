package merchant_responses

import "time"

type MerchantResponse struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	ConversionFactor float64   `json:"conversion_factor"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
