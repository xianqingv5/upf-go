package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	pool := x509.NewCertPool()
	c, err := tls.LoadX509KeyPair("songhq.pem", "songhq.key")
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range c.Certificate {
		z, err := x509.ParseCertificate(v)
		if err != nil {
			fmt.Println(err)
			continue
		}
		pool.AddCert(z)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{ClientCAs: pool, Certificates: []tls.Certificate{c}, ClientAuth: tls.NoClientCert},
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", "https://api.searchads.apple.com/api/v1/search/geo?limit=10", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
}
