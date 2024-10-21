package models

import (
	"time"

	"gorm.io/gorm"
)

type Campaign struct {
	gorm.Model
	MerchantID uint      `gorm:"not null"`
	Merchant   Merchant  `gorm:"foreignKey:MerchantID"`
	BranchID   *uint     `gorm:"index"`
	Branch     Branch    `gorm:"foreignKey:BranchID"`
	StartDate  time.Time `gorm:"not null"`
	EndDate    *time.Time
	Type       string  `gorm:"not null"`
	Value      float64 `gorm:"not null"`
	MinAmount  *float64
}
