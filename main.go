package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/coreos/go-omaha/omaha"
	"github.com/elazarl/goproxy"
	"github.com/fsouza/go-dockerclient"
)

func ParseRequest(req *http.Request) omaha.Request {
	var OmahaRequest omaha.Request

	defer req.Body.Close()
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	if data == nil {
		log.Fatalln("No request body data.")
	}

	err = xml.Unmarshal(data, &OmahaRequest)
	if err != nil {
		log.Fatal(err)
	}
	return OmahaRequest
}

func main() {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(goproxy.UrlIs("public.update.core-os.net/v1/update/")).HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest(goproxy.UrlIs("public.update.core-os.net/v1/update/")).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		OmahaRequest := ParseRequest(req)
		return req, nil
	})
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		resp.Header.Set("X-COREPXE", "corepxe")
		return resp
	})
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(":8080", proxy))

}
