package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/model"
	"room-mate-finance-go-service/payload"
	"room-mate-finance-go-service/utils"
	"time"
	_ "time/tzdata" // Must import when getting current time with timezone
)

type CalculateResult struct {
	PaidToUser int64   `gorm:"column:paid_to_user;"`
	UserToPaid int64   `gorm:"column:user_to_paid;"`
	Amount     float64 `gorm:"column:amount;"`
}

func (h DebitHandler) CalculateDebitOfUser(c *gin.Context) {
	requestPayload := payload.CalculateDebitRequestBody{}

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

	currentUsername, isCurrentUserExist := utils.GetCurrentUsername(c)

	if isCurrentUserExist != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			ReturnResponse(
				c,
				constant.ErrorConstant["UNAUTHORIZED"],
				nil,
				isCurrentUserExist.Error(),
			),
		)
		return
	}

	loc, timeLoadLocationErr := time.LoadLocation("Asia/Ho_Chi_Minh")
	if timeLoadLocationErr != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			ReturnResponse(
				c,
				constant.ErrorConstant["INTERNAL_FAILURE"],
				nil,
				timeLoadLocationErr.Error(),
			),
		)
		return
	}
	currentTimestamp := time.Now().In(loc)

	firstOfMonth := utils.BeginningOfMonth(currentTimestamp)
	lastOfMonth := utils.EndOfMonth(currentTimestamp)

	var calculateResult []CalculateResult
	errorEnum := constant.ErrorConstant["SUCCESS"]

	transactionResult := h.DB.Transaction(func(tx *gorm.DB) error {

		if requestPayload.Request.IsStatisticsAccordingToCurrentUser {

			currentUser := model.Users{}

			tx.
				Where(
					"username = ? AND active is true", *currentUsername,
				).
				Find(&currentUser)

			if currentUser.BaseEntity.Id == 0 {
				errorEnum = constant.ErrorConstant["USER_NOT_EXISTED"]
			}

			tx.Raw(
				`
					select
						du.paid_to_user,
						du.user_to_paid,
						sum(du.amount) as "amount"
					from
						debit_user du
						left join list_of_expenses loe on loe.id = du.expense
					where
						du.paid_to_user = ?
						and du.created_at between to_timestamp(?, 'YYYY-MM-DD HH24:MI:SS')
						and to_timestamp(?, 'YYYY-MM-DD HH24:MI:SS')
						and loe.active is true
						and du.active is true
					group by
						du.paid_to_user,
						du.user_to_paid
				`,
				currentUser.BaseEntity.Id,
				firstOfMonth.Format(constant.YyyyMmDdHhMmSsFormat),
				lastOfMonth.Format(constant.YyyyMmDdHhMmSsFormat),
			).Scan(&calculateResult)
		} else {
			tx.Raw(
				`
					select
						du.paid_to_user,
						du.user_to_paid,
						sum(du.amount) as "amount"
					from
						debit_user du
						left join list_of_expenses loe on loe.id = du.expense
					where
						du.created_at between to_timestamp(?, 'YYYY-MM-DD HH24:MI:SS')
						and to_timestamp(?, 'YYYY-MM-DD HH24:MI:SS')
						and loe.active is true
						and du.active is true
					group by
						du.paid_to_user,
						du.user_to_paid
				`,
				firstOfMonth.Format(constant.YyyyMmDdHhMmSsFormat),
				lastOfMonth.Format(constant.YyyyMmDdHhMmSsFormat),
			).Scan(&calculateResult)
		}

		return nil
	})

	if transactionResult != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				c,
				constant.ErrorConstant["QUERY_ERROR"],
				nil,
				transactionResult.Error(),
			),
		)
		return
	}
	if errorEnum.ErrorCode != 0 {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				c,
				errorEnum,
				nil,
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		ReturnResponse(
			c,
			constant.ErrorConstant["SUCCESS"],
			calculateResult,
		),
	)

}
