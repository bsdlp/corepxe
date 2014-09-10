package corepxe

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func proxySetup(corepxeChan *chan int) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(goproxy.UrlIs("public.update.core-os.net/v1/update/")).HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest(goproxy.UrlIs("public.update.core-os.net/v1/update/")).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		*corepxeChan <- 0
		return req, nil
	})
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		resp.Header.Set("X-COREPXE", "corepxe")
		return resp
	})
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(":8080", proxy))
}

func main() {
	corepxeChan := make(chan int)
	proxySetup(&corepxeChan)
}
