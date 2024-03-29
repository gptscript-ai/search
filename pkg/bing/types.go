package bing

import "github.com/gptscript-ai/search/pkg/common"

// source: https://learn.microsoft.com/en-us/bing/search-apis/bing-web-search/reference/market-codes#country-codes
var SupportedCountries = []string{"AR", "AU", "AT", "BE", "BR", "CA", "CL", "DK", "FI", "FR", "DE", "HK", "IN", "ID", "IT", "JP", "KR", "MY", "MX", "NL", "NZ", "NO", "CN", "PL", "PT", "PH", "RU", "SA", "ZA", "ES", "SE", "CH", "TW", "TR", "GB", "US"}

type params struct {
	Query      string `json:"q"`
	Country    string `json:"country"`
	SearchLang string `json:"search_lang"`
	Offset     string `json:"offset"`
}

type apiResponse struct {
	WebPages struct {
		Value []struct {
			Name    string `json:"name"`
			URL     string `json:"url"`
			Snippet string `json:"snippet"`
		} `json:"value"`
	} `json:"webPages"`
}

func (a apiResponse) toSearchResults() common.WebSearchResults {
	var results common.WebSearchResults
	for _, v := range a.WebPages.Value {
		results.Results = append(results.Results, common.WebResult{
			Title:       v.Name,
			URL:         v.URL,
			Description: v.Snippet,
		})
	}
	return results
}

type imageParams struct {
	Query      string `json:"q"`
	Country    string `json:"country"`
	SearchLang string `json:"search_lang"`
	Offset     string `json:"offset"`
}

type imageAPIResponse struct {
	Value []struct {
		Name        string `json:"name"`
		ContentURL  string `json:"contentUrl"`
		HostPageURL string `json:"hostPageUrl"`
	} `json:"value"`
}

func (a imageAPIResponse) toSearchResults() common.ImageSearchResults {
	var results common.ImageSearchResults
	for _, v := range a.Value {
		results.Results = append(results.Results, common.ImageResult{
			Title:    v.Name,
			ImageURL: v.ContentURL,
			Source:   v.HostPageURL,
		})
	}
	return results
}
