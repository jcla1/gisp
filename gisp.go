package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var DEBUG = false

type tokenType int

const (
	_INVALID tokenType = iota
	_EOF
	_INT
	_SYMBOL
	_LPAREN
	_RPAREN
	_STRING
	_FLOAT
	_BOOL
	_QUOTE
)

func (t tokenType) String() string {
	switch t {
	case _INVALID:
		return "INVALID TOKEN"
	case _EOF:
		return "EOF"
	case _INT:
		return "INT"
	case _SYMBOL:
		return "SYMBOL"
	case _LPAREN:
		return "LEFT_PAREN"
	case _RPAREN:
		return "RIGHT_PAREN"
	case _STRING:
		return "STRING"
	case _FLOAT:
		return "FLOAT"
	case _BOOL:
		return "BOOL"
	case _QUOTE:
		return "'"
	default:
		return "WTF!?"
	}
}

type token struct {
	typ tokenType // The type of this item.
	pos Pos       // The starting position, in bytes, of this item in the input string.
	val string    // The value of this item.
}

func (t token) String() string {
	return fmt.Sprintf("%s", t.val)
}

const eof = -1

type stateFn func(*lexer) stateFn
type Pos int

type lexer struct {
	name       string
	input      string
	state      stateFn
	pos        Pos
	start      Pos
	width      Pos
	lastPos    Pos
	tokens     chan token
	parenDepth int
}

func (l *lexer) run() {
	for l.state = lexWhitespace; l.state != nil; {
		l.state = l.state(l)
	}
}

func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) emit(t tokenType) {
	l.tokens <- token{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.lastPos], "\n")
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{_INVALID, l.start, fmt.Sprintf(format, args...)}
	return nil
}

func (l *lexer) nextToken() token {
	token := <-l.tokens
	l.lastPos = token.pos
	return token
}

// lexes an open parenthesis
func lexOpenParen(l *lexer) stateFn {
	debug("--> lexOpenParen")

	l.emit(_LPAREN)
	l.parenDepth++

	r := l.next()

	switch r {
	case ' ', '\t', '\n', '\r':
		return lexWhitespace
	case '\'':
		return lexQuote
	case '(':
		return lexOpenParen
	case ')':
		return lexCloseParen
	case ';':
		return lexComment
	case '#':
		return lexBool
	}

	if unicode.IsDigit(r) {
		return lexInt
	}

	return lexSymbol
}

func lexBool(l *lexer) stateFn {
	debug("--> lexBool")
	l.accept("tf")

	l.emit(_BOOL)

	r := l.next()

	switch r {
	case ' ', '\t', '\n':
		return lexWhitespace
	case ')':
		return lexCloseParen
	case ';':
		return lexComment
	}

	return l.errorf("unexpected tokens")
}

func lexQuote(l *lexer) stateFn {
	debug("--> lexQuote")
	l.acceptRun(" ")
	l.ignore()
	l.emit(_QUOTE)

	r := l.next()

	switch r {
	case '"':
		return lexString
	case '(':
		return lexOpenParen
	case ')':
		return lexCloseParen
	case '#':
		return lexBool
	case '\'':
		return lexQuote
	}

	if unicode.IsDigit(r) {
		return lexInt
	}

	return lexSymbol
}

func lexWhitespace(l *lexer) stateFn {
	debug("--> lexWhitespace")
	l.ignore()
	r := l.next()

	switch r {
	case ' ', '\t', '\n':
		return lexWhitespace
	case '\'':
		return lexQuote
	case '"':
		return lexString
	case '(':
		return lexOpenParen
	case ')':
		return lexCloseParen
	case ';':
		return lexComment
	case '#':
		return lexBool
	case eof:
		if l.parenDepth > 0 {
			return l.errorf("unclosed paren")
		}
		l.emit(_EOF)
		return nil
	}

	if unicode.IsDigit(r) {
		return lexInt
	}

	return lexSymbol
}

func lexString(l *lexer) stateFn {
	r := l.next()

	switch r {
	case '"':
		l.emit(_STRING)
		return lexWhitespace
	case '\\':
		// l.backup()
		// l.input = append(l.input[:l.pos], l.input[l.pos+1:])
		l.next()
		return lexString
	}

	return lexString
}

// lex an integer.  Once we're on an integer, the only valid characters are
// whitespace, close paren, a period to indicate we want a float, or more
// digits.  Everything else is crap.
func lexInt(l *lexer) stateFn {
	digits := "0123456789"
	l.acceptRun(digits)

	r := l.peek()

	switch r {
	case ' ', '\t', '\n':
		l.emit(_INT)
		l.next()
		return lexWhitespace
	case '.':
		l.next()
		return lexFloat
	case ')':
		l.emit(_INT)
		l.next()
		return lexCloseParen
	case ';':
		l.emit(_INT)
		l.next()
		return lexComment
	}

	return l.errorf("unexpected rune in lexInt: %c", r)
}

// once we're in a float, the only valid values are digits, whitespace or close
// paren.
func lexFloat(l *lexer) stateFn {

	digits := "0123456789"
	l.acceptRun(digits)

	l.emit(_FLOAT)

	r := l.next()

	switch r {
	case ' ', '\t', '\n':
		return lexWhitespace
	case ')':
		return lexCloseParen
	case ';':
		return lexComment
	}

	return l.errorf("unexpected run in lexFloat: %c", r)
}

func lexSymbol(l *lexer) stateFn {

	r := l.peek()

	switch r {
	case ' ', '\t', '\n':
		l.emit(_SYMBOL)
		l.next()
		return lexWhitespace
	case ')':
		l.emit(_SYMBOL)
		l.next()
		return lexCloseParen
	case ';':
		l.emit(_SYMBOL)
		l.next()
		return lexComment
	default:
		l.next()
		return lexSymbol
	}
}

// lex a close parenthesis
func lexCloseParen(l *lexer) stateFn {
	debug("--> lexCloseParen")
	l.emit(_RPAREN)
	l.parenDepth--
	if l.parenDepth < 0 {
		return l.errorf("unexpected close paren")
	}

	r := l.next()
	switch r {
	case ' ', '\t', '\n':
		return lexWhitespace
	case '(':
		return lexOpenParen
	case ')':
		return lexCloseParen
	case ';':
		return lexComment
	}
	return l.errorf("unimplemented")
}

// lexes a comment
func lexComment(l *lexer) stateFn {
	debug("--> lexComment")

	r := l.next()

	switch r {
	case '\n', '\r':
		return lexWhitespace
	}
	return lexComment
}

func lex(input string) *lexer {
	l := &lexer{
		// name:       name,
		input:  input,
		tokens: make(chan token),
	}
	go l.run()
	return l
}

type Pair struct {
	car, cdr interface{}
}

func (p *Pair) String() string {
	return fmt.Sprintf("(%v %v)", p.car, p.cdr)
}

func Cons(car, cdr interface{}) *Pair {
	return &Pair{car, cdr}
}

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
			var v Any
			switch t.typ {
			case _QUOTE:
				nextExp := parse(l, []Any{})
				return append(append(p, []Any{Symbol("quote"), nextExp[0]}), nextExp[1:]...)
			case _INT:
				v, _ = strconv.ParseInt(t.val, 10, 0)
			case _FLOAT:
				v, _ = strconv.ParseFloat(t.val, 64)
			case _STRING:
				v = t.val[1 : len(t.val)-1]
			case _BOOL:
				if t.val == "#t" {
					v = true
				} else {
					v = false
				}
			case _SYMBOL:
				v = Symbol(t.val)
			}
			return parse(l, append(p, v))
		}
	}

	return p
}

func debug(s string) {
	if DEBUG {
		fmt.Println(s)
	}
}

func args() {
	filename := os.Args[1]
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	l := lex(string(b) + "\n")
	p := parse(l, []Any{})

	s := NewRootScope()
	s.EvalAll(p)

}

func main() {
	if len(os.Args) > 1 {
		args()
		return
	}

	r := bufio.NewReader(os.Stdin)
	s := NewRootScope()
	for {
		fmt.Print(">> ")
		line, _, err := r.ReadLine()

		if err != nil {
			fmt.Println("error: ", err)
			continue
		}

		l := lex(string(line) + "\n")
		p := parse(l, []Any{})
		fmt.Println(p)
		res := s.EvalAll(p)
		if len(res) > 0 {
			fmt.Println(res[len(res)-1])
		}
	}
}
