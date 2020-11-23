package ast

import "github.com/NicoNex/monkey/token"

type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) ENode() {}

func (s *StringLiteral) Literal() string {
	return s.Token.Lit
}

func (s *StringLiteral) String() string {
	return s.Token.Lit
}
