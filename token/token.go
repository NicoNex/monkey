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
		return fmt.Sprintf("%s %.10q...", t.Typ, t.Lit)
	}
	return fmt.Sprintf("%s %q", t.Typ, t.Lit)
}

type TokenType int

const (
	EOF TokenType = iota
	ILLEGAL

	// Identifiers and literals.
	IDENT // function names, variable names...
	INT   // Integer

	// Operators.
	ASSIGN
	PLUS
	MINUS
	SLASH
	ASTERISK
	POWER
	EQ
	NOT_EQ
	BANG
	LT
	GT

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
	IF
	ELSE
	RETURN
	TRUE
	FALSE
)

// Useful to get the string representation of the type.
var typemap = map[TokenType]string{
	EOF: "EOF",
	ILLEGAL: "ILLEGAL",

	IDENT: "IDENT",
	INT: "INT",

	ASSIGN: "ASSIGN",
	PLUS: "PLUS",
	MINUS: "MINUS",
	SLASH: "SLASH",
	ASTERISK: "ASTERISK",
	POWER: "POWER",
	EQ: "EQ",
	NOT_EQ: "NOT_EQ",
	BANG: "BANG",
	LT: "LT",
	GT: "GT",

	COMMA: "COMMA",
	SEMICOLON: "SEMICOLON",

	LPAREN: "LPAREN",
	RPAREN: "RPAREN",

	LBRACE: "LBRACE",
	RBRACE: "RBRACE",

	FUNCTION: "FUNCTION",
	LET: "LET",
	IF: "IF",
	ELSE: "ELSE",
	RETURN: "RETURN",
	TRUE: "TRUE",
	FALSE: "FALSE",
}

func (t TokenType) String() string {
	return typemap[t]
}

var keywords = map[string]TokenType {
	"fn": FUNCTION,
	"let": LET,
	"if": IF,
	"else": ELSE,
	"return": RETURN,
	"true": TRUE,
	"false": FALSE,
}

func LookupIdent(ident string) TokenType {
	if typ, ok := keywords[ident]; ok {
		return typ
	}
	return IDENT
}
