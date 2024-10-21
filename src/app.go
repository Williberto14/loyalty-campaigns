package src

import (
	"log"
	"loyalty-campaigns/src/common/configs"
	"loyalty-campaigns/src/common/utils"
	"loyalty-campaigns/src/merchant/merchant_infra/merchant_controller"

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

	merchant_controller.NewMerchantController(router)

	router.Run("127.0.0.1:7070")

}
