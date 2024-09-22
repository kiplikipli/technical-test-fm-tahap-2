package database

import (
	"log"
	"os"

	"github.com/kiplikipli/technical-test-fm-tahap-2/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() {
	var err error // define error here to prevent overshadowing the global DB

	env := os.Getenv("DATABASE_URL")
	DB, err = gorm.Open(sqlite.Open(env), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = DB.AutoMigrate(&entity.User{}, &entity.Transaction{})
	if err != nil {
		log.Fatal(err)
	}

}
