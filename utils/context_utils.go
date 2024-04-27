package utils

import (
	"context"
	"net/http"
	"room-mate-finance-go-service/constant"

	"github.com/gin-gonic/gin"
)

func PrepareContext(c *gin.Context, isBypassCurrentUserCheck ...bool) (context.Context, bool) {
	ctx := context.Background()

	currentUser, isCurrentUserExist := GetCurrentUsername(c)

	if len(isBypassCurrentUserCheck) < 1 || !isBypassCurrentUserCheck[0] {
		if isCurrentUserExist != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				ReturnResponse(
					c,
					constant.Unauthorized,
					nil,
					isCurrentUserExist.Error(),
				),
			)
			return ctx, false
		}
		ctx = context.WithValue(ctx, constant.UsernameLogKey, *currentUser)
	}
	ctx = context.WithValue(ctx, constant.TraceIdLogKey, GetTraceId(c))

	return ctx, true
}

func GetCurrentUsernameFromContextForInsertOrUpdateDataInDb(ctx context.Context) string {
	var currentUsernameInsertOrUpdateData = ""
	var usernameFromContext = ctx.Value(constant.UsernameLogKey)
	if usernameFromContext != nil {
		currentUsernameInsertOrUpdateData = usernameFromContext.(string)
	}
	return currentUsernameInsertOrUpdateData
}
