package main

func test(x int) func() {
	println(&x)

	return func() {
		println(&x, x)
	}
}

func test2() []func() {
	var s []func()
	for i := 0; i < 2; i++ {
		x := i
		s = append(s, func() {
			println(&x, x)
		})
	}

	return s
}

func test3(x int) (func(), func()) {
	return func() {
			println(x)
			x += 10
		}, func() {
			println(x)
		}
}

func main() {
	f := test(0x100)
	f()

	test2()

	a, b := test3(100)
	a()
	b()
}
