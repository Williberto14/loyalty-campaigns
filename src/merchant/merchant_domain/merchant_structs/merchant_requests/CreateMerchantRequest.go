package merchant_requests

type CreateMerchantRequest struct {
	Name             string  `json:"name" binding:"required"`
	ConversionFactor float64 `json:"conversion_factor" binding:"required"`
}
