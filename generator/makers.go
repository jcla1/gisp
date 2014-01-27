package generator

import (
	"../parser"
	h "./helpers"
	"go/ast"
	"go/token"
)

// func makeIfStmtFun(node *parser.CallNode) ast.Expr {
// 	var otherwise ast.Stmt = nil
// 	if len(node.Args) > 2 {
// 		otherwise = makeReturnStmt(h.E(EvalExpr(node.Args[2])))
// 	}

// 	cond, body := EvalExpr(node.Args[0]), makeBlockStmt(h.S(makeReturnStmt(h.E(EvalExpr(node.Args[1])))))

// 	return makeFuncCall(makeFunLitNoArgsSingleStmt(makeIfStmt(cond, body, otherwise)), h.EmptyE())
// }

func makeLetFun(node *parser.CallNode) ast.Expr {
	bindings := makeBindings(node.Args[0].(*parser.VectorNode))

    body := append(h.S(bindings), wrapExprsWithStmt(EvalExprs(node.Args[1:]))...)
    body[len(body)-1] = makeReturnStmt(h.E(body[len(body)-1].(*ast.ExprStmt).X))

    fieldList := makeFieldList([]*ast.Field{makeField(nil, ast.NewIdent("Any"))})
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
