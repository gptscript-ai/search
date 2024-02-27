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

func SearchImage(input string) (common.ImageSearchResults, error) {
	token := os.Getenv("GPTSCRIPT_GOOGLE_SEARCH_TOKEN")
	if token == "" {
		return common.ImageSearchResults{}, fmt.Errorf("GPTSCRIPT_GOOGLE_SEARCH_TOKEN is not set")
	}

	searchEngineID := os.Getenv("GPTSCRIPT_GOOGLE_SEARCH_ENGINE_ID")
	if searchEngineID == "" {
		return common.ImageSearchResults{}, fmt.Errorf("GPTSCRIPT_GOOGLE_SEARCH_ENGINE_ID is not set")
	}

	var params imageParams
	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return common.ImageSearchResults{}, err
	}

	resultsJSON, err := getImageSearchResults(token, searchEngineID, params)
	if err != nil {
		return common.ImageSearchResults{}, err
	}

	var resp imageAPIResponse
	if err := json.Unmarshal([]byte(resultsJSON), &resp); err != nil {
		return common.ImageSearchResults{}, err
	}

	return resp.toSearchResults(), nil
}

func getImageSearchResults(token, engineID string, params imageParams) (string, error) {
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
	if params.ImageSize != "" && !slices.Contains(SupportedImageSizes, params.ImageSize) {
		return "", fmt.Errorf("unsupported image size: %s", params.ImageSize)
	}
	if params.ImageType != "" && !slices.Contains(SupportedImageTypes, params.ImageType) {
		return "", fmt.Errorf("unsupported image type: %s", params.ImageType)
	}

	baseURL := "https://www.googleapis.com/customsearch/v1"
	queryParams := url.Values{}
	queryParams.Set("key", token)
	queryParams.Set("cx", engineID)
	queryParams.Set("q", params.Query)
	queryParams.Set("searchType", "image")
	queryParams.Set("num", count)

	if params.Country != "" {
		queryParams.Set("cr", params.Country)
	}
	if params.SearchLang != "" {
		queryParams.Set("lr", "lang_"+params.SearchLang)
	}
	if params.Offset != "" {
		queryParams.Set("start", params.Offset)
	}
	if params.ImageSize != "" {
		queryParams.Set("imgSize", params.ImageSize)
	}
	if params.ImageType != "" {
		queryParams.Set("imgType", params.ImageType)
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

	return string(body), err
}
