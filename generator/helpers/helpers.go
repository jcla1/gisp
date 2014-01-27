package helpers

import (
	"go/ast"
)

func EmptyS() []ast.Stmt {
	return []ast.Stmt{}
}

func S(stmt ast.Stmt) []ast.Stmt {
	return []ast.Stmt{stmt}
}

func EmptyE() []ast.Expr {
	return []ast.Expr{}
}

func E(expr ast.Expr) []ast.Expr {
	return []ast.Expr{expr}
}

func EmptyI() []*ast.Ident {
	return []*ast.Ident{}
}

func I(ident *ast.Ident) []*ast.Ident {
	return []*ast.Ident{ident}
}
