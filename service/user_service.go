package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/model"
	"room-mate-finance-go-service/payload"
	"room-mate-finance-go-service/utils"
)

func (h UserHandler) GetUsers(ginContext *gin.Context) {

	var user []model.Users

	if result := h.DB.Where("active is not null AND active is true ORDER BY id DESC").Find(&user); result.Error != nil {
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, &payload.Response{
			Trace:        utils.GetTraceId(ginContext),
			ErrorCode:    constant.ErrorConstant["QUERY_ERROR"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["QUERY_ERROR"].ErrorMessage + result.Error.Error(),
		})
		return
	}
	ginContext.JSON(
		http.StatusOK,
		ReturnResponse(
			ginContext,
			constant.ErrorConstant["SUCCESS"],
			user,
		),
	)
}
