package generator

import (
	"go/ast"
	"go/token"
	"../parser"
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
		return evalFuncCall(node)
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
	case parser.NodeNil:
		return makeNil()
	default:
		println(t)
		panic("not implemented yet!")
	}
}

func evalFuncCall(node *parser.CallNode) ast.Expr {
	callee := EvalExpr(node.Callee)
	args := EvalExprs(node.Args)

	return makeFunCall(callee, args)
}

// func makeLitFunCall(body []Any) ast.Expr {
// 	return &ast.CallExpr{
// 		Fun:  makeFuncLit([]Any{}, body),
// 		Args: []ast.Expr{},
// 	}
// }

// func wrapExprsWithStmt(exps []ast.Expr) []ast.Stmt {
// 	out := make([]ast.Stmt, len(exps))
// 	for i, v := range exps {
// 		out[i] = &ast.ExprStmt{X: v}
// 	}
// 	return out
// }

func makeFunCall(callee ast.Expr, args []ast.Expr) ast.Expr {
	return &ast.CallExpr{
		Fun:  callee,
		Args: args,
	}
}

// func makeFuncLit(args, body []Any) *ast.FuncLit {
// 	node := &ast.FuncLit{
// 		Type: &ast.FuncType{

// 			Results: &ast.FieldList{
// 				List: []*ast.Field{
// 					&ast.Field{
// 						Type: makeIdent("Any"),
// 					},
// 				},
// 			},
// 		},
// 		Body: &ast.BlockStmt{
// 			List: returnLast(wrapExprsWithStmt(evalExprs(body))),
// 		},
// 	}

// 	if len(args) > 0 {
// 		node.Type.Params = makeParameterList(args)
// 	}

// 	return node
// }

// func makeParameterList(args []Any) *ast.FieldList {
// 	return &ast.FieldList{
// 		List: []*ast.Field{
// 			&ast.Field{
// 				Type:  makeIdent("Any"),
// 				Names: makeIdentSlice(args),
// 			},
// 		},
// 	}
// }

// func returnLast(stmts []ast.Stmt) []ast.Stmt {
// 	if len(stmts) < 1 {
// 		return stmts
// 	}

// 	stmts[len(stmts)-1] = &ast.ReturnStmt{
// 		Results: []ast.Expr{
// 			stmts[len(stmts)-1].(*ast.ExprStmt).X,
// 		},
// 	}
// 	return stmts
// }

// func makeIdentSlice(args []Any) []*ast.Ident {
// 	out := make([]*ast.Ident, len(args))
// 	for i, v := range args {
// 		out[i] = makeIdent(v.(astToken).Value)
// 	}
// 	return out
// }

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


func makeNil() *ast.Ident {
	return makeIdent("nil")
}