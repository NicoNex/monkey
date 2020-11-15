package ast

import (
	"testing"
	"monkey/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Typ: token.LET, Lit: "let", Pos: 0},
				Name: &Identifier{
					Token: token.Token{Typ: token.IDENT, Lit: "myVar", Pos: 4},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Typ: token.IDENT, Lit: "anotherVar", Pos: 12},
					Value: "anotherVar",
				},
			},
		},
	}

	if str := program.String(); str != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong, got %q", str)
	}
}
