package main

import (
	"encoding/json"
	"os"
)

type Playlist struct {
	RedditData  *RedditData
	SpotifyData *SpotifyData
}

type RedditData struct {
	User        string
	Score       int
	Title       string
	Description string
	SpotifyURL  string
}

type SpotifyData struct {
	Description string
	Followers   int
	ID          string
	ImagesURL   []string
	Name        string
	OwnerName   string
	OwnerID     string
	Tracks      []Track
}

type Track struct {
	Album      Album
	Popularity int
	Artists    []Artist
	Duration   int
	ID         string
	Name       string
	PreviewURL string
}

type Artist struct {
	Name string
	ID   string
}

type Album struct {
	Name string
	ID   string
}

func writeToFile(playlists []Playlist, outputFilename string) error {
	b, err := json.Marshal(playlists)
	if err != nil {
		return err
	}
	fo, err := os.Create(outputFilename)
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
