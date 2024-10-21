package user_ports

import (
	"loyalty-campaigns/src/common/models"
)

type IUserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List() ([]models.User, error)
	SearchByName(name string) ([]models.User, error)
	GetUserWithTransactions(id uint) (*models.User, error)
	GetUserWithRewards(id uint) (*models.User, error)
}
