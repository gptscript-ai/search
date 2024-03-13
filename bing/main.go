package main

import (
	"fmt"
	"os"

	"github.com/gptscript-ai/search/pkg/bing"
	"github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) != 2 {
		logrus.Errorf("Usage: %s <JSON parameters>", os.Args[0])
		os.Exit(1)
	}

	results, err := bing.Search(os.Args[1])
	if err != nil {
		logrus.Errorf(err.Error())
		os.Exit(1)
	}
	fmt.Print(results.ToString())
}
