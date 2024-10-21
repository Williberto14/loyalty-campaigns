package merchant_controller

import (
	"loyalty-campaigns/src/common/configs"
	"loyalty-campaigns/src/merchant/merchant_app"
	"loyalty-campaigns/src/merchant/merchant_domain/merchant_structs/merchant_requests"
	"loyalty-campaigns/src/merchant/merchant_infra/merchant_repository"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type MerchantController struct {
	service merchant_app.IMerchantService
}

var merchantControllerInstance *MerchantController
var merchantControllerOnce sync.Once

func NewMerchantController(router *gin.Engine) *MerchantController {
	merchantControllerOnce.Do(func() {
		merchantControllerInstance = &MerchantController{}
		db := configs.NewDBConnection().GetDB()
		merchantRepository := merchant_repository.NewGormMerchantRepository(db)
		merchantControllerInstance.service = merchant_app.NewMerchantService(merchantRepository)
		merchantControllerInstance.setupMerchantRoutes(router)
	})
	return merchantControllerInstance
}

func (c *MerchantController) setupMerchantRoutes(router *gin.Engine) {
	merchantGroup := router.Group("/api/merchants")
	{
		merchantGroup.POST("", c.CreateMerchant)
		merchantGroup.GET("", c.ListMerchants)
		merchantGroup.GET("/:id", c.GetMerchant)
		merchantGroup.PUT("/:id", c.UpdateMerchant)
		merchantGroup.DELETE("/:id", c.DeleteMerchant)
	}
}

// CreateMerchant godoc
//	@Summary		Create a new merchant
//	@Description	Create a new merchant in the system
//	@Tags			merchants
//	@Accept			json
//	@Produce		json
//	@Param			request	body		merchant_requests.CreateMerchantRequest	true	"Merchant creation request"
//	@Success		201		{object}	merchant_responses.MerchantResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/merchants [post]
func (c *MerchantController) CreateMerchant(ctx *gin.Context) {
	var req merchant_requests.CreateMerchantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.CreateMerchant(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// ListMerchants godoc
//	@Summary		List all merchants
//	@Description	Get a list of all merchants in the system
//	@Tags			merchants
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		merchant_responses.MerchantResponse
//	@Failure		500	{object}	map[string]string
//	@Router			/api/merchants [get]
func (c *MerchantController) ListMerchants(ctx *gin.Context) {
	response, err := c.service.ListMerchants()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetMerchant godoc
//	@Summary		Get a merchant by ID
//	@Description	Get details of a specific merchant
//	@Tags			merchants
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Merchant ID"
//	@Success		200	{object}	merchant_responses.MerchantResponse
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/merchants/{id} [get]
func (c *MerchantController) GetMerchant(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.GetMerchant(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateMerchant godoc
//	@Summary		Update a merchant
//	@Description	Update details of an existing merchant
//	@Tags			merchants
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int										true	"Merchant ID"
//	@Param			request	body		merchant_requests.UpdateMerchantRequest	true	"Merchant update request"
//	@Success		200		{object}	merchant_responses.MerchantResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/merchants/{id} [put]
func (c *MerchantController) UpdateMerchant(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req merchant_requests.UpdateMerchantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.UpdateMerchant(uint(id), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteMerchant godoc
//	@Summary		Delete a merchant
//	@Description	Delete an existing merchant from the system
//	@Tags			merchants
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Merchant ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/merchants/{id} [delete]
func (c *MerchantController) DeleteMerchant(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.service.DeleteMerchant(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
