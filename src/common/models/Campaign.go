package models

import (
	"time"

	"gorm.io/gorm"
)

type Campaign struct {
	gorm.Model
	MerchantID   uint
	Merchant     Merchant
	BranchID     uint
	Branch       Branch
	StartDate    time.Time
	EndDate      time.Time
	RewardTypeID uint
	RewardType   RewardType
	Value        float64
	MinAmount    float64
}
