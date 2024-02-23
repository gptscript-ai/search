# Search Tools for GPTScript (WIP)

This repo contains search tools for GPTScript. The plan is to support Bing, Brave, and Google.

Usage example:

```bash
./search brave '{"q":"What are the most populous cities in North America?","country":"US","search_lang":"en"}'
```

The output is in the following format:

```
Title: the title of the web page
URL: the link to the web page
Description: a short snippet from the web page
```

For usage examples with GPTScript, see the `examples` directory.

> **You are responsible for ensuring that your use of these search APIs with GPTScript does not violate the terms of service of the respective search engines.**

## Bing

For Bing, this tool uses the official [Bing Web Search API](https://www.microsoft.com/en-us/bing/apis/bing-web-search-api).
The environment variable `GPTSCRIPT_BING_SEARCH_TOKEN` must be set to your API key in order for it to work.

### JSON Parameters

- `q` (required): The search query.
- `country` (optional): The country to search from, in ISO 3166-2 format.
- `search_lang` (optional): The language to search in, in RFC 5646 format.
- `offset` (optional): The offset of the first result to return (used for pagination).

## Brave

For Brave, this tool uses the official [Brave Search API](https://brave.com/search/api/).
The environment variable `GPTSCRIPT_BRAVE_SEARCH_TOKEN` must be set to your API key in order for it to work.

### JSON Parameters

- `q` (required): The search query.
- `country` (optional): The country to search from, in ISO 3166-2 format.
- `search_lang` (optional): The language to search in. Use the standard IETF language tag.
  - Exceptions:
    - Chinese must be either `zh-hans` (Simplified) or `zh-hant` (Traditional).
    - Japanese is `jp`.
    - Portuguese must be either `pt-br` (Brazil) or `pt-pt` (Portugal).
- `offset` (optional): The offset of the first result to return (used for pagination).

## Google

For Google, this tool uses the official [Google Custom Search JSON API](https://developers.google.com/custom-search/v1/overview).
To use this API, you will need to create a [Programmable Search Engine](https://developers.google.com/custom-search/docs/tutorial/creatingcse) first.
The environment variables `GPTSCRIPT_GOOGLE_SEARCH_ENGINE_ID` and `GPTSCRIPT_GOOGLE_SEARCH_TOKEN` must be set to your Programmable Search Engine ID and API key, respectively.

### JSON Parameters

- `q` (required): The search query.
- `country` (optional): The country to search from. The list of supported countries can be found in the `tool.gpt` file.
- `search_lang` (optional). The language to search in. Use the standard IETF language tag.
  - Exception:
    - Chinese must be either `zh-CN` (Simplified) or `zh-TW` (Traditional).
- `offset` (optional): The offset of the first result to return (used for pagination).
