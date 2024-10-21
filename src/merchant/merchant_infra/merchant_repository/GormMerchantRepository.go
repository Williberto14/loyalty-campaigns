package merchant_repository

import (
	"loyalty-campaigns/src/common/models"
	"loyalty-campaigns/src/merchant/merchant_domain/merchant_ports"

	"gorm.io/gorm"
)

type GormMerchantRepository struct {
	DB *gorm.DB
}

func NewGormMerchantRepository(db *gorm.DB) merchant_ports.IMerchantRepository {
	return &GormMerchantRepository{DB: db}
}

func (r *GormMerchantRepository) Create(merchant *models.Merchant) error {
	return r.DB.Create(merchant).Error
}

func (r *GormMerchantRepository) GetByID(id uint) (*models.Merchant, error) {
	var merchant models.Merchant
	err := r.DB.First(&merchant, id).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

func (r *GormMerchantRepository) Update(merchant *models.Merchant) error {
	return r.DB.Save(merchant).Error
}

func (r *GormMerchantRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Merchant{}, id).Error
}

func (r *GormMerchantRepository) List() ([]models.Merchant, error) {
	var merchants []models.Merchant
	err := r.DB.Find(&merchants).Error
	return merchants, err
}
