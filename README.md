# Search Tools for GPTScript

This repo contains search tools for GPTScript. We currently support Bing (web and image), Brave (web and image), DuckDuckGo (web only), and Google (web and image).

Web search output is in the following format:

```
Title: the title of the web page
URL: the link to the web page
Description: a short snippet from the web page
```

Image search output is in the following format:

```
Title: the title of the image
Source: the link to the web page where the image came from
Image URL: the link to the image
```

For usage examples with GPTScript, see the `examples` directory.

> **You are responsible for ensuring that your use of these search APIs with GPTScript does not violate the terms of service of the respective search engines.**

## Usage

Each of the search tools live in their own subdirectory. They can be referenced like `github.com/gptscript-ai/search/<tool>`.
For example, `github.com/gptscript-ai/search/duckduckgo` references the DuckDuckGo tool.
The options are `bing`, `bing-image`, `brave`, `brave-image`, `duckduckgo`, `google`, and `google-image`.

```bash
gptscript --cache=false github.com/gptscript-ai/search/duckduckgo '{"q":"best AI models for function calling"}'
```

Specific details and instructions for each search engine follow.

### Bing

The `bing` and `bing-image` tools return search results from the [Bing Web Search API](https://www.microsoft.com/en-us/bing/apis/bing-web-search-api).

The credential tool will ask for your API token.

### Brave

The `brave` and `brave-image` tools return search results from the [Brave Search API](https://brave.com/search/api/).

The credential tool will ask for your API token.

### DuckDuckGo

The `duckduckgo` tool returns search results from the [DuckDuckGo HTML-only Site](https://html.duckduckgo.com).

No API key is required to use this tool.

By default, this tool will make an HTTP request to DuckDuckGo and parse the results. If you do this enough times, it will start to get rate limited.
Rate limits can be more easily avoided by using Google Chrome in headless mode. The tool will do this if the `GPTSCRIPT_USE_CHROME` environment variable is set to `true`.

### Google

The `google` and `google-image` tools return search results from the [Google Custom Search JSON API](https://developers.google.com/custom-search/v1/overview).

The credential tool will ask for your search engine ID and your API token.
