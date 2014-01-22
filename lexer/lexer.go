package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Pos int

type Item struct {
	Typ   ItemType
	Pos   Pos
	Value string
}

type ItemType int

const (
	ItemError ItemType = iota
	ItemEOF

	ItemLeftParen
	ItemRightParen
	ItemLeftVect
	ItemRightVect

	ItemIdent
	ItemString
	ItemChar
	ItemFloat
	ItemInt
	ItemComplex

	ItemQuote
	ItemQuasiQuote
	ItemUnquote
	ItemUnquoteSplice
)

const EOF = -1

type stateFn func(*Lexer) stateFn

type Lexer struct {
	name    string
	input   string
	state   stateFn
	pos     Pos
	start   Pos
	width   Pos
	lastPos Pos
	items   chan Item

	parenDepth int
	vectDepth  int
}

// next returns the next rune in the input.
func (l *Lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return EOF
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *Lexer) backup() {
	l.pos -= l.width
}

// emit passes an Item back to the client.
func (l *Lexer) emit(t ItemType) {
	l.items <- Item{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *Lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- Item{ItemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

func (l *Lexer) NextItem() Item {
	item := <-l.items
	l.lastPos = item.Pos
	return item
}

func Lex(name, input string) *Lexer {
	l := &Lexer{
		name:  name,
		input: input,
		items: make(chan Item),
	}
	go l.run()
	return l
}

func (l *Lexer) run() {
	for l.state = lexWhitespace; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

func lexLeftVect(l *Lexer) stateFn {
	l.emit(ItemLeftVect)

	return lexWhitespace
}

func lexRightVect(l *Lexer) stateFn {
	l.emit(ItemRightVect)

	return lexWhitespace
}

// lexes an open parenthesis
func lexLeftParen(l *Lexer) stateFn {
	l.emit(ItemLeftParen)

	return lexWhitespace
}

func lexWhitespace(l *Lexer) stateFn {
	atLeastOne := false
	for r := l.next(); isSpace(r) || r == '\n'; {
		// just absorb the spaces
		atLeastOne = true
	}

	if atLeastOne {
		l.backup()
		l.ignore()
	}

	switch r := l.next(); {
	case r == EOF:
		return nil
	case r == '(':
		return lexLeftParen
	case r == ')':
		return lexRightParen
	case r == '[':
		return lexLeftVect
	case r == ']':
		return lexRightVect
	case r == '"':
		return lexString
	case r == '+' || r == '-' || ('0' <= r && r <= '9'):
		return lexNumber
	case r == ';':
		return lexComment
	case isAlphaNumeric(r):
		return lexIdentifier
	}

	return lexWhitespace
}

func lexString(l *Lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != EOF {
				break
			}
			fallthrough
		case EOF:
			return l.errorf("unterminated quoted string")
		case '"':
			break Loop
		}
	}

	l.emit(ItemString)
	return lexWhitespace
}

func lexIdentifier(l *Lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb it!
		default:
			l.backup()
			break Loop
		}
	}

	l.emit(ItemIdent)

	return lexWhitespace
}

// lex a close parenthesis
func lexRightParen(l *Lexer) stateFn {
	l.emit(ItemRightParen)

	return lexWhitespace
}

// lex a comment, comment delimiter is known to be already read
func lexComment(l *Lexer) stateFn {
	i := strings.Index(l.input[l.pos:], "\n")
	l.pos += Pos(i)
	l.ignore()
	return lexWhitespace
}

func lexNumber(l *Lexer) stateFn {
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}

	if dot := l.peek(); dot == '.' { // it's a float
		if !l.scanNumber() {
			return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
		}
		l.emit(ItemFloat)
	} else if sign := l.peek(); sign == '+' || sign == '-' {
		// Complex: 1+2i. No spaces, must end in 'i'.
		if !l.scanNumber() || l.input[l.pos-1] != 'i' {
			return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
		}
		l.emit(ItemComplex)
	} else {
		l.emit(ItemInt)
	}

	return lexWhitespace
}

func (l *Lexer) scanNumber() bool {
	// Optional leading sign.
	l.accept("+-")
	// Is it hex?
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)

	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	// Is it imaginary?
	l.accept("i")
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '-' || r == '.' || r == '/' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
