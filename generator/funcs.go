package generator

import (
	"../parser"
	h "./helpers"
	"go/ast"
)

func evalFunCall(node *parser.CallNode) ast.Expr {
	switch {
	case isUnaryOperator(node):
		return makeUnaryExpr(unaryOperatorMap[node.Callee.(*parser.IdentNode).Ident], EvalExpr(node.Args[0]))
	case isBinaryOperator(node):
		return makeNAryBinaryExpr(node)
	case isComparisonOperator(node):
		return makeNAryComparisonExpr(node)
	case checkLetArgs(node):
		return makeLetFun(node)
	case checkIfArgs(node):
		return makeIfStmtFunc(node)
	case checkFuncArgs(node):
		argIdents := getArgIdentsFromVector(node.Args[0].(*parser.VectorNode))

		// TODO: In case of type annotations change the following
		returnField := []*ast.Field{makeField(nil, anyType)}
		results := makeFieldList(returnField)

		var params *ast.FieldList = nil

		if len(argIdents) != 0 {
			params = makeFieldList([]*ast.Field{makeField(argIdents, anyType)})
		}

		fnType := makeFuncType(results, params)
		body := makeFuncBody(EvalExprs(node.Args[1:]))

		return makeFuncLit(fnType, body)
	case checkDefArgs(node):
		panic("you can't have a def within an expression!")
	case checkNSArgs(node):
		panic("you can't define a namespace in an expression!")
	}

	callee := EvalExpr(node.Callee)
	if c, ok := callee.(*ast.Ident); ok {
		callee = makeIdomaticIdent(c.Name)
	}

	args := EvalExprs(node.Args)

	return makeFuncCall(callee, args)
}

func getArgIdentsFromVector(vect *parser.VectorNode) []*ast.Ident {
	args := vect.Nodes
	argIdents := make([]*ast.Ident, len(vect.Nodes))
	for i, arg := range args {
		argIdents[i] = makeIdomaticIdent(arg.(*parser.IdentNode).Ident)
	}

	return argIdents
}

func makeFuncBody(exprs []ast.Expr) *ast.BlockStmt {
	wrapped := wrapExprsWithStmt(exprs)
	wrapped[len(wrapped)-1] = makeReturnStmt(h.E(wrapped[len(wrapped)-1].(*ast.ExprStmt).X))
	return makeBlockStmt(wrapped)
}

func makeFuncLit(typ *ast.FuncType, body *ast.BlockStmt) *ast.FuncLit {
	return &ast.FuncLit{
		Type: typ,
		Body: body,
	}
}

func makeFuncType(results, params *ast.FieldList) *ast.FuncType {
	return &ast.FuncType{
		Results: results,
		Params:  params,
	}
}

func makeFieldList(list []*ast.Field) *ast.FieldList {
	return &ast.FieldList{
		List: list,
	}
}

func makeField(names []*ast.Ident, typ ast.Expr) *ast.Field {
	return &ast.Field{
		Names: names,
		Type:  typ,
	}
}

func makeReturnStmt(exprs []ast.Expr) ast.Stmt {
	return &ast.ReturnStmt{
		Results: exprs,
	}
}

func makeFuncCall(callee ast.Expr, args []ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{
		Fun:  callee,
		Args: args,
	}
}

// Fn type checks (let, fn, def, ns, etc.)

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

func checkFuncArgs(node *parser.CallNode) bool {
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

	// The bindings should be also vectors
	b := bindings.(*parser.VectorNode)
	for _, bind := range b.Nodes {
		if _, ok := bind.(*parser.VectorNode); !ok {
			return false
		}
	}

	// The bound identifiers, should be identifiers
	for _, bind := range b.Nodes {
		bindingVect := bind.(*parser.VectorNode)
		if bindingVect.Nodes[0].Type() != parser.NodeIdent {
			return false
		}
	}

	return true
}
