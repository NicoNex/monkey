package ast

import "monkey/token"

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (r *ReturnStatement) Literal() string {
	return r.Token.Lit
}

func (r *ReturnStatement) SNode() {}
