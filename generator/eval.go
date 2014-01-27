package generator

import (
	"../parser"
	"go/ast"
	"go/token"
	"strings"
)

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
		return makeVector(ast.NewIdent("Any"), EvalExprs(node.Nodes))
	case parser.NodeNumber:
		node := node.(*parser.NumberNode)
		return makeBasicLit(node.NumberType, node.Value)
	case parser.NodeString:
		node := node.(*parser.StringNode)
		return makeBasicLit(token.STRING, node.Value)
	case parser.NodeIdent:
		node := node.(*parser.IdentNode)

		if strings.Contains(node.Ident, "/") {
			parts := strings.Split(node.Ident, "/")
			outerSelector := makeSelectorExpr(ast.NewIdent(parts[0]), ast.NewIdent(goify(parts[1], true)))

			for i := 2; i < len(parts); i++ {
				outerSelector = makeSelectorExpr(outerSelector, ast.NewIdent(goify(parts[i], true)))
			}

			return outerSelector
		}

		return ast.NewIdent(goify(node.Ident, false))
	default:
		println(t)
		panic("not implemented yet!")
	}
}

func evalFunCall(node *parser.CallNode) ast.Expr {
	switch {
	case checkLetArgs(node):
		return makeLetFun(node)
	case checkFunArgs(node):
		nodes := node.Args[0].(*parser.VectorNode).Nodes
		idents := make([]*parser.IdentNode, len(nodes))
		for i := 0; i < len(nodes); i++ {
			idents[i] = nodes[i].(*parser.IdentNode)
		}

		params := makeIdentSlice(idents)
		body := wrapExprsWithStmt(EvalExprs(node.Args[1:]))
		return makeFunLit(params, body)
	case checkDefArgs(node):
		panic("you can't have a def within an expression!")
	case checkNSArgs(node):
		panic("you can't define a namespace in an expression!")
	}

	callee := EvalExpr(node.Callee)
	if c, ok := callee.(*ast.Ident); ok {
		c.Name = goify(c.Name, true)
	}

	args := EvalExprs(node.Args)

	return makeFunCall(callee, args)
}
