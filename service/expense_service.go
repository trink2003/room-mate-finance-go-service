package service

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm/clause"
	"net/http"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/model"
	"room-mate-finance-go-service/payload"
	"room-mate-finance-go-service/utils"
)

func (h *ExpenseHandler) AddNewExpense(c *gin.Context) {
	requestPayload := payload.ExpenseRequestBody{}

	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.ErrorResponse{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["JSON_BINDING_ERROR"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["JSON_BINDING_ERROR"].ErrorMessage + " " + err.Error(),
		})
		return
	}

	if requestPayload.Request.Amount <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.ErrorResponse{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["DATA_FORMAT_ERROR"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["DATA_FORMAT_ERROR"].ErrorMessage + " Amount need to be equal or greater than 0",
		})
		return
	}

	if len(requestPayload.Request.UserToPaid) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.ErrorResponse{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["DATA_FORMAT_ERROR"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["DATA_FORMAT_ERROR"].ErrorMessage + " List of user need to pay must be not empty",
		})
		return
	}

	if requestPayload.Request.Purpose == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.ErrorResponse{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["DATA_FORMAT_ERROR"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["DATA_FORMAT_ERROR"].ErrorMessage + " What is your purpose of this expense?",
		})
		return
	}

	currentUser, isCurrentUserExist := c.Get("auth")

	if isCurrentUserExist == false {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.ErrorResponse{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["UNAUTHORIZED"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["UNAUTHORIZED"].ErrorMessage + " Who are you?",
		})
		return
	}

	claim := currentUser.(jwt.MapClaims)

	boughtUser := model.Users{}

	h.DB.Where(
		&model.Users{
			BaseEntity: model.BaseEntity{
				Active: true,
			},
			Username: claim["sub"].(string),
		},
	).Find(&boughtUser)

	if boughtUser.BaseEntity.Id == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.ErrorResponse{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["USER_NOT_EXISTED"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["USER_NOT_EXISTED"].ErrorMessage + " Who are you?",
		})
		return
	}

	var numberOfActiveUser int64 = 0

	h.DB.Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&model.Users{}).
		Where(
			h.DB.
				Where(
					&model.Users{
						BaseEntity: model.BaseEntity{
							Active: true,
						},
					},
				).
				Where(
					"id NOT IN ?",
					[]int64{boughtUser.BaseEntity.Id},
				),
		).
		Count(&numberOfActiveUser)

	var allActiveUserInList []model.Users

	h.DB.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(
			h.DB.
				Where(
					&model.Users{
						BaseEntity: model.BaseEntity{
							Active: true,
						},
					},
				).
				Where(
					"id IN ?",
					requestPayload.Request.UserToPaid,
				),
		).
		Find(&allActiveUserInList)

	if numberOfActiveUser < 2 || len(allActiveUserInList) < len(requestPayload.Request.UserToPaid) {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.ErrorResponse{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["INVALID_NUMBER_OF_USER"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["INVALID_NUMBER_OF_USER"].ErrorMessage,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": "ok",
	})

}
