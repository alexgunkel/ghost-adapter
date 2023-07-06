package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	"github.com/alexgunkel/ghost-adapter/backend"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	url := flag.String("url", "localhost", "target url")
	key := flag.String("key", "", "")
	tag := flag.String("tag", "software", "")
	listen := flag.String("listen", "127.0.0.1:9013", "")
	flag.Parse()
	dst := fmt.Sprintf("%s/ghost/api/content/posts/?key=%s&filter=tag:%s", *url, *key, *tag)
	storage := backend.NewStorage(dst)

	getListVec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "requests",
			Namespace: "ghost",
			Help:      "collects number of incoming requests",
		},
		[]string{"host"},
	)
	pingVec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "calls",
			Namespace: "ghost",
			Help:      "collects number of calls",
		},
		[]string{"page"},
	)
	prometheus.MustRegister(getListVec, pingVec)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/api/ping", func(writer http.ResponseWriter, request *http.Request) {
		pingVec.WithLabelValues(getReferer(request)).Inc()
		writer.WriteHeader(201)
	})

	http.HandleFunc("/api/posts", func(writer http.ResponseWriter, request *http.Request) {
		getListVec.WithLabelValues(getReferer(request)).Inc()
		posts := storage.Posts()
		body, err := json.Marshal(posts)
		if err != nil {
			writer.WriteHeader(500)
			return
		}

		writer.Header().Set("Content-Type", "application/json")

		_, err = writer.Write(body)
		if err != nil {
			writer.WriteHeader(500)
			return
		}
	})

	err := http.ListenAndServe(*listen, nil)
	if err != nil {
		panic(err)
	}
}

func getReferer(request *http.Request) string {
	return request.Header.Get("Referer")
}
