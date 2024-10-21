package user_requests

type UpdateUserRequest struct {
	Name string `json:"name" binding:"required"`
}
