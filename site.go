package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
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
	tlsEnv := os.Getenv("TLS_ENABLED")
	mtlsEnv := os.Getenv("MTLS")

	tlsEnabled := tlsEnv != "false"
	mtls := mtlsEnv != "false"
	tlsConfig := &tls.Config{}
	if tlsEnabled && mtls {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		certPool, err := x509.SystemCertPool()
		if err != nil {
			panic(err)
		}
		caCert, err := ioutil.ReadFile("ca.crt")
		if err != nil {
			log.Fatal("error reading ca cert: ", err)
		}
		if ok := certPool.AppendCertsFromPEM([]byte(caCert)); !ok {
			panic(errors.New("error adding caCert to cert pool"))
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
