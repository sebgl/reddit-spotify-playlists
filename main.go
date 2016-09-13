package main

import (
	"flag"

	log "github.com/Sirupsen/logrus"
)

var (
	redditUser     = flag.String("reddit-user", "", "Reddit username")
	redditPassword = flag.String("reddit-password", "", "Reddit username")
	subreddit      = flag.String("subreddit", "", "Subreddit to look for playlists")
)

func main() {
	flag.Parse()

	scraper, err := NewPlaylistScraper(*redditUser, *redditPassword, *subreddit)
	if err != nil {
		log.Fatal(err)
	}
	playlists := scraper.ScrapLast(10)
	log.WithField("count", len(playlists)).Info("Successfully scraped playlists")

	spotifyClient := getClient()
	err = getSpotifyPlaylist(spotifyClient, playlists[0].SpotifyURL)
	if err != nil {
		log.Fatal(err)
	}
}
