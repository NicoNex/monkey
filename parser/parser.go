package parser

import (
	"fmt"
	"github.com/NicoNex/monkey/ast"
	"github.com/NicoNex/monkey/token"
	"strconv"
)

type Parser struct {
	cur           token.Token
	peek          token.Token
	tokens        chan token.Token
	errors        []string
	prefixParsers map[token.Type]parsePrefixFn
	infixParsers  map[token.Type]parseInfixFn
}

type (
	parsePrefixFn func() ast.Expression
	parseInfixFn  func(ast.Expression) ast.Expression
)

// Operators' precedence classes.
const (
	LOWEST int = iota
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

// Links each operator to its precedence class.
var precedences = map[token.Type]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LT_EQ:    LESSGREATER,
	token.GT_EQ:    LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.POWER:    PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

func New(tokens chan token.Token) *Parser {
	p := &Parser{
		cur:           <-tokens,
		peek:          <-tokens,
		tokens:        tokens,
		prefixParsers: make(map[token.Type]parsePrefixFn),
		infixParsers:  make(map[token.Type]parseInfixFn),
	}
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LT_EQ, p.parseInfixExpression)
	p.registerInfix(token.GT_EQ, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.POWER, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	return p
}

func (p *Parser) next() {
	p.cur = p.peek
	p.peek = <-p.tokens
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) Parse() *ast.Program {
	var prog = new(ast.Program)

	for !p.cur.Is(token.EOF) {
		if s := p.parseStatement(); s != nil {
			prog.Statements = append(prog.Statements, s)
		}
		p.next()
	}

	return prog
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.cur.Typ {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	var s = &ast.LetStatement{Token: p.cur}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	s.Name = &ast.Identifier{p.cur, p.cur.Lit}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.next()
	s.Value = p.parseExpression(LOWEST)
	if p.peek.Is(token.SEMICOLON) {
		p.next()
	}
	return s
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	var r = &ast.ReturnStatement{Token: p.cur}

	p.next()
	r.Value = p.parseExpression(LOWEST)
	if p.peek.Is(token.SEMICOLON) {
		p.next()
	}
	return r
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	var ret = &ast.ExpressionStatement{
		Token: p.cur,
		Expr:  p.parseExpression(LOWEST),
	}

	if p.peek.Is(token.SEMICOLON) {
		p.next()
	}
	return ret
}

func (p *Parser) noParsePrefixFnError(t token.Type) {
	msg := fmt.Sprintf("no parse prefix function for '%s' found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	if fn, ok := p.prefixParsers[p.cur.Typ]; ok {
		leftExp := fn()

		for !p.peek.Is(token.SEMICOLON) && precedence < p.peekPrecedence() {
			if infixFn, ok := p.infixParsers[p.peek.Typ]; ok {
				p.next()
				leftExp = infixFn(leftExp)
			} else {
				break
			}
		}
		return leftExp
	}
	p.noParsePrefixFnError(p.cur.Typ)
	return nil
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.cur,
		Operator: p.cur.Lit,
	}
	p.next()
	expr.Right = p.parseExpression(PREFIX)
	return expr
}

// Returns the expression obtained by parsin an infix expression.
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.cur,
		Operator: p.cur.Lit,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.next()
	expr.Right = p.parseExpression(precedence)
	return expr
}

// Returns an identifier expression.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.cur, Value: p.cur.Lit}
}

// Returns an integer expression.
func (p *Parser) parseIntegerLiteral() ast.Expression {
	var r = &ast.IntegerLiteral{Token: p.cur}

	i, err := strconv.ParseInt(p.cur.Lit, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as an integer", p.cur.Lit)
		p.errors = append(p.errors, msg)
		return nil
	}
	r.Value = i
	return r
}

// Returns a string expression
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.cur, Value: p.cur.Lit}
}

// Returns a boolean expression.
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.cur, Value: p.cur.Is(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.next()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

// Returns the expression resulting from an if expression.
func (p *Parser) parseIfExpression() ast.Expression {
	var expr = &ast.IfExpression{Token: p.cur}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.next()
	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()

	if p.peek.Is(token.ELSE) {
		p.next()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}

// Returns the expression representing the function.
func (p *Parser) parseFunctionLiteral() ast.Expression {
	var expr = &ast.FunctionLiteral{Token: p.cur}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	expr.Params = p.parseFunctionParams()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()
	return expr
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	return &ast.ArrayLiteral{
		Token:    p.cur,
		Elements: p.parseExpressionList(token.RBRACKET),
	}
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	var list []ast.Expression

	if p.peek.Is(end) {
		p.next()
		return list
	}

	p.next()
	list = append(list, p.parseExpression(LOWEST))

	for p.peek.Is(token.COMMA) {
		p.next()
		p.next()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	return &ast.CallExpression{
		Token: p.cur,
		Func:  fn,
		Args:  p.parseExpressionList(token.RPAREN),
	}
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	var exp = &ast.IndexExpression{Token: p.cur, Left: left}

	p.next()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

// Return a *ast.BlockStatement representing a block enclosed in braces.
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	var block = &ast.BlockStatement{Token: p.cur}

	p.next()

	for !p.cur.Is(token.RBRACE) && !p.cur.Is(token.EOF) {
		if s := p.parseStatement(); s != nil {
			block.Statements = append(block.Statements, s)
		}
		p.next()
	}
	return block
}

// Returns a list of *ast.Identifier consisting in the function parameters.
func (p *Parser) parseFunctionParams() []*ast.Identifier {
	var ret []*ast.Identifier

	if p.peek.Is(token.RPAREN) {
		p.next()
		return ret
	}

	p.next()
	ret = append(ret, &ast.Identifier{Token: p.cur, Value: p.cur.Lit})

	for p.peek.Is(token.COMMA) {
		p.next()
		p.next()
		ret = append(ret, &ast.Identifier{Token: p.cur, Value: p.cur.Lit})
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return ret
}

// Returns true if the peek token is of type 't'.
func (p *Parser) expectPeek(t token.Type) bool {
	if p.peek.Is(t) {
		p.next()
		return true
	}
	p.peekError(t)
	return false
}

// Emits an error if the peek token is not of tipe t.
func (p *Parser) peekError(t token.Type) {
	p.errors = append(
		p.errors,
		fmt.Sprintf(
			"expected next token to be %s, got %s instead",
			t.String(),
			p.peek.Typ.String(),
		),
	)
}

// Returns the precedence value of the type of the peek token.
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peek.Typ]; ok {
		return p
	}
	return LOWEST
}

// Returns the precedence value of the type of the current token.
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.cur.Typ]; ok {
		return p
	}
	return LOWEST
}

// Adds fn to the prefix parsers table with key 'typ'.
func (p *Parser) registerPrefix(typ token.Type, fn parsePrefixFn) {
	p.prefixParsers[typ] = fn
}

// Adds fn to the infix parsers table with key 'typ'.
func (p *Parser) registerInfix(typ token.Type, fn parseInfixFn) {
	p.infixParsers[typ] = fn
}
