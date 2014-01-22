package parser

import (
	"../lexer"
	"go/token"
)

type Node interface {
	Type() NodeType
	Position() Pos
	// String() string
}

type Pos int

func (p Pos) Position() Pos {
	return p
}

type NodeType int

func (t NodeType) Type() NodeType {
	return t
}

type IdentNode struct {
	Pos
	NodeType
	Ident string
}

type StringNode struct {
	Pos
	NodeType
	Value string
}

type NumberNode struct {
	Pos
	NodeType
	Value      string
	NumberType token.Token
}

type VectorNode struct {
	Pos
	NodeType
	Nodes []Node
}

type CallNode struct {
	Pos
	NodeType
	Callee Node
	Args   []Node
}

func ParseFromString(name, program string) []Node {
	return Parse(lexer.Lex(name, program))
}

func Parse(l *lexer.Lexer) []Node {
	return parser(l, make([]Node, 0))
}

func parser(l *lexer.Lexer, tree []Node) []Node {

	return tree
}
