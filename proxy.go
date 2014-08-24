package corepxe

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func reponseHandler(r *http.Response, ctx *goproxy.ProxyCtx, corpxeChan chan []int) *http.Response {
	r.Header.Set("X-COREPXE", "corepxe")
	return r
}

func proxySetup(corpxeChan chan []int) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(UrlIs("public.update.core-os.net/v1/update/")).HandleConnect(goproxy.AlwaysMitm)
	proxy.OnResponse().DoFunc(reponseHandler(corpxeChan))
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
