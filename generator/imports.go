package generator

import (
	"../parser"
	"go/ast"
	"go/token"
)

func getImports(node *parser.CallNode) ast.Decl {
	if len(node.Args) < 2 {
		return nil
	}

	imports := node.Args[1:]
	specs := make([]ast.Spec, len(imports))

	for i, imp := range imports {
		if t := imp.Type(); t == parser.NodeVector {
			specs[i] = makeImportSpecFromVector(imp.(*parser.VectorNode))
		} else if t == parser.NodeString {
			path := makeBasicLit(token.STRING, imp.(*parser.StringNode).Value)
			specs[i] = makeImportSpec(path, nil)
		} else {
			panic("invalid import!")
		}
	}

	decl := makeGeneralDecl(token.IMPORT, specs)
	decl.Lparen = token.Pos(1) // Need this so we can have multiple imports
	return decl
}

func makeImportSpecFromVector(vect *parser.VectorNode) *ast.ImportSpec {
	if len(vect.Nodes) < 3 {
		panic("invalid use of import!")
	}

	if vect.Nodes[0].Type() != parser.NodeString {
		panic("invalid use of import!")
	}

	pathString := vect.Nodes[0].(*parser.StringNode).Value
	path := makeBasicLit(token.STRING, pathString)

	if vect.Nodes[1].Type() != parser.NodeIdent || vect.Nodes[1].(*parser.IdentNode).Ident != ":as" {
		panic("invalid use of import! expecting: \":as\"!!!")
	}
	name := ast.NewIdent(vect.Nodes[2].(*parser.IdentNode).Ident)

	return makeImportSpec(path, name)
}

func makeImportSpec(path *ast.BasicLit, name *ast.Ident) *ast.ImportSpec {
	return &ast.ImportSpec{
		Path: path,
		Name: name,
	}
}
