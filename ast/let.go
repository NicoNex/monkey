package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
)

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) SNode() {}

func (ls *LetStatement) Literal() string {
	return ls.Token.Lit
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf(
		"%s %s = ",
		ls.Literal(),
		ls.Name.String(),
	))

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}
