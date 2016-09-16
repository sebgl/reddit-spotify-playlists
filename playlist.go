package main

import "github.com/zmb3/spotify"

type Playlist struct {
	RedditUser        string
	RedditScore       int
	RedditTitle       string
	RedditDescription string
	SpotifyURL        string
	SpotifyPlaylist   *spotify.FullPlaylist
}
