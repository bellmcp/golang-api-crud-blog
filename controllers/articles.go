package controllers

import (
	"course-go/models"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Articles struct {
	DB *gorm.DB
}

type createArticleForm struct {
	Title   string                `form:"title" binding:"required"` // package: validator; change validate -> binding if use with gin
	Body    string                `form:"body" binding:"required"`
	Excerpt string                `form:"excerpt" binding:"required"`
	Image   *multipart.FileHeader `form:"image" binding:"required"` // file
}

type createdArticleResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Body    string `json:"body"`
	Image   string `json:"image"`
}

func (a *Articles) FindAll(ctx *gin.Context) {

}

func (a *Articles) FindOne(ctx *gin.Context) {

}

func (a *Articles) Create(ctx *gin.Context) {
	var form createArticleForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// form => article
	var article models.Article
	copier.Copy(&article, &form) // must use & to refer to original data

	// article => db
	if err := a.DB.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	a.setArticleImage(ctx, &article)

	serializedArticle := createdArticleResponse{}
	copier.Copy(&serializedArticle, &article)

	ctx.JSON(http.StatusCreated, gin.H{"article": serializedArticle})
}

func (a *Articles) setArticleImage(ctx *gin.Context, article *models.Article) error {
	file, err := ctx.FormFile("image")
	if err != nil || file == nil {
		return err
	}

	if article.Image != "" {
		// http://localhost/upload/articles/<ID>/image.png
		// 1. /upload/articles/<ID>/image.png
		article.Image = strings.Replace(article.Image, os.Getenv("HOST"), "", 1)
		// 2. <WD>/upload/articles/<ID>/image.png
		pwd, _ := os.Getwd()
		// 3. Remove <WD>/upload/articles/<ID>/image.png
		os.Remove(pwd + article.Image)

	}

	path := "uploads/articles/" + strconv.Itoa(int(article.ID))
	os.MkdirAll(path, 0755)
	filename := path + "/" + file.Filename
	if err = ctx.SaveUploadedFile(file, filename); err != nil {
		return err
	}

	article.Image = os.Getenv("HOST") + "/" + filename
	a.DB.Save(article)

	return nil
}