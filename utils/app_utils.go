package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strconv"
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
