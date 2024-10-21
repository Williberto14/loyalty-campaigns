package branch_app

import (
	"loyalty-campaigns/src/branch/branch_domain/branch_ports"
	"loyalty-campaigns/src/branch/branch_domain/branch_structs/branch_requests"
	"loyalty-campaigns/src/branch/branch_domain/branch_structs/branch_responses"
	"loyalty-campaigns/src/common/models"
	"sync"
)

type IBranchService interface {
	CreateBranch(req branch_requests.CreateBranchRequest) (*branch_responses.BranchResponse, error)
	GetBranch(id uint) (*branch_responses.BranchResponse, error)
	UpdateBranch(id uint, req branch_requests.UpdateBranchRequest) (*branch_responses.BranchResponse, error)
	DeleteBranch(id uint) error
	ListBranches() ([]branch_responses.BranchResponse, error)
	GetBranchesByMerchant(merchantID uint) ([]branch_responses.BranchResponse, error)
	GetBranchWithCampaigns(id uint) (*branch_responses.BranchWithCampaignsResponse, error)
}

type branchService struct {
	branchRepo branch_ports.IBranchRepository
}

var (
	branchServiceInstance *branchService
	branchServiceOnce     sync.Once
)

func NewBranchService(branchRepo branch_ports.IBranchRepository) IBranchService {
	branchServiceOnce.Do(func() {
		branchServiceInstance = &branchService{
			branchRepo: branchRepo,
		}
	})
	return branchServiceInstance
}

func (s *branchService) CreateBranch(req branch_requests.CreateBranchRequest) (*branch_responses.BranchResponse, error) {
	branch := &models.Branch{
		Name:       req.Name,
		MerchantID: req.MerchantID,
	}

	err := s.branchRepo.Create(branch)
	if err != nil {
		return nil, err
	}

	return &branch_responses.BranchResponse{
		ID:         branch.ID,
		Name:       branch.Name,
		MerchantID: branch.MerchantID,
	}, nil
}

func (s *branchService) GetBranch(id uint) (*branch_responses.BranchResponse, error) {
	branch, err := s.branchRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &branch_responses.BranchResponse{
		ID:         branch.ID,
		Name:       branch.Name,
		MerchantID: branch.MerchantID,
	}, nil
}

func (s *branchService) UpdateBranch(id uint, req branch_requests.UpdateBranchRequest) (*branch_responses.BranchResponse, error) {
	branch, err := s.branchRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	branch.Name = req.Name
	if req.MerchantID != 0 {
		branch.MerchantID = req.MerchantID
	}

	err = s.branchRepo.Update(branch)
	if err != nil {
		return nil, err
	}

	return &branch_responses.BranchResponse{
		ID:         branch.ID,
		Name:       branch.Name,
		MerchantID: branch.MerchantID,
	}, nil
}

func (s *branchService) DeleteBranch(id uint) error {
	return s.branchRepo.Delete(id)
}

func (s *branchService) ListBranches() ([]branch_responses.BranchResponse, error) {
	branches, err := s.branchRepo.List()
	if err != nil {
		return nil, err
	}

	var responses []branch_responses.BranchResponse
	for _, branch := range branches {
		responses = append(responses, branch_responses.BranchResponse{
			ID:         branch.ID,
			Name:       branch.Name,
			MerchantID: branch.MerchantID,
		})
	}

	return responses, nil
}

func (s *branchService) GetBranchesByMerchant(merchantID uint) ([]branch_responses.BranchResponse, error) {
	branches, err := s.branchRepo.GetByMerchantID(merchantID)
	if err != nil {
		return nil, err
	}

	var responses []branch_responses.BranchResponse
	for _, branch := range branches {
		responses = append(responses, branch_responses.BranchResponse{
			ID:         branch.ID,
			Name:       branch.Name,
			MerchantID: branch.MerchantID,
		})
	}

	return responses, nil
}

func (s *branchService) GetBranchWithCampaigns(id uint) (*branch_responses.BranchWithCampaignsResponse, error) {
	branch, err := s.branchRepo.GetBranchWithCampaigns(id)
	if err != nil {
		return nil, err
	}

	campaignResponses := make([]branch_responses.CampaignResponse, len(branch.Campaigns))
	for i, campaign := range branch.Campaigns {
		campaignResponses[i] = branch_responses.CampaignResponse{
			ID:        campaign.ID,
			StartDate: campaign.StartDate,
			EndDate:   campaign.EndDate,
			Type:      campaign.Type,
			Value:     campaign.Value,
			MinAmount: campaign.MinAmount,
		}
	}

	return &branch_responses.BranchWithCampaignsResponse{
		ID:         branch.ID,
		Name:       branch.Name,
		MerchantID: branch.MerchantID,
		Campaigns:  campaignResponses,
	}, nil
}
