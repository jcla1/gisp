package main

import (
	"bufio"
	"bytes"
	"fmt"
	// "go/ast"
	"go/printer"
	goToken "go/token"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Any interface{}
type Symbol string

func parse(l *lexer, p []Any) []Any {

	for {
		t := l.nextToken()
		if t.typ == _EOF {
			break
		} else if t.typ == _INVALID {
			panic("syntax error")
		}

		if t.typ == _LPAREN {
			p = append(p, parse(l, []Any{}))
			return parse(l, p)
		} else if t.typ == _RPAREN {
			return p
		} else {
			var v astToken
			v.Value = t.val
			switch t.typ {
			// case _UNQUOTESPLICE:
			// 	nextExp := parse(l, []Any{})
			// 	return append(append(p, []Any{Symbol("unquote-splice"), nextExp[0]}), nextExp[1:]...)
			// case _UNQUOTE:
			// 	nextExp := parse(l, []Any{})
			// 	return append(append(p, []Any{Symbol("unquote"), nextExp[0]}), nextExp[1:]...)
			// case _QUASIQUOTE:
			// 	nextExp := parse(l, []Any{})
			// 	return append(append(p, []Any{Symbol("quasiquote"), nextExp[0]}), nextExp[1:]...)
			// case _QUOTE:
			// 	nextExp := parse(l, []Any{})
			// 	return append(append(p, []Any{Symbol("quote"), nextExp[0]}), nextExp[1:]...)
			case _INT:
				v.Type = "INT"
			case _FLOAT:
				v.Type = "FLOAT"
			case _STRING:
				v.Type = "STRING"
			case _SYMBOL:
				v.Type = "IDENT"
			}
			return parse(l, append(p, v))
		}
	}

	return p
}

func args(filename string) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	l := lex(string(b) + "\n")
	p := parse(l, []Any{})
	// fmt.Printf("%#v\n\n\n", p)

	a := generateAST(p)

	fset := goToken.NewFileSet()
	// ast.Print(fset, a)

	var buf bytes.Buffer
	printer.Fprint(&buf, fset, a)
	fmt.Printf("%s\n", buf.String())
}

func main() {
	if len(os.Args) > 1 {
		args(os.Args[1])
		return
	}

	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">> ")
		line, _, _ := r.ReadLine()

		l := lex(string(line) + "\n")
		p := parse(l, []Any{})

		a := generateAST(p)
		fset := goToken.NewFileSet()
		var buf bytes.Buffer
		printer.Fprint(&buf, fset, a)
		fmt.Printf("%s\n", buf.String())
	}
}
