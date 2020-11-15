package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
)

type Parser struct {
	cur    token.Token
	peek   token.Token
	tokens chan token.Token
	errors []string
}

func New(tokens chan token.Token) *Parser {
	return &Parser{
		cur:    <-tokens,
		peek:   <-tokens,
		tokens: tokens,
	}
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
		return nil
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
