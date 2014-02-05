package generator

import (
	"../parser"
	h "./helpers"
	"go/ast"
	"go/token"
)

func evalFuncCall(node *parser.CallNode) ast.Expr {
	switch {
	case isUnaryOperator(node):
		return makeUnaryExpr(unaryOperatorMap[node.Callee.(*parser.IdentNode).Ident], EvalExpr(node.Args[0]))

	case isCallableOperator(node):
		return makeNAryCallableExpr(node)

	case isLogicOperator(node):
		return makeNAryLogicExpr(node)

	case isLoop(node):
		return makeLoop(node)

	case isRecur(node):
		return makeRecur(node)

	case isAssert(node):
		return makeAssert(node)

	case isCoreFunc(node):
		return makeCoreCall(node)

	case checkLetArgs(node):
		return makeLetFun(node)

	case checkIfArgs(node):
		return makeIfStmtFunc(node)

	case checkFuncArgs(node):
		// TODO: In case of type annotations change the following
		returnField := []*ast.Field{makeField(nil, anyType)}
		results := makeFieldList(returnField)

		argIdents, ellipsis := getArgIdentsFromVector(node.Args[0].(*parser.VectorNode))
		params := make([]*ast.Field, 0, len(argIdents))

		if len(argIdents) != 0 {
			params = append(params, makeField(argIdents, anyType))
		}

		if ellipsis != nil {
			params = append(params, makeField(h.I(ellipsis), makeEllipsis(anyType)))
		}

		fnType := makeFuncType(results, makeFieldList(params))
		body := makeFuncBody(EvalExprs(node.Args[1:]))

		return makeFuncLit(fnType, body)

	case checkDefArgs(node):
		panic("you can't have a def within an expression!")

	case checkNSArgs(node):
		panic("you can't define a namespace in an expression!")
	}

	callee := EvalExpr(node.Callee)
	if c, ok := callee.(*ast.Ident); ok {
		callee = makeIdomaticIdent(c.Name)
	}

	args := EvalExprs(node.Args)

	return makeFuncCall(callee, args)
}

func getArgIdentsFromVector(vect *parser.VectorNode) ([]*ast.Ident, *ast.Ident) {
	args := vect.Nodes
	argIdents := make([]*ast.Ident, 0, len(vect.Nodes))

	var ident string
	var ellipsis *ast.Ident

	for i := 0; i < len(args); i++ {
		ident = args[i].(*parser.IdentNode).Ident

		if ident == "&" {
			ellipsis = makeIdomaticIdent(args[i+1].(*parser.IdentNode).Ident)
			break
		}

		argIdents = append(argIdents, makeIdomaticIdent(ident))
	}

	return argIdents, ellipsis
}

func makeFuncBody(exprs []ast.Expr) *ast.BlockStmt {
	wrapped := wrapExprsWithStmt(exprs)
	wrapped[len(wrapped)-1] = makeReturnStmt(h.E(wrapped[len(wrapped)-1].(*ast.ExprStmt).X))
	return makeBlockStmt(wrapped)
}

func makeFuncLit(typ *ast.FuncType, body *ast.BlockStmt) *ast.FuncLit {
	return &ast.FuncLit{
		Type: typ,
		Body: body,
	}
}

func makeFuncType(results, params *ast.FieldList) *ast.FuncType {
	return &ast.FuncType{
		Results: results,
		Params:  params,
	}
}

func makeFieldList(list []*ast.Field) *ast.FieldList {
	return &ast.FieldList{
		List: list,
	}
}

func makeField(names []*ast.Ident, typ ast.Expr) *ast.Field {
	return &ast.Field{
		Names: names,
		Type:  typ,
	}
}

func makeReturnStmt(exprs []ast.Expr) ast.Stmt {
	return &ast.ReturnStmt{
		Results: exprs,
	}
}

func makeFuncCall(callee ast.Expr, args []ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{
		Fun:  callee,
		Args: args,
	}
}

// Fn type checks (let, fn, def, ns, etc.)

func checkIfArgs(node *parser.CallNode) bool {
	if node.Callee.Type() != parser.NodeIdent {
		return false
	}

	if callee := node.Callee.(*parser.IdentNode); callee.Ident != "if" {
		return false
	}

	if len(node.Args) < 2 {
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

func checkFuncArgs(node *parser.CallNode) bool {
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

	// The bindings should be also vectors
	b := bindings.(*parser.VectorNode)
	for _, bind := range b.Nodes {
		if _, ok := bind.(*parser.VectorNode); !ok {
			return false
		}
	}

	// The bound identifiers, should be identifiers
	for _, bind := range b.Nodes {
		bindingVect := bind.(*parser.VectorNode)
		if bindingVect.Nodes[0].Type() != parser.NodeIdent {
			return false
		}
	}

	return true
}

func isLoop(node *parser.CallNode) bool {
	// Need an identifier for it to be "loop"
	if node.Callee.Type() != parser.NodeIdent {
		return false
	}

	// Not a "loop"
	if callee := node.Callee.(*parser.IdentNode); callee.Ident != "loop" {
		return false
	}

	// Bindings should be a vector
	bindings := node.Args[0]
	if bindings.Type() != parser.NodeVector {
		return false
	}

	// The bindings should be also vectors
	b := bindings.(*parser.VectorNode)
	for _, bind := range b.Nodes {
		if _, ok := bind.(*parser.VectorNode); !ok {
			return false
		}
	}

	// The bound identifiers, should be identifiers
	for _, bind := range b.Nodes {
		bindingVect := bind.(*parser.VectorNode)
		if bindingVect.Nodes[0].Type() != parser.NodeIdent {
			return false
		}
	}

	if !searchForRecur(node.Args[1:]) {
		panic("no recur found in loop!")
	}

	return true
}

func isRecur(node *parser.CallNode) bool {
	// Need an identifier for it to be "loop"
	if node.Callee.Type() != parser.NodeIdent {
		return false
	}

	// Not a "loop"
	if callee := node.Callee.(*parser.IdentNode); callee.Ident != "recur" {
		return false
	}

	// Bindings should be a vector
	bindings := node.Args[0]
	if bindings.Type() != parser.NodeVector {
		return false
	}

	// The bindings should be also vectors
	b := bindings.(*parser.VectorNode)
	for _, bind := range b.Nodes {
		if _, ok := bind.(*parser.VectorNode); !ok {
			return false
		}
	}

	// The bound identifiers, should be identifiers
	for _, bind := range b.Nodes {
		bindingVect := bind.(*parser.VectorNode)
		if bindingVect.Nodes[0].Type() != parser.NodeIdent {
			return false
		}
	}

	return true
}

func searchForRecur(nodes []parser.Node) bool {
	for _, node := range nodes {
		if node.Type() == parser.NodeCall {
			n := node.(*parser.CallNode)
			if ident, ok := n.Callee.(*parser.IdentNode); ok && ident.Ident == "recur" {
				return true
			} else if searchForRecur(n.Args) {
				return true
			}
		}
	}

	return false
}

func addNewValuesToBindings(bindingsVector *parser.VectorNode, vals []parser.Node) *parser.VectorNode {
	for i, _ := range bindingsVector.Nodes {
		bind := bindingsVector.Nodes[i].(*parser.VectorNode).Nodes
		bind[len(bind)-1] = vals[i]
	}

	return bindingsVector
}

func addRecurLabelAndBindings(label *parser.IdentNode, bindingsVector *parser.VectorNode, nodes []parser.Node) {
	for _, node := range nodes {
		if node.Type() == parser.NodeCall {
			n := node.(*parser.CallNode)
			if ident, ok := n.Callee.(*parser.IdentNode); ok && ident.Ident == "recur" {
				newValues := make([]parser.Node, len(n.Args))
				copy(newValues, n.Args)

				n.Args = make([]parser.Node, 2)
				n.Args[0] = addNewValuesToBindings(bindingsVector.Copy().(*parser.VectorNode), newValues)
				n.Args[1] = label
			} else {
				addRecurLabelAndBindings(label, bindingsVector, n.Args)
			}
		}
	}
}

func makeLoop(node *parser.CallNode) *ast.CallExpr {
	returnIdent := generateIdent()
	loopIdent := generateIdent()

	fnBody := h.EmptyS()

	bindingsVector := node.Args[0].(*parser.VectorNode)

	addRecurLabelAndBindings(parser.NewIdentNode(loopIdent.String()), bindingsVector.Copy().(*parser.VectorNode), node.Args[1:])

	bindings := makeBindings(bindingsVector, token.DEFINE)
	returnIdentValueSpec := makeValueSpec(h.I(returnIdent), nil, anyType)
	returnIdentDecl := makeDeclStmt(makeGeneralDecl(token.VAR, []ast.Spec{returnIdentValueSpec}))

	fnBody = append(fnBody, bindings...)
	fnBody = append(fnBody, returnIdentDecl)

	init := makeAssignStmt(h.E(loopIdent), h.E(ast.NewIdent("true")), token.DEFINE)
	forBody := h.EmptyS()

	forBody = append(forBody, makeAssignStmt(h.E(loopIdent), h.E(ast.NewIdent("false")), token.ASSIGN))
	forBody = append(forBody, wrapExprsWithStmt(EvalExprs(node.Args[1:len(node.Args)-1]))...)
	forBody = append(forBody, makeAssignStmt(h.E(returnIdent), h.E(EvalExpr(node.Args[len(node.Args)-1])), token.ASSIGN))

	forStmt := makeForStmt(init, nil, loopIdent, makeBlockStmt(forBody))

	fnBody = append(fnBody, forStmt)
	fnBody = append(fnBody, makeReturnStmt(h.E(returnIdent)))

	results := makeFieldList([]*ast.Field{makeField(nil, anyType)})
	fnType := makeFuncType(results, nil)
	fn := makeFuncLit(fnType, makeBlockStmt(fnBody))

	return makeFuncCall(fn, h.EmptyE())
}

func makeRecur(node *parser.CallNode) *ast.CallExpr {
	bindings := makeBindings(node.Args[0].(*parser.VectorNode), token.ASSIGN)
	loopUpdate := makeAssignStmt(h.E(EvalExpr(node.Args[1])), h.E(ast.NewIdent("true")), token.ASSIGN)

	body := append(h.EmptyS(), bindings...)
	body = append(body, loopUpdate, makeReturnStmt(h.E(ast.NewIdent("nil"))))

	resultType := makeFieldList([]*ast.Field{makeField(nil, anyType)})
	fnType := makeFuncType(resultType, nil)
	fn := makeFuncLit(fnType, makeBlockStmt(body))
	return makeFuncCall(fn, h.EmptyE())
}

func isAssert(node *parser.CallNode) bool {
	// Need an identifier for it to be "assert"
	if node.Callee.Type() != parser.NodeIdent {
		return false
	}

	// Not a "loop"
	if callee := node.Callee.(*parser.IdentNode); callee.Ident != "assert" {
		return false
	}

	if len(node.Args) != 2 {
		panic("assert needs 2 arguments")
	}

	if _, ok := node.Args[0].(*parser.IdentNode); !ok {
		panic("assert's first argument needs to be a type")
	}

	return true
}

func makeAssert(node *parser.CallNode) *ast.TypeAssertExpr {
	return makeTypeAssertion(EvalExpr(node.Args[1]), ast.NewIdent(node.Args[0].(*parser.IdentNode).Ident))
}

var coreFuncs = []string{"get"}

func isCoreFunc(node *parser.CallNode) bool {
	// Need an identifier for it to be a func
	if node.Callee.Type() != parser.NodeIdent {
		return false
	}

	ident := node.Callee.(*parser.IdentNode).Ident

	for _, v := range coreFuncs {
		if v == ident {
			return true
		}
	}

	return false
}

// TODO: just a quick and dirty implementation
func makeCoreCall(node *parser.CallNode) ast.Expr {
	ident := node.Callee.(*parser.IdentNode).Ident
	node.Callee.(*parser.IdentNode).Ident = "core/" + ident
	return evalFuncCall(node)
}