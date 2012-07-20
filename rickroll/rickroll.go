package main

import (
	"os"
	"net/http"
	"log"
	"flag"
	"bitbucket.org/saljam/proxy/proxy"
)

func rickRoll(res *http.Response) *http.Response {
	if ct, ok := res.Header["Content-Type"]; !ok || len(ct) == 0{
		return res
	}

	if res.Header["Content-Type"][0] == "video/x-flv" {
		res.Header.Del("Content-Length")
		if file, err := os.Open(*filename); err == nil {
			res.Body = file
			log.Println("win!")
		} else {
			log.Println("failed to open file!")
		}
	}
	return res
}

var filename = flag.String("filename", "", "Path to the replacement .flv file.")

func main() {
	flag.Parse()
	
	p := &proxy.Proxy{
		proxy.ReqManglers{},
		proxy.ResManglers{rickRoll},
	}

	err := http.ListenAndServe(":3128", p)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
