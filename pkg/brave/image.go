package brave

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"

	"github.com/gptscript-ai/search/pkg/common"
)

func SearchImage(input string) (common.ImageSearchResults, error) {
	token := os.Getenv("GPTSCRIPT_BRAVE_SEARCH_TOKEN")
	if token == "" {
		return common.ImageSearchResults{}, fmt.Errorf("GPTSCRIPT_BRAVE_SEARCH_TOKEN is not set")
	}

	var params imageParams
	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return common.ImageSearchResults{}, err
	}

	resultsJSON, err := getImageSearchResults(token, params)
	if err != nil {
		return common.ImageSearchResults{}, err
	}

	var resp imageAPIResponse
	if err := json.Unmarshal([]byte(resultsJSON), &resp); err != nil {
		return common.ImageSearchResults{}, err
	}

	return resp.toSearchResults(), nil
}

func getImageSearchResults(token string, params imageParams) (string, error) {
	// Validate parameters
	if params.Query == "" {
		return "", fmt.Errorf("query is required")
	}
	if params.Country != "" && !slices.Contains(SupportedCountries, params.Country) {
		return "", fmt.Errorf("unsupported country: %s", params.Country)
	}
	if params.SearchLang != "" && !slices.Contains(SupportedLanguages, params.SearchLang) {
		return "", fmt.Errorf("unsupported language: %s", params.SearchLang)
	}

	baseURL := "https://api.search.brave.com/res/v1/images/search"
	queryParams := url.Values{}
	queryParams.Add("q", params.Query)

	if params.Country != "" {
		queryParams.Add("country", params.Country)
	}
	if params.SearchLang != "" {
		queryParams.Add("search_lang", params.SearchLang)
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("X-Subscription-Token", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	gzipReader, err := gzip.NewReader(res.Body)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(gzipReader)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
