package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/NicoNex/monkey/token"
)

type lexer struct {
	input  string
	start  int
	pos    int
	width  int
	tokens chan token.Token
}

type stateFn func(*lexer) stateFn

func (l *lexer) next() rune {
	var r rune
	if l.pos >= len(l.input) {
		l.width = 0
		return 0
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// Consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// Consumes all the runes if they're in the valid set.
func (l *lexer) acceptRun(valid string) bool {
	for strings.IndexRune(valid, l.next()) >= 0 {

	}
	l.backup()
	return true
}

func (l *lexer) emit(t token.TokenType) {
	l.tokens <- token.Token{
		Typ: t,
		Lit: l.input[l.start:l.pos],
		Pos: l.start,
	}
	l.start = l.pos
}

func (l *lexer) current() string {
	return l.input[l.start:l.pos]
}

func (l *lexer) errorf(format string, args ...interface{}) {
	l.tokens <- token.Token{
		token.ILLEGAL,
		fmt.Sprintf(format, args...),
		l.start,
	}
	l.start = l.pos
}

func (l *lexer) run() {
	for state := lexExpression; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func lexOperator(l *lexer) stateFn {
	switch r := l.next(); {
	case r == '+':
		l.emit(token.PLUS)
	case r == '-':
		l.emit(token.MINUS)
	case r == '*':
		if l.next() == '*' {
			l.emit(token.POWER)
		} else {
			l.backup()
			l.emit(token.ASTERISK)
		}
	case r == '/':
		l.emit(token.SLASH)
	case r == '=':
		if l.next() == '=' {
			l.emit(token.EQ)
		} else {
			l.backup()
			l.emit(token.ASSIGN)
		}
	case r == '!':
		if l.next() == '=' {
			l.emit(token.NOT_EQ)
		} else {
			l.backup()
			l.emit(token.BANG)
		}
	case r == '<':
		if l.next() == '=' {
			l.emit(token.LT_EQ)
		} else {
			l.backup()
			l.emit(token.LT)
		}
	case r == '>':
		if l.next() == '=' {
			l.emit(token.GT_EQ)
		} else {
			l.backup()
			l.emit(token.GT)
		}
	default:
		l.errorf("illegal operator: %q", r)
		return nil
	}
	return lexExpression
}

func lexIdentifier(l *lexer) stateFn {
	var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	if l.acceptRun(chars) {
		l.emit(token.LookupIdent(l.current()))
	}
	return lexExpression
}

func lexNumber(l *lexer) stateFn {
	var digits = "0123456789"

	// Optional leading sign
	l.accept("+-")

	// Is it hex?
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}

	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}

	if l.accept("eE") {
		l.accept("+-")
		l.accept("0123456789")
	}

	l.emit(token.INT)
	return lexExpression
}

func lexString(l *lexer) stateFn {
	if l.peek() == '"' {
		l.emit(token.STRING)
		l.next()
		l.ignore()
		return lexExpression
	}
	l.next()
	return lexString
}

func lexExpression(l *lexer) stateFn {
	switch r := l.next(); {

	case isSpace(r):
		l.ignore()

	case isOperator(r):
		l.backup()
		return lexOperator

	case isLetter(r):
		l.backup()
		return lexIdentifier

	case r == '"':
		l.ignore()
		return lexString

	case r == ';':
		l.emit(token.SEMICOLON)

	case r == '(':
		l.emit(token.LPAREN)

	case r == ')':
		l.emit(token.RPAREN)

	case r == ',':
		l.emit(token.COMMA)

	case r == '{':
		l.emit(token.LBRACE)

	case r == '}':
		l.emit(token.RBRACE)

	case r == 0:
		l.emit(token.EOF)
		return nil

	default:
		if isNumber(r) {
			l.backup()
			return lexNumber
		}
		l.errorf("lexer: invalid token %q", r)
	}
	return lexExpression
}

func isLetter(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isOperator(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/' || r == '^' ||
		r == '=' || r == '!' || r == '<' || r == '>'
}

func isNumber(r rune) bool {
	return r == '+' || r == '-' || unicode.IsNumber(r)
}

func Lex(in string) chan token.Token {
	l := &lexer{
		input:  in,
		tokens: make(chan token.Token),
	}
	go l.run()
	return l.tokens
}
