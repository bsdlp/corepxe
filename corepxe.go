package main

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnResponse().DoFunc(
		func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			r.Header.Set("X-COREPXE", "corepxe")
			return r
		})
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
