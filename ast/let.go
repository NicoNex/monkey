package ast

import "monkey/token"

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) SNode() {}

func (ls *LetStatement) Literal() string {
	return ls.Token.Lit
}
