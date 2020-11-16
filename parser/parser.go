package parser

import (
	"fmt"
	// "strconv"
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

// Operators' precedences.
const (
	LOWEST int = iota
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

func New(tokens chan token.Token) *Parser {
	p := &Parser{
		cur:    <-tokens,
		peek:   <-tokens,
		tokens: tokens,
		prefixParsers: make(map[token.TokenType]parsePrefixFn),
		infixParsers: make(map[token.TokenType]parseInfixFn),
	}
	p.registerPrefix(token.IDENT, p.parseIdentifier)
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

func (p *Parser) parseExpression(precedence int) ast.Expression {
	if fn, ok := p.prefixParsers[p.cur.Typ]; ok {
		leftExp := fn()
		return leftExp
	}
	return nil
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.cur, Value: p.cur.Lit}
}

// TODO: remove this if unused.
func (p *Parser) currentIs(t token.TokenType) bool {
	return p.cur.Typ == t
}

// TODO: remove this if unused.
func (p *Parser) peekIs(t token.TokenType) bool {
	return p.peek.Typ == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peek.Is(t) {
		p.next()
		return true
	}
	p.peekError(t)
	return false
}

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

func (p *Parser) registerPrefix(typ token.TokenType, fn parsePrefixFn) {
	p.prefixParsers[typ] = fn
}

func (p *Parser) registerInfix(typ token.TokenType, fn parseInfixFn) {
	p.infixParsers[typ] = fn
}
