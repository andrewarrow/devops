package server

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/foomo/simplecert"
)

func ServeCertFromFile() {
	ReverseProxyWeb = makeReverseProxy(WebPort, false)

	go http.ListenAndServe(":80", http.HandlerFunc(simplecert.Redirect))

	domain := "apps.greyspace.co"
	path1 := fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", domain)
	path2 := fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", domain)
	tlsconf, err := loadTLSConfigFromFile(path1, path2)
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
func loadTLSConfigFromFile(path1, path2 string) (*tls.Config, error) {

	serverCert, err := tls.LoadX509KeyPair(path1, path2)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	tlsconf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}
	return tlsconf, nil
}

func loadTLSConfigFromFileInt(filePath string) (*tls.Config, error) {

	intermediateCert, err := tls.LoadX509KeyPair("/certs/int.crt", "")
	fmt.Println(err)
	serverCert, err := tls.LoadX509KeyPair("/certs/server.crt", "/certs/file.key")
	fmt.Println(err)

	tlsconf := &tls.Config{
		Certificates: []tls.Certificate{serverCert, intermediateCert},
	}

	return tlsconf, nil
}
