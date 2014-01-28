package helpers

import (
	"go/ast"
)

func EmptyS() []ast.Stmt {
	return []ast.Stmt{}
}

func S(stmts ...ast.Stmt) []ast.Stmt {
	return stmts
}

func EmptyE() []ast.Expr {
	return []ast.Expr{}
}

func E(exprs ...ast.Expr) []ast.Expr {
	return exprs
}

func EmptyI() []*ast.Ident {
	return []*ast.Ident{}
}

func I(ident *ast.Ident) []*ast.Ident {
	return []*ast.Ident{ident}
}
