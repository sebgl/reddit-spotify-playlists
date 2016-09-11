package reddit

import (
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/jzelinskie/geddit"
	"github.com/sebgl/reddify/model"
)

var spotifyURLRegex = regexp.MustCompile(`^https:\/\/open\.spotify\.com\/user\/.*\/playlist\/.*$`)

type PlaylistScraper struct {
	Playlists     []model.Playlist
	RedditSession *geddit.LoginSession
	Subreddit     string
}

func NewPlaylistScraper(redditLogin string, redditPassword string, subreddit string) (*PlaylistScraper, error) {
	session, err := geddit.NewLoginSession(redditLogin, redditPassword, "gedditAgent v1")
	if err != nil {
		return nil, err
	}
	return &PlaylistScraper{
		Playlists:     make([]model.Playlist, 0, 10000),
		RedditSession: session,
		Subreddit:     subreddit,
	}, nil
}

func (ps *PlaylistScraper) ScrapLast(n int) {
	batchSize := 100
	params := geddit.ListingOptions{
		Limit: batchSize,
	}
	nLoops := int(n/batchSize) + 1
	for i := 0; i < nLoops; i++ {
		submissions, err := ps.RedditSession.SubredditSubmissions(ps.Subreddit, geddit.NewSubmissions, params)
		if err != nil {
			log.Fatal(err)
		}
		lastSubmission := ps.parseSubmissions(submissions)
		if lastSubmission == nil {
			return // no more submissions to parse
		}
		log.WithField("count", len(submissions)).WithField("last", lastSubmission).Info("Got submissions")
		params.After = lastSubmission.FullID
	}
}

func (ps *PlaylistScraper) parseSubmissions(submissions []*geddit.Submission) (lastSub *geddit.Submission) {
	if len(submissions) == 0 {
		return nil
	}
	for _, s := range submissions {
		isSpotifyPlaylist := spotifyURLRegex.MatchString(s.URL)
		if isSpotifyPlaylist {
			ps.Playlists = append(ps.Playlists, model.Playlist{
				SpotifyURL:        s.URL,
				RedditUser:        s.Author,
				RedditScore:       s.Score,
				RedditTitle:       s.Title,
				RedditDescription: s.Selftext,
			})
		}
	}
	lastSub = submissions[len(submissions)-1]
	return lastSub
}
