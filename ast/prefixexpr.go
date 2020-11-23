package ast

import (
	"fmt"
	"github.com/NicoNex/monkey/token"
)

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) ENode() {}

func (p *PrefixExpression) Literal() string {
	return p.Token.Lit
}

func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", p.Operator, p.Right.String())
}
