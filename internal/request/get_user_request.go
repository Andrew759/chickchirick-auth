package request

type GetUserRequest struct {
	UserUuid string `json:"user_uuid" binding:"required"`
	//TODO: тут жесткие требования по длине пароля от 8 символов.
	// разобраться с кейсом, когда пароль не задан
	Password string `json:"password" binding:"required,min=8,max=256"`
}
