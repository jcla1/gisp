package generator

import (
	"../parser"
	"go/ast"
	"go/token"
)

var (
	comparisonOperatorMap = map[string]token.Token{
		">": token.GTR,
		"<": token.LSS,
		"=": token.EQL,
	}

	binaryOperatorMap = map[string]token.Token{
		"+": token.ADD,
		"-": token.SUB,
		"*": token.MUL,
		"/": token.QUO,
        "and": token.LAND,
        "or": token.LOR,
	}

	unaryOperatorMap = map[string]token.Token{
		"!": token.NOT,
	}
)

func isComparisonOperator(node *parser.CallNode) bool {
	if node.Callee.Type() != parser.NodeIdent {
		return false
	}

	_, ok := comparisonOperatorMap[node.Callee.(*parser.IdentNode).Ident]

	if len(node.Args) < 2 && ok {
		panic("can't use comparison operator with only one argument")
	}

	return ok
}

func makeNAryComparisonExpr(node *parser.CallNode) *ast.BinaryExpr {
	op := comparisonOperatorMap[node.Callee.(*parser.IdentNode).Ident]

	return makeBinaryExpr(op, EvalExpr(node.Args[0]), EvalExpr(node.Args[1]))
}

func chainComparisons(nodes []parser.Node) *ast.BinaryExpr {
    return nil
}

func isBinaryOperator(node *parser.CallNode) bool {
	if node.Callee.Type() != parser.NodeIdent {
		return false
	}

	_, ok := binaryOperatorMap[node.Callee.(*parser.IdentNode).Ident]

	if len(node.Args) < 2 && ok {
		panic("can't use binary operator with only one argument!")
	}

	return ok
}

func makeNAryBinaryExpr(node *parser.CallNode) *ast.BinaryExpr {
	op := binaryOperatorMap[node.Callee.(*parser.IdentNode).Ident]
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

func isUnaryOperator(node *parser.CallNode) bool {
	if node.Callee.Type() != parser.NodeIdent {
		return false
	}

	_, ok := unaryOperatorMap[node.Callee.(*parser.IdentNode).Ident]

	if len(node.Args) != 1 && ok {
		panic("unary expression takes, exactly, one argument!")
	}

	return ok
}

func makeUnaryExpr(op token.Token, x ast.Expr) *ast.UnaryExpr {
	return &ast.UnaryExpr{
		X:  x,
		Op: op,
	}
}
