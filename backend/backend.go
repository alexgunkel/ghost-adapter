package backend

import (
	"io"
	"log"
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
			} else {
				log.Printf("%s: %s", url, err)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return Extract(body)
}
