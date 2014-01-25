package generator

import (
	"../parser"
	"go/ast"
	"go/token"
)

// func generateAST(tree []parser.Node) *ast.File {
// 	return &ast.File{Name: makeIdent("main"), Decls: generateDeclarations(tree)}
// }

// func generateDeclarations(tree []parser.Node) []ast.Decl {
// 	decls := make([]ast.Decl, len(tree))

// 	for i, node := range tree {
// 		switch node.Type() {
// 		case parser.NodeCall:
// 			decls[i] = generateDeclaration(node)
// 		default:
// 			panic("unexpected behaviour!")
// 		}
// 	}

// 	return decls
// }

// func generateDeclaration(node parser.CallNode) ast.Decl {
// 	if node.Callee.(parser.IdentNode).Ident == "def" {
// 		return evalDef(node)
// 	}

// 	return nil
// }

// func evalDef(node parser.CallNode) ast.Decl {
// 	ident := node.Callee.(parser.IdentNode).Ident
// 	val := evalExpr(node.Args[0])

// 	return &ast.GenDecl{
// 		Tok: token.VAR,
// 		Specs: []ast.Spec{
// 			&ast.ValueSpec{
// 				Names:  []*ast.Ident{makeIdent(ident.Value)},
// 				Values: []ast.Expr{val},
// 			},
// 		},
// 	}
// }

func EvalExprs(nodes []parser.Node) []ast.Expr {
	out := make([]ast.Expr, len(nodes))

	for i, node := range nodes {
		out[i] = EvalExpr(node)
	}

	return out
}

func EvalExpr(node parser.Node) ast.Expr {
	switch t := node.Type(); t {
	case parser.NodeCall:
		node := node.(*parser.CallNode)
		return evalFunCall(node)
	case parser.NodeVector:
		node := node.(*parser.VectorNode)
		return makeVector(makeIdent("Any"), EvalExprs(node.Nodes))
	case parser.NodeNumber:
		node := node.(*parser.NumberNode)
		return makeBasicLit(node.NumberType, node.Value)
	case parser.NodeString:
		node := node.(*parser.StringNode)
		return makeBasicLit(token.STRING, node.Value)
	case parser.NodeIdent:
		node := node.(*parser.IdentNode)
		return makeIdent(node.Ident)
	default:
		println(t)
		panic("not implemented yet!")
	}
}

func evalFunCall(node *parser.CallNode) ast.Expr {
	switch {
	case checkLetArgs(node):
		return makeLetFun(node)
	}

	callee := EvalExpr(node.Callee)
	args := EvalExprs(node.Args)

	return makeFunCall(callee, args)
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

func makeLetFun(node *parser.CallNode) ast.Expr {
	bindings := makeBindings(node.Args[0].(*parser.VectorNode))
	// TODO: clean this!
	return makeFunCall(makeFunLit([]*ast.Ident{}, append(bindings, wrapExprsWithStmt(EvalExprs(node.Args[1:]))...)), []ast.Expr{})
}

func makeBindings(bindings *parser.VectorNode) []ast.Stmt {
	vars := make([]*ast.Ident, len(bindings.Nodes)/2)
	for i := 0; i < len(bindings.Nodes); i += 2 {
		vars[i/2] = makeIdent(bindings.Nodes[i].(*parser.IdentNode).Ident)
	}

	vals := make([]ast.Expr, len(bindings.Nodes)/2)
	for i := 1; i < len(bindings.Nodes); i += 2 {
		vals[(i-1)/2] = EvalExpr(bindings.Nodes[i])
	}

	stmts := make([]ast.Stmt, len(vars))
	for i := 0; i < len(vars); i++ {
		stmts[i] = makeAssignStmt(vars[i], vals[i])
	}

	return stmts
}

func makeAssignStmt(name, val ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{name},
		Rhs: []ast.Expr{val},
		// TODO: check if following line can be omitted
		Tok: token.DEFINE,
	}
}

func wrapExprsWithStmt(exps []ast.Expr) []ast.Stmt {
	out := make([]ast.Stmt, len(exps))
	for i, v := range exps {
		out[i] = &ast.ExprStmt{X: v}
	}
	return out
}

func makeFunCall(callee ast.Expr, args []ast.Expr) ast.Expr {
	return &ast.CallExpr{
		Fun:  callee,
		Args: args,
	}
}

func makeFunLit(args []*ast.Ident, body []ast.Stmt) *ast.FuncLit {
	node := &ast.FuncLit{
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Type: makeIdent("Any"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: returnLast(body),
		},
	}

	if len(args) > 0 {
		node.Type.Params = makeParameterList(args)
	}

	return node
}

func makeParameterList(args []*ast.Ident) *ast.FieldList {
	return &ast.FieldList{
		List: []*ast.Field{
			&ast.Field{
				Type:  makeIdent("Any"),
				Names: args,
			},
		},
	}
}

func returnLast(stmts []ast.Stmt) []ast.Stmt {
	if len(stmts) < 1 {
		return stmts
	}

	stmts[len(stmts)-1] = &ast.ReturnStmt{
		Results: []ast.Expr{
			stmts[len(stmts)-1].(*ast.ExprStmt).X,
		},
	}
	return stmts
}

func makeIdentSlice(nodes []*parser.IdentNode) []*ast.Ident {
	out := make([]*ast.Ident, len(nodes))
	for i, node := range nodes {
		out[i] = makeIdent(node.Ident)
	}
	return out
}

func makeIdent(name string) *ast.Ident {
	return ast.NewIdent(name)
}

func makeVector(typ *ast.Ident, elements []ast.Expr) *ast.CompositeLit {
	return makeCompositeLit(&ast.ArrayType{Elt: typ}, elements)
}

func makeCompositeLit(typ ast.Expr, elements []ast.Expr) *ast.CompositeLit {
	return &ast.CompositeLit{
		Type: typ,
		Elts: elements,
	}
}

func makeBasicLit(kind token.Token, value string) *ast.BasicLit {
	return &ast.BasicLit{Kind: kind, Value: value}
}

func makeBlockStmt(statements []ast.Stmt) *ast.BlockStmt {
	return &ast.BlockStmt{List: statements}
}
