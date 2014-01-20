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
		switch sexp := sexp.(type) {
		case []Any:
			decls[i] = generateDeclaration(sexp)
		default:
			panic("unexpected behaviour!")
		}
	}

	return decls
}

func generateDeclaration(sexp []Any) ast.Decl {
	// return &ast.GenDecl{
	//     Tok: goToken.CONST,
	//     Specs: []ast.Spec{
	//         &ast.ValueSpec{
	//             Names: []*ast.Ident{makeIdent(string(sexp[1].(Symbol)))},
	//             Values: []ast.Expr{generateExpression(sexp[2])},
	//         },
	//     },
	// }

	decl := &ast.GenDecl{
		Tok: goToken.CONST,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{makeIdent(string(sexp[1].(Symbol)))},
			},
		},
	}

	f := sexp[0]
	switch f := f.(type) {
	case Symbol:
		switch f {
		case "def":
			switch exp := sexp[2].(type) {
			case []Any:
				(decl.Specs[0].(*ast.ValueSpec)).Values = []ast.Expr{generateCallExpr(sexp[1].(Symbol), exp)}
			case string:

			}
		default:
			panic("unknown!")
		}
	default:
		panic("unknown behaviour!")
	}

	return decl
}

func generateCallExpr(f Symbol, exp []Any) *ast.CallExpr {
	return &ast.CallExpr{
		Fun:  makeIdent(string(f)),
		Args: []ast.Expr{},
	}
}

func generateExpression(val Any) ast.Expr {
	return makeBasicLit(goToken.INT, val.(string))
}

func makeIdent(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

func makeBasicLit(kind goToken.Token, value string) *ast.BasicLit {
	return &ast.BasicLit{Kind: kind, Value: value}
}

// &ast.GenDecl{
//             Tok: goToken.VAR,
//             Specs: []ast.Spec{
//                 &ast.ValueSpec{
//                     Names:  []*ast.Ident{makeIdent(string(sexp.([]Any)[0].(Symbol)))},
//                     Values: []ast.Expr{makeBasicLit(goToken.INT, "10")},
//                 },
//             },
//         }
