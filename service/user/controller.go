package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	h := &UserHandler{
		DB: db,
	}

	routes := router.Group("/auth")
	routes.POST("/register", h.AddNewUser)
}
