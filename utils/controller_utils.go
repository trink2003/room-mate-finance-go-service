package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/log"
	"room-mate-finance-go-service/payload"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Permission struct {
	Url  string   `json:"url"`
	Role []string `json:"role"`
}

type BodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w BodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ErrorHandler(c *gin.Context) {
	CheckAndSetTraceId(c)
	if c.Errors != nil && len(c.Errors.Errors()) != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": c.Errors.Errors()})
	}
}

func RequestLogger(c *gin.Context) {
	CheckAndSetTraceId(c)
	// t := time.Now()
	var buf bytes.Buffer
	tee := io.TeeReader(c.Request.Body, &buf)
	body, _ := io.ReadAll(tee)
	c.Request.Body = io.NopCloser(&buf)
	dst := &bytes.Buffer{}
	if err := json.Compact(dst, body); err != nil && len(body) > 0 {
		panic(err)
	}

	header := map[string][]string(c.Request.Header)

	headerString := ""

	for k, v := range header {
		if IsSensitiveField(k) {
			headerString += fmt.Sprintf("\n\t\t- %s: %s", k, "***")
		} else {
			headerString += fmt.Sprintf("\n\t\t- %s: %s", k, strings.Join(v, ", "))
		}
	}

	message := fmt.Sprintf(
		"Request info:\n\t- header: %s\n\t- url: %s\n\t- method: %s\n\t- proto: %s\n\t- payload:\n\t%s",
		headerString,
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
	var ctx = context.Background()
	ctx = context.WithValue(ctx, constant.UsernameLogKey, currentUser)
	ctx = context.WithValue(ctx, constant.TraceIdLogKey, GetTraceId(c))
	log.WithLevel(
		constant.Info,
		ctx,
		HideSensitiveJsonField(message),
	)
	c.Next()
	// latency := time.Since(t)
	// log.Info("%s %s %s %s\n",
	// 	c.Request.RequestURI,
	// )
}

func ResponseLogger(c *gin.Context) {
	CheckAndSetTraceId(c)
	c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, PUT, GET, OPTIONS, DELETE")
	c.Writer.Header().Set("Access-Control-Max-Age", "3600")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-Requested-With, credential, X-XSRF-TOKEN")
	blw := &BodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	c.Next()

	header := map[string][]string(c.Writer.Header())

	headerString := ""

	for k, v := range header {
		if IsSensitiveField(k) {
			headerString += fmt.Sprintf("\n\t\t- %s: %s", k, "***")
		} else {
			headerString += fmt.Sprintf("\n\t\t- %s: %s", k, strings.Join(v, ", "))
		}
	}

	statusCode := c.Writer.Status()
	message := fmt.Sprintf(
		"Response info:\n\t- status code: %s\n\t- method: %s\n\t- url: %s\n\t- header: %s\n\t- payload:\n\t%s",
		strconv.Itoa(statusCode),
		c.Request.Method,
		c.Request.RequestURI,
		headerString,
		blw.body.String(),
	)
	currentUser := "unknown"
	claimFromGinContext, _ := c.Get("auth")
	if claimFromGinContext != nil {
		claims := claimFromGinContext.(jwt.MapClaims)
		currentUser = claims["sub"].(string)
	}
	var ctx = context.Background()
	ctx = context.WithValue(ctx, constant.UsernameLogKey, currentUser)
	ctx = context.WithValue(ctx, constant.TraceIdLogKey, GetTraceId(c))
	log.WithLevel(
		constant.Info,
		ctx,
		HideSensitiveJsonField(message),
	)

}

func AuthenticationWithAuthorization(listOfRole []string) func(c *gin.Context) {
	return func(c *gin.Context) {
		CheckAndSetTraceId(c)
		traceId := GetTraceId(c)
		ctx := context.Background()
		ctx = context.WithValue(ctx, constant.TraceIdLogKey, traceId)
		token := c.Request.Header.Get("Authorization")
		var mapClaims jwt.MapClaims
		var verifyJwtTokenError error
		if strings.Contains(token, "Bearer") {
			mapClaims, verifyJwtTokenError = VerifyJwtToken(ctx, token[7:])
		} else {
			mapClaims, verifyJwtTokenError = VerifyJwtToken(ctx, token)
		}
		if verifyJwtTokenError != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &payload.Response{
				Trace:        traceId,
				ErrorCode:    constant.Unauthorized.ErrorCode,
				ErrorMessage: constant.Unauthorized.ErrorMessage,
			})
			return
		}
		currentUsername := mapClaims["sub"].(string)
		ctx = context.WithValue(ctx, constant.UsernameLogKey, currentUsername)
		c.Set("auth", mapClaims)
		log.WithLevel(
			constant.Info,
			ctx,
			fmt.Sprintf("Check permission for url: %v", c.Request.RequestURI),
		)
		if listOfRole == nil || len(listOfRole) < 1 {
			c.Next()
			return
		}
		userRolesFromAccessToken := mapClaims["aud"]
		if userRolesFromAccessToken != nil {
			roleListInterface := userRolesFromAccessToken.([]interface{})
			log.WithLevel(
				constant.Info,
				ctx,
				fmt.Sprintf(
					"\n\t- this user has role: %v\n\t- current api require user with role: %v",
					roleListInterface,
					listOfRole,
				),
			)
			for _, roleListInterfaceElement := range roleListInterface {
				if slices.Contains(listOfRole, fmt.Sprintf("%v", roleListInterfaceElement)) {
					c.Next()
					return
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, &payload.Response{
			Trace:        traceId,
			ErrorCode:    constant.Forbidden.ErrorCode,
			ErrorMessage: constant.Forbidden.ErrorMessage,
		})
	}
}

func ReturnResponse(c *gin.Context, errCode constant.ErrorEnums, responseData any, additionalMessage ...string) *payload.Response {
	message := ""
	if len(additionalMessage) > 0 {
		message = additionalMessage[0]
	}

	return &payload.Response{
		Trace:        GetTraceId(c),
		ErrorCode:    errCode.ErrorCode,
		ErrorMessage: strings.Replace(strings.Trim(strings.Trim(errCode.ErrorMessage, " ")+". "+strings.Trim(message, " ")+".", " "), ". .", ".", -1),
		Response:     responseData,
	}
}

func ReturnPageResponse(
	c *gin.Context,
	errCode constant.ErrorEnums,
	totalElement int64,
	totalPage int64,
	responseData any,
	additionalMessage ...string,
) *payload.PageResponse {
	message := ""
	if len(additionalMessage) > 0 {
		message = additionalMessage[0]
	}

	return &payload.PageResponse{
		Trace:        GetTraceId(c),
		ErrorCode:    errCode.ErrorCode,
		ErrorMessage: strings.Replace(strings.Trim(strings.Trim(errCode.ErrorMessage, " ")+". "+strings.Trim(message, " ")+".", " "), ". .", ".", -1),
		TotalElement: totalElement,
		TotalPage:    totalPage,
		Response:     responseData,
	}
}

// func readPermissionJsonFile() *[]Permission {
// 	var result []Permission
// 	filePath := filepath.Join(GetCurrentDirectory(), "permission.json")
// 	// log.Printf(filePath)
// 	jsonFile, err := os.Open(filePath)
// 	if err != nil {
// 		// log.Printf(err.Error())
// 		return &result
// 	}

// 	defer func(jsonFile *os.File) {
// 		deferErr := jsonFile.Close()
// 		if deferErr != nil {
// 			// log.Printf(deferErr.Error())
// 			panic(deferErr)
// 		}
// 	}(jsonFile)
// 	byteValue, readAllErr := io.ReadAll(jsonFile)
// 	if readAllErr != nil {
// 		// log.Printf(readAllErr.Error())
// 		return &result
// 	}
// 	// log.Printf(string(byteValue))
// 	ByteJsonToStruct(byteValue, &result)

// 	return &result

// }
