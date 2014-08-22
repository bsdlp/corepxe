package main

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	proxy.OnRequest().DoFunc(
		func(r *http.Request, _ *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			r.Header.Set("X-COREPXE", "corepxe")
			return r, nil
		})

	log.Fatal(http.ListenAndServe(":8080", proxy))
}
