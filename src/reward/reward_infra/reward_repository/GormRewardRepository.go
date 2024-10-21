package reward_repository

import (
	"loyalty-campaigns/src/common/models"
	"loyalty-campaigns/src/reward/reward_domain/reward_ports"
	"time"

	"gorm.io/gorm"
)

type GormRewardRepository struct {
	DB *gorm.DB
}

func NewGormRewardRepository(db *gorm.DB) reward_ports.IRewardRepository {
	return &GormRewardRepository{DB: db}
}

func (r *GormRewardRepository) Create(reward *models.Reward) error {
	return r.DB.Create(reward).Error
}

func (r *GormRewardRepository) GetByID(id uint) (*models.Reward, error) {
	var reward models.Reward
	err := r.DB.First(&reward, id).Error
	if err != nil {
		return nil, err
	}
	return &reward, nil
}

func (r *GormRewardRepository) Update(reward *models.Reward) error {
	return r.DB.Save(reward).Error
}

func (r *GormRewardRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Reward{}, id).Error
}

func (r *GormRewardRepository) List() ([]models.Reward, error) {
	var rewards []models.Reward
	err := r.DB.Find(&rewards).Error
	return rewards, err
}

func (r *GormRewardRepository) GetByUserID(userID uint) ([]models.Reward, error) {
	var rewards []models.Reward
	err := r.DB.Where("user_id = ?", userID).Find(&rewards).Error
	return rewards, err
}

func (r *GormRewardRepository) GetByMerchantID(merchantID uint) ([]models.Reward, error) {
	var rewards []models.Reward
	err := r.DB.Where("merchant_id = ?", merchantID).Find(&rewards).Error
	return rewards, err
}

func (r *GormRewardRepository) GetByUserAndMerchant(userID, merchantID uint) ([]models.Reward, error) {
	var rewards []models.Reward
	err := r.DB.Where("user_id = ? AND merchant_id = ?", userID, merchantID).Find(&rewards).Error
	return rewards, err
}

func (r *GormRewardRepository) SumRewardsByUser(userID uint, rewardType string) (float64, error) {
	var totalReward float64
	err := r.DB.Model(&models.Reward{}).
		Select("SUM(amount)").
		Where("user_id = ? AND type = ? AND is_redeemed = false", userID, rewardType).
		Scan(&totalReward).Error
	return totalReward, err
}

func (r *GormRewardRepository) GetActiveRewards(userID uint, currentDate time.Time) ([]models.Reward, error) {
	var rewards []models.Reward
	err := r.DB.Where("user_id = ? AND expiry_date > ? AND is_redeemed = false", userID, currentDate).Find(&rewards).Error
	return rewards, err
}

func (r *GormRewardRepository) MarkAsRedeemed(rewardID uint) error {
	return r.DB.Model(&models.Reward{}).Where("id = ?", rewardID).Update("is_redeemed", true).Error
}

func (r *GormRewardRepository) GetExpiredRewards(currentDate time.Time) ([]models.Reward, error) {
	var rewards []models.Reward
	err := r.DB.Where("expiry_date <= ? AND is_redeemed = false", currentDate).Find(&rewards).Error
	return rewards, err
}

func (r *GormRewardRepository) GetTotalRewardsByUser(userID uint) (totalPoints float64, totalCashback float64, err error) {
	var pointsSum, cashbackSum struct {
		Total float64
	}

	err = r.DB.Model(&models.Reward{}).
		Select("SUM(amount) as total").
		Where("user_id = ? AND type = ?", userID, "points").
		Scan(&pointsSum).Error
	if err != nil {
		return 0, 0, err
	}

	err = r.DB.Model(&models.Reward{}).
		Select("SUM(amount) as total").
		Where("user_id = ? AND type = ?", userID, "cashback").
		Scan(&cashbackSum).Error
	if err != nil {
		return 0, 0, err
	}

	return pointsSum.Total, cashbackSum.Total, nil
}

func (r *GormRewardRepository) GetByUserMerchantAndType(userID, merchantID uint, rewardType string) ([]models.Reward, error) {
	var rewards []models.Reward
	err := r.DB.Where("user_id = ? AND merchant_id = ? AND type = ?", userID, merchantID, rewardType).Find(&rewards).Error
	return rewards, err
}
