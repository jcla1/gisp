package main

import "fmt"
// Won't work. Funcs are not constant, but why?
const MyConstFunc = func() string {
    return "This is from your ConstFunc!"
}

func main() {
    fmt.Println(MyConstFunc())
}