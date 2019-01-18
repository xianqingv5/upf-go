package main

import (
	"fmt"
)

func test(x *int) {
	println(&x, x) // 输出形参x的地址
}

func test2(a ...int) {
	fmt.Println(a)
}

func test3(a ...int) {
	for i := range a {
		a[i] += 100
	}
}

// 作为第一类对象，不管是普通函数，还是匿名函数都可作为struct字段，或通过channel传递
func testStruct() {
	type calc struct {
		mul func(x, y int) int
	}

	x := calc{
		mul: func(x, y int) int {
			return x * y
		},
	}

	println(x.mul(2, 3))
}

func testChannel() {
	// c := make(chan func(int, int) int, 2)

}

func main() {
	a := 0x100
	p := &a
	println(&p, p)

	test(p)
	//----------------------------------
	b := [3]int{10, 20, 30}
	test2(b[:]...)
	//----------------------------------
	c := []int{10, 20, 30}
	test3(c...)
	fmt.Println(c)
}
