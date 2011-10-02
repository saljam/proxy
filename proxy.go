package main

import (
	"http"
	"log"
	"io"
)

type proxy struct{}

func (p *proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	
	if req.RawURL[0] == '/' {
		newURL := "http://" + req.Host + req.RawURL
		outReq, _ := http.NewRequest(req.Method, newURL, req.Body)

		for k, vv := range req.Header {
			for _, v := range vv {
				outReq.Header.Add(k, v)
			}
		}

		req = outReq
	}

	log.Println("request", req.Method, req.URL, req.UserAgent())
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Println("proxy fail:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	header := w.Header()
	for k, vv := range res.Header {
		for _, v := range vv {
			header.Add(k, v)
		}
	}

	w.WriteHeader(res.StatusCode)

	if res.Body != nil {
		io.Copy(w, res.Body)
	}
}

func main() {
	p := &proxy{}
	err := http.ListenAndServe(":3128", p)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.String())
	}
}

