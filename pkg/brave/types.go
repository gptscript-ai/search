package brave

import "github.com/gptscript-ai/search/pkg/common"

// These lists come from https://api.search.brave.com/app/documentation/web-search/codes, but you must be logged in in order to view the docs.
var (
	SupportedCountries = []string{"AR", "AU", "AT", "BE", "BR", "CA", "CL", "DK", "FI", "FR", "DE", "HK", "IN", "ID", "IT", "JP", "KR", "MY", "MX", "NL", "NZ", "NO", "CN", "PL", "PT", "PH", "RU", "SA", "ZA", "ES", "SE", "CH", "TW", "TR", "GB", "US"}
	SupportedLanguages = []string{"ar", "eu", "bn", "bg", "ca", "zh-hans", "zh-hant", "hr", "cs", "da", "nl", "en", "en-gb", "et", "fi", "fr", "gl", "de", "gu", "he", "hi", "hu", "is", "it", "jp", "kn", "ko", "lv", "lt", "ms", "ml", "mr", "nb", "pl", "pt-br", "pt-pt", "pa", "ro", "ru", "sr", "sk", "sl", "es", "sv", "ta", "te", "th", "tr", "uk", "vi"}
)

type webParams struct {
	Query      string `json:"q"`
	Country    string `json:"country"`
	SearchLang string `json:"search_lang"`
	Offset     string `json:"offset"`
}

type webAPIResponse struct {
	Web struct {
		Results []struct {
			Title       string `json:"title"`
			URL         string `json:"url"`
			Description string `json:"description"`
		} `json:"results"`
	} `json:"web"`
}

func (r webAPIResponse) toSearchResults() common.WebSearchResults {
	var results common.WebSearchResults
	for _, res := range r.Web.Results {
		results.Results = append(results.Results, common.WebResult{
			Title:       res.Title,
			URL:         res.URL,
			Description: res.Description,
		})
	}
	return results
}

type imageParams struct {
	Query      string `json:"q"`
	Country    string `json:"country"`
	SearchLang string `json:"search_lang"`
}

type imageAPIResponse struct {
	Results []struct {
		Title      string `json:"title"`
		URL        string `json:"url"`
		Properties struct {
			URL string `json:"url"`
		} `json:"properties"`
	} `json:"results"`
}

func (r imageAPIResponse) toSearchResults() common.ImageSearchResults {
	var results common.ImageSearchResults
	for _, res := range r.Results {
		results.Results = append(results.Results, common.ImageResult{
			Title:    res.Title,
			Source:   res.URL,
			ImageURL: res.Properties.URL,
		})
	}
	return results
}
