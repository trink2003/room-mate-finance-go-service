package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"room-mate-finance-go-service/utils"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	userHandler := &UserHandler{DB: db}
	authHandler := &AuthHandler{DB: db}
	expenseHandler := &ExpenseHandler{DB: db}
	debitHandler := &DebitHandler{DB: db}
	roomHandler := &RoomHandler{DB: db}

	router.Use(utils.ErrorHandler)
	router.Use(utils.RequestLogger)
	router.Use(utils.ResponseLogger)

	authRouter := router.Group("/roommate/api/v1/auth")
	authRouter.POST("/register", authHandler.AddNewUser)
	authRouter.POST("/login", authHandler.Login)

	userRouter := router.Group("/roommate/api/v1/user")
	userRouter.POST("/get_all_active_user", utils.AuthenticationWithAuthorization([]string{"ADMIN"}), userHandler.GetUsers)

	expenseRouter := router.Group("/roommate/api/v1/expense")
	expenseRouter.POST("/create_new_expense", utils.AuthenticationWithAuthorization([]string{"USER"}), expenseHandler.AddNewExpense)
	expenseRouter.POST("/get_list_of_expense", utils.AuthenticationWithAuthorization([]string{"ADMIN", "USER"}), expenseHandler.ListExpense)
	expenseRouter.POST("/remove_expense", utils.AuthenticationWithAuthorization([]string{"USER"}), expenseHandler.RemoveExpense)
	expenseRouter.POST("/soft_remove_expense", utils.AuthenticationWithAuthorization([]string{"USER"}), expenseHandler.SoftRemoveExpense)
	expenseRouter.POST("/active_expense", utils.AuthenticationWithAuthorization([]string{"USER"}), expenseHandler.ActiveRemoveExpense)

	debitRouter := router.Group("/roommate/api/v1/debit")
	debitRouter.POST("/calculate", utils.AuthenticationWithAuthorization([]string{"USER"}), debitHandler.CalculateDebitOfUser)

	roomRouter := router.Group("/roommate/api/v1/room")
	roomRouter.POST("/add_new_room", utils.AuthenticationWithAuthorization([]string{"ADMIN"}), roomHandler.AddNewRoom)
}
