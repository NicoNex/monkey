package ast

import (
	"bytes"
	"monkey/token"
)

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (r *ReturnStatement) Literal() string {
	return r.Token.Lit
}

func (r *ReturnStatement) SNode() {}

func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.Literal() + " ")
	if r.Value != nil {
		out.WriteString(r.Value.String())
	}
	out.WriteString(";")
	return out.String()
}