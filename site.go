package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

func main() {
	var err error
	http.HandleFunc("/", HelloServer)
	port := os.Getenv("PORT")
	tlsEnv := os.Getenv("TLS")
	mtlsEnv := os.Getenv("MTLS")

	tlsEnabled := tlsEnv != "false"
	mtlsEnabled := mtlsEnv != "false"
	if !tlsEnabled && mtlsEnabled {
		log.Fatal("invalid config: mtls requires tls")
	}
	tlsConfig := &tls.Config{}
	if tlsEnabled && mtlsEnabled {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		certPool, err := x509.SystemCertPool()
		if err != nil {
			log.Fatalf("opening system cert pool: %s", err)
		}
		caCert, err := ioutil.ReadFile("ca.crt")
		if err != nil {
			log.Fatal("error reading ca cert: ", err)
		}
		if ok := certPool.AppendCertsFromPEM([]byte(caCert)); !ok {
			log.Fatal("error adding caCert to cert pool")
		}
		tlsConfig.ClientCAs = certPool
	}

	httpServer := &http.Server{
		Addr:      fmt.Sprintf(":%s", port),
		TLSConfig: tlsConfig,
	}
	if tlsEnabled {
		err = httpServer.ListenAndServeTLS("server.crt", "server.key")
	} else {
		err = httpServer.ListenAndServe()
	}
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
