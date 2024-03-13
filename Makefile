build:
	CGO_ENABLED=0 go build -o bing/bin/gptscript-go-tool -tags "${GO_TAGS}" -ldflags "-s -w" ./bing
	CGO_ENABLED=0 go build -o bing-image/bin/gptscript-go-tool -tags "${GO_TAGS}" -ldflags "-s -w" ./bing-image
	CGO_ENABLED=0 go build -o brave/bin/gptscript-go-tool -tags "${GO_TAGS}" -ldflags "-s -w" ./brave
	CGO_ENABLED=0 go build -o brave-image/bin/gptscript-go-tool -tags "${GO_TAGS}" -ldflags "-s -w" ./brave-image
	CGO_ENABLED=0 go build -o duckduckgo/bin/gptscript-go-tool -tags "${GO_TAGS}" -ldflags "-s -w" ./duckduckgo
	CGO_ENABLED=0 go build -o google/bin/gptscript-go-tool -tags "${GO_TAGS}" -ldflags "-s -w" ./google
	CGO_ENABLED=0 go build -o google-image/bin/gptscript-go-tool -tags "${GO_TAGS}" -ldflags "-s -w" ./google-image

tidy:
	go mod tidy

GOLANGCI_LINT_VERSION ?= v1.56.1
lint:
	if ! command -v golangci-lint &> /dev/null; then \
  		echo "Could not find golangci-lint, installing version $(GOLANGCI_LINT_VERSION)."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION); \
	fi

	golangci-lint run


validate: tidy lint
	if [ -n "$$(git status --porcelain)" ]; then \
		git status --porcelain; \
		echo "Encountered dirty repo!"; \
		git diff; \
		exit 1 \
	;fi
