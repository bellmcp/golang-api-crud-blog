package routes

import (
	"course-go/controllers"

	"github.com/gin-gonic/gin"
)

func Serve(r *gin.Engine) {
	articlesGroup := r.Group("/api/v1/articles")
	// GET /api/v1/articles
	// GET /api/v1/articles/:id
	// POST /api/v1/articles
	// PATCH /api/v1/articles/:id
	// DELETE /api/v1/articles/:id

	articleController := controllers.Articles{}

	{
		articlesGroup.GET("", articleController.FindAll)
		articlesGroup.GET("/:id", articleController.FindOne)
		articlesGroup.POST("", articleController.Create)
	}
}
