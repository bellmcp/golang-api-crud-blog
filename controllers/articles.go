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
	Title      string                `form:"title" binding:"required"` // package: validator; change validate -> binding if use with gin
	Body       string                `form:"body" binding:"required"`
	Excerpt    string                `form:"excerpt" binding:"required"`
	CategoryID uint                  `form:"categoryId" binding:"required"`
	Image      *multipart.FileHeader `form:"image" binding:"required"` // file
}

type updateArticleForm struct {
	Title      string                `form:"title"` // package: validator; change validate -> binding if use with gin
	Body       string                `form:"body"`
	Excerpt    string                `form:"excerpt"`
	CategoryID uint                  `form:"categoryId"`
	Image      *multipart.FileHeader `form:"image"` // file
}

type articleResponse struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Excerpt    string `json:"excerpt"`
	Body       string `json:"body"`
	Image      string `json:"image"`
	CategoryID uint   `json:"category_id"`
	Category   struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
}

type articlesPaging struct {
	Items  []articleResponse `json:"items"`
	Paging *pagingResult     `json:"paging"`
}

func (a *Articles) FindAll(ctx *gin.Context) {
	var articles []models.Article

	// a.DB.Find(&articles)
	pagination := pagination{
		ctx:     ctx,
		query:   a.DB.Preload("Category").Order("id desc"),
		records: &articles,
	}
	paging := pagination.paginate()

	var serializedArticles []articleResponse
	copier.Copy(&serializedArticles, &articles)
	ctx.JSON(http.StatusOK, gin.H{"articles": articlesPaging{
		Items:  serializedArticles,
		Paging: paging,
	}})
}

func (a *Articles) FindOne(ctx *gin.Context) {
	article, err := a.findArticleByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	serializeArticleResponse := articleResponse{}
	copier.Copy(&serializeArticleResponse, &article)
	ctx.JSON(http.StatusOK, gin.H{"article": serializeArticleResponse})
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

	serializedArticle := articleResponse{}
	copier.Copy(&serializedArticle, &article)

	ctx.JSON(http.StatusCreated, gin.H{"article": serializedArticle})
}

func (a *Articles) Update(ctx *gin.Context) {
	var form updateArticleForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	article, err := a.findArticleByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := a.DB.Model(&article).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	a.setArticleImage(ctx, article)

	serializedArticle := articleResponse{}
	copier.Copy(&serializedArticle, &article)
	ctx.JSON(http.StatusOK, gin.H{"article": serializedArticle})
}

func (a *Articles) Delete(ctx *gin.Context) {
	article, err := a.findArticleByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	a.DB.Unscoped().Delete(&article) // hard delete
	ctx.Status(http.StatusNoContent)
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

func (a *Articles) findArticleByID(ctx *gin.Context) (*models.Article, error) {
	var article models.Article
	id := ctx.Param("id")

	if err := a.DB.Preload("Category").First(&article, id).Error; err != nil {
		return nil, err
	}

	return &article, nil
}
