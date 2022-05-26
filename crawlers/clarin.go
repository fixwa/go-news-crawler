package crawlers

import (
	"github.com/fixwa/go-news-crawler/models"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"time"
)

var (
	clarinArticles map[string]bool
)

func init() {
	existingLinks := models.GetArticlesBySource(models.Infobae)
	clarinArticles = map[string]bool{}

	for _, link := range existingLinks {
		clarinArticles[link] = true
	}
}

func CrawlClarin() {
	log.Println("crawl clarin.com")

	c := colly.NewCollector(
		colly.AllowedDomains("www.clarin.com"),
	)

	q, _ := queue.New(
		1,
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	detailCollector := c.Clone()

	c.OnHTML(".list-news", func(e *colly.HTMLElement) {
		e.ForEach("a.link-new", func(_ int, el *colly.HTMLElement) {
			link := el.Attr("href")
			if _, found := clarinArticles[link]; !found {
				detailCollector.Visit(link)
				clarinArticles[link] = true
			}
		})
	})

	detailCollector.OnHTML(".news.container", func(e *colly.HTMLElement) {
		title := e.ChildText("h1#title")
		date := e.ChildText(".publishedDate")
		thumbnail := e.ChildAttr("img.com-image", "src")
		publishedAt := time.Now()
		content := e.ChildText(".body-nota")

		loc, err := time.LoadLocation("America/Cordoba")
		if err == nil {
			if t, err := time.ParseInLocation("Feb 02, 2022 02:20PM", date, loc); err == nil {
				publishedAt = t
			}
		}

		article := &models.Article{
			Source:      models.Clarin,
			Title:       title,
			Content:     content,
			URL:         e.Request.URL.String(),
			Thumbnail:   thumbnail,
			PublishedAt: publishedAt,
		}

		models.CreateArticle(article)
	})

	q.AddURL("https://www.clarin.com/ultimo-momento/")
	q.Run(c)

	log.Println("Finished clarin.com")
}
