package routes

import (
	"course-go/config"
	"course-go/controllers"

	"github.com/gin-gonic/gin"
)

func Serve(r *gin.Engine) {
	// GET /api/v1/articles
	// GET /api/v1/articles/:id
	// POST /api/v1/articles
	// PATCH /api/v1/articles/:id
	// DELETE /api/v1/articles/:id

	db := config.GetDB()
	v1 := r.Group("/api/v1")

	authGroup := v1.Group("auth")
	authController := controllers.Auth{DB: db}
	{
		authGroup.POST("/sign-up", authController.Signup)
	}

	articlesGroup := v1.Group("articles")
	articleController := controllers.Articles{
		DB: db,
	}
	{
		articlesGroup.GET("", articleController.FindAll)
		articlesGroup.GET("/:id", articleController.FindOne)
		articlesGroup.PATCH("/:id", articleController.Update)
		articlesGroup.DELETE("/:id", articleController.Delete)
		articlesGroup.POST("", articleController.Create)
	}

	categoriesGroup := v1.Group("categories")
	categoryController := controllers.Categories{
		DB: db,
	}
	{
		categoriesGroup.GET("", categoryController.FindAll)
		categoriesGroup.GET("/:id", categoryController.FindOne)
		categoriesGroup.PATCH("/:id", categoryController.Update)
		categoriesGroup.DELETE("/:id", categoryController.Delete)
		categoriesGroup.POST("", categoryController.Create)
	}
}
