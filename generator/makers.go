package generator

import (
	"../parser"
	h "./helpers"
	"go/ast"
	"go/token"
)

func makeIfStmtFunc(node *parser.CallNode) ast.Expr {
	var elseBody ast.Stmt
	if len(node.Args) > 2 {
		elseBody = makeBlockStmt(h.S(makeReturnStmt(h.E(EvalExpr(node.Args[2])))))
	} else {
		elseBody = makeBlockStmt(h.S(makeReturnStmt(h.E(ast.NewIdent("nil")))))
	}

	cond := EvalExpr(node.Args[0])
	ifBody := makeBlockStmt(h.S(makeReturnStmt(h.E(EvalExpr(node.Args[1])))))

	ifStmt := makeIfStmt(cond, ifBody, elseBody)
	fnBody := makeBlockStmt(h.S(ifStmt))

	returnList := makeFieldList([]*ast.Field{makeField(nil, anyType)})
	fnType := makeFuncType(returnList, nil)

	fn := makeFuncLit(fnType, fnBody)

	return makeFuncCall(fn, h.EmptyE())
}

func makeLetFun(node *parser.CallNode) ast.Expr {
	bindings := makeBindings(node.Args[0].(*parser.VectorNode))

	body := append(h.S(bindings), wrapExprsWithStmt(EvalExprs(node.Args[1:]))...)
	body[len(body)-1] = makeReturnStmt(h.E(body[len(body)-1].(*ast.ExprStmt).X))

	fieldList := makeFieldList([]*ast.Field{makeField(nil, anyType)})
	typ := makeFuncType(fieldList, nil)
	fn := makeFuncLit(typ, makeBlockStmt(body))

	return makeFuncCall(fn, h.EmptyE())
	// return makeFuncCall(makeFuncLit(h.EmptyI(), append(h.S(bindings), wrapExprsWithStmt(EvalExprs(node.Args[1:]))...)), h.EmptyE())
}

func makeBindings(bindings *parser.VectorNode) ast.Stmt {
	vars := make([]ast.Expr, len(bindings.Nodes))
	for i, bind := range bindings.Nodes {
		b := bind.(*parser.VectorNode)
		vars[i] = makeIdomaticSelector(b.Nodes[0].(*parser.IdentNode).Ident)
	}

	vals := make([]ast.Expr, len(bindings.Nodes))
	for i, bind := range bindings.Nodes {
		b := bind.(*parser.VectorNode)
		vals[i] = EvalExpr(b.Nodes[1])
	}

	return makeAssignStmt(vars, vals)
}

func mainable(fn *ast.FuncLit) {
	fn.Type.Results = nil

	returnStmt := fn.Body.List[len(fn.Body.List)-1].(*ast.ReturnStmt)
	fn.Body.List[len(fn.Body.List)-1] = makeExprStmt(returnStmt.Results[0])

	// return fn
}

// func makeTypeAssertFromArgList(expr ast.Expr, args []parser.Node) *ast.TypeAssertExpr {

// 	argList

// 	// currently we only support HOF that were written in Gisp
// 	returnList := makeFieldList([]*ast.Field{makeField(nil, anyType)})
// 	fnType := makeFuncType(returnList, argList)

// 	return makeTypeAssertion(expr, fnType)
// }

//////////////////////////////////
// Checked makers from here on! //
//////////////////////////////////

func makeValueSpec(names []*ast.Ident, values []ast.Expr) *ast.ValueSpec {
	return &ast.ValueSpec{
		Names:  names,
		Values: values,
	}
}

func makeFunDeclFromFuncLit(name *ast.Ident, f *ast.FuncLit) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: name,
		Type: f.Type,
		Body: f.Body,
	}
}

func makeGeneralDecl(typ token.Token, specs []ast.Spec) *ast.GenDecl {
	return &ast.GenDecl{
		Tok:   typ,
		Specs: specs,
	}
}

func makeTypeAssertion(expr, typ ast.Expr) *ast.TypeAssertExpr {
	return &ast.TypeAssertExpr{
		X:    expr,
		Type: typ,
	}
}

func makeEllipsis(typ ast.Expr) *ast.Ellipsis {
	return &ast.Ellipsis{
		Elt: typ,
	}
}
