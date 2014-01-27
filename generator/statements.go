package generator

import (
	"go/ast"
	"go/token"
)

func wrapExprsWithStmt(exps []ast.Expr) []ast.Stmt {
	out := make([]ast.Stmt, len(exps))
	for i, v := range exps {
		out[i] = makeExprStmt(v)
	}
	return out
}

func makeBlockStmt(statements []ast.Stmt) *ast.BlockStmt {
	return &ast.BlockStmt{List: statements}
}

func makeExprStmt(exp ast.Expr) ast.Stmt {
	return &ast.ExprStmt{X: exp}
}

func makeIfStmt(cond ast.Expr, body *ast.BlockStmt, otherwise ast.Stmt) *ast.IfStmt {
	return &ast.IfStmt{
		Cond: cond,
		Body: body,
		Else: otherwise,
	}
}

func makeAssignStmt(names, vals []ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: names,
		Rhs: vals,
		// TODO: check if following line can be omitted
		Tok: token.DEFINE,
	}
}
