package bing

import (
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
	count          = "20" // 50 is the max, but 20 should be sufficient
	responseFilter = "Webpages"
)

func Search(input string) (common.SearchResults, error) {
	token := os.Getenv("GPTSCRIPT_BING_SEARCH_TOKEN")
	if token == "" {
		return common.SearchResults{}, fmt.Errorf("GPTSCRIPT_BING_SEARCH_TOKEN is not set")
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
	if params.Country != "" && !slices.Contains(SupportedCountries, params.Country) {
		return "", fmt.Errorf("unsupported country: %s", params.Country)
	}
	if params.Offset != "" {
		if offsetInt, err := strconv.Atoi(params.Offset); err != nil || offsetInt < 0 {
			return "", fmt.Errorf("offset must be a non-negative integer")
		}
	}

	baseURL := "https://api.bing.microsoft.com/v7.0/search"
	queryParams := url.Values{}
	queryParams.Add("q", params.Query)
	queryParams.Add("count", count)
	queryParams.Add("responseFilter", responseFilter)

	if params.Country != "" {
		queryParams.Add("cc", params.Country)
	}
	if params.Offset != "" {
		queryParams.Add("offset", params.Offset)
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", token)
	if params.SearchLang != "" {
		req.Header.Set("Accept-Language", params.SearchLang)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
