package payload

import "room-mate-finance-go-service/model"

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

type ListUsersResponseBody struct {
	Users []model.Users `json:"listOfUsers"`
}

type ErrorResponse struct {
	Trace        string `json:"trace"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}
