package tweetstorage

import (
        "testing"
	"reflect"
)

func TestGetTags(t *testing.T) {
        tests := []struct {
                text         string
                expectedTags []string
        }{
                {"my first tweet", nil},
                {"#cat", []string{"cat"}},
                {"#hype is #overhyped", []string{"hype", "overhyped"}},
                {"#tag nota#tag", []string{"tag"}},
                {"many recurring #tags #tags #tags", []string{"tags"}},
        }

        for _, test := range tests {
                tags := getTags(test.text)
                if !reflect.DeepEqual(tags, test.expectedTags) {
                        t.Errorf("wrong tags for %q. Wanted %q, got %q",
                                test.text, test.expectedTags, tags)
                }
        }
}
