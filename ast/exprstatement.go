package ast

import "monkey/token"

type ExpressionStatement struct {
	Token token.Token
	Value Expression
}

func (e *ExpressionStatement) SNode() {}

func (e *ExpressionStatement) Literal() string {
	return e.Token.Lit
}

func (e *ExpressionStatement) String() string {
	if e.Value != nil {
		return e.Value.String()
	}
	return ""
}
