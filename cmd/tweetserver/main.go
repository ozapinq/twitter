package main

import (
	"github.com/ozapinq/twitter/internal/http"
	"github.com/ozapinq/twitter/tweetstorage"
	"log"
        "os"
        "strings"
)

func main() {
        dbNodes := os.Getenv("DB_NODES")
        if dbNodes == "" {
                log.Fatal("DB_NODES not specified. Comma separated list of hosts required")
        }
	hosts := strings.Split(dbNodes, ",")

        db, err := tweetstorage.NewDB(hosts, "tweetserver")
        if err != nil {
		log.Fatal("unable to connect to database: ", err)
        }

	storage := tweetstorage.NewTweetStorage(db)
	server := http.NewServer(storage)
	server.ListenAndServe()
}
