package main

import (
	"encoding/json"
	"flag"
	"os"

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
	playlists := scraper.ScrapLast(2)
	log.WithField("count", len(playlists)).Info("Successfully scraped playlists")

	spotifyClient := getSpotifyClient()
	for i, p := range playlists {
		sp, err := getSpotifyPlaylist(spotifyClient, p.SpotifyURL)
		if err != nil {
			log.Fatal(err)
		}
		playlists[i].SpotifyPlaylist = sp
	}
	err = ToElasticSearch(playlists)
	if err != nil {
		log.WithError(err).Fatal("Unable to send data to elasticsearch")
	}
}

func writeToFile(playlists []Playlist) error {
	b, err := json.Marshal(playlists)
	if err != nil {
		log.WithError(err).Fatal("Unable to marshal playlist as json")
		return err
	}
	fo, err := os.Create("output.json")
	if err != nil {
		return err
	}
	defer fo.Close()
	_, err = fo.Write(b)
	if err != nil {
		return err
	}
	return nil
}
