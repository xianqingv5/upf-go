package main

import (
	"time"
)

var c int

func counter() int {
	c++
	return c
}

func main2() {
	a := 100
	go func(x, y int) {
		time.Sleep(time.Second)
		println("go:", x, y)
	}(a, counter())

	a += 100
	println("main:", a, counter())

	time.Sleep(time.Second * 3)
}

func main() {
	exit := make(chan struct{})

	go func() {
		time.Sleep(time.Second)
		println("goroutine done.")

		close(exit)
	}()

	println("main ...")
	<-exit
	println("main exit.")
}
