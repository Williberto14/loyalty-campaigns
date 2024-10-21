package user_app

import (
	"loyalty-campaigns/src/common/models"
	"loyalty-campaigns/src/user/user_domain/user_ports"
	"loyalty-campaigns/src/user/user_domain/user_structs/user_requests"
	"loyalty-campaigns/src/user/user_domain/user_structs/user_responses"
	"sync"
)

type IUserService interface {
	CreateUser(req user_requests.CreateUserRequest) (*user_responses.UserResponse, error)
	GetUser(id uint) (*user_responses.UserResponse, error)
	UpdateUser(id uint, req user_requests.UpdateUserRequest) (*user_responses.UserResponse, error)
	DeleteUser(id uint) error
	ListUsers() ([]user_responses.UserResponse, error)
	GetUserWithTransactions(id uint) (*user_responses.UserWithTransactionsResponse, error)
	GetUserWithRewards(id uint) (*user_responses.UserWithRewardsResponse, error)
}

type userService struct {
	userRepo user_ports.IUserRepository
}

var (
	userServiceInstance *userService
	userServiceOnce     sync.Once
)

func NewUserService(userRepo user_ports.IUserRepository) IUserService {
	userServiceOnce.Do(func() {
		userServiceInstance = &userService{
			userRepo: userRepo,
		}
	})
	return userServiceInstance
}

func (s *userService) CreateUser(req user_requests.CreateUserRequest) (*user_responses.UserResponse, error) {
	user := &models.User{
		Name: req.Name,
	}

	err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return &user_responses.UserResponse{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}

func (s *userService) GetUser(id uint) (*user_responses.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &user_responses.UserResponse{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}

func (s *userService) UpdateUser(id uint, req user_requests.UpdateUserRequest) (*user_responses.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.Name = req.Name

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return &user_responses.UserResponse{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}

func (s *userService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *userService) ListUsers() ([]user_responses.UserResponse, error) {
	users, err := s.userRepo.List()
	if err != nil {
		return nil, err
	}

	var responses []user_responses.UserResponse
	for _, user := range users {
		responses = append(responses, user_responses.UserResponse{
			ID:   user.ID,
			Name: user.Name,
		})
	}

	return responses, nil
}

func (s *userService) GetUserWithTransactions(id uint) (*user_responses.UserWithTransactionsResponse, error) {
	user, err := s.userRepo.GetUserWithTransactions(id)
	if err != nil {
		return nil, err
	}

	transactionResponses := make([]user_responses.TransactionResponse, len(user.Transactions))
	for i, transaction := range user.Transactions {
		transactionResponses[i] = user_responses.TransactionResponse{
			ID:       transaction.ID,
			BranchID: transaction.BranchID,
			Amount:   transaction.Amount,
			Date:     transaction.Date,
		}
	}

	return &user_responses.UserWithTransactionsResponse{
		ID:           user.ID,
		Name:         user.Name,
		Transactions: transactionResponses,
	}, nil
}

func (s *userService) GetUserWithRewards(id uint) (*user_responses.UserWithRewardsResponse, error) {
	user, err := s.userRepo.GetUserWithRewards(id)
	if err != nil {
		return nil, err
	}

	rewardResponses := make([]user_responses.RewardResponse, len(user.Rewards))
	for i, reward := range user.Rewards {
		rewardResponses[i] = user_responses.RewardResponse{
			ID:         reward.ID,
			MerchantID: reward.MerchantID,
			Type:       reward.Type,
			Amount:     reward.Amount,
		}
	}

	return &user_responses.UserWithRewardsResponse{
		ID:      user.ID,
		Name:    user.Name,
		Rewards: rewardResponses,
	}, nil
}
