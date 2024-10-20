package models

import (
	"gorm.io/gorm"
)

type Merchant struct {
	gorm.Model
	Name             string
	ConversionFactor float64
	Branches         []Branch
	Campaigns        []Campaign
}
