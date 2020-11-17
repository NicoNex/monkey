package ast

import (
	"fmt"
	"monkey/token"
)

type IfExpression struct {
	Token token.Token
	Condition Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) ENode() {}

func (i *IfExpression) Literal() string {
	return i.Token.Lit
}

func (i *IfExpression) String() string {
	if i.Alternative == nil {
		return fmt.Sprintf("if %s %s", i.Condition, i.Consequence)
	}
	return fmt.Sprintf("if %s %s else %s", i.Condition, i.Consequence, i.Alternative)
}
