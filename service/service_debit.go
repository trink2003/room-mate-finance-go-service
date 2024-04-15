package service

import (
	"context"
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
	PaidToUser int64   `json:"paidToUser" gorm:"column:paid_to_user;"`
	UserToPaid int64   `json:"userToPaid" gorm:"column:user_to_paid;"`
	Amount     float64 `json:"amount" gorm:"column:amount;"`
}

func (h DebitHandler) CalculateDebitOfUser(c *gin.Context) {

	currentUser, isCurrentUserExist := utils.GetCurrentUsername(c)

	ctx := context.Background()

	ctx = context.WithValue(ctx, "username", *currentUser)
	ctx = context.WithValue(ctx, "traceId", utils.GetTraceId(c))

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
	requestPayload := payload.CalculateDebitRequestBody{}

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

	currentUsername, isCurrentUserExist := utils.GetCurrentUsername(c)

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

	loc, timeLoadLocationErr := time.LoadLocation("Asia/Ho_Chi_Minh")
	if timeLoadLocationErr != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.InternalFailure,
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
	errorEnum := constant.Success

	transactionResult := h.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if requestPayload.Request.IsStatisticsAccordingToCurrentUser {

			currentUserModel := model.Users{}

			tx.
				Where(
					"username = ? AND active is true", *currentUsername,
				).
				Find(&currentUserModel)

			if currentUserModel.BaseEntity.Id == 0 {
				errorEnum = constant.UserNotExisted
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
				currentUserModel.BaseEntity.Id,
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
			utils.ReturnResponse(
				c,
				constant.QueryError,
				nil,
				transactionResult.Error(),
			),
		)
		return
	}
	if errorEnum.ErrorCode != 0 {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				errorEnum,
				nil,
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		utils.ReturnResponse(
			c,
			constant.Success,
			calculateResult,
		),
	)

}
