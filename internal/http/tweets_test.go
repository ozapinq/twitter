package http

import (
        "log"
	"errors"
	"github.com/ozapinq/twitter"
	"github.com/ozapinq/twitter/tweetstorage"
	"github.com/ozapinq/twitter/internal/mock"
	"github.com/ozapinq/twitter/auth"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type ErrorReader struct {
	text string
}

func (e ErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New(e.text)
}

func TestCreateTweet(t *testing.T) {
	const tweetsUrl = "/tweets"
	var storage = &mock.TweetStorage{}
	var server = NewServer(storage)

        user := auth.User{"john"}
        defaultToken := "correct_token"
        authenticator.Add(defaultToken, &user)
        defer func() { authenticator.Delete(defaultToken) }()

	type NewFnMock func(string, string) (*twitter.Tweet, error)
	var createTweet = func(b string, fn NewFnMock, token string) *httptest.ResponseRecorder {
		storage.NewFn = fn
		defer func() { storage.NewFn = nil }()

		body := strings.NewReader(b)
		r := httptest.NewRequest(http.MethodPost, tweetsUrl, body)
		w := httptest.NewRecorder()

                if token != "" {
                        r.Header.Set("X-Auth-Token", token)
                }

		server.ServeHTTP(w, r)

		return w
	}

	t.Run("successful tweet creation", func(t *testing.T) {
		var text, author string
                f := func(a, t string) (*twitter.Tweet, error) {
			text = t
                        author = a
			return &twitter.Tweet{ID: 12345}, nil
		}
		defer func() { storage.NewFn = nil }()

                w := createTweet(`{"text":"hello, world"}`, f, defaultToken)

		if w.Code != 201 {
			t.Errorf("POST %s status code %d != 201", tweetsUrl, w.Code)
		}

		location := w.Result().Header.Get("location")
		if expected := "/tweets/12345"; location != expected {
			t.Errorf("POST %s Location field (%q) is incorrect. Expected %q",
				tweetsUrl, location, expected)
		}

                contentType := w.Result().Header.Get("Content-Type")
                if contentType != "application/json" {
                        t.Errorf("Content-Type is incorrect: expected %q, got %q",
                                "application/json", contentType)
                }

		expectedText := "hello, world"
		if text != expectedText {
			t.Errorf("Passed tweet text != saved; expected: %s, got: %s",
				expectedText, text)
		}

                if author != user.Username {
                        t.Errorf("Passed author doesn't match passed token: %s, got %s",
                                user.Username, author)
                }
	})

        t.Run("401 Unauthorized on missing token header", func(t *testing.T) {
                body := `{"text":"hello, world"}`
                w := createTweet(body, nil, "")

                if w.Code != 401 {
                        t.Errorf("Expected %q on request without auth token, got %q",
                                http.StatusText(401), http.StatusText(w.Code))
                }
                expectedBody := `[{"code":"no_auth_token"}]`
                if !strings.Contains(w.Body.String(), expectedBody) {
                        t.Errorf("Expected error message %v, got %v",
                                expectedBody, w.Body.String())
                }
                contentType := w.Result().Header.Get("Content-Type")
                if contentType != "application/json" {
                        t.Errorf("Content-Type is incorrect: expected %q, got %q",
                                "application/json", contentType)
                }
        })

        t.Run("401 Unauthorized on invalid token header", func(t *testing.T) {
                body := `{"text":"hello, world"}`
                w := createTweet(body, nil, "invalid_token")

                if w.Code != 401 {
                        t.Errorf("Expected %q on request without auth token, got %q",
                                http.StatusText(401), http.StatusText(w.Code))
                }
                expectedBody := `[{"code":"invalid_auth_token"}]`
                if !strings.Contains(w.Body.String(), expectedBody) {
                        t.Errorf("Expected error message %v, got %v",
                                expectedBody, w.Body.String())
                }
        })

	t.Run("400 Bad Request on tweet without tags", func(t *testing.T) {
		data := `{"text":"not-a-tag another-not-a-tag"}`

		f := func(_, _ string) (*twitter.Tweet, error) {
			return nil, tweetstorage.ErrNoTags
		}
		w := createTweet(data, f, defaultToken)

		if w.Code != 400 {
			t.Errorf("createTweet without tags didn't return 400 Bad Request")
		}

                expectedBody := `[{"code": "tagless_tweets_not_implemented"}]`
                if w.Body.String() != expectedBody {
                        t.Error("Expected error description not present")
                }
	})

	t.Run("invalid input raises error", func(t *testing.T) {
		input := []string{
			"",
			"not-json-at-all",
			`{"broken":"json`,
		}

		for _, data := range input {
			w := createTweet(data, nil, defaultToken)
			body := w.Body.String()

			if w.Code != 400 {
				t.Errorf("Invalid JSON input didn't return %q: %q",
					http.StatusText(400), data)
			}

			expectedBody := `[{"code": "invalid_json"}]`
			if body != expectedBody {
				t.Error("Expected error description not present")
			}
		}
	})

	t.Run("reject requests exceeding length limit", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, tweetsUrl, nil)
		w := httptest.NewRecorder()

		r.ContentLength = server.tweetPayloadLimit + 1

		server.ServeHTTP(w, r)

		if w.Code != 413 {
			t.Errorf("Huge payload didn't return %q", http.StatusText(413))
		}
	})

	t.Run("return 500 on body read error", func(t *testing.T) {
                lw := &strings.Builder{}
                server.logger = log.New(lw, "", 0)
                defer func() { server.logger = nil }()

		r := httptest.NewRequest(http.MethodPost, tweetsUrl, ErrorReader{"some error"})
		w := httptest.NewRecorder()

                r.Header.Set("X-Auth-Token", defaultToken)
		server.ServeHTTP(w, r)

		if w.Code != 500 {
			t.Errorf("Body read error returned %q, expected %q",
                                http.StatusText(w.Code), http.StatusText(500))
		}

                expectedMsg := "createTweet: body read error: some error"
                if !strings.Contains(lw.String(), expectedMsg) {
                        t.Errorf("body read error wasn't logged. wanted %q, got %q",
                                expectedMsg, lw.String())
                }

                contentType := w.Result().Header.Get("Content-Type")
                if contentType != "application/json" {
                        t.Errorf("Content-Type is incorrect: expected %q, got %q",
                                "application/json", contentType)
                }
	})

	t.Run("return 503 on storage errors", func(t *testing.T) {
                lw := &strings.Builder{}
                server.logger = log.New(lw, "", 0)
                defer func() { server.logger = nil }()

		f := func(_, _ string) (*twitter.Tweet, error) {
			return nil, errors.New("random error")
		}
		w := createTweet(`{"text":"hello, world"}`, f, defaultToken)

		if w.Code != 503 {
			t.Errorf("wanted %q on storage error, got %q",
				http.StatusText(503), http.StatusText(w.Code))
		}

                expectedMsg := "createTweet: storage error: random error"
                if !strings.Contains(lw.String(), expectedMsg) {
                        t.Errorf("storage error wasn't logged. wanted %q, got %q",
                                expectedMsg, lw.String())
                }
	})
}
