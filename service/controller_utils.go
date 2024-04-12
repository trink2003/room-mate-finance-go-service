package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/payload"
	"room-mate-finance-go-service/utils"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
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

	header := map[string][]string(c.Request.Header)

	headerString := ""

	for k, v := range header {
		if utils.IsSensitiveField(k) {
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
	log.Info(
		fmt.Sprintf(
			constant.LogPattern,
			utils.GetTraceId(c),
			currentUser,
			utils.HideSensitiveJsonField(message),
		),
	)
	c.Next()
	// latency := time.Since(t)
	// log.Info("%s %s %s %s\n",
	// 	c.Request.RequestURI,
	// )
}

func ResponseLogger(c *gin.Context) {
	utils.CheckAndSetTraceId(c)
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
		if utils.IsSensitiveField(k) {
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
	log.Info(
		fmt.Sprintf(
			constant.LogPattern,
			utils.GetTraceId(c),
			currentUser,
			utils.HideSensitiveJsonField(message),
		),
	)

}

func Authentication(c *gin.Context) {
	traceId := utils.GetTraceId(c)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "traceId", traceId)
	token := c.Request.Header.Get("Authorization")
	var mapClaims jwt.MapClaims
	var err error
	if strings.Contains(token, "Bearer") {
		mapClaims, err = utils.VerifyJwtToken(ctx, token[7:])
	} else {
		mapClaims, err = utils.VerifyJwtToken(ctx, token)
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &payload.Response{
			Trace:        traceId,
			ErrorCode:    constant.ErrorConstant["UNAUTHORIZED"].ErrorCode,
			ErrorMessage: constant.ErrorConstant["UNAUTHORIZED"].ErrorMessage,
		})
		return
	}
	c.Set("auth", mapClaims)
	permissionList := *readPermissionJsonFile()
	log.Info(
		fmt.Sprintf(
			constant.LogPattern,
			traceId,
			"",
			fmt.Sprintf("Check permission for url: %v", c.Request.RequestURI),
		),
	)
	for _, p := range permissionList {
		log.Info(
			fmt.Sprintf(
				constant.LogPattern,
				traceId,
				"",
				fmt.Sprintf("Current url: %v", p.Url),
			),
		)
		if strings.Compare(strings.ToLower(c.Request.RequestURI), strings.ToLower(p.Url)) == 0 {
			log.Info(
				fmt.Sprintf(
					constant.LogPattern,
					traceId,
					"",
					fmt.Sprintf("This api is accessable with role: %v", p.Role),
				),
			)
			if slices.Contains(p.Role, "*") {
				c.Next()
				return
			}
			userRole := mapClaims["aud"]
			if userRole != nil {

				roleList := userRole.([]interface{})
				log.Info(
					fmt.Sprintf(
						constant.LogPattern,
						traceId,
						"",
						fmt.Sprintf("This user have role: %v", roleList),
					),
				)
				for _, rI := range roleList {
					if slices.Contains(p.Role, fmt.Sprintf("%v", rI)) {
						c.Next()
						return
					}
				}
				// fmt.Printf("\n\n%s - %T\n\n", userRole, userRole)
			}
		} else {
			log.Info(
				fmt.Sprintf(
					constant.LogPattern,
					traceId,
					"",
					"Not match",
				),
			)
		}
	}
	c.AbortWithStatusJSON(http.StatusForbidden, &payload.Response{
		Trace:        traceId,
		ErrorCode:    constant.ErrorConstant["FORBIDDEN"].ErrorCode,
		ErrorMessage: constant.ErrorConstant["FORBIDDEN"].ErrorMessage,
	})
	return
}

func ReturnResponse(c *gin.Context, errCode constant.ErrorEnums, responseData any, additionalMessage ...string) *payload.Response {
	message := ""
	if len(additionalMessage) > 0 {
		message = additionalMessage[0]
	}

	return &payload.Response{
		Trace:        utils.GetTraceId(c),
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
		Trace:        utils.GetTraceId(c),
		ErrorCode:    errCode.ErrorCode,
		ErrorMessage: strings.Replace(strings.Trim(strings.Trim(errCode.ErrorMessage, " ")+". "+strings.Trim(message, " ")+".", " "), ". .", ".", -1),
		TotalElement: totalElement,
		TotalPage:    totalPage,
		Response:     responseData,
	}
}

func readPermissionJsonFile() *[]Permission {
	var result []Permission
	filePath := filepath.Join(utils.GetCurrentDirectory(), "permission.json")
	// log.Printf(filePath)
	jsonFile, err := os.Open(filePath)
	if err != nil {
		// log.Printf(err.Error())
		return &result
	}

	defer func(jsonFile *os.File) {
		deferErr := jsonFile.Close()
		if deferErr != nil {
			// log.Printf(deferErr.Error())
			panic(deferErr)
		}
	}(jsonFile)
	byteValue, readAllErr := io.ReadAll(jsonFile)
	if readAllErr != nil {
		// log.Printf(readAllErr.Error())
		return &result
	}
	// log.Printf(string(byteValue))
	utils.ByteJsonToStruct(byteValue, &result)

	return &result

}
