package main

import (
	"http"
	"log"
	"proxy"
)

func logRequest(r *http.Request) *http.Request {
	log.Println(r.Method, r.URL, r.UserAgent(), r.Header)
	return r
}

func main() {
	p := &proxy.Proxy{
		RequestManglers: []func(*http.Request)*http.Request{logRequest},
	}
	err := http.ListenAndServe(":3128", p)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.String())
	}
}

