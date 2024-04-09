package config

import (
	"fmt"
	"github.com/charmbracelet/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"time"
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
			"Database connect info:\n    - databaseHost: %s\n    - databaseUsername: %s\n    - databaseName: %s\n    - databasePort: %s",
			databaseHost,
			databaseUsername,
			databaseName,
			databasePort,
		),
	)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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

	/*
		usersMigrateErr := db.AutoMigrate(&model.Users{})
		if usersMigrateErr != nil {
			return nil, usersMigrateErr
		}

		debitUserMigrateErr := db.AutoMigrate(&model.DebitUser{})
		if debitUserMigrateErr != nil {
			return nil, debitUserMigrateErr
		}

		listOfExpensesMigrateErr := db.AutoMigrate(&model.ListOfExpenses{})
		if listOfExpensesMigrateErr != nil {
			return nil, listOfExpensesMigrateErr
		}
	*/

	return db, err
}
