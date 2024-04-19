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

	// router.Use(utils.ErrorHandler)
	// router.Use(utils.RequestLogger)
	// router.Use(utils.ResponseLogger)

	authRouter := router.Group(constant.BaseApiPath + "/auth")
	authRouter.POST("/register", utils.AuthenticationWithAuthorization([]string{AdminRole}), utils.RequestLogger, utils.ResponseLogger, handler.AddNewUser, utils.ErrorHandler)
	authRouter.POST("/login", utils.RequestLogger, utils.ResponseLogger, handler.Login, utils.ErrorHandler)

	userRouter := router.Group(constant.BaseApiPath + "/user")
	userRouter.POST("/get_all_active_user", utils.AuthenticationWithAuthorization([]string{AdminRole}), utils.RequestLogger, utils.ResponseLogger, handler.GetAllActiveUser, utils.ErrorHandler)
	userRouter.POST("/get_member_in_room", utils.AuthenticationWithAuthorization([]string{UserRole}), utils.RequestLogger, utils.ResponseLogger, handler.GetMemberInRoom, utils.ErrorHandler)
	userRouter.POST("/get_member_in_a_specific_room_code", utils.AuthenticationWithAuthorization([]string{AdminRole}), utils.RequestLogger, utils.ResponseLogger, handler.GetMemberInASpecificRoomCode, utils.ErrorHandler)

	expenseRouter := router.Group(constant.BaseApiPath + "/expense")
	expenseRouter.POST("/create_new_expense", utils.AuthenticationWithAuthorization([]string{UserRole}), utils.RequestLogger, utils.ResponseLogger, handler.AddNewExpense, utils.ErrorHandler)
	expenseRouter.POST("/get_list_of_expense", utils.AuthenticationWithAuthorization([]string{AdminRole, UserRole}), utils.RequestLogger, utils.ResponseLogger, handler.ListExpense, utils.ErrorHandler)
	expenseRouter.POST("/remove_expense", utils.AuthenticationWithAuthorization([]string{AdminRole}), utils.RequestLogger, utils.ResponseLogger, handler.RemoveExpense, utils.ErrorHandler)
	expenseRouter.POST("/soft_remove_expense", utils.AuthenticationWithAuthorization([]string{UserRole}), utils.RequestLogger, utils.ResponseLogger, handler.SoftRemoveExpense, utils.ErrorHandler)
	expenseRouter.POST("/active_expense", utils.AuthenticationWithAuthorization([]string{UserRole}), utils.RequestLogger, utils.ResponseLogger, handler.ActiveRemoveExpense, utils.ErrorHandler)

	debitRouter := router.Group(constant.BaseApiPath + "/debit")
	debitRouter.POST("/calculate", utils.AuthenticationWithAuthorization([]string{UserRole}), utils.RequestLogger, utils.ResponseLogger, handler.CalculateDebitOfUser, utils.ErrorHandler)

	roomRouter := router.Group(constant.BaseApiPath + "/room")
	roomRouter.POST("/add_new_room", utils.AuthenticationWithAuthorization([]string{AdminRole}), utils.RequestLogger, utils.ResponseLogger, handler.AddNewRoom, utils.ErrorHandler)
	roomRouter.POST("/get_list_of_rooms", utils.AuthenticationWithAuthorization([]string{AdminRole}), utils.RequestLogger, utils.ResponseLogger, handler.GetListOfRooms, utils.ErrorHandler)
	roomRouter.POST("/delete_room", utils.AuthenticationWithAuthorization([]string{AdminRole}), utils.RequestLogger, utils.ResponseLogger, handler.DeleteRoom, utils.ErrorHandler)
	roomRouter.POST("/edit_room_name", utils.AuthenticationWithAuthorization([]string{AdminRole}), utils.RequestLogger, utils.ResponseLogger, handler.EditRoomName, utils.ErrorHandler)
}
