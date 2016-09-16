package main

import (
	"fmt"
	"regexp"

	"github.com/zmb3/spotify"
)

var spotifyURLIDsRegex = regexp.MustCompile(`^https:\/\/open\.spotify\.com\/user\/(.*)\/playlist\/(.*)$`)

func getSpotifyPlaylist(c *spotify.Client, url string) (*spotify.FullPlaylist, error) {
	userID, playlistID, err := IDsFromURL(url)
	if err != nil {
		return nil, err
	}
	playlist, err := c.GetPlaylist(userID, spotify.ID(playlistID))
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

func IDsFromURL(url string) (userID, playlistID string, err error) {
	groups := spotifyURLIDsRegex.FindStringSubmatch(url)
	if len(groups) < 3 {
		return "", "", fmt.Errorf("Failed to parse url %s", url)
	}
	return groups[1], groups[2], nil
}
