package reward_ports

import (
	"loyalty-campaigns/src/common/models"
	"time"
)

type IRewardRepository interface {
	Create(reward *models.Reward) error
	GetByID(id uint) (*models.Reward, error)
	Update(reward *models.Reward) error
	Delete(id uint) error
	List() ([]models.Reward, error)
	GetByUserID(userID uint) ([]models.Reward, error)
	GetTotalRewardsByUser(userID uint) (totalPoints float64, totalCashback float64, err error)
	GetByMerchantID(merchantID uint) ([]models.Reward, error)
	GetByUserAndMerchant(userID, merchantID uint) ([]models.Reward, error)
	SumRewardsByUser(userID uint, rewardType string) (float64, error)
	GetActiveRewards(userID uint, currentDate time.Time) ([]models.Reward, error)
	MarkAsRedeemed(rewardID uint) error
	GetExpiredRewards(currentDate time.Time) ([]models.Reward, error)
	GetByUserMerchantAndType(userID, merchantID uint, rewardType string) ([]models.Reward, error)
}
