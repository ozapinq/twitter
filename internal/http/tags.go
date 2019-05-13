package http

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/ozapinq/twitter"
	"net/http"
	"net/url"
	"strconv"
)

type TweetResponse struct {
	Tweets []twitter.Tweet `json:"tweets"`
	Next   string          `json:"next,omitempty"`
}

// TODO: move to tag handler or something
func (s *Server) tweetsByTag(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")
	err := r.ParseForm()
	if err != nil {
                w.WriteHeader(400)
                w.Write([]byte(`[{"code":"unable_to_parse_form"}]`))
		return
	}

	before, err := positiveInt64(r.FormValue("before"))
	if err != nil {
                w.WriteHeader(400)
                w.Write([]byte(`[{"code":"invalid_before_value"}]`))
		return
	}

	count, err := positiveInt(r.FormValue("count"))
	if err != nil {
                w.WriteHeader(400)
                w.Write([]byte(`[{"code":"invalid_count_value"}]`))
		return
	} else if count > s.tweetListMaxCount {
                w.WriteHeader(400)
                w.Write([]byte(`[{"code":"invalid_count_value"}]`))
		return
	} else if count == 0 {
		count = s.tweetListMaxCount
	}

	tweets, err := s.storage.TweetsByTag(tag, before, count)
	if err != nil {
                s.logger.Print("tweetsByTag: unable to fetch tweets: ", err)
		w.WriteHeader(503)
		return
	}
	if len(tweets) == 0 {
                w.WriteHeader(404)
                w.Write([]byte(`[{"code":"no_tweets_containing_tag"}]`))
		return
	}

        lastTweet := tweets[len(tweets)-1]
	nextUrl := generateNextUrl(r.URL.String(), lastTweet.CreatedAt)

	data, err := json.Marshal(TweetResponse{tweets, nextUrl})
	if err != nil {
                s.logger.Print("tweetsByTag: unable to marshal response: ", err)
		w.WriteHeader(500)
		return
	}

	w.Write(data)
}

var generateNextUrl = func(baseUrl string, before int64) string {
        u, err := url.Parse(baseUrl)
        if err != nil {
                return baseUrl
        }
        q := u.Query()
        q.Set("before", strconv.FormatInt(before, 10))
        u.RawQuery = q.Encode()
        return u.String()
}

func positiveInt64(str string) (int64, error) {
	var result int64 = 0
	if str != "" {
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil || i < 0 {
			return 0, errors.New("not positive integer value")
		}
		result = i
	}
	return result, nil
}

func positiveInt(str string) (int, error) {
        res, err := positiveInt64(str)
        return int(res), err
}
