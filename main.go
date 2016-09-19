package main

import (
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"
)

var (
	redditUser     = flag.String("reddit-user", "", "Reddit username")
	redditPassword = flag.String("reddit-password", "", "Reddit username")
	subreddit      = flag.String("subreddit", "", "Subreddit to look for playlists")
	nSubmissions   = flag.Int("last", 10, "Number of submissions to parse (rounded up to Reddit page size)")
	esURL          = flag.String("elasticsearch", "http://localhost:9200", "Elasticsearch URL")
)

func main() {
	flag.Parse()

	checkEnvVars("SPOTIFY_ID", "SPOTIFY_SECRET")

	// scrap playlist submissions from reddit
	scraper, err := NewPlaylistScraper(*redditUser, *redditPassword, *subreddit)
	if err != nil {
		log.WithError(err).Fatal("Unable to scrap playlists from reddit")
	}
	dataFromReddit := scraper.ScrapLast(*nSubmissions)
	log.WithField("count", len(dataFromReddit)).Info("Successfully scraped reddit submissions")

	// get playlists data from spotify
	playlists := make([]Playlist, 0, len(dataFromReddit))
	spotifyClient := getSpotifyClient()
	for _, p := range dataFromReddit {
		sp, err := getSpotifyData(spotifyClient, p.SpotifyURL)
		if err != nil {
			log.WithError(err).Fatal("Unable to retrieve playlist data from spotify")
		}
		playlists = append(playlists, Playlist{
			RedditData:  p,
			SpotifyData: sp,
		})
	}

	// write playlists data into elasticsearch
	esWriter := NewElasticsearchWriter(*esURL)
	err = esWriter.Write(playlists)
	if err != nil {
		log.WithError(err).Fatal("Unable to send data to elasticsearch")
	}
	log.Info("Data written to elasticsearch")
}

func checkEnvVars(vars ...string) {
	for _, v := range vars {
		if os.Getenv(v) == "" {
			log.Fatalf("$%s must be defined", v)
		}
	}
}
