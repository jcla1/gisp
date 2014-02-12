package generator

import (
	"github.com/jcla1/gisp/parser"
	"go/ast"
	"go/token"
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
		return evalFuncCall(node)

	case parser.NodeVector:
		node := node.(*parser.VectorNode)
		return makeVector(anyType, EvalExprs(node.Nodes))

	case parser.NodeNumber:
		node := node.(*parser.NumberNode)
		return makeBasicLit(node.NumberType, node.Value)

	case parser.NodeString:
		node := node.(*parser.StringNode)
		return makeBasicLit(token.STRING, node.Value)

	case parser.NodeIdent:
		node := node.(*parser.IdentNode)
		return makeIdomaticSelector(node.Ident)

	default:
		println(t)
		panic("not implemented yet!")
	}
}
