package models

import (
	"gorm.io/gorm"
)

type Merchant struct {
	gorm.Model
	Name              string
	ConversionFactor  float64
	DefaultRewardType string
	Branches          []Branch
	Campaigns         []Campaign
}
