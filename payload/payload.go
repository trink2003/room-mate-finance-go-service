package payload

type UserRegisterRequestBodyValue struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegisterRequestBody struct {
	Request UserRegisterRequestBodyValue `json:"request"`
}
