package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func InitDatabaseConnection() (db *gorm.DB, err error) {
	// if databaseUsername, ok := os.LookupEnv("DATABASE_USERNAME"); !ok {
	// 	databaseUsername = "root"
	// }
	var databaseUsername = os.Getenv("DATABASE_USERNAME")
	if databaseUsername == "" {
		databaseUsername = "postgres"
	}
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	if databasePassword == "" {
		databasePassword = "postgres"
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
	// dsn := fmt.Sprintf(
	// 	"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	// 	databaseUsername,
	// 	databasePassword,
	// 	databaseHost,
	// 	databasePort,
	// 	databaseName,
	// )
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		databaseHost,
		databaseUsername,
		databasePassword,
		databaseName,
		databasePort,
	)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	// migrateErr := db.AutoMigrate(&model.User{})
	// if migrateErr != nil {
	// 	return nil, migrateErr
	// }
	return db, err
}
