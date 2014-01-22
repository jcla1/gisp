package main

var square = func(x Any) Any {
	println("Hello, World!")
	times(x, x)
	return func(y Any) Any {
		return id(y)
	}(x)
}

// Output ends here!

type Any interface{}

func times(x, y Any) int {
	return x.(int) * y.(int)
}

func id(y Any) Any {
	return y
}

func main() {
	square(10)
}
