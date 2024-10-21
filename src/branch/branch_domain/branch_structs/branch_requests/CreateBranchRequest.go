package branch_requests

type CreateBranchRequest struct {
	Name       string `json:"name" binding:"required"`
	MerchantID uint   `json:"merchant_id" binding:"required"`
}
