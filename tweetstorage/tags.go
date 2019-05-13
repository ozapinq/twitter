package tweetstorage

import (
	"regexp"
)

// regexp for hashtags.
// matches words starting with '#', ignoring '#' in the middle of the word.
var tagFinderRegexp = regexp.MustCompile(`(?:\s|\A)#(\S*)`)

var getTags = func(text string) []string {
	var result []string
	set := make(map[string]bool)

	r := tagFinderRegexp.FindAllStringSubmatch(text, -1)
	for _, v := range r {
		tag := v[1]
		if !set[tag] {
			result = append(result, v[1])
			set[tag] = true
		}
	}
	return result
}
