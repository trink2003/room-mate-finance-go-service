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

type ExpenseRequestBodyValue struct {
	Purpose         string  `json:"purpose" binding:"required"`
	Amount          float64 `json:"amount" binding:"required"`
	IsParticipating bool    `json:"isParticipating"`
	UserToPaid      []int64 `json:"userToPaid" binding:"required"`
}

type ExpenseRequestBody struct {
	Request ExpenseRequestBodyValue `json:"request" binding:"required"`
}

type RemoveExpenseBody struct {
	Request int64 `json:"request"`
}

type PageRequestBodyValue struct {
	PageNumber int               `json:"pageNumber"`
	PageSize   int               `json:"pageSize"`
	Sort       map[string]string `json:"sort"`
}

type PageRequestBody struct {
	Request PageRequestBodyValue `json:"request"`
}

type CalculateDebitRequestBodyValue struct {
	IsStatisticsAccordingToCurrentUser bool `json:"isStatisticsAccordingToCurrentUser"`
}

type CalculateDebitRequestBody struct {
	Request CalculateDebitRequestBodyValue `json:"request"`
}

type Response struct {
	Trace        string `json:"trace"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Response     any    `json:"response"`
}

type PageResponse struct {
	Trace        string `json:"trace"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	TotalElement int64  `json:"totalElement"`
	TotalPage    int64  `json:"totalPage"`
	Response     any    `json:"response"`
}
