package tweetstorage

import (
	"github.com/ozapinq/twitter"
)

type mockDB struct {
	appendTweetToTagFn   func(string, *twitter.Tweet) error
	removeTweetFromTagFn func(string, *twitter.Tweet) error
	listTweetsByTagFn    func(string, int64, int) ([]twitter.Tweet, error)
}

func (m mockDB) appendTweetToTag(tag string, tweet *twitter.Tweet) error {
	if m.appendTweetToTagFn == nil {
		panic("mockDB.appendTweetToTag called, but not set")
	}
	return m.appendTweetToTagFn(tag, tweet)
}

func (m mockDB) removeTweetFromTag(tag string, tweet *twitter.Tweet) error {
	if m.removeTweetFromTagFn == nil {
		panic("mockDB.removeTweetFromTag called, but not set")
	}
	return m.removeTweetFromTagFn(tag, tweet)
}

func (m mockDB) listTweetsByTag(tag string, before int64, count int) ([]twitter.Tweet, error) {
	if m.listTweetsByTagFn == nil {
		panic("mockDB.listTweetsByTagFn called, but not set")
	}
	return m.listTweetsByTagFn(tag, before, count)
}
