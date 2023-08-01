package routes

import (
	"course-go/config"
	"course-go/controllers"
	"course-go/middleware"

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
	authenticate := middleware.Authenticate().MiddlewareFunc()
	authorize := middleware.Authorize()

	authGroup := v1.Group("auth")
	authController := controllers.Auth{DB: db}
	{
		authGroup.POST("/sign-up", authController.Signup)
		authGroup.POST("/sign-in", middleware.Authenticate().LoginHandler)
		authGroup.GET("/profile", authenticate, authController.GetProfile)
		authGroup.PATCH("/profile", authenticate, authController.UpdateProfile)
	}

	usersController := controllers.Users{DB: db}
	usersGroup := v1.Group("users")
	usersGroup.Use(authenticate, authorize) // apply authenticate and follows by authorize to all routes under user group
	{
		usersGroup.GET("", usersController.FindAll)
		usersGroup.POST("", usersController.Create)
		usersGroup.GET("/:id", usersController.FindOne)
		usersGroup.PATCH("/:id", usersController.Update)
		usersGroup.DELETE("/:id", usersController.Delete)
		usersGroup.PATCH("/:id/promote", usersController.Promote)
		usersGroup.PATCH("/:id/demote", usersController.Demote)
	}

	articleController := controllers.Articles{
		DB: db,
	}
	articlesGroup := v1.Group("articles")
	articlesGroup.Use(authenticate, authorize)
	articlesGroup.GET("", articleController.FindAll)     // skip authenticate and authorize (any user allows)
	articlesGroup.GET("/:id", articleController.FindOne) // skip authenticate and authorize (any user allows)
	{
		articlesGroup.PATCH("/:id", articleController.Update)
		articlesGroup.DELETE("/:id", articleController.Delete)
		articlesGroup.POST("", authenticate, articleController.Create)
	}

	categoryController := controllers.Categories{
		DB: db,
	}
	categoriesGroup := v1.Group("categories")
	categoriesGroup.Use(authenticate, authorize)
	categoriesGroup.GET("", categoryController.FindAll)     // skip authenticate and authorize (any user allows)
	categoriesGroup.GET("/:id", categoryController.FindOne) // skip authenticate and authorize (any user allows)
	{
		categoriesGroup.PATCH("/:id", categoryController.Update)
		categoriesGroup.DELETE("/:id", categoryController.Delete)
		categoriesGroup.POST("", categoryController.Create)
	}
}
