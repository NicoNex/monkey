package token

import "fmt"

type Token struct {
	Typ TokenType
	Lit string
	Pos int
}

func (t Token) String() string {
	switch t.Typ {
	case ILLEGAL:
		return t.Lit
	}
	if len(t.Lit) > 10 {
		return fmt.Sprintf("%.10q...", t.Lit)
	}
	return fmt.Sprintf("%q", t.Lit)
}

type TokenType int

const (
	EOF TokenType = iota
	ILLEGAL

	// Identifiers and literals.
	IDENT // function names, variable names...
	INT   // Integer

	// Operators.
	ASSIGN // = statement
	PLUS
	MINUS
	DIVIDE
	TIMES
	POWER

	// Delimiters.
	COMMA
	SEMICOLON

	LPAREN
	RPAREN

	LBRACE
	RBRACE

	// Keywords.
	FUNCTION
	LET
)

var typemap = map[TokenType]string{
	EOF:       "EOF",
	ILLEGAL:   "ILLEGAL",
	IDENT:     "IDENT",
	INT:       "INT",
	ASSIGN:    "ASSIGN",
	PLUS:      "PLUS",
	COMMA:     "COMMA",
	SEMICOLON: "SEMICOLON",
	LPAREN:    "LPAREN",
	RPAREN:    "RPAREN",
	LBRACE:    "LBRACE",
	RBRACE:    "RBRACE",
	FUNCTION:  "FUNCTION",
	LET:       "LET",
}

func (t TokenType) String() string {
	return typemap[t]
}
