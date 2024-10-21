package reward_controller

import (
	"loyalty-campaigns/src/common/configs"
	"loyalty-campaigns/src/reward/reward_app"
	"loyalty-campaigns/src/reward/reward_domain/reward_structs/reward_requests"
	"loyalty-campaigns/src/reward/reward_infra/reward_repository"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type RewardController struct {
	rewardService reward_app.IRewardService
}

var (
	rewardControllerInstance *RewardController
	rewardControllerOnce     sync.Once
)

func NewRewardController(router *gin.Engine) *RewardController {
	rewardControllerOnce.Do(func() {
		rewardControllerInstance = &RewardController{}
		db := configs.NewDBConnection().GetDB()
		rewardRepository := reward_repository.NewGormRewardRepository(db)
		rewardControllerInstance.rewardService = reward_app.NewRewardService(rewardRepository)
		rewardControllerInstance.setupRewardRoutes(router)
	})
	return rewardControllerInstance
}

func (c *RewardController) setupRewardRoutes(router *gin.Engine) {
	rewardGroup := router.Group("/api/rewards")
	{
		rewardGroup.POST("", c.CreateReward)
		rewardGroup.GET("/:id", c.GetReward)
		rewardGroup.GET("/users/:userID", c.ListRewardsByUser)
		rewardGroup.GET("/users/:userID/total-amount", c.GetTotalRewardsByUser)
	}
}

func (c *RewardController) CreateReward(ctx *gin.Context) {
	var req reward_requests.CreateRewardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.rewardService.CreateReward(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *RewardController) GetReward(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	response, err := c.rewardService.GetReward(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Reward not found"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *RewardController) ListRewardsByUser(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("userID"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	responses, err := c.rewardService.ListRewardsByUser(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, responses)
}

func (c *RewardController) GetTotalRewardsByUser(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("userID"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	response, err := c.rewardService.GetTotalRewardsByUser(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
