package src

import (
	"log"
	"loyalty-campaigns/src/common/configs"
	"loyalty-campaigns/src/common/utils"
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
}
