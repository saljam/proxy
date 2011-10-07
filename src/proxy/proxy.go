package proxy

import (
	"http"
	"log"
	"io"
)

type Proxy struct{
	RequestManglers  []func(*http.Request) *http.Request
	ResponseManglers []func(*http.Response) *http.Response
}

func copyHeader(from, to http.Header) {
	for hdr, items := range from {
		for _, item := range items {
			to.Add(hdr, item)
		}
	}
}

func canonicalizeURL(req *http.Request) *http.Request {
	newURL := "http://" + req.Host + req.RawURL
	outReq, _ := http.NewRequest(req.Method, newURL, req.Body)
	copyHeader(req.Header, outReq.Header)
	return outReq
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Raw[0] == '/' {
		req = canonicalizeURL(req)
	}

	for _, f := range p.RequestManglers {
		req = f(req)
	}

	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Println("proxy fail:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, f := range p.ResponseManglers {
		res = f(res)
	}

	copyHeader(res.Header, rw.Header())

	rw.WriteHeader(res.StatusCode)

	if res.Body != nil {
		io.Copy(rw, res.Body)
	}
}

