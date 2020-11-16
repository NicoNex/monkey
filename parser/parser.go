package parser

import (
	"fmt"
	"strconv"
	"monkey/ast"
	"monkey/token"
)

type Parser struct {
	cur           token.Token
	peek          token.Token
	tokens        chan token.Token
	errors        []string
	prefixParsers map[token.TokenType]parsePrefixFn
	infixParsers  map[token.TokenType]parseInfixFn
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
)

// Links each operator to its precedence class.
var precedences = map[token.TokenType]int{
	token.EQ: EQUALS,
	token.NOT_EQ: EQUALS,
	token.LT: LESSGREATER,
	token.GT: LESSGREATER,
	token.PLUS: SUM,
	token.MINUS: SUM,
	token.SLASH: PRODUCT,
	token.ASTERISK: PRODUCT,
	token.POWER: PRODUCT,
}

func New(tokens chan token.Token) *Parser {
	p := &Parser{
		cur:    <-tokens,
		peek:   <-tokens,
		tokens: tokens,
		prefixParsers: make(map[token.TokenType]parsePrefixFn),
		infixParsers: make(map[token.TokenType]parseInfixFn),
	}
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)

	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.POWER, p.parseInfixExpression)
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

	// skip everything until semicolon (just for the moment)
	for !p.cur.Is(token.SEMICOLON) {
		p.next()
	}

	return s
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	var r = &ast.ReturnStatement{Token: p.cur}
	p.next()

	// skip everything until semicolon (just for the moment)
	if !p.cur.Is(token.SEMICOLON) {
		p.next()
	}
	return r
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	var ret = &ast.ExpressionStatement{
		Token: p.cur,
		Expr: p.parseExpression(LOWEST),
	}

	if p.peek.Is(token.SEMICOLON) {
		p.next()
	}
	return ret
}

func (p *Parser) noParsePrefixFnError(t token.TokenType) {
	msg := fmt.Sprintf("no parse prefix function for '%s' found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	if fn, ok := p.prefixParsers[p.cur.Typ]; ok {
		leftExp := fn()

		for !p.peek.Is(token.SEMICOLON) && precedence < p.peekPrecedence() {
			if infix, ok := p.infixParsers[p.peek.Typ]; ok {
				p.next()
				leftExp = infix(leftExp)
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
		Token: p.cur,
		Operator: p.cur.Lit,
	}
	p.next()
	expr.Right = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token: p.cur,
		Operator: p.cur.Lit,
		Left: left,
	}
	precedence := p.curPrecedence()
	p.next()
	expr.Right = p.parseExpression(precedence)
	return expr
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.cur, Value: p.cur.Lit}
}

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

// TODO: remove this if unused.
func (p *Parser) currentIs(t token.TokenType) bool {
	return p.cur.Typ == t
}

// TODO: remove this if unused.
func (p *Parser) peekIs(t token.TokenType) bool {
	return p.peek.Typ == t
}

// Returns true if the peek token is of type 't'.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peek.Is(t) {
		p.next()
		return true
	}
	p.peekError(t)
	return false
}

// Emits an error if the peek token is not of tipe t.
func (p *Parser) peekError(t token.TokenType) {
	p.errors = append(
		p.errors,
		fmt.Sprintf(
			"expected next token to be %s, got %s instead",
			p.peek.Typ.String(),
			t.String(),
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
func (p *Parser) registerPrefix(typ token.TokenType, fn parsePrefixFn) {
	p.prefixParsers[typ] = fn
}

// Adds fn to the infix parsers table with key 'typ'.
func (p *Parser) registerInfix(typ token.TokenType, fn parseInfixFn) {
	p.infixParsers[typ] = fn
}
