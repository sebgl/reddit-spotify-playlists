package main

import (
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/jzelinskie/geddit"
)

var spotifyURLRegex = regexp.MustCompile(`^https:\/\/open\.spotify\.com\/user\/.*\/playlist\/.*$`)

type PlaylistScraper struct {
	Scraped       []*RedditData
	RedditSession *geddit.LoginSession
	Subreddit     string
}

func NewPlaylistScraper(redditLogin string, redditPassword string, subreddit string) (*PlaylistScraper, error) {
	session, err := geddit.NewLoginSession(redditLogin, redditPassword, "gedditAgent v1")
	if err != nil {
		return nil, err
	}
	return &PlaylistScraper{
		Scraped:       make([]*RedditData, 0, 10000),
		RedditSession: session,
		Subreddit:     subreddit,
	}, nil
}

func (ps *PlaylistScraper) ScrapLast(n int) []*RedditData {
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
			return ps.Scraped // no more submissions to parse
		}
		log.WithField("count", len(submissions)).WithField("last", lastSubmission).Info("Got submissions")
		params.After = lastSubmission.FullID
	}
	return ps.Scraped
}

func (ps *PlaylistScraper) parseSubmissions(submissions []*geddit.Submission) (lastSub *geddit.Submission) {
	if len(submissions) == 0 {
		return nil
	}
	for _, s := range submissions {
		isSpotifyPlaylist := spotifyURLRegex.MatchString(s.URL)
		if isSpotifyPlaylist {
			redditData := RedditData{
				SpotifyURL:  s.URL,
				User:        s.Author,
				Score:       s.Score,
				Title:       s.Title,
				Description: s.Selftext,
			}
			ps.Scraped = append(ps.Scraped, &redditData)
		}
	}
	lastSub = submissions[len(submissions)-1]
	return lastSub
}
