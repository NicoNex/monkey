package ast

import (
	"fmt"
	"strings"
	"monkey/token"
)

type FunctionLiteral struct {
	Token token.Token
	Params []*Identifier
	Body *BlockStatement
}

func (f *FunctionLiteral) ENode() {}

func (f *FunctionLiteral) Literal() string {
	return f.Token.Lit
}

func (f *FunctionLiteral) String() string {
	var params []string

	for _, p := range f.Params {
		params = append(params, p.String())
	}

	return fmt.Sprintf(
		"%s(%s) %s",
		f.Literal(),
		strings.Join(params, ", "),
		f.Body,
	)
}
