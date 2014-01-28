package generator

import (
	"../parser"
	"go/ast"
	"go/token"
)

var operatorMap = map[string]token.Token{
	// "=": token.EQL,
	"+": token.ADD,
	"-": token.SUB,
	"*": token.MUL,
	"/": token.QUO,
}

func isBinaryOperator(node *parser.CallNode) bool {
	if node.Callee.Type() != parser.NodeIdent {
		return false
	}

	_, ok := operatorMap[node.Callee.(*parser.IdentNode).Ident]

	if len(node.Args) < 2 && ok {
		panic("can't use binary operator with only one argument!")
	}

	return ok
}

func makeNAryBinaryExpr(node *parser.CallNode) *ast.BinaryExpr {
	op := operatorMap[node.Callee.(*parser.IdentNode).Ident]
	outer := makeBinaryExpr(op, EvalExpr(node.Args[0]), EvalExpr(node.Args[1]))

	for i := 2; i < len(node.Args); i++ {
		outer = makeBinaryExpr(op, outer, EvalExpr(node.Args[i]))
	}

	return outer
}

func makeBinaryExpr(op token.Token, x, y ast.Expr) *ast.BinaryExpr {
	return &ast.BinaryExpr{
		X:  x,
		Y:  y,
		Op: op,
	}
}

func makeUnaryExpr(op token.Token, x ast.Expr) *ast.UnaryExpr {
	return &ast.UnaryExpr{
		X:  x,
		Op: op,
	}
}
