package payload

import "room-mate-finance-go-service/model"

type AddMemberToARoomBodyValue struct {
	RoomCode string `json:"roomCode" binding:"required"`
	UserId   int64  `json:"userId" binding:"required"`
}

type AddMemberToARoomBody struct {
	Request AddMemberToARoomBodyValue `json:"request"`
}

type RemoveMemberInARoomBodyValue struct {
	RoomCode string `json:"roomCode" binding:"required"`
	UserId   int64  `json:"userId" binding:"required"`
}

type RemoveMemberInARoomBody struct {
	Request RemoveMemberInARoomBodyValue `json:"request"`
}

type GetMemberInARoomRequestBodyValue struct {
	RoomCode string `json:"roomCode" binding:"required"`
}

type GetMemberInARoomRequestBody struct {
	Request GetMemberInARoomRequestBodyValue `json:"request"`
}

type EditRoomNameRequestBodyValue struct {
	RoomCode string `json:"roomCode" binding:"required"`
	RoomName string `json:"roomName" binding:"required"`
}

type EditRoomNameRequestBody struct {
	Request EditRoomNameRequestBodyValue `json:"request"`
}

type DeleteRoomRequestBodyValue struct {
	RoomCode string `json:"roomCode" binding:"required"`
}

type DeleteRoomRequestBody struct {
	Request DeleteRoomRequestBodyValue `json:"request"`
}

type AddNewRoomRequestBodyValue struct {
	RoomCode string `json:"roomCode" binding:"required"`
	RoomName string `json:"roomName"`
}

type AddNewRoomRequestBody struct {
	Request AddNewRoomRequestBodyValue `json:"request"`
}

type UserRegisterRequestBodyValue struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RoomCode string `json:"roomCode" binding:"required"`
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
