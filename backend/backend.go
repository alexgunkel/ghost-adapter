package backend

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"
)

type Storage struct {
	posts []Post
	mtx   sync.RWMutex
}

func NewStorage(url string) *Storage {
	s := new(Storage)
	s.posts = []Post{}
	go func() {
		for {
			posts, err := getPosts(url)
			if err == nil {
				r := regexp.MustCompile("<.*?>")
				for idx, _ := range posts {
					posts[idx].TextBegin = r.ReplaceAllString(posts[idx].Html[:1000], "")
				}
				s.mtx.Lock()
				s.posts = posts
				s.mtx.Unlock()
			}

			<-time.After(5 * time.Minute)
		}
	}()

	return s
}

func (s *Storage) Posts() []Post {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	res := make([]Post, len(s.posts))
	copy(res, s.posts)

	return res
}

func getPosts(url string) ([]Post, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var result Result
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Posts, nil
}

type Post struct {
	Title               string `json:"title"`
	Url                 string `json:"url"`
	Html                string `json:"html"`
	Excerpt             string `json:"excerpt"`
	TextBegin           string `json:"text_begin"`
	FeatureImage        string `json:"feature_image"`
	FeatureImageAlt     string `json:"feature_image_alt"`
	FeatureImageCaption string `json:"feature_image_caption"`
}

type Result struct {
	Posts []Post `json:"posts"`
}
