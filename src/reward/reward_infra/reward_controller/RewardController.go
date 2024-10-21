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

// CreateReward godoc
//	@Summary		Create a new reward
//	@Description	Create a new reward in the system
//	@Tags			rewards
//	@Accept			json
//	@Produce		json
//	@Param			request	body		reward_requests.CreateRewardRequest	true	"Reward creation request"
//	@Success		201		{object}	reward_responses.RewardResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/rewards [post]
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

// GetReward godoc
//	@Summary		Get a reward by ID
//	@Description	Get details of a specific reward
//	@Tags			rewards
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Reward ID"
//	@Success		200	{object}	reward_responses.RewardResponse
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/api/rewards/{id} [get]
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

// ListRewardsByUser godoc
//	@Summary		List rewards for a specific user
//	@Description	Get a list of all rewards for a given user
//	@Tags			rewards
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{array}		reward_responses.RewardResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/rewards/user/{userID} [get]
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

// GetTotalRewardsByUser godoc
//	@Summary		Get total rewards for a user
//	@Description	Calculate the total rewards (points and cashback) for a specific user
//	@Tags			rewards
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	reward_responses.TotalRewardsResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/rewards/user/{userID}/total [get]
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
