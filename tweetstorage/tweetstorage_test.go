package tweetstorage

import (
	"errors"
	"github.com/ozapinq/twitter"
	"reflect"
	"testing"
)


func TestNew(t *testing.T) {
        db := &mockDB{
                appendTweetToTagFn: func(tag string, tw *twitter.Tweet) error {
			return nil
		},
        }
	s := NewTweetStorage(db)

	t.Run("check fields set", func(t *testing.T) {
		author := "john wayne"
		text := "my #first tweet"
		expectedTags := []string{"first"}

		generateTweetIdOrig := generateTweetId
		generateTweetId = func() uint64 { return 54321 }
		defer func() { generateTweetId = generateTweetIdOrig }()

		getTagsOrig := getTags
		getTags = func(text string) []string {
			return expectedTags
		}
		defer func() { getTags = getTagsOrig }()

		tweet, err := s.New(author, text)

		if err != nil {
			t.Fatal("New returned unexpected error: ", err)
		}
		if tweet.ID != 54321 {
			t.Error("Tweet.ID not set")
		}
		if tweet.Author != author {
			t.Errorf("Tweet.Author expected to be %q, got %q",
				author, tweet.Author)
		}
		if tweet.CreatedAt == 0 {
			t.Error("Tweet.CreatedAt is undefined")
		}
		if !reflect.DeepEqual(tweet.Tags, expectedTags) {
			t.Errorf("Tweet.Tags is incorrect: %q, wanted %q",
				tweet.Tags, expectedTags)
		}
	})

	t.Run("adds tweet for each tag", func(t *testing.T) {
		addedTags := []string{}
		s.db = &mockDB{appendTweetToTagFn: func(tag string, tw *twitter.Tweet) error {
			addedTags = append(addedTags, tag)
			return nil
		}}
		s.New("john", "my #first #tweet")

		expectedTags := []string{"first", "tweet"}
		if !reflect.DeepEqual(addedTags, expectedTags) {
			t.Error("New didn't add tweet to every tag")
		}
	})

	t.Run("insert error", func(t *testing.T) {
		s.db = &mockDB{appendTweetToTagFn: func(t string, tw *twitter.Tweet) error {
			return errors.New("some error")
		}}

		_, err := s.New("john wayne", "my first #tweet")
		if err == nil {
			t.Error("New didn't return error on failed insertion")
		}
	})

	t.Run("tweets without tags are not allowed", func(t *testing.T) {
		_, err := s.New("john wayne", "my tweet")
		if err != ErrNoTags {
			t.Error("New didn't return error on tweet without tags")
		}
	})

	t.Run("rollback if one of inserts failed", func(t *testing.T) {
		i := 1
		appended := []string{}
		appendMock := func(t string, tw *twitter.Tweet) error {
			if i == 2 {
				return errors.New("random error")
			}
			appended = append(appended, t)
			i++
			return nil
		}

		removed := []string{}
		removeMock := func(tag string, tw *twitter.Tweet) error {
			removed = append(removed, tag)
			return nil
		}

		s.db = &mockDB{
			appendTweetToTagFn:   appendMock,
			removeTweetFromTagFn: removeMock,
		}
		s.New("john wayne", "#tag1 #tag2 #tag3")

		if !reflect.DeepEqual(removed, appended) {
			t.Error("New didn't remove tweet from tags on insertion error")
		}
	})
}

func TestTweetsByTag(t *testing.T) {
	s := NewTweetStorage(&mockDB{})

	expectedTweets := []twitter.Tweet{
		{ID: 1, Author: "Alice", Text: "hello world"},
		{ID: 2, Author: "Bob", Text: "hello world"},
	}
	s.db = &mockDB{
		listTweetsByTagFn: func(tag string, s int64, c int) ([]twitter.Tweet, error) {
			return expectedTweets, nil
		}}
	tweets, _ := s.TweetsByTag("", 0, 10)
	if !reflect.DeepEqual(tweets, expectedTweets) {
		t.Errorf("TweetsByTag returned unexpected tweets")
	}
}
