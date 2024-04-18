package config

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/log"
	"time"
)

func InitDatabaseConnection() (db *gorm.DB, err error) {
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
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		databaseHost,
		databaseUsername,
		databasePassword,
		databaseName,
		databasePort,
	)
	var ctx = context.Background()
	log.WithLevel(
		constant.Info,
		ctx,
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
	myLog := &dbLogger{}
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: myLog.LogMode(logger.Info),
	})
	if err != nil {
		log.WithLevel(
			constant.Error,
			ctx,
			fmt.Sprintf(
				"An error has been occurred when trying to connect to Database:\n\t- error: %v",
				err,
			),
		)
		return db, err
	}

	sqlDB, getDbObjectError := db.DB()

	if getDbObjectError != nil {
		log.WithLevel(
			constant.Error,
			ctx,
			fmt.Sprintf(
				"An error has been occurred when trying to get Database Object:\n\t- error: %v",
				getDbObjectError,
			),
		)
		return db, getDbObjectError
	}

	// sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	if isDatabaseMigration, isDatabaseMigrationExist := os.LookupEnv("DATABASE_MIGRATION"); isDatabaseMigration == "true" && isDatabaseMigrationExist {
		MigrationAndInsertDate(db)
	}
	if isDatabaseInitializationData, isDatabaseInitializationDataExist := os.LookupEnv("DATABASE_INITIALIZATION_DATA"); isDatabaseInitializationData == "true" && isDatabaseInitializationDataExist {
		InsertData(db)
	}

	return db, err
}

type dbLogger struct{}

func (d *dbLogger) LogMode(level logger.LogLevel) logger.Interface {
	return d
}
func (d *dbLogger) Info(ctx context.Context, s string, i ...interface{}) {
	log.WithLevel(
		constant.Info,
		ctx,
		s,
	)
}
func (d *dbLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	log.WithLevel(
		constant.Warn,
		ctx,
		s,
	)
}
func (d *dbLogger) Error(ctx context.Context, s string, i ...interface{}) {
	log.WithLevel(
		constant.Error,
		ctx,
		s,
	)
}
func (d *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rowsAffected := fc()
	errorMessage := "no error"
	if err != nil {
		errorMessage = err.Error()
		log.WithLevel(
			constant.Error,
			ctx,
			fmt.Sprintf(
				"info:\n    - sql: %s\n    - rowsAffected: %v\n    - begin: %s\n    - error: %s",
				sql,
				rowsAffected,
				begin.Format(constant.YyyyMmDdHhMmSsFormat),
				errorMessage,
			),
		)
		return
	}
	log.WithLevel(
		constant.Info,
		ctx,
		fmt.Sprintf(
			"info:\n    - sql: %s\n    - rowsAffected: %v\n    - begin: %s\n    - error: %s",
			sql,
			rowsAffected,
			begin.Format(constant.YyyyMmDdHhMmSsFormat),
			errorMessage,
		),
	)
}
