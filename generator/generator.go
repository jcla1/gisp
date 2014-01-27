package generator

import (
	"../parser"
	"fmt"
	"go/ast"
	"go/token"
)

func GenerateAST(tree []parser.Node) *ast.File {
	f := &ast.File{Name: ast.NewIdent("main")}
	decls := make([]ast.Decl, 0, len(tree))

	if len(tree) < 1 {
		return f
	}

	// you can only have (ns ...) as the first form
	if isNSDecl(tree[0]) {
		name, imports := getNamespace(tree[0].(*parser.CallNode))

		f.Name = name
		if imports != nil {
			decls = append(decls, imports)
		}

		tree = tree[1:]
	}

	decls = append(decls, generateDecls(tree)...)

	f.Decls = decls
	return f
}

func generateDecls(tree []parser.Node) []ast.Decl {
	decls := make([]ast.Decl, len(tree))

	for i, node := range tree {
		if node.Type() != parser.NodeCall {
			panic("expected call node in root scope!")
		}

		decls[i] = evalDeclNode(node.(*parser.CallNode))
	}

	return decls
}

func evalDeclNode(node *parser.CallNode) ast.Decl {
	// Let's just assume that all top-level functions called will be "def"
	if node.Callee.Type() != parser.NodeIdent {
		panic("expecting call to identifier (i.e. def, defconst, etc.)")
	}

	callee := node.Callee.(*parser.IdentNode)
	switch callee.Ident {
	case "def":
		return evalDef(node)
	}

	return nil
}

func evalDef(node *parser.CallNode) ast.Decl {
	if len(node.Args) < 2 {
		panic(fmt.Sprintf("expecting expression to be assigned to variable: %q", node.Args[0]))
	}

	val := EvalExpr(node.Args[1])
	fn, ok := val.(*ast.FuncLit)

	ident := ast.NewIdent(goify(node.Args[0].(*parser.IdentNode).Ident, true))

	if ok {
		if ident.Name != "Main" {
			return makeFuncDeclFromFuncLit(ident, fn)
		} else {
			ident.Name = "main"
			return mainable(makeFuncDeclFromFuncLit(ident, fn))
		}
	} else {
		return makeGeneralDecl(token.VAR, []ast.Spec{makeValueSpec(ident, val)})
	}
}

func isNSDecl(node parser.Node) bool {
	if node.Type() != parser.NodeCall {
		return false
	}

	call := node.(*parser.CallNode)
	if call.Callee.(*parser.IdentNode).Ident != "ns" {
		return false
	}

	if len(call.Args) < 1 {
		return false
	}

	return true
}

func getNamespace(node *parser.CallNode) (*ast.Ident, ast.Decl) {
	return getPackageName(node), getImports(node)
}

func getPackageName(node *parser.CallNode) *ast.Ident {
	if node.Args[0].Type() != parser.NodeIdent {
		panic("ns package name is not an identifier!")
	}

	return ast.NewIdent(node.Args[0].(*parser.IdentNode).Ident)
}

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

func checkNSArgs(node *parser.CallNode) bool {
	if node.Callee.Type() != parser.NodeIdent {
		return false
	}

	if callee := node.Callee.(*parser.IdentNode); callee.Ident != "ns" {
		return false
	}

	return true
}

func checkIfArgs(node *parser.CallNode) bool {
	if node.Callee.Type() != parser.NodeIdent {
		return false
	}

	if callee := node.Callee.(*parser.IdentNode); callee.Ident != "if" {
		return false
	}

	if len(node.Args) < 2 {
		return false
	}

	return true
}

// Only need this to check if "def" is in
// an expression, which is illegal
func checkDefArgs(node *parser.CallNode) bool {
	if node.Callee.Type() != parser.NodeIdent {
		return false
	}

	if callee := node.Callee.(*parser.IdentNode); callee.Ident != "def" {
		return false
	}

	return true
}

func checkFunArgs(node *parser.CallNode) bool {
	// Need an identifier for it to be "fn"
	if node.Callee.Type() != parser.NodeIdent {
		return false
	}

	if callee := node.Callee.(*parser.IdentNode); callee.Ident != "fn" {
		return false
	}

	// Need argument list and at least one expression
	if len(node.Args) < 2 {
		return false
	}

	// Parameters should be a vector
	params := node.Args[0]
	if params.Type() != parser.NodeVector {
		return false
	}

	p := params.(*parser.VectorNode)
	for _, param := range p.Nodes {
		// TODO: change this in case of variable unpacking
		if param.Type() != parser.NodeIdent {
			return false
		}
	}

	return true
}

func checkLetArgs(node *parser.CallNode) bool {
	// Need an identifier for it to be "let"
	if node.Callee.Type() != parser.NodeIdent {
		return false
	}

	// Not a "let"
	if callee := node.Callee.(*parser.IdentNode); callee.Ident != "let" {
		return false
	}

	// Need _at least_ the bindings & one expression
	if len(node.Args) < 2 {
		return false
	}

	// Bindings should be a vector
	bindings := node.Args[0]
	if bindings.Type() != parser.NodeVector {
		return false
	}

	// There should be an even number of elements in the bindings
	b := bindings.(*parser.VectorNode)
	if len(b.Nodes)%2 != 0 {
		return false
	}

	// The bound identifiers, should be identifiers
	for i := 0; i < len(b.Nodes); i += 2 {
		if b.Nodes[i].Type() != parser.NodeIdent {
			return false
		}
	}

	return true
}
