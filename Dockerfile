# Multistage Dockerfile

FROM golang:1.12.4 as builder

WORKDIR /go/src/github.com/ozapinq/twitter
ADD . .

# for simplicity run tests inside builder image
# in practice CI must be responsible for testing as part of pipeline
RUN go test -short ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o tweetserver ./cmd/tweetserver


FROM alpine:3.9.3

WORKDIR /app
COPY --from=builder /go/src/github.com/ozapinq/twitter/tweetserver .

CMD ["./tweetserver"]
