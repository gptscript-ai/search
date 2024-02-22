package google

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func SearchGoogle(input string) (string, error) {
	token := os.Getenv("GPTSCRIPT_GOOGLE_SEARCH_TOKEN")
	if token == "" {
		return "", fmt.Errorf("GPTSCRIPT_GOOGLE_SEARCH_TOKEN is not set")
	}

	searchEngineID := os.Getenv("GPTSCRIPT_GOOGLE_SEARCH_ENGINE_ID")
	if searchEngineID == "" {
		return "", fmt.Errorf("GPTSCRIPT_GOOGLE_SEARCH_ENGINE_ID is not set")
	}

	var params struct {
		Query string `json:"q"`
	}

	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return "", err
	}

	if params.Query == "" {
		return "", fmt.Errorf("query is empty")
	}

	baseURL := "https://www.googleapis.com/customsearch/v1"
	queryParams := url.Values{}
	queryParams.Add("q", params.Query)
	queryParams.Add("key", token)
	queryParams.Add("cx", searchEngineID)

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
