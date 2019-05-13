package http

import (
	"fmt"
	"log"
	"net/http"
        "os"

	"github.com/go-chi/chi"
	"github.com/ozapinq/twitter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	storage twitter.TweetStorage
	logger  *log.Logger

	// tweetPayloadLimit SHOULD NOT be treated as tweet size limit
	// because of variable-width nature of UTF8
	tweetPayloadLimit int64

	tweetListMaxCount int
}

func NewServer(storage twitter.TweetStorage) *Server {
        logger := log.New(os.Stderr, "", 0)
	var tweetPayloadLimit int64 = 1000
	var tweetListMaxCount int = 50

	return &Server{storage, logger, tweetPayloadLimit, tweetListMaxCount}
}

func (s *Server) ListenAndServe() {
	r := s.Router()
	host := ""
	port := 5555
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Print("starting tweetserver on ", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func (s *Server) Router() *chi.Mux {
	r := chi.NewRouter()
        r.Use(setHeaders)

	tweetLimit := LimitSize(s.tweetPayloadLimit)
	r.Handle("/metrics", promhttp.Handler())

	r.With(
                tweetLimit,
                measure("create_tweet"),
                authenticate,
        ).Post("/tweets", s.createTweet)

	r.With(
                measure("tweets_by_tag"),
        ).Get("/tags/{tag}/tweets", s.tweetsByTag)

	return r
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router().ServeHTTP(w, r)
}
