package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	u "net/url"
)

var ReverseProxyWeb *httputil.ReverseProxy

var WebPort int = 3000

func makeReverseProxy(port int, ws bool) *httputil.ReverseProxy {
	url, _ := u.Parse(fmt.Sprintf("http://localhost:%d", port))
	proxy := httputil.NewSingleHostReverseProxy(url)
	if ws {
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = url.Scheme
			req.URL.Host = url.Host
			if req.Header.Get("Connection") != "Upgrade" {
				req.Header.Set("Connection", "Upgrade")
				req.Header.Set("Upgrade", "websocket")
			}
		}
	}
	return proxy
}

func handleRequest(writer http.ResponseWriter, request *http.Request) {
	//path := request.URL.Path
	request.Header.Set("X-Forwarded-Proto", "https")
	ip := request.RemoteAddr
	tokens := strings.Split(ip, ":")
	if len(tokens) > 1 {
		request.Header.Set("X-Real-Ip", tokens[0])
	}

	//host := request.Host
	ReverseProxyWeb.ServeHTTP(writer, request)
}
