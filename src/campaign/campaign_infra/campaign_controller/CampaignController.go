package campaign_controller

import (
	"loyalty-campaigns/src/campaign/campaign_app"
	"loyalty-campaigns/src/campaign/campaign_domain/campaign_structs/campaign_requests"
	"loyalty-campaigns/src/campaign/campaign_infra/campaign_repository"
	"loyalty-campaigns/src/common/configs"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type CampaignController struct {
	campaignService campaign_app.ICampaignService
}

var (
	campaignControllerInstance *CampaignController
	campaignControllerOnce     sync.Once
)

func NewCampaignController(router *gin.Engine) *CampaignController {
	campaignControllerOnce.Do(func() {
		campaignControllerInstance = &CampaignController{}
		db := configs.NewDBConnection().GetDB()
		campaignRepository := campaign_repository.NewGormCampaignRepository(db)
		campaignControllerInstance.campaignService = campaign_app.NewCampaignService(campaignRepository)
		campaignControllerInstance.setupCampaignRoutes(router)
	})
	return campaignControllerInstance
}

func (c *CampaignController) setupCampaignRoutes(router *gin.Engine) {
	campaignGroup := router.Group("/api/campaigns")
	{
		campaignGroup.POST("", c.CreateCampaign)
		campaignGroup.GET("/:id", c.GetCampaign)
		campaignGroup.PUT("/:id", c.UpdateCampaign)
		campaignGroup.DELETE("/:id", c.DeleteCampaign)
		campaignGroup.GET("", c.ListCampaigns)
		campaignGroup.GET("/active", c.GetActiveCampaigns)
	}
}

// CreateCampaign godoc
//	@Summary		Create a new campaign
//	@Description	Create a new loyalty campaign in the system
//	@Tags			campaigns
//	@Accept			json
//	@Produce		json
//	@Param			request	body		campaign_requests.CreateCampaignRequest	true	"Campaign creation request"
//	@Success		201		{object}	campaign_responses.CampaignResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/campaigns [post]
func (c *CampaignController) CreateCampaign(ctx *gin.Context) {
	var req campaign_requests.CreateCampaignRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.campaignService.CreateCampaign(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// GetCampaign godoc
//	@Summary		Get a campaign by ID
//	@Description	Get details of a specific campaign
//	@Tags			campaigns
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Campaign ID"
//	@Success		200	{object}	campaign_responses.CampaignResponse
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/api/campaigns/{id} [get]
func (c *CampaignController) GetCampaign(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	response, err := c.campaignService.GetCampaign(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateCampaign godoc
//	@Summary		Update a campaign
//	@Description	Update details of an existing campaign
//	@Tags			campaigns
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int										true	"Campaign ID"
//	@Param			request	body		campaign_requests.UpdateCampaignRequest	true	"Campaign update request"
//	@Success		200		{object}	campaign_responses.CampaignResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/campaigns/{id} [put]
func (c *CampaignController) UpdateCampaign(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req campaign_requests.UpdateCampaignRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.campaignService.UpdateCampaign(uint(id), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteCampaign godoc
//	@Summary		Delete a campaign
//	@Description	Delete an existing campaign from the system
//	@Tags			campaigns
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Campaign ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/campaigns/{id} [delete]
func (c *CampaignController) DeleteCampaign(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = c.campaignService.DeleteCampaign(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Campaign deleted successfully"})
}

// ListCampaigns godoc
//	@Summary		List all campaigns
//	@Description	Get a list of all campaigns in the system
//	@Tags			campaigns
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		campaign_responses.CampaignResponse
//	@Failure		500	{object}	map[string]string
//	@Router			/api/campaigns [get]
func (c *CampaignController) ListCampaigns(ctx *gin.Context) {
	responses, err := c.campaignService.ListCampaigns()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, responses)
}

// GetActiveCampaigns godoc
//	@Summary		Get active campaigns
//	@Description	Get a list of active campaigns for a specific merchant, branch, and date
//	@Tags			campaigns
//	@Accept			json
//	@Produce		json
//	@Param			merchantId	query		int		false	"Merchant ID"
//	@Param			branchId	query		int		false	"Branch ID"
//	@Param			date		query		string	false	"Date (RFC3339 format)"	Format(date-time)
//	@Success		200			{array}		campaign_responses.CampaignResponse
//	@Failure		500			{object}	map[string]string
//	@Router			/api/campaigns/active [get]
func (c *CampaignController) GetActiveCampaigns(ctx *gin.Context) {
	merchantID, _ := strconv.ParseUint(ctx.Query("merchantId"), 10, 32)
	branchID, _ := strconv.ParseUint(ctx.Query("branchId"), 10, 32)
	date, _ := time.Parse(time.RFC3339, ctx.Query("date"))

	responses, err := c.campaignService.GetActiveCampaigns(uint(merchantID), uint(branchID), date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, responses)
}
