package configs

import (
	"loyalty-campaigns/src/common/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Merchant{},
		&models.Branch{},
		&models.Campaign{},
		&models.User{},
		&models.Transaction{},
		&models.Reward{},
	)
}
