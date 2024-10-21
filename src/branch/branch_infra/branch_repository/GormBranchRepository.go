package branch_repository

import (
	"loyalty-campaigns/src/branch/branch_domain/branch_ports"
	"loyalty-campaigns/src/common/models"

	"gorm.io/gorm"
)

type GormBranchRepository struct {
	DB *gorm.DB
}

func NewGormBranchRepository(db *gorm.DB) branch_ports.IBranchRepository {
	return &GormBranchRepository{DB: db}
}

func (r *GormBranchRepository) Create(branch *models.Branch) error {
	return r.DB.Create(branch).Error
}

func (r *GormBranchRepository) GetByID(id uint) (*models.Branch, error) {
	var branch models.Branch
	err := r.DB.First(&branch, id).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r *GormBranchRepository) Update(branch *models.Branch) error {
	return r.DB.Save(branch).Error
}

func (r *GormBranchRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Branch{}, id).Error
}

func (r *GormBranchRepository) List() ([]models.Branch, error) {
	var branches []models.Branch
	err := r.DB.Find(&branches).Error
	return branches, err
}

func (r *GormBranchRepository) GetByMerchantID(merchantID uint) ([]models.Branch, error) {
	var branches []models.Branch
	err := r.DB.Where("merchant_id = ?", merchantID).Find(&branches).Error
	return branches, err
}

func (r *GormBranchRepository) GetBranchWithCampaigns(id uint) (*models.Branch, error) {
	var branch models.Branch
	err := r.DB.Preload("Campaigns").First(&branch, id).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}
