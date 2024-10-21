package user_controller

import (
	"loyalty-campaigns/src/common/configs"
	"loyalty-campaigns/src/user/user_app"
	"loyalty-campaigns/src/user/user_domain/user_structs/user_requests"
	"loyalty-campaigns/src/user/user_infra/user_repository"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService user_app.IUserService
}

var (
	userControllerInstance *UserController
	userControllerOnce     sync.Once
)

func NewUserController(router *gin.Engine) *UserController {
	userControllerOnce.Do(func() {
		userControllerInstance = &UserController{}
		db := configs.NewDBConnection().GetDB()
		userRepository := user_repository.NewGormUserRepository(db)
		userControllerInstance.userService = user_app.NewUserService(userRepository)
		userControllerInstance.setupUserRoutes(router)
	})
	return userControllerInstance
}

func (c *UserController) setupUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/users")
	{
		userGroup.POST("", c.CreateUser)
		userGroup.GET("/:id", c.GetUser)
		userGroup.PUT("/:id", c.UpdateUser)
		userGroup.DELETE("/:id", c.DeleteUser)
		userGroup.GET("", c.ListUsers)
		userGroup.GET("/:id/transactions", c.GetUserWithTransactions)
		userGroup.GET("/:id/rewards", c.GetUserWithRewards)
	}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req user_requests.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.userService.CreateUser(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *UserController) GetUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	response, err := c.userService.GetUser(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req user_requests.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.userService.UpdateUser(uint(id), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = c.userService.DeleteUser(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (c *UserController) ListUsers(ctx *gin.Context) {
	responses, err := c.userService.ListUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, responses)
}

func (c *UserController) GetUserWithTransactions(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	response, err := c.userService.GetUserWithTransactions(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *UserController) GetUserWithRewards(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	response, err := c.userService.GetUserWithRewards(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
