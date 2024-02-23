package common

import (
	"fmt"
)

type SearchResults struct {
	Results []Result `json:"results"`
}

type Result struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

func (s SearchResults) ToString() string {
	var result string
	for _, r := range s.Results {
		result += fmt.Sprintf("Title: %s\nURL: %s\nDescription: %s\n\n", r.Title, r.URL, r.Description)
	}
	return result
}
