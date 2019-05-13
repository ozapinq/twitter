package http

import (
	"encoding/json"
	"github.com/ozapinq/twitter/tweetstorage"
	"github.com/ozapinq/twitter/auth"
	"fmt"
	"io/ioutil"
	"net/http"
)

type newTweet struct {
	Text string `json:"text"`
}

// TODO: change method owner to local one
func (s *Server) createTweet(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        username := ctx.Value("user").(*auth.User).Username

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
                s.logger.Print("createTweet: body read error: ", err)
		w.WriteHeader(500)
		return
	}

	var tweet newTweet
	if err := json.Unmarshal(body, &tweet); err != nil {
                w.WriteHeader(400)
                w.Write([]byte(`[{"code": "invalid_json"}]`))
		return
	}

	createdTweet, err := s.storage.New(username, tweet.Text)
        if err == tweetstorage.ErrNoTags {
                // in real world implementation we should save every tweet
                w.WriteHeader(400)
                w.Write([]byte(`[{"code": "tagless_tweets_not_implemented"}]`))
                return
        } else if err != nil {
                s.logger.Print("createTweet: storage error: ", err)
		w.WriteHeader(503)
		return
	}
	url := fmt.Sprintf("/tweets/%d", createdTweet.ID)
	w.Header().Set("Location", url)
	w.WriteHeader(201)
}
