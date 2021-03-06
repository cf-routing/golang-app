package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

type AppInfo struct {
	Name               string   `json:"name"`
	Routes             []string `json:"uris"`
	Guid               string   `json:"application_id"`
	DiegoCellAddress   string
	ContainerNetworkIP string
}

func GetAppInfo() (*AppInfo, error) {
	var appInfo AppInfo
	err := json.Unmarshal([]byte(os.Getenv("VCAP_APPLICATION")), &appInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to parse JSON from VCAP_APPLICATION env var: %s", err)
	}
	appInfo.DiegoCellAddress = os.Getenv("CF_INSTANCE_ADDR")
	appInfo.ContainerNetworkIP = os.Getenv("CF_INSTANCE_INTERNAL_IP")
	return &appInfo, nil
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
	w.Write([]byte(fmt.Sprintf("I see you're connecting from %s\n", req.RemoteAddr)))
	appInfo, err := GetAppInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	json.NewEncoder(w).Encode(appInfo)
}

func DumpRequestServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	reqBytes, err := httputil.DumpRequest(req, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	w.Write([]byte(fmt.Sprintf("Request Info: %s\n", string(reqBytes))))
}

func main() {
	var err error
	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/dump", DumpRequestServer)
	port := os.Getenv("PORT")
	tlsEnv := os.Getenv("TLS")
	mtlsEnv := os.Getenv("MTLS")
	cipherSuite := os.Getenv("CIPHER")

	tlsEnabled := tlsEnv != "false"
	mtlsEnabled := mtlsEnv != "false"
	if !tlsEnabled && mtlsEnabled {
		log.Fatal("invalid config: mtls requires tls")
	}
	tlsConfig := &tls.Config{}
	if tlsEnabled {
		if cipherSuite != "" {
			ciphers := strings.Split(cipherSuite, ":")
			tlsConfig.CipherSuites = cipherValue(ciphers)
		}

	}
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
		Addr:         fmt.Sprintf(":%s", port),
		TLSConfig:    tlsConfig,
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){},
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

func cipherValue(ciphers []string) []uint16 {
	cipherMap := map[string]uint16{
		"RC4-SHA":                                 0x0005, // openssl formatted values
		"DES-CBC3-SHA":                            0x000a,
		"AES128-SHA":                              0x002f,
		"AES256-SHA":                              0x0035,
		"AES128-SHA256":                           0x003c,
		"AES128-GCM-SHA256":                       0x009c,
		"AES256-GCM-SHA384":                       0x009d,
		"ECDHE-ECDSA-RC4-SHA":                     0xc007,
		"ECDHE-ECDSA-AES128-SHA":                  0xc009,
		"ECDHE-ECDSA-AES256-SHA":                  0xc00a,
		"ECDHE-RSA-RC4-SHA":                       0xc011,
		"ECDHE-RSA-DES-CBC3-SHA":                  0xc012,
		"ECDHE-RSA-AES128-SHA":                    0xc013,
		"ECDHE-RSA-AES256-SHA":                    0xc014,
		"ECDHE-ECDSA-AES128-SHA256":               0xc023,
		"ECDHE-RSA-AES128-SHA256":                 0xc027,
		"ECDHE-RSA-AES128-GCM-SHA256":             0xc02f,
		"ECDHE-ECDSA-AES128-GCM-SHA256":           0xc02b,
		"ECDHE-RSA-AES256-GCM-SHA384":             0xc030,
		"ECDHE-ECDSA-AES256-GCM-SHA384":           0xc02c,
		"ECDHE-RSA-CHACHA20-POLY1305":             0xcca8,
		"ECDHE-ECDSA-CHACHA20-POLY1305":           0xcca9,
		"TLS_RSA_WITH_RC4_128_SHA":                0x0005, // RFC formatted values
		"TLS_RSA_WITH_3DES_EDE_CBC_SHA":           0x000a,
		"TLS_RSA_WITH_AES_128_CBC_SHA":            0x002f,
		"TLS_RSA_WITH_AES_256_CBC_SHA":            0x0035,
		"TLS_RSA_WITH_AES_128_CBC_SHA256":         0x003c,
		"TLS_RSA_WITH_AES_128_GCM_SHA256":         0x009c,
		"TLS_RSA_WITH_AES_256_GCM_SHA384":         0x009d,
		"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":        0xc007,
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":    0xc009,
		"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":    0xc00a,
		"TLS_ECDHE_RSA_WITH_RC4_128_SHA":          0xc011,
		"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":     0xc012,
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":      0xc013,
		"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":      0xc014,
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256": 0xc023,
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":   0xc027,
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":   0xc02f,
		"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256": 0xc02b,
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":   0xc030,
		"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384": 0xc02c,
		"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305":    0xcca8,
		"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305":  0xcca9,
	}

	var cipherVals []uint16
	for _, c := range ciphers {
		if val, ok := cipherMap[c]; ok {
			cipherVals = append(cipherVals, val)
		}
	}
	return cipherVals
}
