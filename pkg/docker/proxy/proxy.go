package dockerproxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func getScope(opt HandlerOptions) string {
	var scope string
	if opt.namespace != "" && opt.namespace2 == "" {
		scope = fmt.Sprintf("library/%s", opt.namespace)
	} else {
		scope = fmt.Sprintf("%s/%s", opt.namespace1, opt.namespace2)
	}
	return scope
}

type ReverseProxyDockerHub struct {
	accessToken string
	proxy       *httputil.ReverseProxy
}

func NewReverseProxyDockerHub(upstream string) *ReverseProxyDockerHub {
	_url, err := url.Parse(upstream)
	if err != nil {
		log.Fatal(err)
	}
	p := &ReverseProxyDockerHub{
		accessToken: "",
		proxy:       httputil.NewSingleHostReverseProxy(_url),
	}
	return p
}

func (p *ReverseProxyDockerHub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}
