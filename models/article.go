package models

import (
	"fmt"
	"time"
)

const (
	Infobae  = 1
	LaNacion = 2
	Clarin   = 3
)

type Article struct {
	ID          int `gorm:"primary_key;"`
	Title       string
	Content     string
	Source      int
	URL         string
	Thumbnail   string
	PublishedAt time.Time
	CreatedAt   time.Time
}

var NewsSources map[int]string

func init() {
	NewsSources = map[int]string{
		Infobae:  "infobae.com",
		LaNacion: "lanacion.com.ar",
		Clarin:   "clarin.com",
	}
}

func CreateArticle(input *Article) (*Article, error) {
	article := input

	db := ConnectDatabase()
	defer db.Close()

	// Create the user
	db.Create(article)

	if article.ID <= 0 {
		return nil, fmt.Errorf("Could not create the Article.")
	}

	return article, nil
}

func GetArticlesBySource(sourceID int) []string {
	var links []string

	db := ConnectDatabase()
	defer db.Close()

	db.Table("articles").Where("source = ?", sourceID).Pluck("url", &links)

	return links
}
