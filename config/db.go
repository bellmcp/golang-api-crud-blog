package config

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"log"
)

var db *gorm.DB

func InitDB() {
	var err error
	db, err = gorm.Open("postgres", os.Getenv("DATABASE_CONNECTION"))

	if err != nil {
		log.Fatal(err)
	}

	db.LogMode(gin.Mode() == gin.DebugMode)
}

func CloseDB() {
	db.Close()
}

func GetDB() *gorm.DB {
	return db
}
