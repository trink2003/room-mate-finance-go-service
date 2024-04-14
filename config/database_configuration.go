package config

import (
	"context"
	"fmt"
	"os"
	"room-mate-finance-go-service/constant"
	"time"

	"github.com/charmbracelet/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabaseConnection() (db *gorm.DB, err error) {
	/*
		if databaseUsername, ok := os.LookupEnv("DATABASE_USERNAME"); !ok {
			databaseUsername = "root"
		}
	*/
	var databaseUsername = os.Getenv("DATABASE_USERNAME")
	if databaseUsername == "" {
		databaseUsername = "postgres"
	}
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	if databasePassword == "" {
		databasePassword = "mysecretpassword"
	}
	databaseHost := os.Getenv("DATABASE_HOST")
	if databaseHost == "" {
		databaseHost = "localhost"
	}
	databasePort := os.Getenv("DATABASE_PORT")
	if databasePort == "" {
		databasePort = "5432"
	}
	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		databaseName = "room-mate-finance"
	}
	/*
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			databaseUsername,
			databasePassword,
			databaseHost,
			databasePort,
			databaseName,
		)
	*/
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		databaseHost,
		databaseUsername,
		databasePassword,
		databaseName,
		databasePort,
	)
	log.Info(
		fmt.Sprintf(
			constant.LogPattern,
			"",
			"",
			fmt.Sprintf(
				"Database connect info:\n    - databaseHost: %s\n    - databaseUsername: %s\n    - databaseName: %s\n    - databasePort: %s",
				databaseHost,
				databaseUsername,
				databaseName,
				databasePort,
			),
		),
	)
	/*
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	*/
	myLog := &dbLogger{}
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: myLog.LogMode(logger.Info),
	})
	if err != nil {
		log.Error(err)
		return db, err
	}

	sqlDB, err := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// MigrationAndInsertDate(db)
	// InsertData(db)

	return db, err
}

type dbLogger struct{}

func (d *dbLogger) LogMode(level logger.LogLevel) logger.Interface {
	return d
}
func (d *dbLogger) Info(ctx context.Context, s string, i ...interface{}) {
	usernameFromContext := ctx.Value("username")
	traceIdFromContext := ctx.Value("traceId")
	username := ""
	traceId := ""
	if usernameFromContext != nil {
		username = usernameFromContext.(string)
	}
	if traceIdFromContext != nil {
		traceId = traceIdFromContext.(string)
	}
	log.Info(
		fmt.Sprintf(
			constant.LogPattern,
			traceId,
			username,
			s,
		),
	)
}
func (d *dbLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	usernameFromContext := ctx.Value("username")
	traceIdFromContext := ctx.Value("traceId")
	username := ""
	traceId := ""
	if usernameFromContext != nil {
		username = usernameFromContext.(string)
	}
	if traceIdFromContext != nil {
		traceId = traceIdFromContext.(string)
	}
	log.Warn(
		fmt.Sprintf(
			constant.LogPattern,
			traceId,
			username,
			s,
		),
	)
}
func (d *dbLogger) Error(ctx context.Context, s string, i ...interface{}) {
	usernameFromContext := ctx.Value("username")
	traceIdFromContext := ctx.Value("traceId")
	username := ""
	traceId := ""
	if usernameFromContext != nil {
		username = usernameFromContext.(string)
	}
	if traceIdFromContext != nil {
		traceId = traceIdFromContext.(string)
	}
	log.Error(
		fmt.Sprintf(
			constant.LogPattern,
			traceId,
			username,
			s,
		),
	)
}
func (d *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rowsAffected := fc()
	usernameFromContext := ctx.Value("username")
	traceIdFromContext := ctx.Value("traceId")
	username := ""
	traceId := ""
	if usernameFromContext != nil {
		username = usernameFromContext.(string)
	}
	if traceIdFromContext != nil {
		traceId = traceIdFromContext.(string)
	}
	errorMessage := "no error"
	if err != nil {
		errorMessage = err.Error()
	}
	log.Info(
		fmt.Sprintf(
			constant.LogPattern,
			traceId,
			username,
			fmt.Sprintf(
				"info:\n    - sql: %s\n    - rowsAffected: %v\n    - begin: %s\n    - error: %s",
				sql,
				rowsAffected,
				begin.Format(constant.YyyyMmDdHhMmSsFormat),
				errorMessage,
			),
		),
	)
}
