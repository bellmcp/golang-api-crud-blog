package migrations

import (
	"course-go/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1690811624CreateCategoriesTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1690811624",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Category{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("categories").Error
		},
	}
}
