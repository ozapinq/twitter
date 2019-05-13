// Package twitter is a root package defining main entities and interfaces.
package twitter

type Tweet struct {
	ID        uint64    `json:"id,string"`
	Author    string    `json:"author"`
	Text      string    `json:"text"`
	Tags      []string  `json:"tags"`
	CreatedAt int64     `json:"created_at,string"`
}

type TweetStorage interface {
	New(author, text string) (*Tweet, error)
	TweetsByTag(tag string, before int64, count int) ([]Tweet, error)
}
