package ast

import (
	"fmt"
	"monkey/token"
	"strings"
)

type CallExpression struct {
	Token token.Token
	Func  Expression
	Args  []Expression
}

func (c *CallExpression) ENode() {}

func (c *CallExpression) Literal() string {
	return c.Token.Lit
}

func (c *CallExpression) String() string {
	var args []string

	for _, a := range c.Args {
		args = append(args, a.String())
	}

	return fmt.Sprintf("%s(%s)", c.Func, strings.Join(args, ", "))
}
