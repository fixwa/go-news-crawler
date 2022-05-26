package crawlers

import (
	"github.com/fixwa/go-news-crawler/models"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"strings"
	"sync"
	"time"
)

var (
	infobaeArticles map[string]bool
)

func init() {
	existingLinks := models.GetArticlesBySource(models.Infobae)
	infobaeArticles = map[string]bool{}

	for _, link := range existingLinks {
		infobaeArticles[link] = true
	}
}

func CrawlInfobae(w *sync.WaitGroup) {
	log.Println("crawl infobae.com")

	c := colly.NewCollector(
		colly.AllowedDomains("www.infobae.com"),
	)

	q, _ := queue.New(
		1,
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	detailCollector := c.Clone()

	c.OnHTML(".page-container", func(e *colly.HTMLElement) {
		e.ForEach("a.nd-feed-list-card", func(_ int, el *colly.HTMLElement) {
			link := el.Attr("href")
			if strings.Index(link, "/2022/") == -1 {
				return
			}

			link = "https://www.infobae.com" + link
			if _, found := infobaeArticles[link]; !found {
				detailCollector.Visit(link)
				infobaeArticles[link] = true
			}
		})
	})

	detailCollector.OnHTML(".article-section", func(e *colly.HTMLElement) {
		title := e.ChildText(".article-headline")
		date := e.ChildText(".byline-datetime")
		thumbnail := e.ChildAttr(".visual__image > img", "src")
		publishedAt := time.Now()

		var paragraphs []string
		e.ForEach("p.paragraph", func(_ int, el *colly.HTMLElement) {
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
			Source:      models.Infobae,
			Title:       title,
			Content:     content,
			URL:         e.Request.URL.String(),
			Thumbnail:   thumbnail,
			PublishedAt: publishedAt,
		}

		models.CreateArticle(article)
	})

	q.AddURL("https://www.infobae.com/ultimas-noticias/")

	// Consume
	q.Run(c)
	log.Println("Finished infobae.com")
	w.Done()
}

// ConvertInfobaeDate @todo this is exported to allow testing (make private).
func ConvertInfobaeDate(str string) string {
	search := []string{"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"}
	replace := []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}

	for i, rpl := range search {
		r := strings.NewReplacer(rpl, replace[i])
		str = r.Replace(str)
	}

	dateParts := strings.Split(str, " ")

	// add a 0 to the month
	if len(dateParts[2]) < 2 {
		dateParts[2] = "0" + dateParts[2]
	}

	// add a 0 to the day
	if len(dateParts[0]) < 2 {
		dateParts[0] = "0" + dateParts[0]
	}

	return dateParts[4] + "-" + dateParts[2] + "-" + dateParts[0]
}
