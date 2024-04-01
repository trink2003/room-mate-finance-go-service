package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"room-mate-finance-go-service/model"
)

func (h handler) GetUser(ginContext *gin.Context) {
	var user []model.User

	if result := h.DB.Find(&user); result.Error != nil {
		ginContext.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	ginContext.JSON(http.StatusOK, &user)
}
