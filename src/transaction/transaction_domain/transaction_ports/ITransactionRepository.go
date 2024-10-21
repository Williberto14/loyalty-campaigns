package transaction_ports

import (
	"loyalty-campaigns/src/common/models"
	"time"
)

type ITransactionRepository interface {
	Create(transaction *models.Transaction) error
	GetByID(id uint) (*models.Transaction, error)
	Update(transaction *models.Transaction) error
	Delete(id uint) error
	List() ([]models.Transaction, error)
	GetByUserID(userID uint) ([]models.Transaction, error)
	GetByBranchID(branchID uint) ([]models.Transaction, error)
	GetByDateRange(startDate, endDate time.Time) ([]models.Transaction, error)
	GetByUserAndDateRange(userID uint, startDate, endDate time.Time) ([]models.Transaction, error)
	GetTotalAmountByUserAndDateRange(userID uint, startDate, endDate time.Time) (float64, error)
}
