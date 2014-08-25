package corepxe

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func cpReqHandler(r *http.Request, ctx *goproxy.ProxyCtx, corepxeChan chan []int) *http.Request {
	corepxeChan <- 0
	return r
}

func cpRespHandler(r *http.Response, ctx *goproxy.ProxyCtx, corepxeChan chan []int) *http.Response {
	r.Header.Set("X-COREPXE", "corepxe")
	return r
}

func proxySetup(corpxeChan chan []int) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(goproxy.UrlIs("public.update.core-os.net/v1/update/")).HandleConnect(goproxy.AlwaysMitm).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) *http.Request {
		return cpReqHandler(req, ctx, corpxeChan)
	})
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Reponse {
		return cpRespHandler(resp, ctx, corpxeChan)
	})
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
