package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func DatabaseInit() {
	host := "localhost"
	user := "postgres"
	password := "secret"
	dbName := "db_crowdfunding"
	port := 5432

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbName, port)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	return db
}