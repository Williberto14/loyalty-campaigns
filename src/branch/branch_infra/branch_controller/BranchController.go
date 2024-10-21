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

func (c *BranchController) ListBranches(ctx *gin.Context) {
	responses, err := c.branchService.ListBranches()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, responses)
}

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
