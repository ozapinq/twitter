package tweetstorage

// I've decided not to test DB integration in application tests.
// It's almost impossible to write good DB integration tests,
// especially if database is distributed and eventually-consistent.
// In my opinion the better way is to write System tests, that simulate
// as much as possible potential errors (availability, latency-related, etc).
//
// To simplify things, I prefer to isolate DB-related activities in simple,
// predictable functions, which I can mock in my unittests.


import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/ozapinq/twitter"
)

type cassandraDB struct {
        cluster *gocql.ClusterConfig
	session *gocql.Session
}

func NewDB(hosts []string, keyspace string) (*cassandraDB, error) {
        db := &cassandraDB{}

	db.cluster = gocql.NewCluster(hosts...)
	db.cluster.Keyspace = keyspace

	s, err := db.cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("unable to create session: %v", err)
	}
        db.session = s

        return db, nil
}

func (m cassandraDB) appendTweetToTag(tag string, tw *twitter.Tweet) error {
	err := m.session.Query(
		"INSERT INTO tweets_by_tag "+
			"(tag, created_at, tid, text, author) "+
			"VALUES (?, ?, ?, ?, ?)",
		tag, tw.CreatedAt, tw.ID, tw.Text, tw.Author).Exec()
	return err
}

func (m cassandraDB) removeTweetFromTag(tag string, tw *twitter.Tweet) error {
	err := m.session.Query(
		"DELETE FROM tweets_by_tag "+
			"WHERE tag = ? AND created_at = ?",
		tag, tw.CreatedAt).Exec()
	return err
}

func (m cassandraDB) listTweetsByTag(tag string, before int64, count int) ([]twitter.Tweet, error) {
	query := m.session.Query(
		"SELECT created_at, author, text, tid "+
			"FROM tweets_by_tag "+
			"WHERE tag = ? AND created_at > ? "+
			"LIMIT ?",
		tag, before, count).Iter()
	tweets := make([]twitter.Tweet, 0, count)
	var createdAt int64
	var author string
	var text string
	var id uint64

	for query.Scan(&createdAt, &author, &text, &id) {
		tweet := twitter.Tweet{
			id,
			author,
			text,
			[]string{tag},
			createdAt,
		}
		tweets = append(tweets, tweet)
	}
	return tweets, nil
}
