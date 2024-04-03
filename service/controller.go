package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

type AuthHandler struct {
	DB *gorm.DB
}

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	userHandler := &UserHandler{
		DB: db,
	}

	authHandler := &AuthHandler{
		DB: db,
	}

	router.Use(ErrorHandler)
	router.Use(RequestLogger)
	router.Use(ResponseLogger)

	authRouter := router.Group("/auth")
	authRouter.POST("/register", authHandler.AddNewUser)
	authRouter.POST("/login", authHandler.Login)

	userRouter := router.Group("/user")
	userRouter.POST("/get_all_active_user", Authentication, userHandler.GetUsers)
}
