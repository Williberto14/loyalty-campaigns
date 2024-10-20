package models

import (
	"gorm.io/gorm"
)

type Reward struct {
	gorm.Model
	UserID     uint
	User       User
	MerchantID uint
	Merchant   Merchant
	Type       string
	Amount     float64
}
