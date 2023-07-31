package controllers

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type Articles struct {
}

type createArticleForm struct {
	Title string                `form:"title" binding:"required"` // package: validator; change validate -> binding if use with gin
	Body  string                `form:"body" binding:"required"`
	Image *multipart.FileHeader `form:"image" binding:"required"` // file
}

func (a *Articles) FindAll(ctx *gin.Context) {

}

func (a *Articles) FindOne(ctx *gin.Context) {

}

func (a *Articles) Create(ctx *gin.Context) {

}
