package transaction_controller

import (
	"loyalty-campaigns/src/common/configs"
	"loyalty-campaigns/src/transaction/transaction_app"
	"loyalty-campaigns/src/transaction/transaction_domain/transaction_structs/transaction_requests"
	"loyalty-campaigns/src/transaction/transaction_infra/transaction_repository"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	transactionService transaction_app.ITransactionService
}

var (
	transactionControllerInstance *TransactionController
	transactionControllerOnce     sync.Once
)

func NewTransactionController(router *gin.Engine) *TransactionController {
	transactionControllerOnce.Do(func() {
		transactionControllerInstance = &TransactionController{}
		db := configs.NewDBConnection().GetDB()
		transactionRepository := transaction_repository.NewGormTransactionRepository(db)
		transactionControllerInstance.transactionService = transaction_app.NewTransactionService(transactionRepository)
		transactionControllerInstance.setupTransactionRoutes(router)
	})
	return transactionControllerInstance
}

func (c *TransactionController) setupTransactionRoutes(router *gin.Engine) {
	transactionGroup := router.Group("/api/transactions")
	{
		transactionGroup.POST("", c.CreateTransaction)
		transactionGroup.GET("/:id", c.GetTransaction)
		transactionGroup.GET("/users/:userID", c.ListTransactionsByUser)
		transactionGroup.GET("/users/:userID/total-amount", c.GetTotalAmountByUserAndDateRange)
	}
}

func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	var req transaction_requests.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.transactionService.CreateTransaction(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *TransactionController) GetTransaction(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	response, err := c.transactionService.GetTransaction(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *TransactionController) ListTransactionsByUser(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("userID"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	responses, err := c.transactionService.ListTransactionsByUser(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, responses)
}

func (c *TransactionController) GetTotalAmountByUserAndDateRange(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("userID"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	startDate, err := time.Parse(time.RFC3339, ctx.Query("startDate"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, ctx.Query("endDate"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date"})
		return
	}

	totalAmount, err := c.transactionService.GetTotalAmountByUserAndDateRange(uint(userID), startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"totalAmount": totalAmount})
}
