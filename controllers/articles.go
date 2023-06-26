package controllers

import (
	"course-go/models"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Articles struct {
}

type createArticleForm struct {
	Title string                `form:"title" binding:"required"` // package: validator; change validate -> binding if use with gin
	Body  string                `form:"body" binding:"required"`
	Image *multipart.FileHeader `form:"image" binding:"required"` // file
}

var articles []models.Article = []models.Article{
	{ID: 1, Title: "Title#1", Body: "Body#1"},
	{ID: 2, Title: "Title#2", Body: "Body#2"},
	{ID: 3, Title: "Title#3", Body: "Body#3"},
	{ID: 4, Title: "Title#4", Body: "Body#4"},
	{ID: 5, Title: "Title#5", Body: "Body#5"},
}

func (a *Articles) FindAll(ctx *gin.Context) {
	result := articles
	// ?limit=3
	if limit := ctx.Query("limit"); limit != "" {
		n, _ := strconv.Atoi(limit) // str -> int

		result = articles[0:n] // [0,1,2]
	}
	ctx.JSON(http.StatusOK, gin.H{"articles": result})
}

func (a *Articles) FindOne(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	for _, item := range articles {
		if item.ID == uint(id) {
			ctx.JSON(http.StatusOK, gin.H{"articles": item})
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
}

func (a *Articles) Create(ctx *gin.Context) {
	var form createArticleForm

	if err := ctx.ShouldBind(&form); err != nil { // ShouldBindJSON -> ShouldBind; Form Data
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article := models.Article{
		ID:    uint(len(articles) + 1),
		Title: form.Title,
		Body:  form.Body,
	}

	// get file
	file, _ := ctx.FormFile("image")

	// create path
	// ID -> 8, uploads/articles/8/image.png
	path := "uploads/articles/" + strconv.Itoa(int(article.ID))
	os.MkdirAll(path, 0755)

	// upload file
	destination := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, destination); err != nil {
		fmt.Println(err)
	}

	// attach file to article
	article.Image = os.Getenv("HOST") + "/" + destination

	articles = append(articles, article)

	ctx.JSON(http.StatusCreated, gin.H{"article": article})
}
