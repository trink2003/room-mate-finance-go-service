package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"log"
	"net/http"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/payload"
	"room-mate-finance-go-service/utils"
	"strconv"
	"strings"
)

type BodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w BodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ErrorHandler(c *gin.Context) {
	utils.CheckAndSetTraceId(c)
	if c.Errors != nil && len(c.Errors.Errors()) != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": c.Errors.Errors()})
	}
}

func RequestLogger(c *gin.Context) {
	utils.CheckAndSetTraceId(c)
	// t := time.Now()
	var buf bytes.Buffer
	tee := io.TeeReader(c.Request.Body, &buf)
	body, _ := io.ReadAll(tee)
	c.Request.Body = io.NopCloser(&buf)
	dst := &bytes.Buffer{}
	if err := json.Compact(dst, body); err != nil && len(body) > 0 {
		panic(err)
	}
	message := fmt.Sprintf(
		"Request info:\n\t- header: %s\n\t- url: %s\n\t- method: %s\n\t- proto: %s\n\t- payload:\n\t%s",
		c.Request.Header,
		c.Request.RequestURI,
		c.Request.Method,
		c.Request.Proto,
		dst.String(),
	)
	currentUser := "unknown"
	claimFromGinContext, _ := c.Get("auth")
	if claimFromGinContext != nil {
		claims := claimFromGinContext.(jwt.MapClaims)
		currentUser = claims["sub"].(string)
	}
	log.Printf(
		constant.LogPattern,
		utils.GetTraceId(c),
		currentUser,
		message,
	)
	c.Next()
	// latency := time.Since(t)
	// log.Printf("%s %s %s %s\n",
	// 	c.Request.RequestURI,
	// )
}

func ResponseLogger(c *gin.Context) {
	utils.CheckAndSetTraceId(c)
	c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	blw := &BodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	c.Next()

	statusCode := c.Writer.Status()
	message := fmt.Sprintf(
		"Response info:\n\t- status code: %s\n\t- method: %s\n\t- url: %s\n\t- payload:\n\t%s",
		strconv.Itoa(statusCode),
		c.Request.Method,
		c.Request.RequestURI,
		blw.body.String(),
	)
	currentUser := "unknown"
	claimFromGinContext, _ := c.Get("auth")
	if claimFromGinContext != nil {
		claims := claimFromGinContext.(jwt.MapClaims)
		currentUser = claims["sub"].(string)
	}
	log.Printf(
		constant.LogPattern,
		utils.GetTraceId(c),
		currentUser,
		message,
	)

}

func Authentication(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	var mapClaims jwt.MapClaims
	var err error
	if strings.Contains(token, "Bearer") {
		mapClaims, err = utils.VerifyJwtToken(token[7:])
	} else {
		mapClaims, err = utils.VerifyJwtToken(token)
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &payload.Response{
			Trace:        utils.GetTraceId(c),
			ErrorCode:    constant.ErrorConstant["UNAUTHORIZED"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["UNAUTHORIZED"].ErrorMessage,
		})
		return
	}
	c.Set("auth", mapClaims)
	c.Next()
}
