package loyalty_app

import (
	"errors"
	"loyalty-campaigns/src/campaign/campaign_app"
	"loyalty-campaigns/src/common/utils"
	"loyalty-campaigns/src/merchant/merchant_app"
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
	merchantService    merchant_app.IMerchantService
	logger             utils.ILogger
}

func NewLoyaltyService(
	transactionService transaction_app.ITransactionService,
	campaignService campaign_app.ICampaignService,
	rewardService reward_app.IRewardService,
	merchantService merchant_app.IMerchantService,
) ILoyaltyService {
	return &loyaltyService{
		transactionService: transactionService,
		campaignService:    campaignService,
		rewardService:      rewardService,
		merchantService:    merchantService,
		logger:             utils.NewLogger(),
	}
}

func (s *loyaltyService) ProcessTransaction(userID, merchantID, branchID uint, amount float64, date time.Time) error {
	// 1. Crear la transacción
	_, err := s.transactionService.CreateTransaction(transaction_requests.CreateTransactionRequest{
		UserID:   userID,
		BranchID: branchID,
		Amount:   amount,
		Date:     date,
	})
	if err != nil {
		s.logger.Error("Error al crear transacción", err)
		return err
	}

	// Obtener el merchant
	merchant, err := s.merchantService.GetMerchant(merchantID)
	if err != nil {
		s.logger.Error("Error al obtener merchant", err)
		return err
	}

	// Calcular recompensa base
	baseReward := amount * merchant.ConversionFactor

	// Obtener campañas activas
	activeCampaigns, err := s.campaignService.GetActiveCampaigns(merchantID, &branchID, date)
	if err != nil {
		s.logger.Error("Error al obtener campañas activas", err)
		return err
	}

	// Procesar recompensas
	if len(activeCampaigns) > 0 {
		// Hay campañas activas
		for _, campaign := range activeCampaigns {
			if campaign.MinAmount == nil || amount >= *campaign.MinAmount {
				finalReward := baseReward * campaign.Value
				_, err = s.rewardService.CreateReward(reward_requests.CreateRewardRequest{
					UserID:     userID,
					MerchantID: merchantID,
					Type:       campaign.Type,
					Amount:     finalReward,
				})
				if err != nil {
					s.logger.Error("Error al crear recompensa de campaña", err)
					return err
				}
			}
		}
	} else {
		// No hay campañas activas, otorgar la recompensa base según el tipo predeterminado del merchant
		_, err = s.rewardService.CreateReward(reward_requests.CreateRewardRequest{
			UserID:     userID,
			MerchantID: merchantID,
			Type:       merchant.DefaultRewardType,
			Amount:     baseReward,
		})
		if err != nil {
			s.logger.Error("Error al crear recompensa base", err)
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
