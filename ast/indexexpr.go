package ast

import (
	"fmt"
	"github.com/NicoNex/monkey/token"
)

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (i *IndexExpression) ENode() {}

func (i *IndexExpression) Literal() string {
	return i.Token.Lit
}

func (i *IndexExpression) String() string {
	return fmt.Sprintf("(%s[%s])", i.Left.String(), i.Index.String())
}
