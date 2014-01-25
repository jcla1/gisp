package generator

import (
	"../parser"
	"go/ast"
	"go/token"
)

func GenerateAST(tree []parser.Node) *ast.File {
	f := &ast.File{Name: makeIdent("main")}
	decls := make([]ast.Decl, 0, len(tree))

	if len(tree) < 1 {
		return f
	}

	if isNSDecl(tree[0]) {
		name, imports := getNamespace(tree[0].(*parser.CallNode))

		f.Name = name
		if imports != nil {
			decls = append(decls, imports)
		}

		tree = tree[1:]
	}

	f.Decls = decls
	return f
}

func isNSDecl(node parser.Node) bool {
	if node.Type() != parser.NodeCall {
		return false
	}

	call := node.(*parser.CallNode)
	if call.Callee.(*parser.IdentNode).Ident != "ns" {
		return false
	}

	if len(call.Args) < 1 {
		return false
	}

	return true
}

func getNamespace(node *parser.CallNode) (*ast.Ident, ast.Decl) {
	return getPackageName(node), getImports(node)
}

func getPackageName(node *parser.CallNode) *ast.Ident {
	if node.Args[0].Type() != parser.NodeIdent {
		panic("ns package name is not an identifier!")
	}

	return makeIdent(node.Args[0].(*parser.IdentNode).Ident)
}

func getImports(node *parser.CallNode) ast.Decl {
	if len(node.Args) < 2 {
		return nil
	}

	imports := node.Args[1:]
	specs := make([]ast.Spec, len(imports))

	for i, imp := range imports {
		if t := imp.Type(); t == parser.NodeVector {
			specs[i] = makeImportSpecFromVector(imp.(*parser.VectorNode))
		} else if t == parser.NodeString {
			path := makeBasicLit(token.STRING, imp.(*parser.StringNode).Value)
			specs[i] = makeImportSpec(path, nil)
		} else {
			panic("invalid import!")
		}
	}

	decl := makeGeneralDecl(token.IMPORT, specs)
	decl.Lparen = token.Pos(1) // Need this so we can have multiple imports
	return decl
}

func EvalExprs(nodes []parser.Node) []ast.Expr {
	out := make([]ast.Expr, len(nodes))

	for i, node := range nodes {
		out[i] = EvalExpr(node)
	}

	return out
}

func EvalExpr(node parser.Node) ast.Expr {
	switch t := node.Type(); t {
	case parser.NodeCall:
		node := node.(*parser.CallNode)
		return evalFunCall(node)
	case parser.NodeVector:
		node := node.(*parser.VectorNode)
		return makeVector(makeIdent("Any"), EvalExprs(node.Nodes))
	case parser.NodeNumber:
		node := node.(*parser.NumberNode)
		return makeBasicLit(node.NumberType, node.Value)
	case parser.NodeString:
		node := node.(*parser.StringNode)
		return makeBasicLit(token.STRING, node.Value)
	case parser.NodeIdent:
		node := node.(*parser.IdentNode)
		return makeIdent(node.Ident)
	default:
		println(t)
		panic("not implemented yet!")
	}
}

func evalFunCall(node *parser.CallNode) ast.Expr {
	switch {
	case checkLetArgs(node):
		return makeLetFun(node)
	case checkFunArgs(node):
		nodes := node.Args[0].(*parser.VectorNode).Nodes
		idents := make([]*parser.IdentNode, len(nodes))
		for i := 0; i < len(nodes); i++ {
			idents[i] = nodes[i].(*parser.IdentNode)
		}

		params := makeIdentSlice(idents)
		body := wrapExprsWithStmt(EvalExprs(node.Args[1:]))
		return makeFunLit(params, body)
	case checkDefArgs(node):
		panic("you can't have a def within an expression!")
	case checkNSArgs(node):
		panic("you can't define a namespace in an expression!")
	}

	callee := EvalExpr(node.Callee)
	args := EvalExprs(node.Args)

	return makeFunCall(callee, args)
}

func checkNSArgs(node *parser.CallNode) bool {
	if node.Callee.Type() != parser.NodeIdent {
		return false
	}

	if callee := node.Callee.(*parser.IdentNode); callee.Ident != "ns" {
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

func checkFunArgs(node *parser.CallNode) bool {
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

	// There should be an even number of elements in the bindings
	b := bindings.(*parser.VectorNode)
	if len(b.Nodes)%2 != 0 {
		return false
	}

	// The bound identifiers, should be identifiers
	for i := 0; i < len(b.Nodes); i += 2 {
		if b.Nodes[i].Type() != parser.NodeIdent {
			return false
		}
	}

	return true
}

func makeLetFun(node *parser.CallNode) ast.Expr {
	bindings := makeBindings(node.Args[0].(*parser.VectorNode))
	// TODO: clean this!
	return makeFunCall(makeFunLit([]*ast.Ident{}, append(bindings, wrapExprsWithStmt(EvalExprs(node.Args[1:]))...)), []ast.Expr{})
}

func makeBindings(bindings *parser.VectorNode) []ast.Stmt {
	vars := make([]*ast.Ident, len(bindings.Nodes)/2)
	for i := 0; i < len(bindings.Nodes); i += 2 {
		vars[i/2] = makeIdent(bindings.Nodes[i].(*parser.IdentNode).Ident)
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
		out[i] = &ast.ExprStmt{X: v}
	}
	return out
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
						Type: makeIdent("Any"),
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
				Type:  makeIdent("Any"),
				Names: args,
			},
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

func makeIdentSlice(nodes []*parser.IdentNode) []*ast.Ident {
	out := make([]*ast.Ident, len(nodes))
	for i, node := range nodes {
		out[i] = makeIdent(node.Ident)
	}
	return out
}

func makeIdent(name string) *ast.Ident {
	return ast.NewIdent(name)
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
	name := makeIdent(vect.Nodes[2].(*parser.IdentNode).Ident)

	return makeImportSpec(path, name)
}
