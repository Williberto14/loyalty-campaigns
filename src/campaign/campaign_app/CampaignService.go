package campaign_app

import (
	"loyalty-campaigns/src/campaign/campaign_domain/campaign_ports"
	"loyalty-campaigns/src/campaign/campaign_domain/campaign_structs/campaign_requests"
	"loyalty-campaigns/src/campaign/campaign_domain/campaign_structs/campaign_responses"
	"loyalty-campaigns/src/common/models"
	"loyalty-campaigns/src/common/utils"
	"sync"
	"time"
)

type ICampaignService interface {
	CreateCampaign(req campaign_requests.CreateCampaignRequest) (*campaign_responses.CampaignResponse, error)
	GetCampaign(id uint) (*campaign_responses.CampaignResponse, error)
	UpdateCampaign(id uint, req campaign_requests.UpdateCampaignRequest) (*campaign_responses.CampaignResponse, error)
	DeleteCampaign(id uint) error
	ListCampaigns() ([]campaign_responses.CampaignResponse, error)
	GetActiveCampaigns(merchantID, branchID uint, date time.Time) ([]campaign_responses.CampaignResponse, error)
}

type campaignService struct {
	campaignRepo campaign_ports.ICampaignRepository
	logger       utils.ILogger
}

var (
	campaignServiceInstance *campaignService
	campaignServiceOnce     sync.Once
)

func NewCampaignService(campaignRepo campaign_ports.ICampaignRepository) ICampaignService {
	campaignServiceOnce.Do(func() {
		campaignServiceInstance = &campaignService{
			campaignRepo: campaignRepo,
			logger:       utils.NewLogger(),
		}
	})
	return campaignServiceInstance
}

func (s *campaignService) CreateCampaign(req campaign_requests.CreateCampaignRequest) (*campaign_responses.CampaignResponse, error) {
	campaign := &models.Campaign{
		MerchantID: req.MerchantID,
		BranchID:   req.BranchID,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Type:       req.Type,
		Value:      req.Value,
		MinAmount:  req.MinAmount,
	}

	err := s.campaignRepo.Create(campaign)
	if err != nil {
		s.logger.Error("Error al crear campaña", err)
		return nil, err
	}

	return &campaign_responses.CampaignResponse{
		ID:         campaign.ID,
		MerchantID: campaign.MerchantID,
		BranchID:   campaign.BranchID,
		StartDate:  campaign.StartDate,
		EndDate:    campaign.EndDate,
		Type:       campaign.Type,
		Value:      campaign.Value,
		MinAmount:  campaign.MinAmount,
	}, nil
}

func (s *campaignService) GetCampaign(id uint) (*campaign_responses.CampaignResponse, error) {
	campaign, err := s.campaignRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Error al obtener campaña", err)
		return nil, err
	}

	return &campaign_responses.CampaignResponse{
		ID:         campaign.ID,
		MerchantID: campaign.MerchantID,
		BranchID:   campaign.BranchID,
		StartDate:  campaign.StartDate,
		EndDate:    campaign.EndDate,
		Type:       campaign.Type,
		Value:      campaign.Value,
		MinAmount:  campaign.MinAmount,
	}, nil
}

func (s *campaignService) UpdateCampaign(id uint, req campaign_requests.UpdateCampaignRequest) (*campaign_responses.CampaignResponse, error) {
	campaign, err := s.campaignRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Error al obtener campaña para actualizar", err)
		return nil, err
	}

	campaign.StartDate = req.StartDate
	campaign.EndDate = req.EndDate
	campaign.Type = req.Type
	campaign.Value = req.Value
	campaign.MinAmount = req.MinAmount

	err = s.campaignRepo.Update(campaign)
	if err != nil {
		s.logger.Error("Error al actualizar campaña", err)
		return nil, err
	}

	return &campaign_responses.CampaignResponse{
		ID:         campaign.ID,
		MerchantID: campaign.MerchantID,
		BranchID:   campaign.BranchID,
		StartDate:  campaign.StartDate,
		EndDate:    campaign.EndDate,
		Type:       campaign.Type,
		Value:      campaign.Value,
		MinAmount:  campaign.MinAmount,
	}, nil
}

func (s *campaignService) DeleteCampaign(id uint) error {
	err := s.campaignRepo.Delete(id)
	if err != nil {
		s.logger.Error("Error al eliminar campaña", err)
	}
	return err
}

func (s *campaignService) ListCampaigns() ([]campaign_responses.CampaignResponse, error) {
	campaigns, err := s.campaignRepo.List()
	if err != nil {
		s.logger.Error("Error al listar campañas", err)
		return nil, err
	}

	var responses []campaign_responses.CampaignResponse
	for _, campaign := range campaigns {
		responses = append(responses, campaign_responses.CampaignResponse{
			ID:         campaign.ID,
			MerchantID: campaign.MerchantID,
			BranchID:   campaign.BranchID,
			StartDate:  campaign.StartDate,
			EndDate:    campaign.EndDate,
			Type:       campaign.Type,
			Value:      campaign.Value,
			MinAmount:  campaign.MinAmount,
		})
	}

	return responses, nil
}

func (s *campaignService) GetActiveCampaigns(merchantID, branchID uint, date time.Time) ([]campaign_responses.CampaignResponse, error) {
	campaigns, err := s.campaignRepo.GetActiveCampaigns(merchantID, branchID, date)
	if err != nil {
		s.logger.Error("Error al obtener campañas activas", err)
		return nil, err
	}

	var responses []campaign_responses.CampaignResponse
	for _, campaign := range campaigns {
		responses = append(responses, campaign_responses.CampaignResponse{
			ID:         campaign.ID,
			MerchantID: campaign.MerchantID,
			BranchID:   campaign.BranchID,
			StartDate:  campaign.StartDate,
			EndDate:    campaign.EndDate,
			Type:       campaign.Type,
			Value:      campaign.Value,
			MinAmount:  campaign.MinAmount,
		})
	}

	return responses, nil
}
