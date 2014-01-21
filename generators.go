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
	if sexp[0].(astToken).Value == "def" {
		return evalDef(sexp)
	}

	return nil
}

func evalDef(sexp []Any) ast.Decl {
	ident := sexp[1].(astToken)
	val := evalExpr(sexp[2])

	return &ast.GenDecl{
		Tok: goToken.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names:  []*ast.Ident{makeIdent(ident.Value)},
				Values: []ast.Expr{val},
			},
		},
	}
}

func evalExprs(sexp []Any) []ast.Expr {
	out := make([]ast.Expr, len(sexp))

	for i, v := range sexp {
		out[i] = evalExpr(v)
	}

	return out
}

func evalExpr(sexp Any) ast.Expr {
	switch sexp := sexp.(type) {
	case []Any:
		return evalFuncCall(sexp)
	case Any:
		tok := sexp.(astToken)
		switch tok.Type {
		case "INT":
			return makeBasicLit(goToken.INT, tok.Value)
		case "FLOAT":
			return makeBasicLit(goToken.FLOAT, tok.Value)
		case "STRING":
			return makeBasicLit(goToken.STRING, tok.Value)
		case "IDENT":
			return makeIdent(tok.Value)
		default:
			panic("other datatype not implemented yet!")
		}
	default:
		panic("oops!")
	}
}

func evalFuncCall(sexp []Any) ast.Expr {
	tok, ok := sexp[0].(astToken)

	if tok.Value == "fn" && ok {
		return makeFuncLit(sexp[1].([]Any), sexp[2:])
	}

	return makeFunCall(evalExpr(sexp[0]), sexp[1:])
}

func makeLitFunCall(body []Any) ast.Expr {
	return &ast.CallExpr{
		Fun:  makeFuncLit([]Any{}, body),
		Args: []ast.Expr{},
	}
}

func wrapExprsWithStmt(exps []ast.Expr) []ast.Stmt {
	out := make([]ast.Stmt, len(exps))
	for i, v := range exps {
		out[i] = &ast.ExprStmt{X: v}
	}
	return out
}

func makeFunCall(name ast.Expr, args []Any) ast.Expr {
	return &ast.CallExpr{
		Fun:  name,
		Args: evalExprs(args),
	}
}

func makeFuncLit(args, body []Any) *ast.FuncLit {
	return &ast.FuncLit{
		Type: &ast.FuncType{
			Params: &ast.FieldList{

				List: []*ast.Field{
					&ast.Field{
						Type:  makeIdent("Any"),
						Names: makeIdentSlice(args),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Type: makeIdent("Any"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: returnLast(wrapExprsWithStmt(evalExprs(body))),
		},
	}
}

func returnLast(stmts []ast.Stmt) []ast.Stmt {
	if len(stmts) < 1 {
		return stmts
	}

	stmts[len(stmts)-1] = &ast.ReturnStmt{
		Results: []ast.Expr{
			stmts[len(stmts)-1].(*ast.ExprStmt).X,
		},
	}
	return stmts
}

func makeIdentSlice(args []Any) []*ast.Ident {
	out := make([]*ast.Ident, len(args))
	for i, v := range args {
		out[i] = makeIdent(v.(astToken).Value)
	}
	return out
}

func makeIdent(name string) *ast.Ident {
	return ast.NewIdent(name)
}

func makeVector(typ string, elements []Any) ast.CompositeLit {
    return makeCompositeLit(&ast.ArrayType{Elt: makeIdent(typ)}, evalExprs(elements))
}

func makeCompositeLit(typ ast.Expr, elements []ast.Expr) *ast.CompositeLit {
    return &ast.CompositeLit{
        Type: typ,
        Elts: elements,
    }
}

func makeBasicLit(kind goToken.Token, value string) *ast.BasicLit {
	return &ast.BasicLit{Kind: kind, Value: value}
}
