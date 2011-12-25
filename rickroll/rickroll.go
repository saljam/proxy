package main

import (
	"os"
	"net/http"
	"log"
	"fmt"
	"flag"
	"bitbucket.org/saljam/proxy/proxy"
)

func rickRoll(res *http.Response) *http.Response {
	if ct, ok := res.Header["Content-Type"]; !ok || len(ct) == 0{
		return res
	}

	if res.Header["Content-Type"][0] == "video/x-flv" {
		res.Header.Del("Content-Length")
		if file, err := os.Open(filename); err == nil {
			res.Body = file
			log.Println("win!")
		} else {
			log.Println("failed to open file!")
		}
	}
	return res
}

var filename string

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: rickroll filename")
		os.Exit(2)
	}

	filename = args[0]

	p := &proxy.Proxy{
		proxy.ReqManglers{},
		proxy.ResManglers{rickRoll},
	}

	err := http.ListenAndServe(":3128", p)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
