package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(config MainConfig) *gorm.DB {
	dbConfig := config.Database

	connection := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", dbConfig.Host, dbConfig.DBUser, dbConfig.DBName, dbConfig.Port, dbConfig.DBPass)

	db, err := gorm.Open(postgres.Open(connection), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}

	return db
}
