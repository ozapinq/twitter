package http

import (
        "net/http"
        "time"
        "strconv"
        "context"
	"github.com/ozapinq/twitter/auth"
)

var authenticator *auth.TokenAuthenticator

func init() {
        authenticator = auth.New()
}

// LimitSize checks if Content-Type of request is below size
func LimitSize(size int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > size {
				w.WriteHeader(413)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, 200}
}

func (l *loggingResponseWriter) WriteHeader(status int) {
	l.statusCode = status
	l.ResponseWriter.WriteHeader(status)
}

// measure instruments http.Handler with prometheus counters and histogram.
// passed handlerName will be used as label.
func measure(handlerName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lwr := newLoggingResponseWriter(w)

			http_requests_total.WithLabelValues(
				handlerName).Inc()
			start := time.Now()
			next.ServeHTTP(lwr, r)
			duration := time.Since(start).Seconds()

			http_responses_total.WithLabelValues(
				handlerName, strconv.Itoa(lwr.statusCode)).Inc()
			http_response_time_hist.WithLabelValues(
				handlerName).Observe(duration)
		})
	}
}

// authenticate checks request looking for X-Auth-Token header and
// checks it's value through authenticator
func authenticate(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                token := r.Header.Get("X-Auth-Token")
                if token == "" {
                        w.WriteHeader(401)
                        w.Write([]byte(`[{"code":"no_auth_token"}]`))
                        return
                }
                user, err := authenticator.Check(token)
                if err != nil {
                        w.WriteHeader(401)
                        w.Write([]byte(`[{"code":"invalid_auth_token"}]`))
                        return
                }
                ctx := context.WithValue(r.Context(), "user", user)
                next.ServeHTTP(w, r.WithContext(ctx))
        })
}


// setHeaders is a simple middleware responsible for correct Content-Type
// In real-world applications you usually don't want to reinvent the wheel
// and prefer to use existing tools, developed for higher-level frameworks
func setHeaders(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Content-Type", "application/json")
                w.Header().Set("Access-Control-Allow-Origin", "*")
                next.ServeHTTP(w, r)
        })
}
