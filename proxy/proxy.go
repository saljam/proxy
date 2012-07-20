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

// ReqManglers and ResManglers are a slice of fucntions which modify requests
// and replies respectively.
// These functions must have the signature:
//	func(*http.Request) *http.Request
// or:
//	func(*http.Response) *http.Response
type ReqManglers []func(*http.Request) *http.Request
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

func canonicalizeURL(req *http.Request) *http.Request {
	newURL := "http://" + req.Host + req.URL.Path
	outReq, _ := http.NewRequest(req.Method, newURL, req.Body)
	copyHeader(req.Header, outReq.Header)
	return outReq
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

