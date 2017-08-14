package main

import (
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
	http.HandleFunc("/", HelloServer)
	port := os.Getenv("PORT")
	certBits := os.Getenv("APP_CERT")
	keyBits := os.Getenv("APP_KEY")

	err := ioutil.WriteFile("server.crt", []byte(certBits), 0400)
	if err != nil {
		log.Fatal("Error while writing cert: ", err)
	}
	err = ioutil.WriteFile("server.key", []byte(keyBits), 0400)
	if err != nil {
		log.Fatal("Error while writing key ", err)
	}

	fileBytes, err := ioutil.ReadFile("server.crt")
	log.Printf("cert contents %s", string(fileBytes))
	fileBytes, err = ioutil.ReadFile("server.key")
	log.Printf("key contents %s", string(fileBytes))

	err = http.ListenAndServeTLS(fmt.Sprintf(":%s", port), "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
