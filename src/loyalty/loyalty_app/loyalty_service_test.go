package loyalty_app_test

import (
	"loyalty-campaigns/src/campaign/campaign_domain/campaign_structs/campaign_requests"
	"loyalty-campaigns/src/campaign/campaign_domain/campaign_structs/campaign_responses"
	"loyalty-campaigns/src/loyalty/loyalty_app"
	"loyalty-campaigns/src/merchant/merchant_domain/merchant_structs/merchant_requests"
	"loyalty-campaigns/src/merchant/merchant_domain/merchant_structs/merchant_responses"
	"loyalty-campaigns/src/reward/reward_domain/reward_structs/reward_requests"
	"loyalty-campaigns/src/reward/reward_domain/reward_structs/reward_responses"
	"loyalty-campaigns/src/transaction/transaction_domain/transaction_structs/transaction_requests"
	"loyalty-campaigns/src/transaction/transaction_domain/transaction_structs/transaction_responses"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("LoyaltyService", func() {
	var (
		loyaltyService  loyalty_app.ILoyaltyService
		mockTransaction *mockTransactionService
		mockMerchant    *mockMerchantService
		mockCampaign    *mockCampaignService
		mockReward      *mockRewardService
		userID          uint
		merchantID      uint
		branchID        uint
		amount          float64
		date            time.Time
	)

	BeforeEach(func() {
		mockTransaction = new(mockTransactionService)
		mockMerchant = new(mockMerchantService)
		mockCampaign = new(mockCampaignService)
		mockReward = new(mockRewardService)

		loyaltyService = loyalty_app.NewLoyaltyService(
			mockTransaction,
			mockCampaign,
			mockReward,
			mockMerchant,
		)

		userID = 1
		merchantID = 2
		branchID = 3
		amount = 100.0
		date = time.Now()
	})

	Describe("ProcessTransaction", func() {
		Context("When there are no active campaigns", func() {
			BeforeEach(func() {
				mockTransaction.On("CreateTransaction", mock.AnythingOfType("transaction_requests.CreateTransactionRequest")).Return(&transaction_responses.TransactionResponse{}, nil)
				mockMerchant.On("GetMerchant", merchantID).Return(&merchant_responses.MerchantResponse{
					ID:                merchantID,
					ConversionFactor:  0.1,
					DefaultRewardType: "points",
				}, nil)
				mockCampaign.On("GetActiveCampaigns", merchantID, &branchID, mock.AnythingOfType("time.Time")).Return([]campaign_responses.CampaignResponse{}, nil)
				mockReward.On("CreateReward", mock.AnythingOfType("reward_requests.CreateRewardRequest")).Return(&reward_responses.RewardResponse{}, nil)
			})

			It("should process the transaction and create a default reward", func() {
				err := loyaltyService.ProcessTransaction(userID, merchantID, branchID, amount, date)

				Expect(err).To(BeNil())
				mockTransaction.AssertExpectations(GinkgoT())
				mockMerchant.AssertExpectations(GinkgoT())
				mockCampaign.AssertExpectations(GinkgoT())
				mockReward.AssertExpectations(GinkgoT())

				mockReward.AssertCalled(GinkgoT(), "CreateReward", reward_requests.CreateRewardRequest{
					UserID:     userID,
					MerchantID: merchantID,
					Type:       "points",
					Amount:     10.0, // 100 * 0.1
				})
			})
		})

		Context("When there is an active campaign", func() {
			BeforeEach(func() {
				mockTransaction.On("CreateTransaction", mock.AnythingOfType("transaction_requests.CreateTransactionRequest")).Return(&transaction_responses.TransactionResponse{}, nil)
				mockMerchant.On("GetMerchant", merchantID).Return(&merchant_responses.MerchantResponse{
					ID:                merchantID,
					ConversionFactor:  0.1,
					DefaultRewardType: "points",
				}, nil)
				mockCampaign.On("GetActiveCampaigns", merchantID, &branchID, mock.AnythingOfType("time.Time")).Return([]campaign_responses.CampaignResponse{
					{
						Type:  "points",
						Value: 2.0,
					},
				}, nil)
				mockReward.On("CreateReward", mock.AnythingOfType("reward_requests.CreateRewardRequest")).Return(&reward_responses.RewardResponse{}, nil)
			})

			It("should process the transaction and create a campaign reward", func() {
				err := loyaltyService.ProcessTransaction(userID, merchantID, branchID, amount, date)

				Expect(err).To(BeNil())
				mockReward.AssertCalled(GinkgoT(), "CreateReward", reward_requests.CreateRewardRequest{
					UserID:     userID,
					MerchantID: merchantID,
					Type:       "points",
					Amount:     20.0, // (100 * 0.1) * 2
				})
			})
		})

	})

	Describe("RedeemRewards", func() {
		Context("When user has sufficient rewards", func() {
			BeforeEach(func() {
				mockReward.On("ListRewardsByUser", userID).Return([]reward_responses.RewardResponse{
					{
						MerchantID: merchantID,
						Type:       "points",
						Amount:     50.0,
					},
				}, nil)
				// Cambia esta línea para que coincida con el monto que estás probando
				mockReward.On("DeductRewards", userID, merchantID, float64(30), "points").Return(nil)
			})

			It("should redeem the rewards successfully", func() {
				err := loyaltyService.RedeemRewards(userID, merchantID, 30.0, "points")

				Expect(err).To(BeNil())
				mockReward.AssertExpectations(GinkgoT())
			})
		})

		Context("When user has insufficient rewards", func() {
			BeforeEach(func() {
				mockReward.On("ListRewardsByUser", userID).Return([]reward_responses.RewardResponse{
					{
						MerchantID: merchantID,
						Type:       "points",
						Amount:     20.0,
					},
				}, nil)
			})

			It("should return an error", func() {
				err := loyaltyService.RedeemRewards(userID, merchantID, 30.0, "points")

				Expect(err).To(MatchError("insufficient rewards"))
				mockReward.AssertNotCalled(GinkgoT(), "DeductRewards")
			})
		})

	})
})

// Mock implementations
type mockTransactionService struct {
	mock.Mock
}

func (m *mockTransactionService) CreateTransaction(req transaction_requests.CreateTransactionRequest) (*transaction_responses.TransactionResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*transaction_responses.TransactionResponse), args.Error(1)
}

func (m *mockTransactionService) GetTotalAmountByUserAndDateRange(userID uint, startDate, endDate time.Time) (float64, error) {
	args := m.Called(userID, startDate, endDate)
	return args.Get(0).(float64), args.Error(1)
}

func (m *mockTransactionService) GetTransaction(id uint) (*transaction_responses.TransactionResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*transaction_responses.TransactionResponse), args.Error(1)
}

func (m *mockTransactionService) GetTransactionsByDateRange(startDate, endDate time.Time) ([]transaction_responses.TransactionResponse, error) {
	args := m.Called(startDate, endDate)
	return args.Get(0).([]transaction_responses.TransactionResponse), args.Error(1)
}

func (m *mockTransactionService) ListTransactionsByBranch(branchID uint) ([]transaction_responses.TransactionResponse, error) {
	args := m.Called(branchID)
	return args.Get(0).([]transaction_responses.TransactionResponse), args.Error(1)
}

func (m *mockTransactionService) ListTransactionsByUser(userID uint) ([]transaction_responses.TransactionResponse, error) {
	args := m.Called(userID)
	return args.Get(0).([]transaction_responses.TransactionResponse), args.Error(1)
}

type mockCampaignService struct {
	mock.Mock
}

func (m *mockCampaignService) GetActiveCampaigns(merchantID uint, branchID *uint, date time.Time) ([]campaign_responses.CampaignResponse, error) {
	args := m.Called(merchantID, branchID, date)
	return args.Get(0).([]campaign_responses.CampaignResponse), args.Error(1)
}

func (m *mockCampaignService) CreateCampaign(req campaign_requests.CreateCampaignRequest) (*campaign_responses.CampaignResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*campaign_responses.CampaignResponse), args.Error(1)
}

func (m *mockCampaignService) DeleteCampaign(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockCampaignService) GetCampaign(id uint) (*campaign_responses.CampaignResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*campaign_responses.CampaignResponse), args.Error(1)
}

func (m *mockCampaignService) ListCampaigns() ([]campaign_responses.CampaignResponse, error) {
	args := m.Called()
	return args.Get(0).([]campaign_responses.CampaignResponse), args.Error(1)
}

func (m *mockCampaignService) UpdateCampaign(id uint, req campaign_requests.UpdateCampaignRequest) (*campaign_responses.CampaignResponse, error) {
	args := m.Called(id, req)
	return args.Get(0).(*campaign_responses.CampaignResponse), args.Error(1)
}

type mockRewardService struct {
	mock.Mock
}

func (m *mockRewardService) CreateReward(req reward_requests.CreateRewardRequest) (*reward_responses.RewardResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*reward_responses.RewardResponse), args.Error(1)
}

func (m *mockRewardService) ListRewardsByUser(userID uint) ([]reward_responses.RewardResponse, error) {
	args := m.Called(userID)
	return args.Get(0).([]reward_responses.RewardResponse), args.Error(1)
}

func (m *mockRewardService) DeductRewards(userID, merchantID uint, amount float64, rewardType string) error {
	args := m.Called(userID, merchantID, amount, rewardType)
	return args.Error(0)
}

func (m *mockRewardService) GetReward(id uint) (*reward_responses.RewardResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*reward_responses.RewardResponse), args.Error(1)
}

func (m *mockRewardService) GetTotalRewardsByUser(id uint) (*reward_responses.TotalRewardsResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*reward_responses.TotalRewardsResponse), args.Error(1)
}

type mockMerchantService struct {
	mock.Mock
}

func (m *mockMerchantService) GetMerchant(id uint) (*merchant_responses.MerchantResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*merchant_responses.MerchantResponse), args.Error(1)
}

func (m *mockMerchantService) CreateMerchant(req merchant_requests.CreateMerchantRequest) (*merchant_responses.MerchantResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*merchant_responses.MerchantResponse), args.Error(1)
}

func (m *mockMerchantService) UpdateMerchant(id uint, req merchant_requests.UpdateMerchantRequest) (*merchant_responses.MerchantResponse, error) {
	args := m.Called(id, req)
	return args.Get(0).(*merchant_responses.MerchantResponse), args.Error(1)
}

func (m *mockMerchantService) DeleteMerchant(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockMerchantService) ListMerchants() ([]*merchant_responses.MerchantResponse, error) {
	args := m.Called()
	return args.Get(0).([]*merchant_responses.MerchantResponse), args.Error(1)
}
