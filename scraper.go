package main

import (
	"context"
	"database/sql"
	"github.com/blkcor/Rss-aggregator/internal/database"
	"github.com/google/uuid"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

/*
Run a long task to auto update feed resources
*/
func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v go routine every %v duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	//waiting for ticker
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("Fail to fetch feeds: %v", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapFeed(db, wg, feed)
		}
		wg.Wait()
	}
}
func scrapFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Fail to mark feed as fetched: %v", err)
		return
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		if !(err == io.EOF) {
			log.Printf("Error fetching feed: %v", err)
			return
		}
	}
	for _, item := range rssFeed.Channel.Item {
		log.Printf("Found post:%v,On feed %s\n", item.Title, feed.Name)
		//save to the database
		description := sql.NullString{String: item.Description}
		parsedTime, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			log.Printf("Fail to parse the time %v, err:%v", item.PubDate, err)
			continue
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: parsedTime,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if !strings.Contains(err.Error(), "duplicate key") {
				log.Printf("Fail to create post: %v", err)
			}
			continue
		}
	}
	log.Printf("feed %s collected, %v post found", feed.Name, len(rssFeed.Channel.Item))
}
