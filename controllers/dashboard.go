package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Dashboard struct {
	DB *gorm.DB
}

type dashboardArticle struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Image   string `json:"image"`
}

type dashboardResponse struct {
	LatestArticles []dashboardArticle `json:"latestArticles"`
	UsersCount     []struct {
		Role  string `json:"role"`
		Count uint   `json:"count"`
	} `json:"usersCount"`
	CategoriesCount uint `json:"categoriesCount"`
	ArticlesCount   uint `json:"articlesCount"`
}

func (d *Dashboard) GetInfo(ctx *gin.Context) {
	res := dashboardResponse{}
	d.DB.Table("articles").Order("id desc").Limit(5).Find(&res.LatestArticles) // 5 latest articles
	d.DB.Table("users").Select("role, count(*)").Group("role").Scan(&res.UsersCount)
	d.DB.Table("categories").Count(&res.CategoriesCount) // count all categories
	d.DB.Table("articles").Count(&res.ArticlesCount)     // count all articles

	ctx.JSON(http.StatusOK, gin.H{"dashboard": &res})
}
