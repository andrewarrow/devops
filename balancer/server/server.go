package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	u "net/url"

	"github.com/foomo/simplecert"
	"github.com/foomo/tlsconfig"
)

var ReverseProxyBackend *httputil.ReverseProxy
var ReverseProxyWeb *httputil.ReverseProxy

var BackendPort int = 8080
var WebPort int = 3000

var BalancerGuid = os.Getenv("BALANCER_GUID")

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
	path := request.URL.Path

	if strings.HasPrefix(path, "/"+BalancerGuid) {
		tokens := strings.Split(path, "/")
		last := tokens[len(tokens)-1]

		if last == "backend" {
			writer.Write([]byte(fmt.Sprintf("%d", BackendPort)))
			return
		} else if last == "web" {
			writer.Write([]byte(fmt.Sprintf("%d", WebPort)))
			return
		}

		if last == "8080" {
			if BackendPort == 8080 {
				BackendPort++
			} else {
				BackendPort--
			}
			ReverseProxyBackend = makeReverseProxy(BackendPort, false)
		} else if last == "3000" {
			if WebPort == 3000 {
				WebPort++
			} else {
				WebPort--
			}
			ReverseProxyWeb = makeReverseProxy(WebPort, false)
		}

		return
	}

	host := request.Host
	if strings.Contains(host, "api") {
		ReverseProxyBackend.ServeHTTP(writer, request)
	} else if strings.Contains(host, "web") {
		ReverseProxyWeb.ServeHTTP(writer, request)
	} else {
		ReverseProxyWeb.ServeHTTP(writer, request)
	}
}

func Serve() {
	domainList := os.Getenv("BALANCER_DOMAINS")
	ReverseProxyBackend = makeReverseProxy(BackendPort, false)
	ReverseProxyWeb = makeReverseProxy(WebPort, false)

	cfg := simplecert.Default
	cfg.Domains = strings.Split(domainList, ",")
	cfg.CacheDir = "/certs"
	cfg.SSLEmail = os.Getenv("BALANCER_EMAIL")
	certReloader, err := simplecert.Init(cfg, nil)
	fmt.Println("err", err)

	go http.ListenAndServe(":80", http.HandlerFunc(simplecert.Redirect))
	go http.ListenAndServe(":8082", http.HandlerFunc(handleLocal))

	tlsconf := tlsconfig.NewServerTLSConfig(tlsconfig.TLSModeServerStrict)
	tlsconf.GetCertificate = certReloader.GetCertificateFunc()

	handler := http.HandlerFunc(handleRequest)

	s := &http.Server{
		Addr:      ":443",
		Handler:   handler,
		TLSConfig: tlsconf,
	}

	s.ListenAndServeTLS("", "")

	for {
		time.Sleep(time.Second)
	}
}

func handleLocal(writer http.ResponseWriter, request *http.Request) {
	service := request.Header.Get("AA-Service")

	if service == "backend" {
		ReverseProxyBackend.ServeHTTP(writer, request)
	} else if service == "web" {
		ReverseProxyWeb.ServeHTTP(writer, request)
	}
}
