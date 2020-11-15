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

func (i *Identifier) String() string {
	return i.Value
}
