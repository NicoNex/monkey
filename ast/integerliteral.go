package ast

import "github.com/NicoNex/monkey/token"

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) ENode() {}

func (i *IntegerLiteral) Literal() string {
	return i.Token.Lit
}

func (i *IntegerLiteral) String() string {
	return i.Token.Lit
}
