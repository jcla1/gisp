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
var MyExportedFunc = func(myArg, f Any) Any {
	return f.(func(Any)Any)(10)
}

var ExportedVar = 10
func main() {
	x := "Hello, World!"
	println(x)
	{
		x := 123
		MyExportedFunc(x)
	}
}
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	// (f.Decls[0].(*ast.GenDecl)).Specs[0].Name.Obj = nil
	// ((f.Decls[0].(*ast.GenDecl)).Specs[0].(*ast.TypeSpec).Name.Obj) = nil
	ast.Print(fset, f)

	// Print the AST.
	var buf bytes.Buffer
    printer.Fprint(&buf, fset, f)
	fmt.Println(buf.String())

}

