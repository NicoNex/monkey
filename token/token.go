package token

import "fmt"

type Type int

type Token struct {
	Typ Type
	Lit string
	Pos int
}

const (
	EOF Type = iota
	ILLEGAL

	// Identifiers and literals.
	IDENT // function names, variable names...
	INT   // Integer
	STRING

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
	LT_EQ
	GT_EQ

	// Delimiters.
	COMMA
	SEMICOLON
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET

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
var typemap = map[Type]string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",

	IDENT:  "IDENT",
	INT:    "INT",
	STRING: "STRING",

	ASSIGN:   "=",
	PLUS:     "+",
	MINUS:    "-",
	SLASH:    "/",
	ASTERISK: "*",
	POWER:    "**",
	EQ:       "==",
	NOT_EQ:   "!=",
	BANG:     "!",
	LT:       ">",
	GT:       "<",
	LT_EQ:    "<=",
	GT_EQ:    ">=",

	COMMA:     ",",
	SEMICOLON: ";",

	LPAREN:   "(",
	RPAREN:   ")",
	LBRACE:   "{",
	RBRACE:   "}",
	LBRACKET: "[",
	RBRACKET: "]",

	FUNCTION: "FUNCTION",
	LET:      "LET",
	IF:       "IF",
	ELSE:     "ELSE",
	RETURN:   "RETURN",
	TRUE:     "TRUE",
	FALSE:    "FALSE",
}

var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
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

func (t Token) Is(tt Type) bool {
	return t.Typ == tt
}

func (t Type) String() string {
	return typemap[t]
}

func LookupIdent(ident string) Type {
	if typ, ok := keywords[ident]; ok {
		return typ
	}
	return IDENT
}
