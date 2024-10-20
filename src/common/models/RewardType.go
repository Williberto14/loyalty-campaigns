package models

import (
	"gorm.io/gorm"
)

type RewardType struct {
	gorm.Model
	Name string
}
