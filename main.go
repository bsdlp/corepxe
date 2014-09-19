package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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

func ParseResponse(res *http.Response) omaha.Response {
	var OmahaResponse omaha.Response
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	if data == nil {
		log.Fatalln("No response body data.")
	}

	err = xml.Unmarshal(data, &OmahaResponse)
	if err != nil {
		log.Fatal(err)
	}
	return OmahaResponse
}

func PackageToPXE(pkg omaha.Package, urlbase string, c chan int) {
	pkgURL := strings.Join(urlbase, pkg.Name)
}

func main() {
	var OriginalRequest http.Request

	c := make(chan int)

	dockerSocket := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(dockerSocket)
	if err != nil {
		log.Fatal(err)
	}

	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(goproxy.UrlIs("public.update.core-os.net/v1/update/")).HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest(goproxy.UrlIs("public.update.core-os.net/v1/update/")).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		OriginalRequest = req
		OmahaRequest := ParseRequest(req)
		return req, nil
	})
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		OmahaResponse := ParseResponse(resp)

		for i := range OmahaResponse.Apps {
			if OmahaResponse.Apps[i].UpdateCheck.Status == "ok" {
				for pkg := range packages {
					go PackageToPXE(pkg, OmahaReseponse.Apps[i].UpdateCheck.Urls[0], c)
				}
			}
		}

		goproxy.NewResponse(OriginalRequest, goproxy.ContentTypeHtml, http.StatusOK, CorePXEResponse)

		resp.Header.Set("X-COREPXE", "corepxe")
		return resp
	})
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
