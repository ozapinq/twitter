package http

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	http_requests_total = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tweetserver_http_requests_total",
		Help: "Total count of all requests (partitioned by handler)",
	}, []string{"handler"})
	http_responses_total = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tweetserver_http_responses_total",
		Help: "Total count of all responses (partitioned by handler and code)",
	}, []string{"handler", "code"})
	http_response_time_hist = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "tweetserver_http_response_time_hist",
		Help: "HTTP response time histogram (partitioned by handler)",
	}, []string{"handler"})
)
