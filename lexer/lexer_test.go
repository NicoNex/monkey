package lexer

import (
	"testing"

	"monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expTyp token.TokenType
		expLit string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}

	tokens := Lex(input)

	i := 0
	for tok := range tokens {
		t.Log(tok)
		if i >= len(tests) {
			break
		}

		tt := tests[i]
		if tok.Typ != tt.expTyp {
			t.Fatalf("tests[%d] - wrong token type: expected=%s, got=%s", i, tt.expTyp, tok.Typ)
		}
		if tok.Lit != tt.expLit {
			t.Fatalf("tests[%d] - wrong token literal: expected=%q, got=%q", i, tt.expLit, tok.Lit)
		}
		i++
	}
}
