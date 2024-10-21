package transaction_repository

import (
	"loyalty-campaigns/src/common/models"
	"loyalty-campaigns/src/transaction/transaction_domain/transaction_ports"
	"time"

	"gorm.io/gorm"
)

type GormTransactionRepository struct {
	DB *gorm.DB
}

func NewGormTransactionRepository(db *gorm.DB) transaction_ports.TransactionRepository {
	return &GormTransactionRepository{DB: db}
}

func (r *GormTransactionRepository) Create(transaction *models.Transaction) error {
	return r.DB.Create(transaction).Error
}

func (r *GormTransactionRepository) GetByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.DB.First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *GormTransactionRepository) Update(transaction *models.Transaction) error {
	return r.DB.Save(transaction).Error
}

func (r *GormTransactionRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Transaction{}, id).Error
}

func (r *GormTransactionRepository) List() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.DB.Find(&transactions).Error
	return transactions, err
}

func (r *GormTransactionRepository) GetByUserID(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.DB.Where("user_id = ?", userID).Find(&transactions).Error
	return transactions, err
}

func (r *GormTransactionRepository) GetByBranchID(branchID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.DB.Where("branch_id = ?", branchID).Find(&transactions).Error
	return transactions, err
}

func (r *GormTransactionRepository) GetByDateRange(startDate, endDate time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.DB.Where("date BETWEEN ? AND ?", startDate, endDate).Find(&transactions).Error
	return transactions, err
}

func (r *GormTransactionRepository) GetByUserAndDateRange(userID uint, startDate, endDate time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.DB.Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).Find(&transactions).Error
	return transactions, err
}

func (r *GormTransactionRepository) SumAmountByUserAndDateRange(userID uint, startDate, endDate time.Time) (float64, error) {
	var totalAmount float64
	err := r.DB.Model(&models.Transaction{}).
		Select("SUM(amount)").
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).
		Scan(&totalAmount).Error
	return totalAmount, err
}
