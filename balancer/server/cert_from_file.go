package server

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/foomo/simplecert"
)

func ServeCertFromFile() {
	//domainList := os.Getenv("BALANCER_DOMAINS")
	ReverseProxyBackend = makeReverseProxy(BackendPort, false)
	ReverseProxyWeb = makeReverseProxy(WebPort, false)

	go http.ListenAndServe(":80", http.HandlerFunc(simplecert.Redirect))
	go http.ListenAndServe(":8082", http.HandlerFunc(handleLocal))

	path := "/certs/file.cert"
	tlsconf, err := loadTLSConfigFromFile(path)
	fmt.Println(err)

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

func loadTLSConfigFromFile(filePath string) (*tls.Config, error) {
	tlsconf := &tls.Config{}

	tlsconf.Certificates = make([]tls.Certificate, 1)
	cert, err := tls.LoadX509KeyPair(filePath, filePath)
	if err != nil {
		return nil, err
	}
	tlsconf.Certificates[0] = cert

	return tlsconf, nil
}
