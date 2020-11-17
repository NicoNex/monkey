package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
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
		t.Errorf("il not *ast.IntegerLiteral, got %T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d, got %d", value, integ.Value)
		return false
	}

	if integ.Literal() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d, got %s", value,
			integ.Literal())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier, got %T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s, got %s", value, ident.Value)
		return false
	}

	if ident.Literal() != value {
		t.Errorf("ident.TokenLiteral not %s, got %s", value, ident.Literal())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean, got %T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t, got %t", value, bo.Value)
		return false
	}

	if bo.Literal() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t, got %s", value, bo.Literal())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
		operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression, got %T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s', got %q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled, got %T", exp)
	return false
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

func TestBooleanExpression(t *testing.T) {
	input := `true;`

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

	literal, ok := stmt.Expr.(*ast.Boolean)
	if !ok {
		t.Fatalf("exp not *ast.Boolean, got %T", stmt.Expr)
	}

	if literal.Value != true {
		t.Errorf("literal.Value not %t, got %t", true, literal.Value)
	}

	if literal.Literal() != "true" {
		t.Errorf("literal.Literal() not %s, got %s", "5", literal.Literal())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		// {"!foobar;", "!", "foobar"},
		// {"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
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
			t.Fatalf("stmt is not ast.PrefixExpression, got %T", stmt.Expr)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s', got %s", tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
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
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		toks := lexer.Lex(tt.input)
		p := New(toks)
		program := p.Parse()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements, got %d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expr, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
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
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
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

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	toks := lexer.Lex(input)
	p := New(toks)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got %d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T",
			program.Statements[0])
	}

	exp, ok := stmt.Expr.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expr is not ast.IfExpression, got %T", stmt.Expr)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements, got %d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement, got %T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expr, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil, got %+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	toks := lexer.Lex(input)
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

	exp, ok := stmt.Expr.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expr is not ast.IfExpression. got=%T", stmt.Expr)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expr, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expr, "y") {
		return
	}
}
