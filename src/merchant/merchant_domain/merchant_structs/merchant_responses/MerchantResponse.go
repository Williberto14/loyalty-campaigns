package merchant_responses

type MerchantResponse struct {
	ID                uint    `json:"id"`
	Name              string  `json:"name"`
	ConversionFactor  float64 `json:"conversion_factor"`
	DefaultRewardType string  `json:"defaultRewardType"`
}
