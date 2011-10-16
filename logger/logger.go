package main

import (
	"http"
	"log"
	"bitbucket.org/saljam/proxy/proxy"
)

func logRequest(r *http.Request) *http.Request {
	log.Println(r.Method, r.URL, r.UserAgent(), r.Header)
	return r
}

func main() {
	p := &proxy.Proxy{
		proxy.ReqManglers{logRequest},
		proxy.ResManglers{},
	}
	err := http.ListenAndServe(":3128", p)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.String())
	}
}

