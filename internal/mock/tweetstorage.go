package mock

import "github.com/ozapinq/twitter"

type TweetStorage struct {
	NewFn         func(string, string) (*twitter.Tweet, error)
	TweetsByTagFn func(string, int64, int) ([]twitter.Tweet, error)
}

func NewTweetStorage() *TweetStorage {
	return &TweetStorage{}
}

func (t *TweetStorage) New(author, text string) (*twitter.Tweet, error) {
	if t.NewFn == nil {
		panic("mock.TweetStorage.NewFn called, but not set")
	}
	return t.NewFn(author, text)
}

func (t *TweetStorage) TweetsByTag(tag string, before int64, count int) ([]twitter.Tweet, error) {
	if t.TweetsByTagFn == nil {
		panic("mock.TweetStorage.TweetsByTagFn called, but not set")
	}
	return t.TweetsByTagFn(tag, before, count)
}
