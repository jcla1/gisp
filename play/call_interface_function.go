package main

import "fmt"

type Any interface{}

var MyToplevelFunc = func() {
	fmt.Println("In: MyToplevelFunc")
}

var MyMap = func(sequence, f Any) {
	for _, v := range sequence.([]Any) {
		fmt.Println(f.(func(Any) Any)(v))
	}
}

func main() {
	MyToplevelFunc()
	MyMap([]Any{1, 2, 3, 4}, func(val Any) Any {
		return 10 * val.(int)
	})
}
