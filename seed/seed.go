package seed

import (
	"course-go/config"
	"course-go/migrations"
	"course-go/models"
	"log"
	"math/rand"
	"strconv"

	"github.com/bxcodec/faker/v3"
)

func Load() {
	db := config.GetDB()

	// Clean database
	db.DropTableIfExists("articles", "categories", "migrations")

	migrations.Migrate()

	// Add categories
	log.Print("Creating categories")

	numOfCategories := 20
	categories := make([]models.Category, 0, numOfCategories)

	for i := 1; i <= numOfCategories; i++ {
		category := models.Category{
			Name: faker.Word(),
			Desc: faker.Paragraph(),
		}

		db.Create(&category)
		categories = append(categories, category)
	}

	// Add articles
	log.Print("Creating articles")
	numOfArticles := 50
	articles := make([]models.Article, 0, numOfArticles)

	for i := 1; i <= numOfArticles; i++ {
		article := models.Article{
			Title:      faker.Word(),
			Excerpt:    faker.Sentence(),
			Body:       faker.Paragraph(),
			Image:      "https://source.unsplash.com/random/300x300?" + strconv.Itoa(i),
			CategoryID: uint(rand.Intn(numOfArticles) + 1),
		}

		db.Create(&article)
		articles = append(articles, article)
	}
}
