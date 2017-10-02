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

// var cert = flag.String("cert", "", "cert for application to terminate TLS")
// var key = flag.String("key", "", "private key for application")

func main() {
	var err error
	http.HandleFunc("/", HelloServer)
	port := os.Getenv("PORT")
	serverCert := os.Getenv("APP_CERT")
	serverKey := os.Getenv("APP_KEY")
	tlsEnv := os.Getenv("TLS_ENABLED")
	mtlsEnv := os.Getenv("MTLS")
	caCert := os.Getenv("CA_CERT")

	tlsEnabled := tlsEnv != "false"
	if tlsEnabled {
		err = ioutil.WriteFile("server.crt", []byte(serverCert), 0400)
		if err != nil {
			log.Fatal("Error while writing cert: ", err)
		}
		err = ioutil.WriteFile("server.key", []byte(serverKey), 0400)
		if err != nil {
			log.Fatal("Error while writing key ", err)
		}
	}

	mtls := mtlsEnv != "false"
	tlsConfig := &tls.Config{}
	if tlsEnabled && mtls {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		certPool, err := x509.SystemCertPool()
		if err != nil {
			panic(err)
		}
		if caCert != "" {
			if ok := certPool.AppendCertsFromPEM([]byte(caCert)); !ok {
				panic(errors.New("error adding caCert to cert pool"))
			}
			tlsConfig.ClientCAs = certPool
		}
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
