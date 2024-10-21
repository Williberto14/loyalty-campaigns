package branch_ports

import "loyalty-campaigns/src/common/models"

type IBranchRepository interface {
	Create(branch *models.Branch) error
	GetByID(id uint) (*models.Branch, error)
	Update(branch *models.Branch) error
	Delete(id uint) error
	List() ([]models.Branch, error)
	GetByMerchantID(merchantID uint) ([]models.Branch, error)
	GetBranchWithCampaigns(id uint) (*models.Branch, error)
}
