package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/printer"
	"fmt"
	"bytes"
)

func main() {
	// src is the input for which we want to print the AST.
	src := `
package main
type Any interface{}
var MyExportedFunc = func(myArg Any) Any {
	return 10
}
func main() {
	println("Hello, World!")
	MyExportedFunc(123)
}
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	f.Scope = nil

	ast.Print(fset, f)

	// Print the AST.
	var buf bytes.Buffer
    	printer.Fprint(&buf, fset, f)
	fmt.Println(buf.String())

}

