# Search Tools for GPTScript (WIP)

This repo contains (or will contain) search tools for GPTScript. The plan is to support Bing, Brave, and Google.

Usage example:

```bash
./search brave '{"q":"What are the most populous cities in North America?","country":"US","search_lang":"en"}'
```

## Brave

For Brave, this tool uses the official [Brave Search API](https://brave.com/search/api/).
The environment variable `GPTSCRIPT_BRAVE_SEARCH_TOKEN` must be set to your API key in order for it to work.

> **Be sure that your plan allows you to use the Brave Search API for AI inference!** 

### JSON Parameters

- `q` (required): The search query.
- `country` (optional): The country to search from, in ISO 3166-2 format.
- `search_lang` (optional): The language to search in. Use the standard IETF language tag.
  - Exceptions:
    - Chinese must be either `zh-hans` (Simplified) or `zh-hant` (Traditional).
    - Japanese is `jp`.
    - Portuguese must be either `pt-br` (Brazil) or `pt-pt` (Portugal).
- `offset` (optional): The offset of the first result to return (used for pagination). Each query will return up to 20 results.
