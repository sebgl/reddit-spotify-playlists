package main

import (
	"fmt"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/zmb3/spotify"
)

var spotifyURLIDsRegex = regexp.MustCompile(`^https:\/\/open\.spotify\.com\/user\/(.*)\/playlist\/(.*)$`)

func getSpotifyPlaylist(c *spotify.Client, url string) error {
	userID, playlistID, err := IDsFromURL(url)
	if err != nil {
		return err
	}
	playlist, err := c.GetPlaylist(userID, spotify.ID(playlistID))
	if err != nil {
		return err
	}
	log.Infof("Got playlist %+v", playlist)
	return nil
}

func IDsFromURL(url string) (userID, playlistID string, err error) {
	groups := spotifyURLIDsRegex.FindStringSubmatch(url)
	if len(groups) < 3 {
		return "", "", fmt.Errorf("Failed to parse url %s", url)
	}
	return groups[1], groups[2], nil
}
