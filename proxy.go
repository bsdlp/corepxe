package corepxe

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func proxySetup(corpxeChan chan []int) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(goproxy.UrlIs("public.update.core-os.net/v1/update/")).HandleConnect(goproxy.AlwaysMitm).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) *http.Request {
		corepxeChan <- 0
		return r
	})
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Reponse {
		r.Header.Set("X-COREPXE", "corepxe")
		return r
	})
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
