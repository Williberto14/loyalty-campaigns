package merchant_requests

type UpdateMerchantRequest struct {
	Name             string  `json:"name" binding:"required"`
	ConversionFactor float64 `json:"conversion_factor" binding:"required"`
}
