// Package tweetstorage is responsible for integration with tweet storage database.
// Current implementation can create new tweet and retrieve tweets by tag.
// NOTE: Current implementation can't handle tweets without tags due to schema limitations.
//
// Dynamo-like databases are query-centric, not entity-relation centric, so
// we have to create single-purpose tables, designed to solve isolated set of problems.


package tweetstorage

import (
	"github.com/ozapinq/twitter"
	"fmt"
	"time"
	"math/rand"
)

type DB interface {
	appendTweetToTag(tag string, tweet *twitter.Tweet) error
	removeTweetFromTag(tag string, tweet *twitter.Tweet) error
	listTweetsByTag(tag string, before int64, count int) ([]twitter.Tweet, error)
}


type TweetStorage struct {
	db      DB
}

// NewTweetStorage creates instance of tweetStorage backed with specified DB.
func NewTweetStorage(db DB) *TweetStorage {
	return &TweetStorage{db}
}

// New creates new tweet with specified author and text.
// Returns created Tweet and any creation error encountered.
// TODO: add tweet id collision detection / avoidance
func (t *TweetStorage) New(author, text string) (*twitter.Tweet, error) {
	id := generateTweetId()
	createdAt := time.Now().UnixNano()
	tags := getTags(text)
	if tags == nil {
		return nil, ErrNoTags
	}

	newTweet := twitter.Tweet{id, author, text, tags, createdAt}

	appendedTags := make([]string, 0, len(tags))
	var err error
	for _, tag := range tags {
                // sequential writes are ok for PoC, but should be
                // reimplemented using goroutines
		err = t.db.appendTweetToTag(tag, &newTweet)
		if err == nil {
			appendedTags = append(appendedTags, tag)
		} else {
			break
		}
	}

	// make cleanup in case of insert error
	if err != nil {
		for _, tag := range appendedTags {
			t.db.removeTweetFromTag(tag, &newTweet)
		}
		return nil, fmt.Errorf("unable to create new tweet: %v", err)
	}
	return &newTweet, nil
}

// TweetsByTag finds tweets with specified tag, created before some point in time
// It returns slice of tweets, limited by count
func (t *TweetStorage) TweetsByTag(tag string, before int64, count int) ([]twitter.Tweet, error) {
	return t.db.listTweetsByTag(tag, before, count)
}

var generateTweetId = func() uint64 {
	return rand.Uint64()
}
