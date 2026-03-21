package request

// TODO: дубль логики с GetUserRequest
type CreateTokenRequest struct {
	UserUuid string `json:"user_uuid" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=256"`
}
