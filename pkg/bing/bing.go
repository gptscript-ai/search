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
)

var (
	// source: https://learn.microsoft.com/en-us/bing/search-apis/bing-web-search/reference/market-codes
	bingSupportedCountries = []string{"AR", "AU", "AT", "BE", "BR", "CA", "CL", "DK", "FI", "FR", "DE", "HK", "IN", "ID", "IT", "JP", "KR", "MY", "MX", "NL", "NZ", "NO", "CN", "PL", "PT", "PH", "RU", "SA", "ZA", "ES", "SE", "CH", "TW", "TR", "GB", "US"}
)

func SearchBing(input string) (string, error) {
	token := os.Getenv("GPTSCRIPT_BING_SEARCH_TOKEN")
	if token == "" {
		return "", fmt.Errorf("GPTSCRIPT_BING_SEARCH_TOKEN is not set")
	}

	var params struct {
		Query      string `json:"q"`
		Country    string `json:"country"`
		SearchLang string `json:"search_lang"`
		Offset     string `json:"offset"`
	}

	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return "", err
	}

	// Validate parameters
	if params.Country != "" && !slices.Contains(bingSupportedCountries, params.Country) {
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
