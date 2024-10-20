package models

import (
	"gorm.io/gorm"
)

type Reward struct {
	gorm.Model
	Amount       float64
	UserID       uint
	User         User
	MerchantID   uint
	Merchant     Merchant
	RewardType   RewardType
	RewardTypeID uint
}
