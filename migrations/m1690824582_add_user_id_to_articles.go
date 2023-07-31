package migrations

import (
	"course-go/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1690824582AddUserIDToArticles() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1690824582",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Article{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Model(&models.Article{}).DropColumn("user_id").Error
		},
	}
}
