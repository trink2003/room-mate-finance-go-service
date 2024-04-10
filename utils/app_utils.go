package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/big"
	"os"
	"room-mate-finance-go-service/constant"
	"strconv"
	"strings"
	"time"
)

func EncryptPassword(password string) (encryptedPassword string, error error) {
	encryptedPassword = ""
	bytePassword := []byte(password)
	hashedPassword, generateFromPasswordErr := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if generateFromPasswordErr == nil {
		encryptedPassword = string(hashedPassword)
	} else {
		error = generateFromPasswordErr
	}
	return encryptedPassword, error
}

func EncryptPasswordPointer(password *string) (error error) {
	bytePassword := []byte(*password)
	hashedPassword, generateFromPasswordErr := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if generateFromPasswordErr == nil {
		*password = string(hashedPassword)
	}
	error = generateFromPasswordErr
	return error
}

func ComparePassword(inputPassword string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}

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

func GenerateJwtToken(username string, role ...string) string {

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "Q8OzIHRo4buDIGfhu41pIGFuaCBsw6AgxJHhurlwIHRyYWkgbmjhuqV0IFZp4buHdCBOYW0"
	}

	tokenExpireTime := os.Getenv("JWT_EXPIRE_TIME")
	if tokenExpireTime == "" {
		tokenExpireTime = "10"
	}

	expireTime, err := strconv.Atoi(tokenExpireTime)

	if err != nil {
		panic(err)
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                                                       // Subject (user identifier)
		"iss": "room-mate-finance-go-service",                                 // Issuer
		"aud": role,                                                           // Audience (user role)
		"exp": time.Now().Add(time.Duration(expireTime) * time.Minute).Unix(), // Expiration time
		"iat": time.Now().Unix(),                                              // Issued at
	})
	tokenString, signedStringError := claims.SignedString([]byte(secretKey))
	if signedStringError != nil {
		panic(signedStringError)
	}
	return tokenString
}

func VerifyJwtToken(token string) (jwt.MapClaims, error) {

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "Q8OzIHRo4buDIGfhu41pIGFuaCBsw6AgxJHhurlwIHRyYWkgbmjhuqV0IFZp4buHdCBOYW0"
	}

	parsedToken, tokenParseError := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if tokenParseError != nil {
		log.Print(tokenParseError)
		return nil, tokenParseError
	}

	if !parsedToken.Valid {
		return nil, errors.New("token invalid")
	}

	return parsedToken.Claims.(jwt.MapClaims), tokenParseError
}

func GetCurrentUsername(c *gin.Context) (username *string, err error) {

	currentUser, isCurrentUserExist := c.Get("auth")

	emptyString := constant.EmptyString

	if isCurrentUserExist == false {
		return &emptyString, errors.New("can not get current username")
	}

	claim := currentUser.(jwt.MapClaims)

	currentUsername := claim["sub"].(string)

	return &currentUsername, nil
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

func StructToJson(anyStruct any) string {
	result, err := json.Marshal(anyStruct)
	if err != nil {
		return ""
	}
	return string(result)
}

func JsonToStruct[T any](jsonString string, anyStruct *T) {
	err := json.Unmarshal([]byte(jsonString), anyStruct)
	if err != nil {
		return
	}
}

func BeginningOfMonth(date time.Time) time.Time {
	result := date.AddDate(0, 0, -date.Day()+1)
	fmt.Printf("BeginningOfMonth time zone is %s\n", result.Location())
	return time.Date(result.Year(), result.Month(), result.Day(), 0, 0, 0, result.Nanosecond(), result.Location())
}

func EndOfMonth(date time.Time) time.Time {
	result := date.AddDate(0, 1, -date.Day())
	fmt.Printf("EndOfMonth time zone is %s\n", result.Location().String())
	return time.Date(result.Year(), result.Month(), result.Day(), 23, 59, 59, result.Nanosecond(), result.Location())
}

func SortMapToString(inputMap map[string]string) string {
	result := ""
	for k, v := range inputMap {
		sort := ""
		if v != constant.AscKeyword && v != constant.DescKeyword {
			sort = constant.DescKeyword
		} else {
			sort = v
		}
		result += k + " " + sort + ", "
	}
	return strings.TrimSuffix(result, ", ")
}

func HideSensitiveJsonField(inputJson string) string {
	element := strings.Split(inputJson, "\"")
	for i := range element {
		currentField := element[i]
		var colon string
		if (len(element) == 0) || (i+1 > len(element)-1) {
			continue
		}
		colon = element[i+1]
		if isSensitive(currentField) {
			if strings.Contains(strings.Trim(colon, " "), ":") {
				element[i+2] = "***"
			}
		} else if i+2 < len(element) && len(element[i+2]) > 1000 {
			element[i+2] = element[i+2][:50] + "..." + element[i+2][len(element[i+2])-50:]
		}
	}
	return strings.Join(element, "\"")
}

func isSensitive(input string) bool {
	for _, e := range constant.SensitiveField {
		if strings.Contains(strings.ToLower(e), strings.ToLower(input)) || strings.Contains(strings.ToLower(input), strings.ToLower(e)) {
			return true
		}
	}
	return false
}
