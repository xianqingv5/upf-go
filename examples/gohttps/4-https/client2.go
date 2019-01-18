package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	certFile = flag.String("cert", "songhq.pem", "A PEM eoncoded certificate file.")
	keyFile  = flag.String("key", "songhq.key", "A PEM encoded private key file.")
	caFile   = flag.String("CA", "songhq.pem", "A PEM eoncoded CA's certificate file.")
)

func main() {
	flag.Parse()

	// Load client cert 加载客户端证书
	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(*caFile)
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: &tls.Config{ClientCAs: caCertPool, Certificates: []tls.Certificate{cert}, ClientAuth: tls.NoClientCert}}
	client := &http.Client{Transport: transport}

	// Do GET something
	resp, err := client.Get("https://api.searchads.apple.com/api/v1/acls")
	//client.Post()
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Dump response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(data))
}
