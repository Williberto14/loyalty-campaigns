package loyalty_controller

import (
	"loyalty-campaigns/src/campaign/campaign_app"
	"loyalty-campaigns/src/campaign/campaign_infra/campaign_repository"
	"loyalty-campaigns/src/common/configs"
	"loyalty-campaigns/src/loyalty/loyalty_app"
	"loyalty-campaigns/src/loyalty/loyalty_domain/loyalty_structs/loyalty_requests"
	"loyalty-campaigns/src/reward/reward_app"
	"loyalty-campaigns/src/reward/reward_infra/reward_repository"
	"loyalty-campaigns/src/transaction/transaction_app"
	"loyalty-campaigns/src/transaction/transaction_infra/transaction_repository"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type LoyaltyController struct {
	loyaltyService loyalty_app.ILoyaltyService
}

var (
	loyaltyControllerInstance *LoyaltyController
	loyaltyControllerOnce     sync.Once
)

func NewLoyaltyController(router *gin.Engine) *LoyaltyController {
	loyaltyControllerOnce.Do(func() {
		loyaltyControllerInstance = &LoyaltyController{}

		db := configs.NewDBConnection().GetDB()

		transactionRepository := transaction_repository.NewGormTransactionRepository(db)
		campaignRepository := campaign_repository.NewGormCampaignRepository(db)
		rewardRepository := reward_repository.NewGormRewardRepository(db)

		transactionService := transaction_app.NewTransactionService(transactionRepository)
		campaignService := campaign_app.NewCampaignService(campaignRepository)
		rewardService := reward_app.NewRewardService(rewardRepository)

		loyaltyControllerInstance.loyaltyService = loyalty_app.NewLoyaltyService(
			transactionService,
			campaignService,
			rewardService,
		)

		loyaltyControllerInstance.setupRoutes(router)
	})
	return loyaltyControllerInstance
}

func (c *LoyaltyController) setupRoutes(router *gin.Engine) {
	loyaltyGroup := router.Group("/api/loyalty")
	{
		loyaltyGroup.POST("/process-transaction", c.ProcessTransaction)
		loyaltyGroup.POST("/redeem-rewards", c.RedeemRewards)
	}
}

// ProcessTransaction godoc
//	@Summary		Process a transaction and award loyalty points or cashback
//	@Description	Process a user transaction and award loyalty points or cashback based on active campaigns
//	@Tags			loyalty
//	@Accept			json
//	@Produce		json
//	@Param			request	body		loyalty_requests.ProcessTransactionRequest	true	"Transaction details"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/loyalty/process-transaction [post]
func (c *LoyaltyController) ProcessTransaction(ctx *gin.Context) {
	var req loyalty_requests.ProcessTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.loyaltyService.ProcessTransaction(req.UserID, req.MerchantID, req.BranchID, req.Amount, req.Date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Transaction processed successfully"})
}

// RedeemRewards godoc
//	@Summary		Redeem user rewards
//	@Description	Redeem a user's loyalty points or cashback
//	@Tags			loyalty
//	@Accept			json
//	@Produce		json
//	@Param			request	body		loyalty_requests.RedeemRewardsRequest	true	"Redemption details"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/loyalty/redeem-rewards [post]
func (c *LoyaltyController) RedeemRewards(ctx *gin.Context) {
	var req loyalty_requests.RedeemRewardsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.loyaltyService.RedeemRewards(req.UserID, req.MerchantID, req.Amount, req.RewardType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Rewards redeemed successfully"})
}
