package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"room-mate-finance-go-service/utils"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	handler := &Handler{DB: db}

	router.Use(utils.ErrorHandler)
	router.Use(utils.RequestLogger)
	router.Use(utils.ResponseLogger)

	var basePath = "/roommate/api/v1"

	authRouter := router.Group(basePath + "/auth")
	authRouter.POST("/register", utils.AuthenticationWithAuthorization([]string{"ADMIN"}), handler.AddNewUser)
	authRouter.POST("/login", handler.Login)

	userRouter := router.Group(basePath + "/user")
	userRouter.POST("/get_all_active_user", utils.AuthenticationWithAuthorization([]string{"ADMIN"}), handler.GetUsers)
	userRouter.POST("/get_member_in_room", utils.AuthenticationWithAuthorization([]string{"USER"}), handler.GetMemberInRoom)
	userRouter.POST("/get_member_in_a_specific_room_code", utils.AuthenticationWithAuthorization([]string{"ADMIN"}), handler.GetMemberInASpecificRoomCode)

	expenseRouter := router.Group(basePath + "/expense")
	expenseRouter.POST("/create_new_expense", utils.AuthenticationWithAuthorization([]string{"USER"}), handler.AddNewExpense)
	expenseRouter.POST("/get_list_of_expense", utils.AuthenticationWithAuthorization([]string{"ADMIN", "USER"}), handler.ListExpense)
	expenseRouter.POST("/remove_expense", utils.AuthenticationWithAuthorization([]string{"USER"}), handler.RemoveExpense)
	expenseRouter.POST("/soft_remove_expense", utils.AuthenticationWithAuthorization([]string{"USER"}), handler.SoftRemoveExpense)
	expenseRouter.POST("/active_expense", utils.AuthenticationWithAuthorization([]string{"USER"}), handler.ActiveRemoveExpense)

	debitRouter := router.Group(basePath + "/debit")
	debitRouter.POST("/calculate", utils.AuthenticationWithAuthorization([]string{"USER"}), handler.CalculateDebitOfUser)

	roomRouter := router.Group(basePath + "/room")
	roomRouter.POST("/add_new_room", utils.AuthenticationWithAuthorization([]string{"ADMIN"}), handler.AddNewRoom)
	roomRouter.POST("/get_list_of_rooms", utils.AuthenticationWithAuthorization([]string{"ADMIN"}), handler.GetListOfRooms)
	roomRouter.POST("/delete_room", utils.AuthenticationWithAuthorization([]string{"ADMIN"}), handler.DeleteRoom)
	roomRouter.POST("/edit_room_name", utils.AuthenticationWithAuthorization([]string{"ADMIN"}), handler.EditRoomName)
}
