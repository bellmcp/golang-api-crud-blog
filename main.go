package main

import (
	"course-go/config"
	"course-go/migrations"
	"course-go/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// read env variable
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDB()
	defer config.CloseDB()
	migrations.Migrate()

	// articles
	r := gin.Default()

	// serve static file
	// http://localhost:8080/uploads/articles/6/bell2.png"
	r.Static("/uploads", "./uploads")

	uploadDirs := [...]string{"articles", "users"}
	for _, dir := range uploadDirs {
		os.MkdirAll("uploads/"+dir, 0755) // user r/w/e
	}

	routes.Serve(r)

	r.Run(":" + os.Getenv("PORT")) // specific port
}
