package main

import (
	"fmt"
	"os"

	"github.com/gptscript-ai/search/pkg/bing"
	"github.com/gptscript-ai/search/pkg/brave"
	"github.com/gptscript-ai/search/pkg/common"
	"github.com/gptscript-ai/search/pkg/duckduckgo"
	"github.com/gptscript-ai/search/pkg/google"
	"github.com/sirupsen/logrus"
)

const (
	TypeImage = "image"
	TypeWeb   = "web"

	EngineBing   = "bing"
	EngineBrave  = "brave"
	EngineDDG    = "duckduckgo"
	EngineGoogle = "google"
)

func main() {
	if len(os.Args) != 4 {
		logrus.Errorf("Usage: %s <web | image> <search engine> <JSON parameters>", os.Args[0])
		os.Exit(1)
	}

	searchType, engine, input := os.Args[1], os.Args[2], os.Args[3]

	switch searchType {
	case TypeImage:
		image(engine, input)
	case TypeWeb:
		web(engine, input)
	default:
		logrus.Errorf("Unsupported search type: %s", os.Args[1])
		os.Exit(1)
	}
}

func web(engine, input string) {
	var (
		results common.WebSearchResults
		err     error
	)
	switch engine {
	case EngineBing:
		results, err = bing.Search(input)
	case EngineBrave:
		results, err = brave.Search(input)
	case EngineDDG:
		results, err = duckduckgo.Search(input)
	case EngineGoogle:
		results, err = google.Search(input)
	default:
		logrus.Errorf("Unsupported search engine for web search: %s", engine)
		os.Exit(1)
	}

	if err != nil {
		logrus.Errorf(err.Error())
		os.Exit(1)
	}

	fmt.Print(results.ToString())
}

func image(engine, input string) {
	var (
		results common.ImageSearchResults
		err     error
	)
	switch engine {
	case EngineBrave:
		results, err = brave.SearchImage(input)
	case EngineGoogle:
		results, err = google.SearchImage(input)
	default:
		logrus.Errorf("Unsupported search engine for image search: %s", engine)
		os.Exit(1)
	}

	if err != nil {
		logrus.Errorf(err.Error())
		os.Exit(1)
	}

	fmt.Print(results.ToString())
}
