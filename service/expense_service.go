package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"math/big"
	"net/http"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/model"
	"room-mate-finance-go-service/payload"
	"room-mate-finance-go-service/utils"
	"slices"
	"strconv"
	"time"
)

func (h *ExpenseHandler) AddNewExpense(c *gin.Context) {
	requestPayload := payload.ExpenseRequestBody{}

	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.Response{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["JSON_BINDING_ERROR"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["JSON_BINDING_ERROR"].ErrorMessage + " " + err.Error(),
		})
		return
	}

	if requestPayload.Request.Amount <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.Response{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["DATA_FORMAT_ERROR"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["DATA_FORMAT_ERROR"].ErrorMessage + " Amount need to be equal or greater than 0",
		})
		return
	}

	if len(requestPayload.Request.UserToPaid) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.Response{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["DATA_FORMAT_ERROR"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["DATA_FORMAT_ERROR"].ErrorMessage + " List of user need to pay must be not empty",
		})
		return
	}

	if requestPayload.Request.Purpose == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.Response{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["DATA_FORMAT_ERROR"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["DATA_FORMAT_ERROR"].ErrorMessage + " What is your purpose of this expense?",
		})
		return
	}

	currentUser, isCurrentUserExist := utils.GetCurrentUsername(c)

	ctx := context.Background()

	ctx = context.WithValue(ctx, "username", *currentUser)

	if isCurrentUserExist != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.Response{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["UNAUTHORIZED"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["UNAUTHORIZED"].ErrorMessage + " " + isCurrentUserExist.Error(),
		})
		return
	}

	boughtUser := model.Users{}

	h.DB.Where(
		&model.Users{
			BaseEntity: model.BaseEntity{
				Active: true,
			},
			Username: *currentUser,
		},
	).Find(&boughtUser)

	if boughtUser.BaseEntity.Id == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.Response{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["USER_NOT_EXISTED"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["USER_NOT_EXISTED"].ErrorMessage + " Who are you?",
		})
		return
	}

	var numberOfActiveUser int64 = 0

	h.DB. /*Clauses(clause.Locking{Strength: "UPDATE"}).*/
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
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.Response{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["INVALID_NUMBER_OF_USER"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["INVALID_NUMBER_OF_USER"].ErrorMessage,
		})
		return
	}

	if slices.Contains(requestPayload.Request.UserToPaid, boughtUser.BaseEntity.Id) {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.Response{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["INVALID_USER_TO_PAID_LIST"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["INVALID_USER_TO_PAID_LIST"].ErrorMessage,
		})
		return
	}
	log.Printf(
		constant.LogPattern,
		utils.GetTraceId(c),
		*currentUser,
		fmt.Sprintf("bought user is user with id: %s", strconv.FormatInt(boughtUser.BaseEntity.Id, 10)),
	)

	expense := model.ListOfExpenses{
		Purpose:        requestPayload.Request.Purpose,
		Amount:         requestPayload.Request.Amount,
		BoughtByUserID: boughtUser.BaseEntity.Id,
	}

	var equallyDividedAmount *big.Float

	// Calculate the divisor based on participation
	divisor := new(big.Float).SetInt64(int64(len(requestPayload.Request.UserToPaid)))
	if requestPayload.Request.IsParticipating {
		log.Printf(
			constant.LogPattern,
			utils.GetTraceId(c),
			*currentUser,
			"this user will participate to this expense",
		)
		divisor.Add(divisor, big.NewFloat(1))
	}

	requestAmount := new(big.Float).SetFloat64(requestPayload.Request.Amount)

	// Perform the division
	equallyDividedAmount = new(big.Float).Quo(requestAmount, divisor)
	// equallyDividedAmount.SetPrec(2) // Set precision to 2 decimal places

	scaledNumber := new(big.Float).Mul(equallyDividedAmount, big.NewFloat(100))

	roundedNumber, _ := scaledNumber.Int(nil)

	finalEquallyDividedAmount := new(big.Float).Quo(new(big.Float).SetInt(roundedNumber), big.NewFloat(100))
	log.Printf(
		constant.LogPattern,
		utils.GetTraceId(c),
		*currentUser,
		fmt.Sprintf(
			"amount info:\n    - finalEquallyDividedAmount: %s\n\t- requestAmount: %s\n\t- equallyDividedAmount: %s\n\t- divisor: %s\n\t- finalEquallyDividedAmount: %s",
			finalEquallyDividedAmount,
			requestAmount,
			equallyDividedAmount,
			divisor,
			finalEquallyDividedAmount,
		),
	)

	expenseTransactionError := h.DB.Transaction(
		func(tx *gorm.DB) error {
			if saveNewExpenseErr := SaveNewExpense(tx, &expense, ctx); saveNewExpenseErr.Error != nil {
				return saveNewExpenseErr.Error
			}
			paidAmount, _ := finalEquallyDividedAmount.Float64()
			for _, user := range allActiveUserInList {
				debitOfCurrentUser := model.DebitUser{
					UserToPaidID:     user.BaseEntity.Id,
					PaidToUserID:     boughtUser.BaseEntity.Id,
					ListOfExpensesID: expense.BaseEntity.Id,
					Amount:           paidAmount,
				}
				if saveNewDebitUserErr := SaveNewDebitUser(tx, &debitOfCurrentUser, ctx); saveNewDebitUserErr.Error != nil {
					return saveNewDebitUserErr.Error
				}
			}
			return nil
		},
	)
	if expenseTransactionError != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.Response{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["QUERY_ERROR"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["QUERY_ERROR"].ErrorMessage + " " + expenseTransactionError.Error(),
		})
		return
	}

	var savedExpense []model.ListOfExpenses

	h.DB.Preload("Users").Preload("DebitUser").Where(
		model.ListOfExpenses{
			BaseEntity: model.BaseEntity{
				Id: expense.BaseEntity.Id,
			},
		},
	).Find(&savedExpense)

	c.JSON(http.StatusOK, &payload.Response{
		Trace:        utils.GetTraceId(c),
		ErrorCode:    constant.ErrorConstant["SUCCESS"].ErrorCode,
		ErrorMessage: constant.ErrorConstant["SUCCESS"].ErrorMessage,
		Response:     savedExpense,
	})

}

func (h *ExpenseHandler) RemoveExpense(c *gin.Context) {

	requestPayload := payload.RemoveExpenseBody{}
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.Response{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["JSON_BINDING_ERROR"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["JSON_BINDING_ERROR"].ErrorMessage + " " + err.Error(),
		})
		return
	}

	var expense model.ListOfExpenses

	h.DB.Preload("Users").Preload("DebitUser").Where(
		model.ListOfExpenses{
			BaseEntity: model.BaseEntity{
				Id: requestPayload.Request,
			},
		},
	).Find(&expense)

	if expense.BaseEntity.Id == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.Response{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["EXPENSE_NOT_SUCCESS"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["EXPENSE_NOT_SUCCESS"].ErrorMessage,
		})
		return
	}

	transactionResult := h.DB.Transaction(func(tx *gorm.DB) error {
		debitUserRemoveResult := tx.Delete(expense.DebitUser)
		if debitUserRemoveResult.Error != nil {
			return debitUserRemoveResult.Error
		}
		expenseRemoveResult := tx.Delete(expense)
		if expenseRemoveResult.Error != nil {
			return expenseRemoveResult.Error
		}
		return nil
	})
	if transactionResult != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.Response{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["QUERY_ERROR"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["QUERY_ERROR"].ErrorMessage + " " + transactionResult.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &payload.Response{
		Trace:        utils.GetTraceId(c),
		ErrorCode:    constant.ErrorConstant["SUCCESS"].ErrorCode,
		ErrorMessage: constant.ErrorConstant["SUCCESS"].ErrorMessage,
	})
}

func (h *ExpenseHandler) ListExpense(c *gin.Context) {

	requestPayload := payload.PageRequestBody{}
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.Response{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["JSON_BINDING_ERROR"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["JSON_BINDING_ERROR"].ErrorMessage + " " + err.Error(),
		})
		return
	}

	if requestPayload.Request.PageSize == 0 || requestPayload.Request.PageNumber == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, &payload.Response{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["DATA_FORMAT_ERROR"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["DATA_FORMAT_ERROR"].ErrorMessage + " " + "Page number and page size can not be 0",
		})
		return
	}

	limit := requestPayload.Request.PageNumber - 1
	offset := requestPayload.Request.PageSize * limit

	var expense []model.ListOfExpenses

	var total int64

	h.DB.Preload("Users").Preload("DebitUser").Count(&total)

	h.DB.Preload("Users").Preload("DebitUser").Limit(limit).Offset(offset).Find(&expense)

	c.JSON(http.StatusOK, &payload.Response{
		Trace:        utils.GetTraceId(c),
		ErrorCode:    constant.ErrorConstant["SUCCESS"].ErrorCode,
		ErrorMessage: constant.ErrorConstant["SUCCESS"].ErrorMessage,
		Response:     expense,
	})
}

func SaveNewExpense(db *gorm.DB, model *model.ListOfExpenses, ctx context.Context) *gorm.DB {
	model.BaseEntity.Active = true
	model.BaseEntity.CreatedAt = time.Now()
	model.BaseEntity.UpdatedAt = time.Now()
	model.BaseEntity.CreatedBy = ctx.Value("username").(string)
	model.BaseEntity.UpdatedBy = ctx.Value("username").(string)

	return db.Create(model)
}

func SaveNewDebitUser(db *gorm.DB, model *model.DebitUser, ctx context.Context) *gorm.DB {
	model.BaseEntity.Active = true
	model.BaseEntity.CreatedAt = time.Now()
	model.BaseEntity.UpdatedAt = time.Now()
	model.BaseEntity.CreatedBy = ctx.Value("username").(string)
	model.BaseEntity.UpdatedBy = ctx.Value("username").(string)

	return db.Create(model)
}
