package reward_app

import (
	"loyalty-campaigns/src/common/models"
	"loyalty-campaigns/src/common/utils"
	"loyalty-campaigns/src/reward/reward_domain/reward_ports"
	"loyalty-campaigns/src/reward/reward_domain/reward_structs/reward_requests"
	"loyalty-campaigns/src/reward/reward_domain/reward_structs/reward_responses"
	"sync"
)

type IRewardService interface {
	CreateReward(req reward_requests.CreateRewardRequest) (*reward_responses.RewardResponse, error)
	GetReward(id uint) (*reward_responses.RewardResponse, error)
	ListRewardsByUser(userID uint) ([]reward_responses.RewardResponse, error)
	GetTotalRewardsByUser(userID uint) (*reward_responses.TotalRewardsResponse, error)
}

type rewardService struct {
	rewardRepo reward_ports.IRewardRepository
	logger     utils.ILogger
}

var (
	rewardServiceInstance *rewardService
	rewardServiceOnce     sync.Once
)

func NewRewardService(rewardRepo reward_ports.IRewardRepository) IRewardService {
	rewardServiceOnce.Do(func() {
		rewardServiceInstance = &rewardService{
			rewardRepo: rewardRepo,
			logger:     utils.NewLogger(),
		}
	})
	return rewardServiceInstance
}

func (s *rewardService) CreateReward(req reward_requests.CreateRewardRequest) (*reward_responses.RewardResponse, error) {
	reward := &models.Reward{
		UserID:     req.UserID,
		MerchantID: req.MerchantID,
		Type:       req.Type,
		Amount:     req.Amount,
	}

	err := s.rewardRepo.Create(reward)
	if err != nil {
		s.logger.Error("Error al crear recompensa", err)
		return nil, err
	}

	return mapRewardToResponse(reward), nil
}

func (s *rewardService) GetReward(id uint) (*reward_responses.RewardResponse, error) {
	reward, err := s.rewardRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Error al obtener recompensa", err)
		return nil, err
	}

	return mapRewardToResponse(reward), nil
}

func (s *rewardService) ListRewardsByUser(userID uint) ([]reward_responses.RewardResponse, error) {
	rewards, err := s.rewardRepo.GetByUserID(userID)
	if err != nil {
		s.logger.Error("Error al listar recompensas del usuario", err)
		return nil, err
	}

	return mapRewardsToResponses(rewards), nil
}

func (s *rewardService) GetTotalRewardsByUser(userID uint) (*reward_responses.TotalRewardsResponse, error) {
	totalPoints, totalCashback, err := s.rewardRepo.GetTotalRewardsByUser(userID)
	if err != nil {
		s.logger.Error("Error al obtener total de recompensas del usuario", err)
		return nil, err
	}

	return &reward_responses.TotalRewardsResponse{
		UserID:        userID,
		TotalPoints:   totalPoints,
		TotalCashback: totalCashback,
	}, nil
}

func mapRewardToResponse(reward *models.Reward) *reward_responses.RewardResponse {
	return &reward_responses.RewardResponse{
		ID:         reward.ID,
		UserID:     reward.UserID,
		MerchantID: reward.MerchantID,
		Type:       reward.Type,
		Amount:     reward.Amount,
	}
}

func mapRewardsToResponses(rewards []models.Reward) []reward_responses.RewardResponse {
	responses := make([]reward_responses.RewardResponse, len(rewards))
	for i, reward := range rewards {
		responses[i] = *mapRewardToResponse(&reward)
	}
	return responses
}
