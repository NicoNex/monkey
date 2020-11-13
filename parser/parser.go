package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
	// "monkey/lexer"
)

type Parser struct {
	cur    token.Token
	peek   token.Token
	tokens chan token.Token
}

func New(tokens chan token.Token) *Parser {
	return &Parser{<-tokens, <-tokens, tokens}
}

func (p *Parser) next() {
	p.cur = p.peek
	p.peek = <-p.tokens
}

func (p *Parser) Parse() *ast.Program {
	var prog *ast.Program

	for p.cur.Typ != token.EOF {
		fmt.Println(p.cur)
		if s := p.parseStatement(); s != nil {
			fmt.Println(s)
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
	return false
}
