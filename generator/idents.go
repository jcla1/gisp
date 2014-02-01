package generator

import (
	"../parser"
	"bytes"
	"go/ast"
	"regexp"
	"strings"
	"strconv"
)

func makeIdentSlice(nodes []*parser.IdentNode) []*ast.Ident {
	out := make([]*ast.Ident, len(nodes))
	for i, node := range nodes {
		out[i] = ast.NewIdent(node.Ident)
	}
	return out
}

func makeSelectorExpr(x ast.Expr, sel *ast.Ident) *ast.SelectorExpr {
	return &ast.SelectorExpr{
		X:   x,
		Sel: sel,
	}
}

func makeIdomaticSelector(src string) ast.Expr {
	strs := strings.Split(src, "/")
	var expr ast.Expr = makeIdomaticIdent(strs[0])

	for i := 1; i < len(strs); i++ {
		ido := CamelCase(strs[i], true)
		expr = makeSelectorExpr(expr, ast.NewIdent(ido))
	}

	return expr
}

func makeIdomaticIdent(src string) *ast.Ident {
	return ast.NewIdent(CamelCase(src, false))
}

var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

func CamelCase(src string, capit bool) string {
	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx > 0 || capit {
			chunks[idx] = bytes.Title(val)
		}
	}
	return string(bytes.Join(chunks, nil))
}

var gensyms = func() <-chan string {
	syms := make(chan string)
	go func() {
		i := 0
		for {
			syms <- "GEN_" + strconv.Itoa(i)
			i++
		}
	}()
	return syms
}()

func generateIdent() *ast.Ident {
	return ast.NewIdent(<-gensyms)
}
