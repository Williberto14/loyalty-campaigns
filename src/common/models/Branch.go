package models

import (
	"gorm.io/gorm"
)

type Branch struct {
	gorm.Model
	Name         string
	MerchantID   uint
	Merchant     Merchant
	Campaigns    []Campaign
	Transactions []Transaction
}
