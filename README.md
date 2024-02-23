# Search Tools for GPTScript

This repo contains search tools for GPTScript. We currently support Bing (web only), Brave (web and image), and Google (web and image).

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

## Tools

All tools are currently implemented in the tool.gpt file.

### Bing

The `bing` tool returns search results from the [Bing Web Search API](https://www.microsoft.com/en-us/bing/apis/bing-web-search-api).

The environment variable `GPTSCRIPT_BING_SEARCH_TOKEN` must be set to your API key in order for it to work.

### Brave

The `brave` and `brave-image` tools return search results from the [Brave Search API](https://brave.com/search/api/).

The environment variable `GPTSCRIPT_BRAVE_SEARCH_TOKEN` must be set to your API key in order for it to work.

### Google

The `google` and `google-image` tools return search results from the [Google Custom Search JSON API](https://developers.google.com/custom-search/v1/overview).

The environment variable `GPTSCRIPT_GOOGLE_SEARCH_TOKEN` must be set to your API key, and `GPTSCRIPT_GOOGLE_SEARCH_ENGINE_ID` must be set to your search engine ID in order for it to work.
