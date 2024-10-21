package merchant_app

import (
	"loyalty-campaigns/src/common/models"
	"loyalty-campaigns/src/merchant/merchant_domain/merchant_ports"
	"loyalty-campaigns/src/merchant/merchant_domain/merchant_structs/merchant_requests"
	"loyalty-campaigns/src/merchant/merchant_domain/merchant_structs/merchant_responses"
)

type IMerchantService interface {
	CreateMerchant(req merchant_requests.CreateMerchantRequest) (*merchant_responses.MerchantResponse, error)
	ListMerchants() ([]*merchant_responses.MerchantResponse, error)
	GetMerchant(id uint) (*merchant_responses.MerchantResponse, error)
	UpdateMerchant(id uint, req merchant_requests.UpdateMerchantRequest) (*merchant_responses.MerchantResponse, error)
	DeleteMerchant(id uint) error
}

type MerchantService struct {
	repo merchant_ports.IMerchantRepository
}

func NewMerchantService(repo merchant_ports.IMerchantRepository) IMerchantService {
	return &MerchantService{repo: repo}
}

func (s *MerchantService) CreateMerchant(req merchant_requests.CreateMerchantRequest) (*merchant_responses.MerchantResponse, error) {
	merchant := &models.Merchant{
		Name:             req.Name,
		ConversionFactor: req.ConversionFactor,
	}

	err := s.repo.Create(merchant)
	if err != nil {
		return nil, err
	}

	return &merchant_responses.MerchantResponse{
		ID:               merchant.ID,
		Name:             merchant.Name,
		ConversionFactor: merchant.ConversionFactor,
		CreatedAt:        merchant.CreatedAt,
		UpdatedAt:        merchant.UpdatedAt,
	}, nil
}

func (s *MerchantService) ListMerchants() ([]*merchant_responses.MerchantResponse, error) {
	merchants, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	var response []*merchant_responses.MerchantResponse
	for _, merchant := range merchants {
		response = append(response, &merchant_responses.MerchantResponse{
			ID:               merchant.ID,
			Name:             merchant.Name,
			ConversionFactor: merchant.ConversionFactor,
			CreatedAt:        merchant.CreatedAt,
			UpdatedAt:        merchant.UpdatedAt,
		})
	}

	return response, nil
}

func (s *MerchantService) GetMerchant(id uint) (*merchant_responses.MerchantResponse, error) {
	merchant, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &merchant_responses.MerchantResponse{
		ID:               merchant.ID,
		Name:             merchant.Name,
		ConversionFactor: merchant.ConversionFactor,
		CreatedAt:        merchant.CreatedAt,
		UpdatedAt:        merchant.UpdatedAt,
	}, nil
}

func (s *MerchantService) UpdateMerchant(id uint, req merchant_requests.UpdateMerchantRequest) (*merchant_responses.MerchantResponse, error) {
	merchant, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	merchant.Name = req.Name
	merchant.ConversionFactor = req.ConversionFactor

	err = s.repo.Update(merchant)
	if err != nil {
		return nil, err
	}

	return &merchant_responses.MerchantResponse{
		ID:               merchant.ID,
		Name:             merchant.Name,
		ConversionFactor: merchant.ConversionFactor,
		CreatedAt:        merchant.CreatedAt,
		UpdatedAt:        merchant.UpdatedAt,
	}, nil
}

func (s *MerchantService) DeleteMerchant(id uint) error {
	return s.repo.Delete(id)
}
