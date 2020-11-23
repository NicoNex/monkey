package ast

import (
	"fmt"
	"github.com/NicoNex/monkey/token"
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
	if r.Value != nil {
		fmt.Sprintf("%s %s;", r.Literal(), r.Value)
	}
	return fmt.Sprintf("%s;", r.Literal())
}
