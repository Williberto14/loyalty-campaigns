package merchant_ports

import "loyalty-campaigns/src/common/models"

type IMerchantRepository interface {
	Create(merchant *models.Merchant) error
	GetByID(id uint) (*models.Merchant, error)
	Update(merchant *models.Merchant) error
	Delete(id uint) error
	List() ([]models.Merchant, error)
}
