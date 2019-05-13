package http

import (
	"encoding/json"
	"errors"
	"github.com/ozapinq/twitter"
	"github.com/ozapinq/twitter/internal/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
        "strings"
        "log"
)

func TestTweetsByTag(t *testing.T) {
	storage := &mock.TweetStorage{}
	server := NewServer(storage)

	type TweetsByTagFn func(string, int64, int) ([]twitter.Tweet, error)
	var tweetsByTag = func(path string, fn TweetsByTagFn) *httptest.ResponseRecorder {
		u, _ := url.Parse(path)
		path = u.Path

		storage.TweetsByTagFn = fn
		defer func() { storage.TweetsByTagFn = nil }()

		r := httptest.NewRequest(http.MethodGet, path, nil)
		r.URL.RawQuery = u.RawQuery

		w := httptest.NewRecorder()
		server.ServeHTTP(w, r)

		return w
	}

	t.Run("tweets with specified tag exist", func(t *testing.T) {
		expectedTweets := []twitter.Tweet{
			twitter.Tweet{
				3, "Marina", "third tweet",
				[]string{"my_tag"}, 1300,
			},
			twitter.Tweet{
				2, "Bob", "another one",
				[]string{"my_tag"}, 1200,
			},
			twitter.Tweet{
				1, "Alice", "First tweet",
				[]string{"my_tag"}, 1100,
			},
		}
		f := func(tag string, _ int64, _ int) ([]twitter.Tweet, error) {
			return expectedTweets, nil
		}
		u := "/tags/my_tag/tweets"
		w := tweetsByTag(u, f)

		if w.Code != 200 {
			t.Errorf("GET %s returned %d, expected 200", u, w.Code)
		}

		var response TweetResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal("Unable to unmarshal received tweets: ", err)
		}

		if !reflect.DeepEqual(response.Tweets, expectedTweets) {
			t.Errorf("Received unexpected tweets.\n\twanted: %v\n\tgot: %v",
				expectedTweets, response.Tweets)
		}

                expectedNext := "/tags/my_tag/tweets?before=1100"
                if response.Next != expectedNext {
                        t.Errorf("Next url is unexpected. wanted %q, got %q",
                                expectedNext, response.Next)
                }

                contentType := w.Result().Header.Get("Content-Type")
                if contentType != "application/json" {
                        t.Errorf("Content-Type is incorrect: expected %q, got %q",
                                "application/json", contentType)
                }
	})

	t.Run("check before and count set by default", func(t *testing.T) {
		var before int64
		var count int
		f := func(tag string, s int64, c int) ([]twitter.Tweet, error) {
			before = s
			count = c
			return []twitter.Tweet{}, nil
		}
		tweetsByTag("/tags/tag/tweets", f)

		if before != 0 {
			t.Errorf("before is unexpected: wanted 0, got %d", before)
		}

		if count != server.tweetListMaxCount {
			t.Error("default count != tweetListMaxCount")
		}
	})

	t.Run("check before and count set from querystring", func(t *testing.T) {
		var before int64
		var count int
		f := func(tag string, s int64, c int) ([]twitter.Tweet, error) {
			before = s
			count = c
			return []twitter.Tweet{}, nil
		}
		tweetsByTag("/tags/tag/tweets?before=3&count=10", f)

		if before != 3 {
			t.Errorf("before is unexpected: wanted 3, got %v", before)
		}

		if count != 10 {
			t.Errorf("count is unexpected: wanted 10, got %d", count)
		}

		w := tweetsByTag("/tags/tag/tweets?before=3&count=100000", f)

		if w.Code != 400 {
			t.Errorf("count exceeded tweetListMaxCount: got %q, wanted %q  ",
				http.StatusText(w.Code), http.StatusText(400))
		}
	})

	t.Run("incorrect before", func(t *testing.T) {
		urls := []string{
			"/tags/tag/tweets?before=random",
			"/tags/tag/tweets?before=-1",
		}

		for _, u := range urls {
			f := func(_ string, _ int64, _ int) ([]twitter.Tweet, error) {
				return nil, nil
			}
			w := tweetsByTag(u, f)

			if w.Code != 400 {
				t.Errorf("GET %s returned %q, expected %q",
					u, http.StatusText(w.Code), http.StatusText(400))
			}

			body := w.Body.String()
			expected := `[{"code":"invalid_before_value"}]`
			if body != expected {
				t.Errorf("GET %q invalid body. expected '%s', got '%s'",
					u, expected, body)
			}
		}
	})

	t.Run("no tweets matching tag", func(t *testing.T) {
		f := func(tag string, _ int64, _ int) ([]twitter.Tweet, error) {
			return []twitter.Tweet{}, nil
		}
		w := tweetsByTag("/tags/my_tag/tweets", f)

		if w.Code != 404 {
			t.Errorf("empty tweetsByTag returned %d, expected 404",
                                w.Code)
		}
	})

	t.Run("return 503 on storage errors", func(t *testing.T) {
                lw := &strings.Builder{}
                server.logger = log.New(lw, "", 0)
                defer func() { server.logger = nil }()

		f := func(tag string, _ int64, _ int) ([]twitter.Tweet, error) {
			return nil, errors.New("random error")
		}
		w := tweetsByTag("/tags/my_tag/tweets", f)

		if w.Code != 503 {
			t.Errorf("wanted %q on storage error, got %q",
				http.StatusText(503), http.StatusText(w.Code))
		}
                expectedMsg := "tweetsByTag: unable to fetch tweets: random error"
                if !strings.Contains(lw.String(), expectedMsg) {
                        t.Errorf("storage error wasn't logged. wanted %q, got %q",
                                expectedMsg, lw.String())
                }
                contentType := w.Result().Header.Get("Content-Type")
                if contentType != "application/json" {
                        t.Errorf("Content-Type is incorrect: expected %q, got %q",
                                "application/json", contentType)
                }
	})
}

func TestGenerateNextUrl(t *testing.T) {
        tests := []struct{
                url string
                before int64
                expected string
        }{
                {"/tags/tag/tweets", 10, "/tags/tag/tweets?before=10"},
                {"/tags/tag/tweets?before=5", 20, "/tags/tag/tweets?before=20"},
                {"/tags/tag/tweets?before=50&c=1", 3, "/tags/tag/tweets?before=3&c=1"},
        }
        for _, test := range tests {
                got := generateNextUrl(test.url, test.before)
                if got != test.expected {
                        t.Errorf("generateNextUrl(%q) != %q. got %q",
                                test.url, test.expected, got)
                }
        }
}
