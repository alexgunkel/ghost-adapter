package backend

import (
	"encoding/json"
	"regexp"
	"unicode/utf8"
)

func Extract(body []byte) (res []Post, err error) {
	var result Result

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	res = make([]Post, len(result.Posts))

	for idx, p := range result.Posts {
		res[idx] = p.Post
		res[idx].Teaser = Convert(p.Html)
	}

	return res, nil
}

func Convert(in string) string {
	res := regexp.MustCompile("<.*?>").ReplaceAllString(in, "")

	if strLen := utf8.RuneCountInString(res); strLen > 700 {
		for pos := 700; pos < strLen; pos++ {
			if res[pos] == '.' {
				res = res[:(pos + 1)]
				break
			}
		}
	}

	return res
}

type Post struct {
	Title               string `json:"title"`
	Url                 string `json:"url"`
	Teaser              string `json:"teaser"`
	FeatureImage        string `json:"feature_image"`
	FeatureImageAlt     string `json:"feature_image_alt"`
	FeatureImageCaption string `json:"feature_image_caption"`
}

type RemotePost struct {
	Post
	Html    string `json:"html"`
	Excerpt string `json:"excerpt"`
}

type Result struct {
	Posts []RemotePost `json:"posts"`
}
