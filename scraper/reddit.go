package scraper

import (
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/jzelinskie/geddit"
)

var spotifyURLRegex = regexp.MustCompile(`^https:\/\/open\.spotify\.com\/user\/.*\/playlist\/.*$`)

type RedditScraper struct {
	Scraped       []*RedditData
	RedditSession *geddit.LoginSession
	Subreddit     string
}

func NewRedditScraper(redditLogin string, redditPassword string, subreddit string) (*RedditScraper, error) {
	session, err := geddit.NewLoginSession(redditLogin, redditPassword, "gedditAgent v1")
	if err != nil {
		return nil, err
	}
	return &RedditScraper{
		Scraped:       make([]*RedditData, 0, 10000),
		RedditSession: session,
		Subreddit:     subreddit,
	}, nil
}

func (ps *RedditScraper) ScrapLast(n int) []*RedditData {
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

func (ps *RedditScraper) parseSubmissions(submissions []*geddit.Submission) (lastSub *geddit.Submission) {
	if len(submissions) == 0 {
		return nil
	}
	for _, s := range submissions {
		isSpotifyPlaylist := spotifyURLRegex.MatchString(s.URL)
		if isSpotifyPlaylist {
			comments, err := ps.getComments(s)
			if err != nil {
				log.WithError(err).Info("Fail to get submission comments")
				comments = []Comment{}
			}
			redditData := RedditData{
				SpotifyURL:  s.URL,
				User:        s.Author,
				UpVotes:     s.Ups,
				Title:       s.Title,
				Description: s.Selftext,
				Date:        int(s.DateCreated),
				Comments:    comments,
			}
			ps.Scraped = append(ps.Scraped, &redditData)
		}
	}
	lastSub = submissions[len(submissions)-1]
	return lastSub
}

// get comments from first-level only
// TODO: recursively parse comment's comments
func (ps *RedditScraper) getComments(submission *geddit.Submission) ([]Comment, error) {
	redditComments, err := ps.RedditSession.Comments(submission)
	if err != nil {
		return nil, err
	}
	comments := make([]Comment, len(redditComments))
	for i, c := range redditComments {
		comments[i].Author = c.Author
		comments[i].Text = c.Body
		comments[i].UpVotes = int(c.UpVotes)
	}
	return comments, nil
}
