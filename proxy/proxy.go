/*
Proxy is a simple transperent web proxy written in go.

Basic usage:

	p := &proxy.Proxy{
		proxy.ReqManglers{logRequest},
		proxy.ResManglers{flipImageRes},
	}
	log.Fatal(http.ListenAndServe(":3128", p))

Where the members of the manglers are functions which implement the signatures:

	func(*http.Request) *http.Request
	func(*http.Response) *http.Response

*/
package proxy

import (
	"net/http"
	"log"
	"io"
)

// ReqManglers are a slice of fucntions which modify requests.
// These functions must have the signature:
//	func(*http.Request) *http.Request
type ReqManglers []func(*http.Request) *http.Request

// ResManglers are a slice of fucntions which modify responses.
// These functions must have the signature:
//	func(*http.Response) *http.Response
type ResManglers []func(*http.Response) *http.Response

// Proxy implements http.Handler.
type Proxy struct{
	RequestManglers  ReqManglers
	ResponseManglers ResManglers
}

func copyHeader(from, to http.Header) {
	for hdr, items := range from {
		for _, item := range items {
			to.Add(hdr, item)
		}
	}
}

// ServeHTTP proxies the request given and writes the response to w.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, f := range p.RequestManglers {
		req = f(req)
	}

	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Println("proxy fail:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, f := range p.ResponseManglers {
		res = f(res)
	}

	copyHeader(res.Header, w.Header())

	w.WriteHeader(res.StatusCode)

	if res.Body != nil {
		io.Copy(w, res.Body)
	}
}

