package ast

import "github.com/NicoNex/monkey/token"

type ExpressionStatement struct {
	Token token.Token
	Expr  Expression
}

func (e *ExpressionStatement) SNode() {}

func (e *ExpressionStatement) Literal() string {
	return e.Token.Lit
}

func (e *ExpressionStatement) String() string {
	if e.Expr != nil {
		return e.Expr.String()
	}
	return ""
}
