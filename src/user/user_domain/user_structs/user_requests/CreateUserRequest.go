package user_requests

type CreateUserRequest struct {
	Name string `json:"name" binding:"required"`
}
