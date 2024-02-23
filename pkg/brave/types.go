package brave

import "github.com/gptscript-ai/search/pkg/common"

// These lists come from https://api.search.brave.com/app/documentation/web-search/codes, but you must be logged in in order to view the docs.
var (
	SupportedCountries = []string{"AR", "AU", "AT", "BE", "BR", "CA", "CL", "DK", "FI", "FR", "DE", "HK", "IN", "ID", "IT", "JP", "KR", "MY", "MX", "NL", "NZ", "NO", "CN", "PL", "PT", "PH", "RU", "SA", "ZA", "ES", "SE", "CH", "TW", "TR", "GB", "US"}
	SupportedLanguages = []string{"ar", "eu", "bn", "bg", "ca", "zh-hans", "zh-hant", "hr", "cs", "da", "nl", "en", "en-gb", "et", "fi", "fr", "gl", "de", "gu", "he", "hi", "hu", "is", "it", "jp", "kn", "ko", "lv", "lt", "ms", "ml", "mr", "nb", "pl", "pt-br", "pt-pt", "pa", "ro", "ru", "sr", "sk", "sl", "es", "sv", "ta", "te", "th", "tr", "uk", "vi"}
)

type params struct {
	Query      string `json:"q"`
	Country    string `json:"country"`
	SearchLang string `json:"search_lang"`
	Offset     string `json:"offset"`
}

type apiResponse struct {
	Web struct {
		Results []struct {
			Title       string `json:"title"`
			URL         string `json:"url"`
			Description string `json:"description"`
		} `json:"results"`
	} `json:"web"`
}

func (r apiResponse) toSearchResults() common.SearchResults {
	var results common.SearchResults
	for _, res := range r.Web.Results {
		results.Results = append(results.Results, common.Result{
			Title:       res.Title,
			URL:         res.URL,
			Description: res.Description,
		})
	}
	return results
}
