package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID   uint
	User     User
	BranchID uint
	Branch   Branch
	Amount   float64
	Date     time.Time
}
