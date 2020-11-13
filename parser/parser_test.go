package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	tokens := lexer.Lex(input)
	p := New(tokens)
	program := p.Parse()

	if program == nil {
		t.Fatal("Parse() returned nil")
	}
	if l := len(program.Statements); l != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got %d", l)
	}

	tests := []struct {
		expIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if lit := s.Literal(); lit != "let" {
		t.Errorf("s.Literal not 'let', got %q", lit)
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement, got %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s', got '%s'", name, letStmt.Name.Value)
		return false
	}

	if lit := letStmt.Name.Literal(); lit != name {
		t.Errorf("letStmt.Name.Literal() not '%s', got '%s'", name, lit)
		return false
	}

	return true
}
