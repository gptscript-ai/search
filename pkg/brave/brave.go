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
	"strconv"

	"github.com/gptscript-ai/search/pkg/common"
)

const (
	count        = "20" // 20 is the max allowed by Brave
	resultFilter = "web"
)

func Search(input string) (common.SearchResults, error) {
	token := os.Getenv("GPTSCRIPT_BRAVE_SEARCH_TOKEN")
	if token == "" {
		return common.SearchResults{}, fmt.Errorf("GPTSCRIPT_BRAVE_SEARCH_TOKEN is not set")
	}

	var params params
	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return common.SearchResults{}, err
	}

	resultsJSON, err := getSearchResults(token, params)
	if err != nil {
		return common.SearchResults{}, err
	}

	var resp apiResponse
	if err := json.Unmarshal([]byte(resultsJSON), &resp); err != nil {
		return common.SearchResults{}, err
	}

	return resp.toSearchResults(), nil
}

func getSearchResults(token string, params params) (string, error) {
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
	if params.Offset != "" {
		if offsetInt, err := strconv.Atoi(params.Offset); err != nil || offsetInt < 0 {
			return "", fmt.Errorf("offset must be a non-negative integer")
		}
	}

	baseURL := "https://api.search.brave.com/res/v1/web/search"
	queryParams := url.Values{}
	queryParams.Add("q", params.Query)
	queryParams.Add("count", count)
	queryParams.Add("result_filter", resultFilter)

	if params.Country != "" {
		queryParams.Add("country", params.Country)
	}
	if params.SearchLang != "" {
		queryParams.Add("search_lang", params.SearchLang)
	}
	if params.Offset != "" {
		queryParams.Add("offset", params.Offset)
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
