package main

type Any interface{}

func main() {
	myFunc(func(x, y int, rest ...int) {
		println(x)
		println(y)
		println(rest)
	})
}

func myFunc(fn interface{}) {
	fn.(func(Any, Any, ...Any))(1, 2, 3, 4, 5)
}
