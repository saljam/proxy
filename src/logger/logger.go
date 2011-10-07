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
		[]func(*http.Request)*http.Request{logRequest, logRequest},
		[]func(*http.Response)*http.Response{},
	}
	err := http.ListenAndServe(":3128", p)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.String())
	}
}

