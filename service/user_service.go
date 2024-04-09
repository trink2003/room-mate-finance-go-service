package service

import (
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

	requestPayload := payload.PageRequestBody{}
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				c,
				constant.ErrorConstant["JSON_BINDING_ERROR"],
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

		tx.Model(&model.Users{}).
			Where(
				tx.Where(
					model.Users{
						BaseEntity: model.BaseEntity{
							Active: utils.GetPointerOfAnyValue(true),
						},
					},
				).Where("active is not null"),
			).Count(&total)

		tx.Limit(limit).
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
			ReturnResponse(
				c,
				constant.ErrorConstant["QUERY_ERROR"],
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
			ReturnPageResponse(
				c,
				constant.ErrorConstant["SUCCESS"],
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
				ErrorCode:    constant.ErrorConstant["QUERY_ERROR"].ErrorCode,
				ErrorMessage: constant.ErrorConstant["QUERY_ERROR"].ErrorMessage + result.Error.Error(),
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
		ReturnPageResponse(
			c,
			constant.ErrorConstant["SUCCESS"],
			total,
			totalPageInt.Int64(),
			user,
		),
	)
}
