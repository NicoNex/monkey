package parser

import (
	"fmt"
	"testing"
	"monkey/ast"
	"monkey/lexer"
)

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

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.Literal() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.Literal())
		return false
	}

	return true
}

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;`

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

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	toks := lexer.Lex(input)
	p := New(toks)
	prog := p.Parse()
	checkParserErrors(t, p)

	if l := len(prog.Statements); l != 1 {
		t.Fatalf("program has not enough statements, got %d", l)
	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("prog.Statements[0] is not an ast.ExpressionStatement, got %T", prog.Statements[0])
	}

	ident, ok := stmt.Expr.(*ast.Identifier)
	if !ok {
		t.Fatalf("expr not *ast.Identifier, got %T", stmt.Expr)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s, got %s", "foobar", ident.Value)
	}

	if lit := ident.Literal(); lit != "foobar" {
		t.Errorf("ident.Literal not %s, got %s", "foobar", lit)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`

	tokens := lexer.Lex(input)
	p := New(tokens)
	prog := p.Parse()
	checkParserErrors(t, p)

	if l := len(prog.Statements); l != 1 {
		t.Fatalf("program has not enough statements, got %d", l)
	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement, got %T", prog.Statements[0])
	}

	literal, ok := stmt.Expr.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral, got %T", stmt.Expr)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d, got %d", 5, literal.Value)
	}

	if literal.Literal() != "5" {
		t.Errorf("literal.Literal() not %s, got %s", "5", literal.Literal())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		toks := lexer.Lex(tt.input)
		p := New(toks)
		program := p.Parse()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got %d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T",
				program.Statements[0])
		}

		exp, ok := stmt.Expr.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expr)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 ** 5;", 5, "**", 5},
		// {"foobar + barfoo;", "foobar", "+", "barfoo"},
		// {"foobar - barfoo;", "foobar", "-", "barfoo"},
		// {"foobar * barfoo;", "foobar", "*", "barfoo"},
		// {"foobar / barfoo;", "foobar", "/", "barfoo"},
		// {"foobar > barfoo;", "foobar", ">", "barfoo"},
		// {"foobar < barfoo;", "foobar", "<", "barfoo"},
		// {"foobar == barfoo;", "foobar", "==", "barfoo"},
		// {"foobar != barfoo;", "foobar", "!=", "barfoo"},
		// {"true == true", true, "==", true},
		// {"true != false", true, "!=", false},
		// {"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		toks := lexer.Lex(tt.input)
		p := New(toks)
		program := p.Parse()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expr.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression, got %T", stmt.Expr)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s', got '%s'", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}

		// if !testInfixExpression(t, stmt.Expr, tt.leftValue, tt.operator, tt.rightValue) {
		// 	return
		// }
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		// {
		// 	"true",
		// 	"true",
		// },
		// {
		// 	"false",
		// 	"false",
		// },
		// {
		// 	"3 > 5 == false",
		// 	"((3 > 5) == false)",
		// },
		// {
		// 	"3 < 5 == true",
		// 	"((3 < 5) == true)",
		// },
		// {
		// 	"1 + (2 + 3) + 4",
		// 	"((1 + (2 + 3)) + 4)",
		// },
		// {
		// 	"(5 + 5) * 2",
		// 	"((5 + 5) * 2)",
		// },
		// {
		// 	"2 / (5 + 5)",
		// 	"(2 / (5 + 5))",
		// },
		// {
		// 	"(5 + 5) * 2 * (5 + 5)",
		// 	"(((5 + 5) * 2) * (5 + 5))",
		// },
		// {
		// 	"-(5 + 5)",
		// 	"(-(5 + 5))",
		// },
		// {
		// 	"!(true == true)",
		// 	"(!(true == true))",
		// },
		// {
		// 	"a + add(b * c) + d",
		// 	"((a + add((b * c))) + d)",
		// },
		// {
		// 	"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
		// 	"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		// },
		// {
		// 	"add(a + b + c * d / f + g)",
		// 	"add((((a + b) + ((c * d) / f)) + g))",
		// },
	}

	for _, tt := range tests {
		toks := lexer.Lex(tt.input)
		p := New(toks)
		program := p.Parse()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
