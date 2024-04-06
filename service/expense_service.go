package service

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/model"
	"room-mate-finance-go-service/payload"
	"room-mate-finance-go-service/utils"
	"slices"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (h *ExpenseHandler) AddNewExpense(c *gin.Context) {
	requestPayload := payload.ExpenseRequestBody{}

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

	if requestPayload.Request.Amount <= 0 {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				c,
				constant.ErrorConstant["DATA_FORMAT_ERROR"],
				nil,
				"Amount need to be equal or greater than 0",
			),
		)
		return
	}

	if len(requestPayload.Request.UserToPaid) == 0 {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				c,
				constant.ErrorConstant["DATA_FORMAT_ERROR"],
				nil,
				"List of user need to pay must be not empty",
			),
		)
		return
	}

	if requestPayload.Request.Purpose == "" {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				c,
				constant.ErrorConstant["DATA_FORMAT_ERROR"],
				nil,
				"What is your purpose of this expense?",
			),
		)
		return
	}

	currentUser, isCurrentUserExist := utils.GetCurrentUsername(c)

	ctx := context.Background()

	ctx = context.WithValue(ctx, "username", *currentUser)

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

	boughtUser := model.Users{}

	h.DB.Where(
		&model.Users{
			BaseEntity: model.BaseEntity{
				Active: utils.GetPointerOfAnyValue(true),
			},
			Username: *currentUser,
		},
	).Find(&boughtUser)

	if boughtUser.BaseEntity.Id == 0 {
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			ReturnResponse(
				c,
				constant.ErrorConstant["USER_NOT_EXISTED"],
				nil,
				"Who are you?",
			),
		)
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
							Active: utils.GetPointerOfAnyValue(true),
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
							Active: utils.GetPointerOfAnyValue(true),
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
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				c,
				constant.ErrorConstant["INVALID_NUMBER_OF_USER"],
				nil,
			),
		)
		return
	}

	if slices.Contains(requestPayload.Request.UserToPaid, boughtUser.BaseEntity.Id) {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				c,
				constant.ErrorConstant["INVALID_USER_TO_PAID_LIST"],
				nil,
			),
		)
		return
	}
	log.Info(
		fmt.Sprintf(
			constant.LogPattern,
			utils.GetTraceId(c),
			*currentUser,
			fmt.Sprintf("bought user is user with id: %s", strconv.FormatInt(boughtUser.BaseEntity.Id, 10)),
		),
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
		log.Info(
			fmt.Sprintf(
				constant.LogPattern,
				utils.GetTraceId(c),
				*currentUser,
				"this user will participate to this expense",
			),
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
	log.Info(
		fmt.Sprintf(
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
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				c,
				constant.ErrorConstant["QUERY_ERROR"],
				nil,
				expenseTransactionError.Error(),
			),
		)
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

	c.JSON(
		http.StatusOK,
		ReturnResponse(
			c,
			constant.ErrorConstant["SUCCESS"],
			savedExpense,
		),
	)

}

func (h *ExpenseHandler) RemoveExpense(c *gin.Context) {

	requestPayload := payload.RemoveExpenseBody{}
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

	var expense model.ListOfExpenses

	h.DB.Preload("Users").Preload("DebitUser").Where(
		model.ListOfExpenses{
			BaseEntity: model.BaseEntity{
				Id: requestPayload.Request,
			},
		},
	).Find(&expense)

	if expense.BaseEntity.Id == 0 {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				c,
				constant.ErrorConstant["EXPENSE_DELETE_NOT_SUCCESS"],
				nil,
			),
		)
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

	c.JSON(
		http.StatusOK,
		ReturnResponse(
			c,
			constant.ErrorConstant["SUCCESS"],
			nil,
		),
	)
}

func (h *ExpenseHandler) SoftRemoveExpense(c *gin.Context) {

	currentUser, isCurrentUserExist := utils.GetCurrentUsername(c)

	ctx := context.Background()

	ctx = context.WithValue(ctx, "username", *currentUser)

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

	requestPayload := payload.RemoveExpenseBody{}
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

	var errorEnum = constant.ErrorEnums{}

	transactionResult := h.DB.Transaction(func(tx *gorm.DB) error {

		var expense model.ListOfExpenses

		tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Users").Preload("DebitUser").Where(
			model.ListOfExpenses{
				BaseEntity: model.BaseEntity{
					Id:     requestPayload.Request,
					Active: utils.GetPointerOfAnyValue(true),
				},
			},
		).Find(&expense)

		if expense.BaseEntity.Id == 0 {
			errorEnum = constant.ErrorConstant["EXPENSE_DELETE_NOT_SUCCESS"]
			return nil
		}

		debitUser := expense.DebitUser
		for _, du := range debitUser {
			log.Info(
				fmt.Sprintf(
					constant.LogPattern,
					utils.GetTraceId(c),
					*currentUser,
					"changing active status of debit user id "+strconv.FormatInt(du.BaseEntity.Id, 10)+" to false",
				),
			)
			*du.BaseEntity.Active = false
		}
		debitUserSoftRemoveResult := tx.Save(&debitUser)
		if debitUserSoftRemoveResult.Error != nil {
			return debitUserSoftRemoveResult.Error
		}
		*expense.BaseEntity.Active = false
		log.Info(
			fmt.Sprintf(
				constant.LogPattern,
				utils.GetTraceId(c),
				*currentUser,
				"changing active status of expense id "+strconv.FormatInt(expense.BaseEntity.Id, 10)+" to false",
			),
		)
		expenseSoftRemoveResult := tx.Save(&expense)
		if expenseSoftRemoveResult.Error != nil {
			return expenseSoftRemoveResult.Error
		}
		return nil
	})
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

	c.JSON(
		http.StatusOK,
		ReturnResponse(
			c,
			constant.ErrorConstant["SUCCESS"],
			nil,
		),
	)
}

func (h *ExpenseHandler) ActiveRemoveExpense(c *gin.Context) {

	currentUser, isCurrentUserExist := utils.GetCurrentUsername(c)

	ctx := context.Background()

	ctx = context.WithValue(ctx, "username", *currentUser)

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

	requestPayload := payload.RemoveExpenseBody{}
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

	var errorEnum = constant.ErrorEnums{}

	transactionResult := h.DB.Transaction(func(tx *gorm.DB) error {

		var expense model.ListOfExpenses

		tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Users").Preload("DebitUser").Where(
			model.ListOfExpenses{
				BaseEntity: model.BaseEntity{
					Id:     requestPayload.Request,
					Active: utils.GetPointerOfAnyValue(false),
				},
			},
		).Find(&expense)

		if expense.BaseEntity.Id == 0 {
			errorEnum = constant.ErrorConstant["EXPENSE_ACTIVE_NOT_SUCCESS"]
			return nil
		}

		debitUser := expense.DebitUser
		for _, du := range debitUser {
			log.Info(
				fmt.Sprintf(
					constant.LogPattern,
					utils.GetTraceId(c),
					*currentUser,
					"changing active status of debit user id "+strconv.FormatInt(du.BaseEntity.Id, 10)+" to true",
				),
			)
			*du.BaseEntity.Active = true
		}
		debitUserSoftRemoveResult := tx.Save(&debitUser)
		if debitUserSoftRemoveResult.Error != nil {
			return debitUserSoftRemoveResult.Error
		}
		*expense.BaseEntity.Active = true
		log.Info(
			fmt.Sprintf(
				constant.LogPattern,
				utils.GetTraceId(c),
				*currentUser,
				"changing active status of expense id "+strconv.FormatInt(expense.BaseEntity.Id, 10)+" to true",
			),
		)
		expenseSoftRemoveResult := tx.Save(&expense)
		if expenseSoftRemoveResult.Error != nil {
			return expenseSoftRemoveResult.Error
		}
		return nil
	})
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

	c.JSON(
		http.StatusOK,
		ReturnResponse(
			c,
			constant.ErrorConstant["SUCCESS"],
			nil,
		),
	)
}

func (h *ExpenseHandler) ListExpense(c *gin.Context) {

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

	if requestPayload.Request.PageSize == 0 || requestPayload.Request.PageNumber == 0 {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				c,
				constant.ErrorConstant["DATA_FORMAT_ERROR"],
				nil,
				"Page number or page size can not be 0",
			),
		)
		return
	}

	limit := requestPayload.Request.PageSize
	offset := requestPayload.Request.PageSize * (requestPayload.Request.PageNumber - 1)

	var expense []model.ListOfExpenses

	var total int64 = 0

	h.DB.Model(&model.ListOfExpenses{}).Preload("Users").Preload("DebitUser").
		Where(
			model.ListOfExpenses{
				BaseEntity: model.BaseEntity{
					Active: utils.GetPointerOfAnyValue(true),
				},
			},
		).
		Count(&total)

	h.DB.Preload("Users").Preload("DebitUser").Limit(limit).Offset(offset).
		Where(
			model.ListOfExpenses{
				BaseEntity: model.BaseEntity{
					Active: utils.GetPointerOfAnyValue(true),
				},
			},
		).
		Find(&expense)

	totalRecordsSelected := len(expense)

	if totalRecordsSelected == 0 {
		c.JSON(
			http.StatusOK,
			ReturnPageResponse(
				c,
				constant.ErrorConstant["SUCCESS"],
				0,
				0,
				expense,
			),
		)
		return
	}

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
			expense,
		),
	)
}

func SaveNewExpense(db *gorm.DB, model *model.ListOfExpenses, ctx context.Context) *gorm.DB {
	model.BaseEntity.Active = utils.GetPointerOfAnyValue(true)
	model.BaseEntity.CreatedAt = time.Now()
	model.BaseEntity.UpdatedAt = time.Now()
	model.BaseEntity.CreatedBy = ctx.Value("username").(string)
	model.BaseEntity.UpdatedBy = ctx.Value("username").(string)

	return db.Create(model)
}

func SaveNewDebitUser(db *gorm.DB, model *model.DebitUser, ctx context.Context) *gorm.DB {
	model.BaseEntity.Active = utils.GetPointerOfAnyValue(true)
	model.BaseEntity.CreatedAt = time.Now()
	model.BaseEntity.UpdatedAt = time.Now()
	model.BaseEntity.CreatedBy = ctx.Value("username").(string)
	model.BaseEntity.UpdatedBy = ctx.Value("username").(string)

	return db.Create(model)
}
