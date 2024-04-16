package service

import (
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	AdminRole = "ADMIN"
	UserRole  = "USER"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	handler := &Handler{DB: db}

	router.Use(utils.ErrorHandler)
	router.Use(utils.RequestLogger)
	router.Use(utils.ResponseLogger)

	authRouter := router.Group(constant.BaseApiPath + "/auth")
	authRouter.POST("/register", utils.AuthenticationWithAuthorization([]string{AdminRole}), handler.AddNewUser)
	authRouter.POST("/login", handler.Login)

	userRouter := router.Group(constant.BaseApiPath + "/user")
	userRouter.POST("/get_all_active_user", utils.AuthenticationWithAuthorization([]string{AdminRole}), handler.GetUsers)
	userRouter.POST("/get_member_in_room", utils.AuthenticationWithAuthorization([]string{UserRole}), handler.GetMemberInRoom)
	userRouter.POST("/get_member_in_a_specific_room_code", utils.AuthenticationWithAuthorization([]string{AdminRole}), handler.GetMemberInASpecificRoomCode)

	expenseRouter := router.Group(constant.BaseApiPath + "/expense")
	expenseRouter.POST("/create_new_expense", utils.AuthenticationWithAuthorization([]string{UserRole}), handler.AddNewExpense)
	expenseRouter.POST("/get_list_of_expense", utils.AuthenticationWithAuthorization([]string{AdminRole, UserRole}), handler.ListExpense)
	expenseRouter.POST("/remove_expense", utils.AuthenticationWithAuthorization([]string{AdminRole}), handler.RemoveExpense)
	expenseRouter.POST("/soft_remove_expense", utils.AuthenticationWithAuthorization([]string{UserRole}), handler.SoftRemoveExpense)
	expenseRouter.POST("/active_expense", utils.AuthenticationWithAuthorization([]string{UserRole}), handler.ActiveRemoveExpense)

	debitRouter := router.Group(constant.BaseApiPath + "/debit")
	debitRouter.POST("/calculate", utils.AuthenticationWithAuthorization([]string{UserRole}), handler.CalculateDebitOfUser)

	roomRouter := router.Group(constant.BaseApiPath + "/room")
	roomRouter.POST("/add_new_room", utils.AuthenticationWithAuthorization([]string{AdminRole}), handler.AddNewRoom)
	roomRouter.POST("/get_list_of_rooms", utils.AuthenticationWithAuthorization([]string{AdminRole}), handler.GetListOfRooms)
	roomRouter.POST("/delete_room", utils.AuthenticationWithAuthorization([]string{AdminRole}), handler.DeleteRoom)
	roomRouter.POST("/edit_room_name", utils.AuthenticationWithAuthorization([]string{AdminRole}), handler.EditRoomName)
}
