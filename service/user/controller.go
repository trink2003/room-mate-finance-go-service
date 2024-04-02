package user

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"room-mate-finance-go-service/utils"
	"time"
)

type UserHandler struct {
	DB *gorm.DB
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	h := &UserHandler{
		DB: db,
	}

	router.Use(ErrorHandler)
	router.Use(RequestLogger)
	router.Use(ResponseLogger)

	routes := router.Group("/auth")
	routes.POST("/register", h.AddNewUser)
}

func ErrorHandler(c *gin.Context) {
	utils.CheckAndSetTraceId(c)
	if c.Errors != nil && len(c.Errors.Errors()) != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": c.Errors.Errors()})
	}
}

func RequestLogger(c *gin.Context) {
	utils.CheckAndSetTraceId(c)
	t := time.Now()
	var buf bytes.Buffer
	tee := io.TeeReader(c.Request.Body, &buf)
	body, _ := io.ReadAll(tee)
	c.Request.Body = io.NopCloser(&buf)
	log.Printf(utils.GetTraceId(c) + " - Request body: " + string(body))
	log.Printf(utils.GetTraceId(c)+" - "+"Request info: %s", c.Request.Header)
	c.Next()
	latency := time.Since(t)
	log.Printf("%s %s %s %s\n",
		c.Request.Method,
		c.Request.RequestURI,
		c.Request.Proto,
		latency,
	)
}

func ResponseLogger(c *gin.Context) {
	utils.CheckAndSetTraceId(c)
	c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	c.Next()

	statusCode := c.Writer.Status()
	log.Printf(utils.GetTraceId(c)+" - "+"%d %s %s\n",
		statusCode,
		c.Request.Method,
		c.Request.RequestURI,
	)
	log.Printf(utils.GetTraceId(c)+" - "+"Response body: %s", blw.body.String())

}
