package src

import (
	"log"
	"loyalty-campaigns/src/branch/branch_infra/branch_controller"
	"loyalty-campaigns/src/campaign/campaign_infra/campaign_controller"
	"loyalty-campaigns/src/common/configs"
	"loyalty-campaigns/src/common/utils"
	"loyalty-campaigns/src/loyalty/loyalty_infra/loyalty_controller"
	"loyalty-campaigns/src/merchant/merchant_infra/merchant_controller"
	"loyalty-campaigns/src/reward/reward_infra/reward_controller"
	"loyalty-campaigns/src/transaction/transaction_infra/transaction_controller"
	"loyalty-campaigns/src/user/user_infra/user_controller"

	_ "loyalty-campaigns/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

var logger = utils.NewLogger()

func Run() {
	dbConnection := configs.NewDBConnection()
	defer dbConnection.Close()

	err := configs.Migrate(dbConnection.GetDB())
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	logger.Info("[OK] Migrations completed")

	gin.SetMode("debug")
	router := gin.Default()

	// Register controllers
	user_controller.NewUserController(router)
	merchant_controller.NewMerchantController(router)
	branch_controller.NewBranchController(router)
	campaign_controller.NewCampaignController(router)
	reward_controller.NewRewardController(router)
	transaction_controller.NewTransactionController(router)
	loyalty_controller.NewLoyaltyController(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run("127.0.0.1:7070")

}
