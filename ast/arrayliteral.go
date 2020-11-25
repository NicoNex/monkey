package ast

import (
	"fmt"
	"github.com/NicoNex/monkey/token"
	"strings"
)

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (a *ArrayLiteral) ENode() {}

func (a *ArrayLiteral) Literal() string {
	return a.Token.Lit
}

func (a *ArrayLiteral) String() string {
	var elements []string

	for _, e := range a.Elements {
		elements = append(elements, e.String())
	}

	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}
