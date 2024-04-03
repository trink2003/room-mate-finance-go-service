package payload

type UserRegisterRequestBodyValue struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginRequestBodyValue struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegisterRequestBody struct {
	Request UserRegisterRequestBodyValue `json:"request"`
}

type UserLoginRequestBody struct {
	Request UserLoginRequestBodyValue `json:"request"`
}

type UserLoginResponseBody struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Trace        string `json:"trace"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}
