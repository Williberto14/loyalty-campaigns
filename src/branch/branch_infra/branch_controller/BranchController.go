package branch_controller

import (
	"loyalty-campaigns/src/branch/branch_app"
	"loyalty-campaigns/src/branch/branch_domain/branch_structs/branch_requests"
	"loyalty-campaigns/src/branch/branch_infra/branch_repository"
	"loyalty-campaigns/src/common/configs"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type BranchController struct {
	branchService branch_app.IBranchService
}

var (
	branchControllerInstance *BranchController
	branchControllerOnce     sync.Once
)

func NewBranchController(router *gin.Engine) *BranchController {
	branchControllerOnce.Do(func() {
		branchControllerInstance = &BranchController{}
		db := configs.NewDBConnection().GetDB()
		branchRepository := branch_repository.NewGormBranchRepository(db)
		branchControllerInstance.branchService = branch_app.NewBranchService(branchRepository)
		branchControllerInstance.setupBranchRoutes(router)
	})
	return branchControllerInstance
}

func (c *BranchController) setupBranchRoutes(router *gin.Engine) {
	branchGroup := router.Group("/api/branches")
	{
		branchGroup.POST("", c.CreateBranch)
		branchGroup.GET("/:id", c.GetBranch)
		branchGroup.PUT("/:id", c.UpdateBranch)
		branchGroup.DELETE("/:id", c.DeleteBranch)
		branchGroup.GET("", c.ListBranches)
		branchGroup.GET("/merchant/:merchantID", c.GetBranchesByMerchant)
		branchGroup.GET("/:id/campaigns", c.GetBranchWithCampaigns)
	}
}

// CreateBranch godoc
//	@Summary		Create a new branch
//	@Description	Create a new branch for a merchant in the system
//	@Tags			branches
//	@Accept			json
//	@Produce		json
//	@Param			request	body		branch_requests.CreateBranchRequest	true	"Branch creation request"
//	@Success		201		{object}	branch_responses.BranchResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/branches [post]
func (c *BranchController) CreateBranch(ctx *gin.Context) {
	var req branch_requests.CreateBranchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.branchService.CreateBranch(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// GetBranch godoc
//	@Summary		Get a branch by ID
//	@Description	Get details of a specific branch
//	@Tags			branches
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Branch ID"
//	@Success		200	{object}	branch_responses.BranchResponse
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/api/branches/{id} [get]
func (c *BranchController) GetBranch(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	response, err := c.branchService.GetBranch(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateBranch godoc
//	@Summary		Update a branch
//	@Description	Update details of an existing branch
//	@Tags			branches
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int									true	"Branch ID"
//	@Param			request	body		branch_requests.UpdateBranchRequest	true	"Branch update request"
//	@Success		200		{object}	branch_responses.BranchResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/branches/{id} [put]
func (c *BranchController) UpdateBranch(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req branch_requests.UpdateBranchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.branchService.UpdateBranch(uint(id), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteBranch godoc
//	@Summary		Delete a branch
//	@Description	Delete an existing branch from the system
//	@Tags			branches
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Branch ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/branches/{id} [delete]
func (c *BranchController) DeleteBranch(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = c.branchService.DeleteBranch(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Branch deleted successfully"})
}

// ListBranches godoc
//	@Summary		List all branches
//	@Description	Get a list of all branches in the system
//	@Tags			branches
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		branch_responses.BranchResponse
//	@Failure		500	{object}	map[string]string
//	@Router			/api/branches [get]
func (c *BranchController) ListBranches(ctx *gin.Context) {
	responses, err := c.branchService.ListBranches()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, responses)
}

// GetBranchesByMerchant godoc
//	@Summary		Get branches by merchant
//	@Description	Get a list of branches for a specific merchant
//	@Tags			branches
//	@Accept			json
//	@Produce		json
//	@Param			merchantID	path		int	true	"Merchant ID"
//	@Success		200			{array}		branch_responses.BranchResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/api/branches/merchant/{merchantID} [get]
func (c *BranchController) GetBranchesByMerchant(ctx *gin.Context) {
	merchantID, err := strconv.ParseUint(ctx.Param("merchantID"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid merchant ID"})
		return
	}

	responses, err := c.branchService.GetBranchesByMerchant(uint(merchantID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, responses)
}

// GetBranchWithCampaigns godoc
//	@Summary		Get a branch with its campaigns
//	@Description	Get details of a specific branch along with its associated campaigns
//	@Tags			branches
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Branch ID"
//	@Success		200	{object}	branch_responses.BranchWithCampaignsResponse
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/branches/{id}/campaigns [get]
func (c *BranchController) GetBranchWithCampaigns(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	response, err := c.branchService.GetBranchWithCampaigns(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
