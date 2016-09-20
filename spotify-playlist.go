package main

import (
	"fmt"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/zmb3/spotify"
)

var spotifyURLIDsRegex = regexp.MustCompile(`^https:\/\/open\.spotify\.com\/user\/(.*)\/playlist\/(.*)$`)

type SpotifyScraper struct {
	c *spotify.Client
}

func (sp SpotifyScraper) GetSpotifyData(url string) (*SpotifyData, error) {
	userID, playlistID, err := IDsFromURL(url)
	if err != nil {
		return nil, err
	}
	log.WithField("ID", playlistID).Info("Getting Spotify playlist data")
	playlist, err := sp.c.GetPlaylist(userID, spotify.ID(playlistID))
	if err != nil {
		return nil, err
	}
	imagesURL := sp.extractImagesURL(playlist)
	tracks := sp.extractTracks(playlist)
	spotifyData := SpotifyData{
		Description: playlist.Description,
		Followers:   int(playlist.Followers.Count),
		ID:          playlist.ID.String(),
		ImagesURL:   imagesURL,
		Name:        playlist.Name,
		OwnerID:     playlist.Owner.ID,
		OwnerName:   playlist.Owner.DisplayName,
		Tracks:      tracks,
	}
	return &spotifyData, nil
}

func (sp SpotifyScraper) extractImagesURL(playlist *spotify.FullPlaylist) []string {
	imagesURL := make([]string, len(playlist.Images))
	for i, img := range playlist.Images {
		imagesURL[i] = img.URL
	}
	return imagesURL
}

func (sp SpotifyScraper) extractTracks(playlist *spotify.FullPlaylist) []Track {
	tracks := make([]Track, playlist.Tracks.Total)
	// TODO: get next page tracks
	// sadly, not available in zmb3/spotify
	for i, t := range playlist.Tracks.Tracks {
		tracks[i] = Track{
			Album: Album{
				Name: t.Track.Album.Name,
				ID:   t.Track.Album.ID.String(),
			},
			Popularity: t.Track.Popularity,
			Artists:    extractArtistsFromTrack(t),
			Duration:   t.Track.Duration,
			ID:         t.Track.ID.String(),
			Name:       t.Track.Name,
			PreviewURL: t.Track.PreviewURL,
		}
	}
	return tracks
}

func extractArtistsFromTrack(t spotify.PlaylistTrack) []Artist {
	artists := make([]Artist, len(t.Track.Artists))
	for i, artist := range t.Track.Artists {
		artists[i] = Artist{
			Name: artist.Name,
			ID:   artist.ID.String(),
		}
	}
	return artists
}

func IDsFromURL(url string) (userID, playlistID string, err error) {
	groups := spotifyURLIDsRegex.FindStringSubmatch(url)
	if len(groups) < 3 {
		return "", "", fmt.Errorf("Failed to parse url %s", url)
	}
	return groups[1], groups[2], nil
}
