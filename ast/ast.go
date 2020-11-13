package ast

import "monkey/token"

type Node interface {
	Literal() string
}

type Statement interface {
	Node
	SNode()
}

type Expression interface {
	Node
	ENode()
}


type Program struct {
	Statements []Statement
}

func (p *Program) Literal() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].Literal()
	}
	return ""
}


type LetStatement struct {
	Token token.Token
	Name *Identifier
	Value Expression
}

func (ls *LetStatement) SNode() {}

func (ls *LetStatement) Literal() string {
	return ls.Token.Lit
}


type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) ENode() {}

func (i *Identifier) Literal() string {
	return i.Token.Lit
}
