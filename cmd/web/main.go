package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/otaviokr/check-my-speed/web"
)

var (
	url string
)

func init() {
	flag.StringVar(&url, "url", ":8088", "address to listen. Default: localhost:8088")
}

func main() {
	flag.Parse()
	router := web.NewRouter()
	log.Fatal(http.ListenAndServe(url, router))
}
