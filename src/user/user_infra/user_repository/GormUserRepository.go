package user_repository

import (
	"loyalty-campaigns/src/common/models"
	"loyalty-campaigns/src/user/user_domain/user_ports"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	DB *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) user_ports.IUserRepository {
	return &GormUserRepository{DB: db}
}

func (r *GormUserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *GormUserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}

func (r *GormUserRepository) Delete(id uint) error {
	return r.DB.Delete(&models.User{}, id).Error
}

func (r *GormUserRepository) List() ([]models.User, error) {
	var users []models.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *GormUserRepository) SearchByName(name string) ([]models.User, error) {
	var users []models.User
	err := r.DB.Where("name LIKE ?", "%"+name+"%").Find(&users).Error
	return users, err
}

func (r *GormUserRepository) GetUserWithTransactions(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.Preload("Transactions").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetUserWithRewards(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.Preload("Rewards").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
