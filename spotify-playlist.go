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
	// retrieve next pages of tracks
	// TODO: this would be better implemented in the spotify library
	MaxNumberOfTracksToRetrieve := 5000 // set an upper-limit for XXXL playlists
	currentOffset := playlist.Tracks.Limit
	for currentOffset < playlist.Tracks.Total && currentOffset < MaxNumberOfTracksToRetrieve {
		nextTracks, err := sp.c.GetPlaylistTracksOpt(playlist.Owner.ID, playlist.ID,
			&spotify.Options{Offset: &currentOffset}, "")
		if err != nil {
			log.WithError(err).Errorf("Fail to retrieve tracks from offset %d", playlist.Tracks.Limit)
		}
		newOffset := currentOffset + len(nextTracks.Tracks)
		log.WithField("currentOffset", currentOffset).WithField("newOffset", newOffset).WithField("total", playlist.Tracks.Total).
			Info("Retrieved additional tracks")
		playlist.Tracks.Tracks = append(playlist.Tracks.Tracks, nextTracks.Tracks...)
		currentOffset = newOffset
	}

	// extract relevant tracks data
	tracks := make([]Track, len(playlist.Tracks.Tracks))
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
