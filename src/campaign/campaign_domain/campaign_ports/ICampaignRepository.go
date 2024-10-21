package campaign_ports

import (
	"loyalty-campaigns/src/common/models"
	"time"
)

type CampaignRepository interface {
	Create(campaign *models.Campaign) error
	GetByID(id uint) (*models.Campaign, error)
	Update(campaign *models.Campaign) error
	Delete(id uint) error
	List() ([]models.Campaign, error)
	GetActiveCampaigns(merchantID, branchID uint, date time.Time) ([]models.Campaign, error)
}
