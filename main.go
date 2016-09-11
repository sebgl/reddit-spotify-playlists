package main

import (
	"flag"

	log "github.com/Sirupsen/logrus"

	"github.com/sebgl/reddify/reddit"
)

var (
	redditUser     = flag.String("reddit-user", "", "Reddit username")
	redditPassword = flag.String("reddit-password", "", "Reddit username")
)

func main() {
	flag.Parse()

	subreddit := "spotify"
	scraper, err := reddit.NewPlaylistScraper(*redditUser, *redditPassword, subreddit)
	if err != nil {
		log.Fatal(err)
	}
	scraper.ScrapLast(1000)
	log.WithField("count", len(scraper.Playlists)).Info("Successfully scraped playlists")
}
