package branch_requests

type UpdateBranchRequest struct {
	Name       string `json:"name" binding:"required"`
	MerchantID uint   `json:"merchant_id"`
}
