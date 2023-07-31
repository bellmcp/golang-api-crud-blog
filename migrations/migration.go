package migrations

import (
	"course-go/config"
	"log"

	"gopkg.in/gormigrate.v1"
)

func Migrate() {
	db := config.GetDB()
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		m1690787686CreateArticlesTable(),
		m1690811624CreateCategoriesTable(),
		m1690813635AddCategoryIDToArticles(),
		m1690818012CreateUsersTable(),
		m1690824582AddUserIDToArticles(),
	})

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
}
