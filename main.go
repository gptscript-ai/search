package main

import (
	"fmt"
	"os"

	"github.com/gptscript-ai/search/pkg/bing"
	"github.com/gptscript-ai/search/pkg/brave"
	"github.com/gptscript-ai/search/pkg/common"
	"github.com/gptscript-ai/search/pkg/google"
	"github.com/sirupsen/logrus"
)

const (
	EngineBing   = "bing"
	EngineBrave  = "brave"
	EngineGoogle = "google"
)

func main() {
	if len(os.Args) != 3 {
		logrus.Errorf("Usage: %s <search engine> <JSON parameters>", os.Args[0])
		os.Exit(1)
	}

	var (
		results common.SearchResults
		err     error
	)
	switch os.Args[1] {
	case EngineBing:
		results, err = bing.Search(os.Args[2])
	case EngineBrave:
		results, err = brave.Search(os.Args[2])
	case EngineGoogle:
		results, err = google.Search(os.Args[2])
	default:
		logrus.Errorf("Unsupported search engine: %s", os.Args[1])
		os.Exit(1)
	}

	if err != nil {
		logrus.Errorf(err.Error())
		os.Exit(1)
	}

	fmt.Print(results.ToString())
}
