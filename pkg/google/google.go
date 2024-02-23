package google

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

const count = "10" // 10 is the max allowed by Google

func Search(input string) (common.SearchResults, error) {
	token := os.Getenv("GPTSCRIPT_GOOGLE_SEARCH_TOKEN")
	if token == "" {
		return common.SearchResults{}, fmt.Errorf("GPTSCRIPT_GOOGLE_SEARCH_TOKEN is not set")
	}

	searchEngineID := os.Getenv("GPTSCRIPT_GOOGLE_SEARCH_ENGINE_ID")
	if searchEngineID == "" {
		return common.SearchResults{}, fmt.Errorf("GPTSCRIPT_GOOGLE_SEARCH_ENGINE_ID is not set")
	}

	var params params
	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return common.SearchResults{}, err
	}

	resultsJSON, err := getSearchResults(token, searchEngineID, params)
	if err != nil {
		return common.SearchResults{}, err
	}

	var resp apiResponse
	if err := json.Unmarshal([]byte(resultsJSON), &resp); err != nil {
		return common.SearchResults{}, err
	}

	return resp.toSearchResults(), nil
}

func getSearchResults(token, engineID string, params params) (string, error) {
	if params.Query == "" {
		return "", fmt.Errorf("query is empty")
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

	baseURL := "https://www.googleapis.com/customsearch/v1"
	queryParams := url.Values{}
	queryParams.Add("q", params.Query)
	queryParams.Add("key", token)
	queryParams.Add("cx", engineID)
	queryParams.Add("num", count)

	if params.Country != "" {
		queryParams.Add("cr", "country"+params.Country)
	}
	if params.SearchLang != "" {
		queryParams.Add("lr", "lang_"+params.SearchLang)
	}
	if params.Offset != "" {
		queryParams.Add("start", params.Offset)
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return "", err
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
