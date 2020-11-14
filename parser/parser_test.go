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

let b 15;
let = 10;
let 666;
`

	tokens := lexer.Lex(input)
	p := New(tokens)
	program := p.Parse()
	checkParserErrors(t, p)

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

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 12345;
`
	toks := lexer.Lex(input)
	p := New(toks)
	program := p.Parse()
	checkParserErrors(t, p)

	if l := len(program.Statements); l != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got %d", l)
	}

	for _, stmt := range program.Statements {
		retStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement, got %T", stmt)
			continue
		}
		if lit := retStmt.Literal(); lit != "return" {
			t.Errorf("retStmt.Literal not 'return', got %q", lit)
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	var errs = p.Errors()
	if len(errs) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errs))
	for _, msg := range errs {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
