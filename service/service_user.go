package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/big"
	"net/http"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/model"
	"room-mate-finance-go-service/payload"
	"room-mate-finance-go-service/utils"
)

func (h UserHandler) GetUsers(c *gin.Context) {

	currentUser, isCurrentUserExist := utils.GetCurrentUsername(c)

	if isCurrentUserExist != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			utils.ReturnResponse(
				c,
				constant.Unauthorized,
				nil,
				isCurrentUserExist.Error(),
			),
		)
		return
	}

	ctx := context.Background()

	ctx = context.WithValue(ctx, "username", *currentUser)
	ctx = context.WithValue(ctx, "traceId", utils.GetTraceId(c))

	requestPayload := payload.PageRequestBody{}
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.JsonBindingError,
				nil,
				err.Error(),
			),
		)
		return
	}

	limit := requestPayload.Request.PageSize
	offset := requestPayload.Request.PageSize * (requestPayload.Request.PageNumber - 1)
	var user []model.Users
	var total int64 = 0

	transactionResult := h.DB.Transaction(func(tx *gorm.DB) error {

		tx.WithContext(ctx).Model(&model.Users{}).
			Where(
				tx.Where(
					model.Users{
						BaseEntity: model.BaseEntity{
							Active: utils.GetPointerOfAnyValue(true),
						},
					},
				).Where("active is not null"),
			).Count(&total)

		tx.WithContext(ctx).Limit(limit).
			Offset(offset).
			Order(utils.SortMapToString(requestPayload.Request.Sort)).
			Where(
				tx.Where(
					model.Users{
						BaseEntity: model.BaseEntity{
							Active: utils.GetPointerOfAnyValue(true),
						},
					},
				).Where("active is not null"),
			).Find(&user)
		return nil
	})
	if transactionResult != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.QueryError,
				nil,
				transactionResult.Error(),
			),
		)
		return
	}

	totalRecordsSelected := len(user)

	if totalRecordsSelected == 0 {
		c.JSON(
			http.StatusOK,
			utils.ReturnPageResponse(
				c,
				constant.Success,
				0,
				0,
				user,
			),
		)
		return
	}

	/*
		if result := h.DB.Where("active is not null AND active is true ORDER BY id DESC").Find(&user); result.Error != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &payload.Response{
				Trace:        utils.GetTraceId(c),
				ErrorCode:    constant.QueryError.ErrorCode,
				ErrorMessage: constant.QueryError.ErrorMessage + result.Error.Error(),
			})
			return
		}
	*/

	totalBigNumber := new(big.Float).SetInt64(total)
	totalRecordsSelectedBigNumber := new(big.Float).SetInt64(int64(totalRecordsSelected))
	totalPage := new(big.Float).Quo(totalRecordsSelectedBigNumber, totalBigNumber)
	utils.RoundHalfUpBigFloat(totalPage)
	totalPageInt, _ := totalPage.Int(nil)

	c.JSON(
		http.StatusOK,
		utils.ReturnPageResponse(
			c,
			constant.Success,
			total,
			totalPageInt.Int64(),
			user,
		),
	)
}
