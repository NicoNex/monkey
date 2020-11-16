package ast

import (
	"fmt"
	"monkey/token"
)

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) ENode() {}

func (i *InfixExpression) Literal() string {
	return i.Token.Lit
}

func (i *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left, i.Operator, i.Right)
}
