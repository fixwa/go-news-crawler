package crawlers

import (
	"github.com/fixwa/go-news-crawler/models"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"strings"
	"time"
)

var (
	laNacionArticles map[string]bool
)

func init() {
	existingLinks := models.GetArticlesBySource(models.Infobae)
	laNacionArticles = map[string]bool{}

	for _, link := range existingLinks {
		laNacionArticles[link] = true
	}
}

func CrawlLaNacion() {
	log.Println("crawl lanacion.com.ar")

	c := colly.NewCollector(
		colly.AllowedDomains("www.lanacion.com.ar"),
	)

	q, _ := queue.New(
		1,
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	detailCollector := c.Clone()

	c.OnHTML("#content-main", func(e *colly.HTMLElement) {
		e.ForEach("a.com-link", func(_ int, el *colly.HTMLElement) {

			link := el.Attr("href")

			if strings.Index(link, "2022/") == -1 {
				return
			}

			link = "https://www.lanacion.com.ar" + link

			if _, found := laNacionArticles[link]; !found {
				detailCollector.Visit(link)
				laNacionArticles[link] = true
			}
		})

	})

	detailCollector.OnHTML(".nota", func(e *colly.HTMLElement) {
		title := e.ChildText("h1.com-title")
		date := e.ChildText(".com-date")
		thumbnail := e.ChildAttr("img.com-image", "src")
		publishedAt := time.Now()

		var paragraphs []string
		e.ForEach("p.com-paragraph", func(_ int, el *colly.HTMLElement) {
			paragraphs = append(paragraphs, el.Text)
		})
		content := strings.Join(paragraphs, "\n\n")

		loc, err := time.LoadLocation("America/Cordoba")
		if err == nil {
			if t, err := time.ParseInLocation("Feb 02, 2022 02:20PM", ConvertInfobaeDate(date), loc); err == nil {
				publishedAt = t
			}
		}

		article := &models.Article{
			Source:      models.LaNacion,
			Title:       title,
			Content:     content,
			URL:         e.Request.URL.String(),
			Thumbnail:   thumbnail,
			PublishedAt: publishedAt,
		}

		models.CreateArticle(article)
	})

	q.AddURL("https://www.lanacion.com.ar/ultimas-noticias/")
	q.Run(c)

	log.Println("Finished lanacion.com.ar")
}
