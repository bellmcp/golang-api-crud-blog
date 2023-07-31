package migrations

import (
	"course-go/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1690813635AddCategoryIDToArticles() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1690813635",
		Migrate: func(tx *gorm.DB) error {
			err := tx.AutoMigrate(&models.Article{}).Error

			var articles []models.Article
			tx.Unscoped().Find(&articles)
			for _, article := range articles {
				article.CategoryID = 2
				tx.Save(&article)
			}

			return err
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Model(&models.Article{}).DropColumn("category_id").Error
		},
	}
}
