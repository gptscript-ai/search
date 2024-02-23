package google

import "github.com/gptscript-ai/search/pkg/common"

const count = "10" // 10 is the max allowed by Google

// List sources:
// countries: https://developers.google.com/custom-search/docs/json_api_reference#countryCollections
// everything else: https://developers.google.com/custom-search/v1/reference/rest/v1/cse/list#request
var (
	SupportedCountries           = []string{"AF", "AL", "DZ", "AS", "AD", "AO", "AI", "AQ", "AG", "AR", "AM", "AW", "AU", "AT", "AZ", "BS", "BH", "BD", "BB", "BY", "BE", "BZ", "BJ", "BM", "BT", "BO", "BA", "BW", "BV", "BR", "IO", "BN", "BG", "BF", "BI", "KH", "CM", "CA", "CV", "KY", "CF", "TD", "CL", "CN", "CX", "CC", "CO", "KM", "CG", "CD", "CK", "CR", "CI", "HR", "CU", "CY", "CZ", "DK", "DJ", "DM", "DO", "TP", "EC", "EG", "SV", "GQ", "ER", "EE", "ET", "EU", "FK", "FO", "FJ", "FI", "FR", "FX", "GF", "PF", "TF", "GA", "GM", "GE", "DE", "GH", "GI", "GR", "GL", "GD", "GP", "GU", "GT", "GN", "GW", "GY", "HT", "HM", "VA", "HN", "HK", "HU", "IS", "IN", "ID", "IR", "IQ", "IE", "IL", "IT", "JM", "JP", "JO", "KZ", "KE", "KI", "KP", "KR", "KW", "KG", "LA", "LV", "LB", "LS", "LR", "LY", "LI", "LT", "LU", "MO", "MK", "MG", "MW", "MY", "MV", "ML", "MT", "MH", "MQ", "MR", "MU", "YT", "MX", "FM", "MD", "MC", "MN", "MS", "MA", "MZ", "MM", "NA", "NR", "NP", "NL", "AN", "NC", "NZ", "NI", "NE", "NG", "NU", "NF", "MP", "NO", "OM", "PK", "PW", "PS", "PA", "PG", "PY", "PE", "PH", "PN", "PL", "PT", "PR", "QA", "RE", "RO", "RU", "RW", "SH", "KN", "LC", "PM", "VC", "WS", "SM", "ST", "SA", "SN", "CS", "SC", "SL", "SG", "SK", "SI", "SB", "SO", "ZA", "GS", "ES", "LK", "SD", "SR", "SJ", "SZ", "SE", "CH", "SY", "TW", "TJ", "TZ", "TH", "TG", "TK", "TO", "TT", "TN", "TR", "TM", "TC", "TV", "UG", "UA", "AE", "UK", "US", "UM", "UY", "UZ", "VU", "VE", "VN", "VG", "VI", "WF", "EH", "YE", "YU", "ZM", "ZW"}
	SupportedLanguages           = []string{"ar", "bg", "ca", "cs", "da", "de", "el", "en", "es", "et", "fi", "fr", "hr", "hu", "id", "is", "it", "iw", "ja", "ko", "lt", "lv", "nl", "no", "pl", "pt", "ro", "ru", "sk", "sl", "sr", "sv", "tr", "zh-CN", "zh-TW"}
	SupportedImageSizes          = []string{"huge", "icon", "large", "medium", "small", "xlarge", "xxlarge"}
	SupportedImageTypes          = []string{"clipart", "face", "lineart", "stock", "photo", "animated"}
	SupportedImageDominantColors = []string{"black", "blue", "brown", "gray", "green", "pink", "purple", "teal", "white", "yellow"}
	SupportedImageColorTypes     = []string{"color", "gray", "mono", "trans"}
)

type webParams struct {
	Query      string `json:"q"`
	Country    string `json:"country"`
	SearchLang string `json:"search_lang"`
	Offset     string `json:"offset"`
}

type webAPIResponse struct {
	Items []struct {
		Title   string `json:"title"`
		Link    string `json:"link"`
		Snippet string `json:"snippet"`
	} `json:"items"`
}

func (r webAPIResponse) toSearchResults() common.WebSearchResults {
	var results common.WebSearchResults
	for _, item := range r.Items {
		results.Results = append(results.Results, common.WebResult{
			Title:       item.Title,
			URL:         item.Link,
			Description: item.Snippet,
		})
	}
	return results
}

type imageParams struct {
	Query              string `json:"q"`
	Country            string `json:"country"`
	SearchLang         string `json:"search_lang"`
	Offset             string `json:"offset"`
	ImageSize          string `json:"imgSize"`
	ImageType          string `json:"imgType"`
	ImageDominantColor string `json:"imgDominantColor"`
	ImageColorType     string `json:"imgColorType"`
}

type imageAPIResponse struct {
	Items []struct {
		Title string `json:"title"`
		Link  string `json:"link"`
		Image struct {
			ContextLink string `json:"contextLink"`
		} `json:"image"`
	} `json:"items"`
}

func (r imageAPIResponse) toSearchResults() common.ImageSearchResults {
	var results common.ImageSearchResults
	for _, item := range r.Items {
		results.Results = append(results.Results, common.ImageResult{
			Title:    item.Title,
			Source:   item.Image.ContextLink,
			ImageURL: item.Link,
		})
	}
	return results
}
