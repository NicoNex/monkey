package ast

import "github.com/NicoNex/monkey/token"

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) ENode() {}

func (b *Boolean) Literal() string {
	return b.Token.Lit
}

func (b *Boolean) String() string {
	return b.Token.Lit
}
