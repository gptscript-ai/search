package common

import (
	"fmt"
)

type WebSearchResults struct {
	Results []WebResult `json:"results"`
}

type WebResult struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

func (w WebSearchResults) ToString() string {
	var result string
	for _, r := range w.Results {
		result += fmt.Sprintf("Title: %s\nURL: %s\nDescription: %s\n\n", r.Title, r.URL, r.Description)
	}
	return result
}

type ImageSearchResults struct {
	Results []ImageResult `json:"results"`
}

type ImageResult struct {
	Title    string `json:"title"`
	Source   string `json:"source"`
	ImageURL string `json:"image_url"`
}

func (i ImageSearchResults) ToString() string {
	var result string
	for _, r := range i.Results {
		result += fmt.Sprintf("Title: %s\nSource: %s\nImage URL: %s\n\n", r.Title, r.Source, r.ImageURL)
	}
	return result
}
