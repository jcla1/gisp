package main

import (
	"./parser"
	"bufio"
	// "bytes"
	"fmt"
	// "go/ast"
	// "go/printer"
	// "go/token"
	// "io/ioutil"
	"os"
)

type Any interface{}
type Symbol string

// func args(filename string) {
// 	b, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		panic(err)
// 	}

// 	p := parser.ParseFromString(filename, string(b)+"\n")

// 	a := generateAST(p)

// 	fset := token.NewFileSet()
// 	ast.Print(fset, a)

// 	var buf bytes.Buffer
// 	printer.Fprint(&buf, fset, a)
// 	fmt.Printf("%s\n", buf.String())
// }

func main() {
	// if len(os.Args) > 1 {
	// 	args(os.Args[1])
	// 	return
	// }

	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">> ")
		line, _, _ := r.ReadLine()
		p := parser.ParseFromString("<REPL>", string(line)+"\n")
		fmt.Println(p)

		// a := generateAST(p)
		// fset := goToken.NewFileSet()
		// var buf bytes.Buffer
		// printer.Fprint(&buf, fset, a)
		// fmt.Printf("%s\n", buf.String())
	}
}
