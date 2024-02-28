package duckduckgo

import (
	"compress/flate"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/gptscript-ai/search/pkg/common"
	"golang.org/x/net/html"
)

func Search(input string) (common.WebSearchResults, error) {
	var params webParams
	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return common.WebSearchResults{}, err
	}

	if strings.EqualFold(os.Getenv("GPTSCRIPT_USE_CHROME"), "true") {
		return getSearchResults(params)
	}
	return getSearchResultsNoChrome(params)
}

func getSearchResults(params webParams) (common.WebSearchResults, error) {
	if params.Query == "" {
		return common.WebSearchResults{}, nil
	}

	baseURL := "https://html.duckduckgo.com/html"
	queryParams := url.Values{}
	queryParams.Add("q", params.Query)

	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	var nodes []*cdp.Node
	ctx, stop := chromedp.NewContext(context.Background())
	defer stop()
	err := chromedp.Run(ctx,
		chromedp.Navigate(fullURL),
		chromedp.Nodes(".result__body", &nodes, chromedp.NodeVisible, chromedp.Populate(-1, true)))
	if err != nil {
		return common.WebSearchResults{}, err
	}

	var results common.WebSearchResults
	for _, node := range nodes {
		if node == nil {
			continue
		}

		var title, link, description string
		for _, child := range node.Children {
			if child == nil {
				continue
			}

			if child.AttributeValue("class") == "result__title" && len(child.Children) > 0 {
				link = child.Children[0].AttributeValue("href")

				// This shouldn't happen, but I want to prevent a panic in case it does
				if len(child.Children[0].Children) == 0 {
					continue
				}

				title = child.Children[0].Children[0].NodeValue
			} else if child.AttributeValue("class") == "result__snippet" {
				description = combineNodeValues(child.Children)
			}
		}

		results.Results = append(results.Results, common.WebResult{
			Title:       title,
			URL:         fixURL(link),
			Description: description,
		})
	}

	return results, nil
}

func combineNodeValues(nodes []*cdp.Node) string {
	var result string
	for _, node := range nodes {
		if node == nil {
			continue
		}

		if len(node.Children) > 0 {
			result += combineNodeValues(node.Children)
		}

		result += node.NodeValue
	}

	return result
}

func fixURL(link string) string {
	// All the URLs we get look like this: //duckduckgo.com/l/?uddg=<actual URL>
	link = "https:" + link
	urlObj, err := url.Parse(link)
	if err != nil {
		return link
	}

	return urlObj.Query().Get("uddg")
}

func getSearchResultsNoChrome(params webParams) (common.WebSearchResults, error) {
	if params.Query == "" {
		return common.WebSearchResults{}, nil
	}

	baseURL := "https://html.duckduckgo.com/html"
	queryParams := url.Values{}
	queryParams.Add("q", params.Query)

	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return common.WebSearchResults{}, err
	}

	req.Header.Add("Accept", "text/html")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Cache-Control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return common.WebSearchResults{}, err
	}

	if res.StatusCode != http.StatusOK {
		return common.WebSearchResults{}, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var body []byte
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		gzipReader, err := gzip.NewReader(res.Body)
		if err != nil {
			return common.WebSearchResults{}, err
		}
		defer func() {
			_ = gzipReader.Close()
		}()

		body, err = io.ReadAll(gzipReader)
		if err != nil {
			return common.WebSearchResults{}, err
		}
	case "deflate":
		deflateReader := flate.NewReader(res.Body)
		defer func() {
			_ = deflateReader.Close()
		}()

		body, err = io.ReadAll(deflateReader)
		if err != nil {
			return common.WebSearchResults{}, err
		}
	default:
		if res.Uncompressed {
			body, err = io.ReadAll(res.Body)
			if err != nil {
				return common.WebSearchResults{}, err
			}
		} else {
			return common.WebSearchResults{}, fmt.Errorf("unsupported content encoding: %s", res.Header.Get("Content-Encoding"))
		}
	}

	reader := strings.NewReader(string(body))
	doc, err := html.Parse(reader)
	if err != nil {
		return common.WebSearchResults{}, err
	}

	return parseNode(doc), nil
}

func parseNode(node *html.Node) common.WebSearchResults {
	var results common.WebSearchResults
	if strings.Contains(getAttr(node.Attr, "class"), "result__body") {
		results.Results = append(results.Results, parseResultBodyNode(node))
	}

	if node.NextSibling != nil {
		results.Results = append(results.Results, parseNode(node.NextSibling).Results...)
	}

	if node.FirstChild != nil {
		results.Results = append(results.Results, parseNode(node.FirstChild).Results...)
	}

	return results
}

func parseResultBodyNode(node *html.Node) common.WebResult {
	var title, link, description string
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if strings.Contains(getAttr(child.Attr, "class"), "result__title") && child.FirstChild != nil && child.FirstChild.NextSibling != nil {
			link = getAttr(child.FirstChild.NextSibling.Attr, "href")
			if child.FirstChild.NextSibling.FirstChild != nil {
				title = child.FirstChild.NextSibling.FirstChild.Data
			}
		} else if strings.Contains(getAttr(child.Attr, "class"), "result__snippet") {
			description = getDataRecursive(child)
		}
	}

	return common.WebResult{
		Title:       title,
		URL:         fixURL(link),
		Description: description,
	}
}

func getAttr(attrs []html.Attribute, key string) string {
	for _, attr := range attrs {
		if attr.Key == key {
			return attr.Val
		}
	}

	return ""
}

func getDataRecursive(node *html.Node) string {
	var result string
	if node.FirstChild != nil {
		result += getDataRecursive(node.FirstChild)
	}

	if node.Type == html.TextNode && !regexp.MustCompile(`^\s*\n\s*$`).MatchString(node.Data) {
		result += node.Data
	}

	if node.NextSibling != nil {
		result += getDataRecursive(node.NextSibling)
	}

	return result
}
