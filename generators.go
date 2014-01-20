package main

import (
	"go/ast"
	goToken "go/token"
)

func generateAST(parsed []Any) *ast.File {
	return &ast.File{Name: makeIdent("main"), Decls: generateDeclarations(parsed)}
}

func generateDeclarations(parsed []Any) []ast.Decl {
    decls := make([]ast.Decl, len(parsed))

    for i, sexp := range parsed {
        decls[i] = &ast.GenDecl{
            Tok: goToken.VAR,
            Specs: []ast.Spec{
                &ast.ValueSpec{
                    Names:  []*ast.Ident{makeIdent(string(sexp.([]Any)[0].(Symbol)))},
                    Values: []ast.Expr{makeBasicLit(goToken.INT, "10")},
                },
            },
        }
    }

    return decls
}

func makeIdent(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

func makeBasicLit(kind goToken.Token, value string) *ast.BasicLit {
	return &ast.BasicLit{Kind: kind, Value: value}
}
