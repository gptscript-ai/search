package duckduckgo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/gptscript-ai/search/pkg/common"
)

func Search(input string) (common.WebSearchResults, error) {
	var params webParams
	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return common.WebSearchResults{}, err
	}

	return getSearchResults(params)
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
	ctx, _ := chromedp.NewContext(context.Background())
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
