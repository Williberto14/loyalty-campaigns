package loyalty_app

import (
	"errors"
	"loyalty-campaigns/src/campaign/campaign_app"
	"loyalty-campaigns/src/common/utils"
	"loyalty-campaigns/src/reward/reward_app"
	"loyalty-campaigns/src/reward/reward_domain/reward_structs/reward_requests"
	"loyalty-campaigns/src/transaction/transaction_app"
	"loyalty-campaigns/src/transaction/transaction_domain/transaction_structs/transaction_requests"
	"time"
)

type ILoyaltyService interface {
	ProcessTransaction(userID, merchantID, branchID uint, amount float64, date time.Time) error
	RedeemRewards(userID, merchantID uint, amount float64, rewardType string) error
}

type loyaltyService struct {
	transactionService transaction_app.ITransactionService
	campaignService    campaign_app.ICampaignService
	rewardService      reward_app.IRewardService
	logger             utils.ILogger
}

func NewLoyaltyService(
	transactionService transaction_app.ITransactionService,
	campaignService campaign_app.ICampaignService,
	rewardService reward_app.IRewardService,
) ILoyaltyService {
	return &loyaltyService{
		transactionService: transactionService,
		campaignService:    campaignService,
		rewardService:      rewardService,
		logger:             utils.NewLogger(),
	}
}

func (s *loyaltyService) ProcessTransaction(userID, merchantID, branchID uint, amount float64, date time.Time) error {
	// 1. Create the transaction
	_, err := s.transactionService.CreateTransaction(transaction_requests.CreateTransactionRequest{
		UserID:   userID,
		BranchID: branchID,
		Amount:   amount,
		Date:     date,
	})
	if err != nil {
		s.logger.Error("Error creating transaction", err)
		return err
	}

	// 2. Get active campaigns
	campaigns, err := s.campaignService.GetActiveCampaigns(merchantID, branchID, date)
	if err != nil {
		s.logger.Error("Error getting active campaigns", err)
		return err
	}

	// 3. Calculate rewards based on campaigns
	for _, campaign := range campaigns {
		var rewardAmount float64
		if campaign.Type == "points" {
			rewardAmount = amount * campaign.Value
		} else if campaign.Type == "cashback" {
			rewardAmount = amount * (campaign.Value / 100)
		}

		// 4. Create reward
		_, err := s.rewardService.CreateReward(reward_requests.CreateRewardRequest{
			UserID:     userID,
			MerchantID: merchantID,
			Type:       campaign.Type,
			Amount:     rewardAmount,
		})
		if err != nil {
			s.logger.Error("Error creating reward", err)
			return err
		}
	}

	return nil
}

func (s *loyaltyService) RedeemRewards(userID, merchantID uint, amount float64, rewardType string) error {
	// 1. Get user's rewards
	rewards, err := s.rewardService.ListRewardsByUser(userID)
	if err != nil {
		s.logger.Error("Error getting user rewards", err)
		return err
	}

	// 2. Calculate total available rewards
	var totalRewards float64
	for _, reward := range rewards {
		if reward.MerchantID == merchantID && reward.Type == rewardType {
			totalRewards += reward.Amount
		}
	}

	// 3. Check if user has enough rewards
	if totalRewards < amount {
		return errors.New("insufficient rewards")
	}

	// 4. Deduct rewards
	err = s.rewardService.DeductRewards(userID, merchantID, amount, rewardType)
	if err != nil {
		s.logger.Error("Error deducting rewards", err)
		return err
	}

	return nil
}
