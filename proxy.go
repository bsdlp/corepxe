package corepxe

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func cpRespHandler(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	r.Header.Set("X-COREPXE", "corepxe")
	return r
}

func proxySetup(corpxeChan chan []int) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(goproxy.UrlIs("public.update.core-os.net/v1/update/")).HandleConnect(goproxy.AlwaysMitm)
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Reponse {
		return cpRespHandler(resp, ctx)
	})
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
