package backend

import (
	"encoding/json"
	"io"
	"net/http"
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
	Title        string `json:"title"`
	Url          string `json:"url"`
	Excerpt      string `json:"excerpt"`
	FeatureImage string `json:"feature_image"`
}

type Result struct {
	Posts []Post `json:"posts"`
}
