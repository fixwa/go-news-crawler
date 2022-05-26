package main

import (
	"github.com/fixwa/go-news-crawler/crawlers"
	"sync"
)

func main() {
	var waiter sync.WaitGroup

	waiter.Add(1)
	go crawlers.CrawlInfobae(&waiter)

	waiter.Add(1)
	go crawlers.CrawlLaNacion(&waiter)

	waiter.Add(1)
	go crawlers.CrawlClarin(&waiter)

	waiter.Wait()
}
