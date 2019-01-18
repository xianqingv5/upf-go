package main

import (
	"fmt"
	"os/exec"
)

func main() {
	f, err := exec.Command("/bin/sh", "-c", "curl --cert /Users/qingfeng/Downloads/Certificates/songhq.p12 --pass 1 https://api.searchads.apple.com/api/v1/acls").Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(f))
	fmt.Println("nihao")
	fmt.Println("nnn")
	fmt.Println("")
	fmt.Println("songhq")
}
