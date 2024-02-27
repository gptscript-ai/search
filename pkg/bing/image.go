package bing

import (
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
	token := os.Getenv("GPTSCRIPT_BING_SEARCH_TOKEN")
	if token == "" {
		return common.ImageSearchResults{}, fmt.Errorf("GPTSCRIPT_BING_SEARCH_TOKEN is not set")
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

	baseURL := "https://api.bing.microsoft.com/v7.0/images/search"
	queryParams := url.Values{}
	queryParams.Add("q", params.Query)
	queryParams.Add("count", count)

	if params.Country != "" {
		queryParams.Add("cc", params.Country)
	}
	if params.Offset != "" {
		queryParams.Add("offset", params.Offset)
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Ocp-Apim-Subscription-Key", token)
	if params.SearchLang != "" {
		req.Header.Add("Accept-Language", params.SearchLang)
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
