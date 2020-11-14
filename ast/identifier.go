package ast

import "monkey/token"

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) ENode() {}

func (i *Identifier) Literal() string {
	return i.Token.Lit
}
