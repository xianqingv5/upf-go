package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
)

func client(rootCa, rootKey string) *http.Client {
	var tr *http.Transport
	certs, err := tls.LoadX509KeyPair(rootCa, rootKey)
	if err != nil {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		ca, err := x509.ParseCertificate(certs.Certificate[0])
		if err != nil {
			return &http.Client{Transport: tr}
		}
		pool := x509.NewCertPool()
		pool.AddCert(ca)

		tr = &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: pool},
		}

	}
	return &http.Client{Transport: tr}
}

func Wget(url, name, rootCa, rootKey string) (bool, error) {
	fmt.Println("1")
	resp, err := client(rootCa, rootKey).Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	File, err := os.Create(name)
	if err != nil {
		return false, err
	}
	io.Copy(File, resp.Body)
	File.Close()
	fmt.Println("ok")
	return true, nil
}

func main() {
	bool, err := Wget("https://api.searchads.apple.com/api/v1/acls", "test", "songhq.pem", "songhq.key")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(bool)
}
