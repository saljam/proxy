package main

import (
	"os"
	"net/http"
	"log"
	"bitbucket.org/saljam/proxy/proxy"
	"bytes"
	"image"
	"image/jpeg"
	"image/gif"
	"image/png"
)

type imgBuffer struct{
	*bytes.Buffer
}

func (b imgBuffer) Close() os.Error {
	return nil
}

func flipImage(m image.Image) *image.RGBA {
	b := m.Bounds()
	s := image.Point{b.Dx(), b.Dy()}
	n := image.NewRGBA(b)
	for i :=0; i < s.X; i++ {
		for j :=0; j < s.Y; j++ {
			n.Set(i, j, m.At(i, s.Y-j-1))
		}
	}
	return n
}

func flipImageRes(res *http.Response) *http.Response {
	switch res.Header["Content-Type"][0] {
	case "image/jpeg":
		m, err := jpeg.Decode(res.Body)
		if err != nil {
			log.Println("image error:", err)
			return res
		}
		buf := imgBuffer{bytes.NewBuffer([]byte{})}
		n := flipImage(m)
		jpeg.Encode(buf, n, nil)
		res.Header.Del("Content-Length")
		res.Body = buf
	case "image/gif":
		m, err := gif.Decode(res.Body)
		if err != nil {
			log.Println("image error:", err)
			return res
		}
		buf := imgBuffer{bytes.NewBuffer([]byte{})}
		n := flipImage(m)
		png.Encode(buf, n)
		res.Header.Del("Content-Length")
		res.Header.Del("Content-Type")
		res.Header.Add("Content-Type", "image/png")
		res.Body = buf
	case "image/png":
		m, err := png.Decode(res.Body)
		if err != nil {
			log.Println("image error:", err)
			return res
		}
		buf := imgBuffer{bytes.NewBuffer([]byte{})}
		n := flipImage(m)
		png.Encode(buf, n)
		res.Header.Del("Content-Length")
		res.Header.Del("Content-Type")
		res.Header.Add("Content-Type", "image/png")
		res.Body = buf
	}
	return res
}

func logRequest(r *http.Request) *http.Request {
	log.Println(r.Method, r.URL, r.UserAgent(), r.Header)
	return r
}

func main() {
	p := &proxy.Proxy{
		proxy.ReqManglers{logRequest},
		proxy.ResManglers{flipImageRes},
	}
	err := http.ListenAndServe(":3128", p)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

