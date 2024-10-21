package campaign_repository

import (
	"loyalty-campaigns/src/campaign/campaign_domain/campaign_ports"
	"loyalty-campaigns/src/common/models"
	"time"

	"gorm.io/gorm"
)

type GormCampaignRepository struct {
	DB *gorm.DB
}

func NewGormCampaignRepository(db *gorm.DB) campaign_ports.ICampaignRepository {
	return &GormCampaignRepository{DB: db}
}

func (r *GormCampaignRepository) Create(campaign *models.Campaign) error {
	return r.DB.Create(campaign).Error
}

func (r *GormCampaignRepository) GetByID(id uint) (*models.Campaign, error) {
	var campaign models.Campaign
	err := r.DB.First(&campaign, id).Error
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (r *GormCampaignRepository) Update(campaign *models.Campaign) error {
	return r.DB.Save(campaign).Error
}

func (r *GormCampaignRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Campaign{}, id).Error
}

func (r *GormCampaignRepository) List() ([]models.Campaign, error) {
	var campaigns []models.Campaign
	err := r.DB.Find(&campaigns).Error
	return campaigns, err
}

func (r *GormCampaignRepository) GetActiveCampaigns(merchantID, branchID uint, date time.Time) ([]models.Campaign, error) {
	var campaigns []models.Campaign
	err := r.DB.Where("merchant_id = ? AND branch_id = ? AND start_date <= ? AND end_date >= ?",
		merchantID, branchID, date, date).Find(&campaigns).Error
	return campaigns, err
}
