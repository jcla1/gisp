package generator

import (
	"../parser"
	"bytes"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

func makeIfStmtFun(node *parser.CallNode) ast.Expr {
    var otherwise ast.Stmt = nil
    if len(node.Args) > 2 {
        otherwise = makeReturnStmt(EvalExpr(node.Args[2]))
    }

    cond, body := EvalExpr(node.Args[0]), makeBlockStmt([]ast.Stmt{makeReturnStmt(EvalExpr(node.Args[1]))})

    return makeFunCall(makeFunLit([]*ast.Ident{}, []ast.Stmt{makeIfStmt(cond, body, otherwise)}), []ast.Expr{})
}

func makeIfStmt(cond ast.Expr, body *ast.BlockStmt, otherwise ast.Stmt) *ast.IfStmt {
    return &ast.IfStmt{
        Cond: cond,
        Body: body,
        Else: otherwise,
    }
}

func makeLetFun(node *parser.CallNode) ast.Expr {
	bindings := makeBindings(node.Args[0].(*parser.VectorNode))
	// TODO: clean this!
	return makeFunCall(makeFunLit([]*ast.Ident{}, append(bindings, wrapExprsWithStmt(EvalExprs(node.Args[1:]))...)), []ast.Expr{})
}

func makeBindings(bindings *parser.VectorNode) []ast.Stmt {
	vars := make([]*ast.Ident, len(bindings.Nodes)/2)
	for i := 0; i < len(bindings.Nodes); i += 2 {
		vars[i/2] = ast.NewIdent(bindings.Nodes[i].(*parser.IdentNode).Ident)
	}

	vals := make([]ast.Expr, len(bindings.Nodes)/2)
	for i := 1; i < len(bindings.Nodes); i += 2 {
		vals[(i-1)/2] = EvalExpr(bindings.Nodes[i])
	}

	stmts := make([]ast.Stmt, len(vars))
	for i := 0; i < len(vars); i++ {
		stmts[i] = makeAssignStmt(vars[i], vals[i])
	}

	return stmts
}

func makeAssignStmt(name, val ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{name},
		Rhs: []ast.Expr{val},
		// TODO: check if following line can be omitted
		Tok: token.DEFINE,
	}
}

func wrapExprsWithStmt(exps []ast.Expr) []ast.Stmt {
	out := make([]ast.Stmt, len(exps))
	for i, v := range exps {
		out[i] = makeExprStmt(v)
	}
	return out
}

func makeExprStmt(exp ast.Expr) ast.Stmt {
	return &ast.ExprStmt{X: exp}
}

func makeFunCall(callee ast.Expr, args []ast.Expr) ast.Expr {
	return &ast.CallExpr{
		Fun:  callee,
		Args: args,
	}
}

func makeFunLit(args []*ast.Ident, body []ast.Stmt) *ast.FuncLit {
	node := &ast.FuncLit{
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Type: ast.NewIdent("Any"),
					},
				},
			},
		},
		Body: makeBlockStmt(returnLast(body)),
	}

	if len(args) > 0 {
		node.Type.Params = makeParameterList(args)
	}

	return node
}

func makeParameterList(args []*ast.Ident) *ast.FieldList {
	return &ast.FieldList{
		List: []*ast.Field{
			&ast.Field{
				Type:  ast.NewIdent("Any"),
				Names: args,
			},
		},
	}
}

func returnLast(stmts []ast.Stmt) []ast.Stmt {
	if len(stmts) < 1 {
		return stmts
	}

	stmts[len(stmts)-1] = makeReturnStmt(stmts[len(stmts)-1].(*ast.ExprStmt).X)

	return stmts
}

func makeReturnStmt(expr ast.Expr) ast.Stmt {
    return &ast.ReturnStmt{
        Results: []ast.Expr{expr},
    }
}

func makeIdentSlice(nodes []*parser.IdentNode) []*ast.Ident {
	out := make([]*ast.Ident, len(nodes))
	for i, node := range nodes {
		out[i] = ast.NewIdent(node.Ident)
	}
	return out
}

func makeVector(typ *ast.Ident, elements []ast.Expr) *ast.CompositeLit {
	return makeCompositeLit(&ast.ArrayType{Elt: typ}, elements)
}

func makeCompositeLit(typ ast.Expr, elements []ast.Expr) *ast.CompositeLit {
	return &ast.CompositeLit{
		Type: typ,
		Elts: elements,
	}
}

func makeBasicLit(kind token.Token, value string) *ast.BasicLit {
	return &ast.BasicLit{Kind: kind, Value: value}
}

func makeBlockStmt(statements []ast.Stmt) *ast.BlockStmt {
	return &ast.BlockStmt{List: statements}
}

func makeGeneralDecl(typ token.Token, specs []ast.Spec) *ast.GenDecl {
	return &ast.GenDecl{
		Tok:   typ,
		Specs: specs,
	}
}

func makeImportSpec(path *ast.BasicLit, name *ast.Ident) *ast.ImportSpec {
	spec := &ast.ImportSpec{Path: path}

	if name != nil {
		spec.Name = name
	}

	return spec
}

func makeImportSpecFromVector(vect *parser.VectorNode) *ast.ImportSpec {
	if len(vect.Nodes) < 3 {
		panic("invalid use of import!")
	}

	if vect.Nodes[0].Type() != parser.NodeString {
		panic("invalid use of import!")
	}

	pathString := vect.Nodes[0].(*parser.StringNode).Value
	path := makeBasicLit(token.STRING, pathString)

	if vect.Nodes[1].Type() != parser.NodeIdent || vect.Nodes[1].(*parser.IdentNode).Ident != ":as" {
		panic("invalid use of import! expecting: \":as\"!!!")
	}
	name := ast.NewIdent(vect.Nodes[2].(*parser.IdentNode).Ident)

	return makeImportSpec(path, name)
}

func mainable(fn *ast.FuncDecl) *ast.FuncDecl {
	fn.Type.Results = nil

	returnStmt := fn.Body.List[len(fn.Body.List)-1].(*ast.ReturnStmt)
	fn.Body.List[len(fn.Body.List)-1] = makeExprStmt(returnStmt.Results[0])

	return fn
}

func makeFuncDeclFromFuncLit(name *ast.Ident, f *ast.FuncLit) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: name,
		Type: f.Type,
		Body: f.Body,
	}
}

func makeValueSpec(name *ast.Ident, value ast.Expr) *ast.ValueSpec {
	return &ast.ValueSpec{
		Names:  []*ast.Ident{name},
		Values: []ast.Expr{value},
	}
}

func makeSelectorExpr(x ast.Expr, sel *ast.Ident) *ast.SelectorExpr {
	return &ast.SelectorExpr{
		X:   x,
		Sel: sel,
	}
}

var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

func goify(src string, capitalizeFirst bool) string {
	src = strings.Replace(src, "/", ".", -1)
	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx > 0 || capitalizeFirst {
			chunks[idx] = bytes.Title(val)
		}
	}
	return string(bytes.Join(chunks, nil))
}
