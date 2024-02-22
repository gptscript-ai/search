package main

import (
	"fmt"
	"os"

	"github.com/gptscript-ai/search/pkg/bing"
	"github.com/gptscript-ai/search/pkg/brave"
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
		result string
		err    error
	)
	switch os.Args[1] {
	case EngineBing:
		result, err = bing.SearchBing(os.Args[2])
	case EngineBrave:
		result, err = brave.SearchBrave(os.Args[2])
	case EngineGoogle:
		result, err = google.SearchGoogle(os.Args[2])
	default:
		logrus.Errorf("Unsupported search engine: %s", os.Args[1])
		os.Exit(1)
	}

	if err != nil {
		logrus.Errorf(err.Error())
		os.Exit(1)
	}

	fmt.Print(result)
}
