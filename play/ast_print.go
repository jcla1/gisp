package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
)

func main() {
	// src is the input for which we want to print the AST.
	src := `
package main

func main() {
fmt.Println(myNode.Node.Type.Hello)
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
	// f.Imports = nil
	ast.Print(fset, f)

	// Print the AST.
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, f)
	fmt.Println(buf.String())

}
