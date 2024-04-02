package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"room-mate-finance-go-service/model"
)

type UserRegisterRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h handler) GetUser(ginContext *gin.Context) {
	var user []model.User

	if result := h.DB.Find(&user); result.Error != nil {
		ginContext.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	ginContext.JSON(http.StatusOK, &user)
}

func (h handler) AddNewUser(ginContext *gin.Context) {

}
