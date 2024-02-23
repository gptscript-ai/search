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

func (a apiResponse) toSearchResults() common.SearchResults {
	var results common.SearchResults
	for _, v := range a.WebPages.Value {
		results.Results = append(results.Results, common.Result{
			Title:       v.Name,
			URL:         v.URL,
			Description: v.Snippet,
		})
	}
	return results
}
