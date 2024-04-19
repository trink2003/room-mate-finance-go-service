package utils

import (
	"math/big"
	"room-mate-finance-go-service/constant"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CheckAndSetTraceId(c *gin.Context) {
	if traceId, _ := c.Get("traceId"); traceId == nil || traceId == "" {
		c.Set("traceId", uuid.New().String())
	}
}

func GetTraceId(c *gin.Context) string {
	if traceId, _ := c.Get("traceId"); traceId == nil || traceId == "" {
		return ""
	} else {
		return traceId.(string)
	}
}

func RoundHalfUpBigFloat(input *big.Float) {
	delta := constant.DeltaPositive

	if input.Sign() < 0 {
		delta = constant.DeltaNegative
	}
	input.Add(input, new(big.Float).SetFloat64(delta))
}

func GetPointerOfAnyValue[T any](a T) *T {
	return &a
}
