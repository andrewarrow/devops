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

	intermediateCert, err := tls.LoadX509KeyPair("/certs/int.crt", "")
	fmt.Println(err)
	serverCert, err := tls.LoadX509KeyPair("/certs/server.crt", "/certs/file.key")
	fmt.Println(err)

	tlsconf := &tls.Config{
		Certificates: []tls.Certificate{serverCert, intermediateCert},
	}

	return tlsconf, nil
}
