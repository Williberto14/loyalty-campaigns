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

func (c *MerchantController) ListMerchants(ctx *gin.Context) {
	response, err := c.service.ListMerchants()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

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
